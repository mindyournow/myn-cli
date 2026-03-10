package auth

import (
	"testing"
)

// Verify TokenStore interface is properly defined
func TestTokenStore_Interface(t *testing.T) {
	// Verify Keyring implements TokenStore
	var _ TokenStore = (*Keyring)(nil)

	// Verify mockTokenStore implements TokenStore (from oauth_test.go)
	var _ TokenStore = (*mockTokenStore)(nil)
}

func TestTokenStore_KeyringImplementation(t *testing.T) {
	tmpDir := t.TempDir()
	kr := NewKeyring(tmpDir)

	// Test the interface methods
	err := kr.SaveRefreshToken("test-token")
	if err != nil {
		t.Fatalf("SaveRefreshToken() error = %v", err)
	}

	token, err := kr.LoadRefreshToken()
	if err != nil {
		t.Fatalf("LoadRefreshToken() error = %v", err)
	}

	if token != "test-token" {
		t.Errorf("LoadRefreshToken() = %q, want %q", token, "test-token")
	}

	err = kr.Clear()
	if err != nil {
		t.Fatalf("Clear() error = %v", err)
	}

	_, err = kr.LoadRefreshToken()
	if err == nil {
		t.Error("LoadRefreshToken() should return error after Clear()")
	}
}
