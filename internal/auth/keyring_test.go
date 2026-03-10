package auth

import (
	"encoding/base64"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestKeyring_SaveAndLoadRefreshToken(t *testing.T) {
	tmpDir := t.TempDir()
	kr := NewKeyring(tmpDir)

	testToken := "test-refresh-token-12345"

	// Save token
	err := kr.SaveRefreshToken(testToken)
	if err != nil {
		t.Fatalf("SaveRefreshToken() error = %v", err)
	}

	// Load token
	loadedToken, err := kr.LoadRefreshToken()
	if err != nil {
		t.Fatalf("LoadRefreshToken() error = %v", err)
	}

	if loadedToken != testToken {
		t.Errorf("LoadRefreshToken() = %q, want %q", loadedToken, testToken)
	}
}

func TestKeyring_LoadRefreshToken_NotFound(t *testing.T) {
	tmpDir := t.TempDir()
	kr := NewKeyring(tmpDir)

	_, err := kr.LoadRefreshToken()
	if err == nil {
		t.Error("LoadRefreshToken() should return error when token not found")
	}
}

func TestKeyring_SaveRefreshToken_Empty(t *testing.T) {
	tmpDir := t.TempDir()
	kr := NewKeyring(tmpDir)

	err := kr.SaveRefreshToken("")
	if err == nil {
		t.Error("SaveRefreshToken() should return error for empty token")
	}
}

func TestKeyring_Clear(t *testing.T) {
	tmpDir := t.TempDir()
	kr := NewKeyring(tmpDir)

	// Save a token
	err := kr.SaveRefreshToken("test-token")
	if err != nil {
		t.Fatalf("SaveRefreshToken() error = %v", err)
	}

	// Verify token exists
	_, err = kr.LoadRefreshToken()
	if err != nil {
		t.Fatalf("Token should exist before clear")
	}

	// Clear credentials
	err = kr.Clear()
	if err != nil {
		t.Fatalf("Clear() error = %v", err)
	}

	// Verify token is gone
	_, err = kr.LoadRefreshToken()
	if err == nil {
		t.Error("LoadRefreshToken() should return error after Clear()")
	}
}

func TestKeyring_Clear_NoErrorOnEmpty(t *testing.T) {
	tmpDir := t.TempDir()
	kr := NewKeyring(tmpDir)

	// Clear when nothing exists should not error
	err := kr.Clear()
	if err != nil {
		t.Errorf("Clear() on empty keyring error = %v", err)
	}
}

func TestKeyring_Encryption(t *testing.T) {
	tmpDir := t.TempDir()
	kr := NewKeyring(tmpDir)

	testToken := "secret-token"
	err := kr.SaveRefreshToken(testToken)
	if err != nil {
		t.Fatalf("SaveRefreshToken() error = %v", err)
	}

	// Read the raw file and verify it's encrypted (not plaintext)
	tokenPath := filepath.Join(tmpDir, credentialsDir, tokenFile)
	data, err := os.ReadFile(tokenPath)
	if err != nil {
		t.Fatalf("Failed to read token file: %v", err)
	}

	// Parse the JSON structure
	var tokenData tokenFileData
	if err := json.Unmarshal(data, &tokenData); err != nil {
		t.Fatalf("Failed to unmarshal token data: %v", err)
	}

	// Verify salt is present and base64 encoded
	if tokenData.Salt == "" {
		t.Error("Salt should not be empty")
	}
	saltBytes, err := base64.StdEncoding.DecodeString(tokenData.Salt)
	if err != nil {
		t.Errorf("Salt should be base64 encoded: %v", err)
	}
	if len(saltBytes) != saltSize {
		t.Errorf("Salt length = %d, want %d", len(saltBytes), saltSize)
	}

	// Verify encrypted data is present and base64 encoded
	if tokenData.Encrypted == "" {
		t.Error("Encrypted data should not be empty")
	}
	encryptedBytes, err := base64.StdEncoding.DecodeString(tokenData.Encrypted)
	if err != nil {
		t.Errorf("Encrypted data should be base64 encoded: %v", err)
	}

	// Verify the encrypted data doesn't contain the plaintext token
	if string(data) == testToken {
		t.Error("Token should be encrypted, not stored as plaintext")
	}

	// The encrypted data should include nonce + ciphertext, so it should be larger than the token
	if len(encryptedBytes) < len(testToken) {
		t.Error("Encrypted data should include nonce and be larger than plaintext")
	}
}

func TestKeyring_CredentialsDirPermissions(t *testing.T) {
	tmpDir := t.TempDir()
	kr := NewKeyring(tmpDir)

	err := kr.SaveRefreshToken("test-token")
	if err != nil {
		t.Fatalf("SaveRefreshToken() error = %v", err)
	}

	// Check directory permissions
	credDir := filepath.Join(tmpDir, credentialsDir)
	info, err := os.Stat(credDir)
	if err != nil {
		t.Fatalf("Failed to stat credentials dir: %v", err)
	}

	// On Unix systems, verify the directory is not world-readable
	mode := info.Mode().Perm()
	if mode&004 != 0 || mode&002 != 0 || mode&001 != 0 {
		t.Errorf("Credentials directory permissions = %o, should not be world accessible", mode)
	}
}

func TestEncryptDecrypt(t *testing.T) {
	key := make([]byte, keySize)
	for i := range key {
		key[i] = byte(i)
	}

	plaintext := "test message for encryption"

	encrypted, err := encrypt(plaintext, key)
	if err != nil {
		t.Fatalf("encrypt() error = %v", err)
	}

	decrypted, err := decrypt(encrypted, key)
	if err != nil {
		t.Fatalf("decrypt() error = %v", err)
	}

	if decrypted != plaintext {
		t.Errorf("decrypt(encrypt(plaintext)) = %q, want %q", decrypted, plaintext)
	}
}

func TestDecrypt_WrongKey(t *testing.T) {
	key1 := make([]byte, keySize)
	key2 := make([]byte, keySize)
	key2[0] = 1 // Different key

	plaintext := "test message"
	encrypted, _ := encrypt(plaintext, key1)

	_, err := decrypt(encrypted, key2)
	if err == nil {
		t.Error("decrypt() should fail with wrong key")
	}
}

func TestDecrypt_CiphertextTooShort(t *testing.T) {
	key := make([]byte, keySize)
	_, err := decrypt([]byte{1, 2, 3}, key)
	if err == nil {
		t.Error("decrypt() should fail with short ciphertext")
	}
}
