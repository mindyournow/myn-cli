package auth

// TokenStore manages credential persistence.
// Uses Linux Secret Service API (GNOME Keyring / KDE Wallet) when available,
// with fallback to encrypted file storage.
type TokenStore interface {
	SaveRefreshToken(token string) error
	LoadRefreshToken() (string, error)
	Clear() error
}
