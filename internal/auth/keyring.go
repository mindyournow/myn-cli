package auth

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"

	"github.com/zalando/go-keyring"
)

const (
	serviceName = "mynow"
	apiKeyAccount = "api-key"
	refreshTokenAccount = "refresh-token"
)

// Keyring provides credential storage using the OS keyring or encrypted file fallback.
type Keyring struct {
	useFileFallback bool
	filePath        string
}

// NewKeyring creates a new Keyring instance.
func NewKeyring() *Keyring {
	k := &Keyring{}
	// Test if keyring is available
	if err := k.testKeyring(); err != nil {
		k.useFileFallback = true
		k.filePath = k.defaultFilePath()
	}
	return k
}

// SaveAPIKey stores an API key.
func (k *Keyring) SaveAPIKey(key string) error {
	if k.useFileFallback {
		return k.saveToFile("api-key", key)
	}
	return keyring.Set(serviceName, apiKeyAccount, key)
}

// LoadAPIKey retrieves the stored API key.
func (k *Keyring) LoadAPIKey() (string, error) {
	if k.useFileFallback {
		return k.loadFromFile("api-key")
	}
	return keyring.Get(serviceName, apiKeyAccount)
}

// SaveRefreshToken stores an OAuth refresh token.
func (k *Keyring) SaveRefreshToken(token string) error {
	if k.useFileFallback {
		return k.saveToFile("refresh-token", token)
	}
	return keyring.Set(serviceName, refreshTokenAccount, token)
}

// LoadRefreshToken retrieves the stored refresh token.
func (k *Keyring) LoadRefreshToken() (string, error) {
	if k.useFileFallback {
		return k.loadFromFile("refresh-token")
	}
	return keyring.Get(serviceName, refreshTokenAccount)
}

// Clear removes all stored credentials.
func (k *Keyring) Clear() error {
	var errs []error

	if k.useFileFallback {
		_ = os.Remove(k.filePath)
	} else {
		_ = keyring.Delete(serviceName, apiKeyAccount)
		_ = keyring.Delete(serviceName, refreshTokenAccount)
	}

	if len(errs) > 0 {
		return fmt.Errorf("failed to clear some credentials")
	}
	return nil
}

// testKeyring checks if the system keyring is available.
func (k *Keyring) testKeyring() error {
	testValue := "test" + randomString(8)
	if err := keyring.Set(serviceName, "test-account", testValue); err != nil {
		return err
	}
	_ = keyring.Delete(serviceName, "test-account")
	return nil
}

// defaultFilePath returns the path to the encrypted credentials file.
func (k *Keyring) defaultFilePath() string {
	if dir := os.Getenv("XDG_CONFIG_HOME"); dir != "" {
		return filepath.Join(dir, "mynow", "credentials.enc")
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "mynow", "credentials.enc")
}

// saveToFile encrypts and saves data to the fallback file.
func (k *Keyring) saveToFile(key, value string) error {
	data := map[string]string{key: value}

	// Load existing data if present
	if existing, err := k.loadAllFromFile(); err == nil {
		for k, v := range existing {
			if _, ok := data[k]; !ok {
				data[k] = v
			}
		}
	}

	// Serialize and encrypt
	plaintext := ""
	for k, v := range data {
		if plaintext != "" {
			plaintext += "\n"
		}
		plaintext += k + "=" + base64.StdEncoding.EncodeToString([]byte(v))
	}

	encrypted, err := k.encrypt([]byte(plaintext))
	if err != nil {
		return fmt.Errorf("encrypting credentials: %w", err)
	}

	// Ensure directory exists
	_ = os.MkdirAll(filepath.Dir(k.filePath), 0755)

	if err := os.WriteFile(k.filePath, encrypted, 0600); err != nil {
		return fmt.Errorf("writing credentials file: %w", err)
	}

	return nil
}

// loadFromFile loads a specific key from the encrypted file.
func (k *Keyring) loadFromFile(key string) (string, error) {
	data, err := k.loadAllFromFile()
	if err != nil {
		return "", err
	}
	val, ok := data[key]
	if !ok {
		return "", fmt.Errorf("key not found: %s", key)
	}
	return val, nil
}

// loadAllFromFile decrypts and loads all data from the fallback file.
func (k *Keyring) loadAllFromFile() (map[string]string, error) {
	encrypted, err := os.ReadFile(k.filePath)
	if err != nil {
		return nil, err
	}

	plaintext, err := k.decrypt(encrypted)
	if err != nil {
		return nil, fmt.Errorf("decrypting credentials: %w", err)
	}

	data := make(map[string]string)
	lines := splitLines(string(plaintext))
	for _, line := range lines {
		if idx := indexByte(line, '='); idx > 0 {
			key := line[:idx]
			encoded := line[idx+1:]
			if decoded, err := base64.StdEncoding.DecodeString(encoded); err == nil {
				data[key] = string(decoded)
			}
		}
	}

	return data, nil
}

// encrypt encrypts plaintext using AES-256-GCM with a key derived from machine ID.
func (k *Keyring) encrypt(plaintext []byte) ([]byte, error) {
	key := k.deriveKey()
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nil
}

// decrypt decrypts ciphertext using AES-256-GCM.
func (k *Keyring) decrypt(ciphertext []byte) ([]byte, error) {
	key := k.deriveKey()
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}

// deriveKey derives an encryption key from machine-specific data.
func (k *Keyring) deriveKey() []byte {
	// Try to get a machine-specific identifier
	machineID := ""
	if id, err := os.ReadFile("/etc/machine-id"); err == nil {
		machineID = string(id)
	} else if id, err := os.ReadFile("/var/lib/dbus/machine-id"); err == nil {
		machineID = string(id)
	} else {
		// Fallback to hostname + home dir
		hostname, _ := os.Hostname()
		home, _ := os.UserHomeDir()
		machineID = hostname + home
	}

	hash := sha256.Sum256([]byte(machineID + "mynow-salt-v1"))
	return hash[:]
}

// randomString generates a random alphanumeric string.
func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	rand.Read(b)
	for i := range b {
		b[i] = letters[int(b[i])%len(letters)]
	}
	return string(b)
}

// splitLines splits a string into lines.
func splitLines(s string) []string {
	var lines []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			lines = append(lines, s[start:i])
			start = i + 1
		}
	}
	if start < len(s) {
		lines = append(lines, s[start:])
	}
	return lines
}

// indexByte returns the index of the first occurrence of c in s, or -1.
func indexByte(s string, c byte) int {
	for i := 0; i < len(s); i++ {
		if s[i] == c {
			return i
		}
	}
	return -1
}
