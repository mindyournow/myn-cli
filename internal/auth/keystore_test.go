package auth

import (
	"testing"
)

func newFileKeyStore(t *testing.T) *KeyStore {
	t.Helper()
	kr := NewKeyring(t.TempDir())
	return NewKeyStore(kr, "file")
}

func TestKeyStore_SaveAndLoad_RefreshToken(t *testing.T) {
	ks := newFileKeyStore(t)

	want := "refresh-token-abc123"
	if err := ks.SaveRefreshToken(want); err != nil {
		t.Fatalf("SaveRefreshToken() error = %v", err)
	}

	got, err := ks.LoadRefreshToken()
	if err != nil {
		t.Fatalf("LoadRefreshToken() error = %v", err)
	}
	if got != want {
		t.Errorf("LoadRefreshToken() = %q, want %q", got, want)
	}
}

func TestKeyStore_SaveAndLoad_APIKey(t *testing.T) {
	ks := newFileKeyStore(t)

	want := "api-key-xyz789"
	if err := ks.SaveAPIKey(want); err != nil {
		t.Fatalf("SaveAPIKey() error = %v", err)
	}

	got, err := ks.LoadAPIKey()
	if err != nil {
		t.Fatalf("LoadAPIKey() error = %v", err)
	}
	if got != want {
		t.Errorf("LoadAPIKey() = %q, want %q", got, want)
	}
}

func TestKeyStore_Clear(t *testing.T) {
	ks := newFileKeyStore(t)

	if err := ks.SaveRefreshToken("some-refresh-token"); err != nil {
		t.Fatalf("SaveRefreshToken() error = %v", err)
	}
	if err := ks.SaveAPIKey("some-api-key"); err != nil {
		t.Fatalf("SaveAPIKey() error = %v", err)
	}

	if err := ks.Clear(); err != nil {
		t.Fatalf("Clear() error = %v", err)
	}

	if _, err := ks.LoadRefreshToken(); err == nil {
		t.Error("LoadRefreshToken() should return error after Clear()")
	}
	if _, err := ks.LoadAPIKey(); err == nil {
		t.Error("LoadAPIKey() should return error after Clear()")
	}
}

func TestKeyStore_LoadRefreshToken_NotFound(t *testing.T) {
	ks := newFileKeyStore(t)

	_, err := ks.LoadRefreshToken()
	if err == nil {
		t.Error("LoadRefreshToken() should return error on empty store")
	}
}

func TestKeyStore_LoadAPIKey_NotFound(t *testing.T) {
	ks := newFileKeyStore(t)

	_, err := ks.LoadAPIKey()
	if err == nil {
		t.Error("LoadAPIKey() should return error on empty store")
	}
}

func TestKeyStore_FileBackend_AlwaysUsed(t *testing.T) {
	kr := NewKeyring(t.TempDir())
	ks := NewKeyStore(kr, "file")

	if ks.useOSKeyring() {
		t.Error("useOSKeyring() should return false when backend is \"file\"")
	}
}
