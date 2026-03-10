package auth

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"golang.org/x/crypto/pbkdf2"
)

const (
	// Encryption settings for PBKDF2
	saltSize       = 32
	keySize        = 32
	iterations     = 100000
	credentialsDir = "credentials"
	tokenFile      = "refresh_token.enc"
)

// Keyring provides secure storage for refresh tokens.
// It uses PBKDF2 for key derivation and AES-GCM for encryption.
type Keyring struct {
	configDir string
}

// NewKeyring creates a new Keyring instance.
func NewKeyring(configDir string) *Keyring {
	return &Keyring{configDir: configDir}
}

// SaveRefreshToken encrypts and saves the refresh token to disk.
func (k *Keyring) SaveRefreshToken(token string) error {
	if token == "" {
		return fmt.Errorf("token cannot be empty")
	}

	// Ensure credentials directory exists with restricted permissions
	credDir := k.credDir()
	if err := os.MkdirAll(credDir, 0700); err != nil {
		return fmt.Errorf("failed to create credentials directory: %w", err)
	}

	// Derive encryption key using PBKDF2
	key, salt, err := k.deriveKey()
	if err != nil {
		return fmt.Errorf("failed to derive encryption key: %w", err)
	}

	// Encrypt the token
	encrypted, err := encrypt(token, key)
	if err != nil {
		return fmt.Errorf("failed to encrypt token: %w", err)
	}

	// Store salt + encrypted data
	tokenData := &tokenFileData{
		Salt:      base64.StdEncoding.EncodeToString(salt),
		Encrypted: base64.StdEncoding.EncodeToString(encrypted),
	}

	data, err := json.Marshal(tokenData)
	if err != nil {
		return fmt.Errorf("failed to marshal token data: %w", err)
	}

	tokenPath := filepath.Join(credDir, tokenFile)
	if err := os.WriteFile(tokenPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write token file: %w", err)
	}

	return nil
}

// LoadRefreshToken loads and decrypts the refresh token from disk.
func (k *Keyring) LoadRefreshToken() (string, error) {
	tokenPath := filepath.Join(k.credDir(), tokenFile)

	data, err := os.ReadFile(tokenPath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("no refresh token found")
		}
		return "", fmt.Errorf("failed to read token file: %w", err)
	}

	var tokenData tokenFileData
	if err := json.Unmarshal(data, &tokenData); err != nil {
		return "", fmt.Errorf("failed to unmarshal token data: %w", err)
	}

	salt, err := base64.StdEncoding.DecodeString(tokenData.Salt)
	if err != nil {
		return "", fmt.Errorf("failed to decode salt: %w", err)
	}

	encrypted, err := base64.StdEncoding.DecodeString(tokenData.Encrypted)
	if err != nil {
		return "", fmt.Errorf("failed to decode encrypted data: %w", err)
	}

	// Derive key with the stored salt
	secret, err := k.machineSecret()
	if err != nil {
		return "", fmt.Errorf("failed to obtain machine secret: %w", err)
	}
	key := pbkdf2.Key(secret, salt, iterations, keySize, sha256.New)

	token, err := decrypt(encrypted, key)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt token: %w", err)
	}

	return token, nil
}

// Clear removes all stored credentials.
func (k *Keyring) Clear() error {
	credDir := k.credDir()

	// Remove the token file
	tokenPath := filepath.Join(credDir, tokenFile)
	if err := os.Remove(tokenPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove token file: %w", err)
	}

	// Remove any other files in the credentials directory
	entries, err := os.ReadDir(credDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("failed to read credentials directory: %w", err)
	}

	var errs []error
	for _, entry := range entries {
		path := filepath.Join(credDir, entry.Name())
		if err := os.RemoveAll(path); err != nil {
			errs = append(errs, fmt.Errorf("failed to remove %s: %w", entry.Name(), err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors clearing credentials: %v", errs)
	}

	return nil
}

func (k *Keyring) credDir() string {
	return filepath.Join(k.configDir, credentialsDir)
}

// deriveKey derives an encryption key using PBKDF2.
// Returns the key, the salt used, and any error.
func (k *Keyring) deriveKey() ([]byte, []byte, error) {
	salt := make([]byte, saltSize)
	if _, err := rand.Read(salt); err != nil {
		return nil, nil, fmt.Errorf("failed to generate salt: %w", err)
	}

	secret, err := k.machineSecret()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to obtain machine secret: %w", err)
	}
	key := pbkdf2.Key(secret, salt, iterations, keySize, sha256.New)
	return key, salt, nil
}

// machineSecret returns a machine-specific secret for key derivation.
// On Linux, this attempts to use the D-Bus machine ID.
// Falls back to a randomly-generated per-machine key stored with 0600 permissions (HIGH-1 fix).
func (k *Keyring) machineSecret() ([]byte, error) {
	// Try to get machine ID on Linux
	if runtime.GOOS == "linux" {
		// Try systemd machine ID first
		if id, err := os.ReadFile("/etc/machine-id"); err == nil && len(id) >= 16 {
			return id, nil
		}
		// Fallback to D-Bus machine ID
		if id, err := os.ReadFile("/var/lib/dbus/machine-id"); err == nil && len(id) >= 16 {
			return id, nil
		}
	}

	// Universal fallback: load or generate a random per-machine key file.
	// This is cryptographically strong (unlike hostname-based derivation) and
	// persists across CLI invocations so tokens remain decryptable.
	return k.loadOrCreateMachineKey()
}

// loadOrCreateMachineKey loads an existing random machine key from disk, or
// generates and stores a new one if none exists. The key file is stored with
// 0600 permissions so only the owning user can read it (HIGH-1 fix).
func (k *Keyring) loadOrCreateMachineKey() ([]byte, error) {
	keyFile := filepath.Join(k.credDir(), "machine.key")

	// Try to load an existing 32-byte key
	if data, err := os.ReadFile(keyFile); err == nil && len(data) == 32 {
		return data, nil
	}

	// Generate a new random 256-bit key
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		return nil, fmt.Errorf("failed to generate machine key: %w", err)
	}

	// Persist the key so future decryption works.
	// If either MkdirAll or WriteFile fails, return an error rather than silently
	// continuing — without persistence the key is lost on restart and any token
	// encrypted with it becomes permanently undecryptable.
	if err := os.MkdirAll(k.credDir(), 0700); err != nil {
		return nil, fmt.Errorf("failed to create credentials directory for machine key: %w", err)
	}
	if err := os.WriteFile(keyFile, key, 0600); err != nil {
		return nil, fmt.Errorf("failed to persist machine key (token will be unreadable after restart): %w", err)
	}

	return key, nil
}

type tokenFileData struct {
	Salt      string `json:"salt"`
	Encrypted string `json:"encrypted"`
}

// encrypt encrypts plaintext using AES-GCM with the provided key.
func encrypt(plaintext string, key []byte) ([]byte, error) {
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

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return ciphertext, nil
}

// decrypt decrypts ciphertext using AES-GCM with the provided key.
func decrypt(ciphertext []byte, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
