package auth

import (
	"errors"
	"fmt"
	"os"
	"time"

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
		} else {
			fmt.Fprintf(os.Stderr, "Warning: OS keyring unavailable (%v), falling back to file store\n", err)
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
		} else {
			fmt.Fprintf(os.Stderr, "Warning: OS keyring unavailable (%v), falling back to file store\n", err)
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

// SaveAccessToken saves the access token with its expiry for cross-process use.
func (k *KeyStore) SaveAccessToken(token string, expiresAt time.Time) error {
	return k.fileKeyring.SaveAccessToken(token, expiresAt)
}

// LoadAccessToken loads the cached access token and its expiry.
func (k *KeyStore) LoadAccessToken() (string, time.Time, error) {
	return k.fileKeyring.LoadAccessToken()
}

// SaveClientID saves the OAuth client ID for use across processes.
func (k *KeyStore) SaveClientID(clientID string) error {
	return k.fileKeyring.SaveClientID(clientID)
}

// LoadClientID loads the saved OAuth client ID.
func (k *KeyStore) LoadClientID() (string, error) {
	return k.fileKeyring.LoadClientID()
}

// Clear removes all stored credentials from both stores.
func (k *KeyStore) Clear() error {
	if k.useOSKeyring() {
		_ = gokeyring.Delete(KeyringService, KeyringAccountRefresh)
		_ = gokeyring.Delete(KeyringService, KeyringAccountAPIKey)
	}
	return k.fileKeyring.Clear()
}
