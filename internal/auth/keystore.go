package auth

import (
	"errors"

	gokeyring "github.com/zalando/go-keyring"
)

// KeyStore is the unified credential store.
// It tries the OS keyring (GNOME Keyring, KDE Wallet, macOS Keychain)
// first, then falls back to the encrypted file-based Keyring.
type KeyStore struct {
	fileKeyring *Keyring
	backend     string // "auto" | "gnome" | "kde" | "pass" | "file"
}

// NewKeyStore creates a KeyStore backed by the OS keyring with a file fallback.
func NewKeyStore(fileKeyring *Keyring, backend string) *KeyStore {
	return &KeyStore{
		fileKeyring: fileKeyring,
		backend:     backend,
	}
}

// useOSKeyring returns true when the OS keyring should be attempted.
func (k *KeyStore) useOSKeyring() bool {
	return k.backend != "file"
}

// SaveRefreshToken stores the OAuth refresh token.
func (k *KeyStore) SaveRefreshToken(token string) error {
	if k.useOSKeyring() {
		if err := gokeyring.Set(KeyringService, KeyringAccountRefresh, token); err == nil {
			return nil
		}
	}
	return k.fileKeyring.SaveRefreshToken(token)
}

// LoadRefreshToken retrieves the OAuth refresh token.
func (k *KeyStore) LoadRefreshToken() (string, error) {
	if k.useOSKeyring() {
		token, err := gokeyring.Get(KeyringService, KeyringAccountRefresh)
		if err == nil {
			return token, nil
		}
		if !errors.Is(err, gokeyring.ErrNotFound) {
			// Non-fatal: OS keyring query failed, fall through to file store
			_ = err
		}
	}
	return k.fileKeyring.LoadRefreshToken()
}

// SaveAPIKey stores the API key credential.
func (k *KeyStore) SaveAPIKey(key string) error {
	if k.useOSKeyring() {
		if err := gokeyring.Set(KeyringService, KeyringAccountAPIKey, key); err == nil {
			return nil
		}
	}
	return k.fileKeyring.saveRawCredential("api_key.enc", key)
}

// LoadAPIKey retrieves the stored API key.
func (k *KeyStore) LoadAPIKey() (string, error) {
	if k.useOSKeyring() {
		key, err := gokeyring.Get(KeyringService, KeyringAccountAPIKey)
		if err == nil {
			return key, nil
		}
		if !errors.Is(err, gokeyring.ErrNotFound) {
			_ = err
		}
	}
	return k.fileKeyring.loadRawCredential("api_key.enc")
}

// Clear removes all stored credentials from both stores.
func (k *KeyStore) Clear() error {
	if k.useOSKeyring() {
		_ = gokeyring.Delete(KeyringService, KeyringAccountRefresh)
		_ = gokeyring.Delete(KeyringService, KeyringAccountAPIKey)
	}
	return k.fileKeyring.Clear()
}
