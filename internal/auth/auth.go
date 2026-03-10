package auth

// TokenStore manages credential persistence.
// Uses Linux Secret Service API (GNOME Keyring / KDE Wallet) when available,
// with fallback to encrypted file storage.
type TokenStore interface {
	SaveRefreshToken(token string) error
	LoadRefreshToken() (string, error)
	Clear() error
}

// CredentialStore extends TokenStore with API key support.
type CredentialStore interface {
	TokenStore
	SaveAPIKey(key string) error
	LoadAPIKey() (string, error)
}

// AuthMethod represents the configured authentication method.
type AuthMethod string

const (
	AuthMethodOAuth  AuthMethod = "oauth"
	AuthMethodAPIKey AuthMethod = "api-key"
	AuthMethodDevice AuthMethod = "device"
)

// KeyringService and account names used with go-keyring.
const (
	KeyringService        = "mynow"
	KeyringAccountRefresh = "refresh-token"
	KeyringAccountAPIKey  = "api-key"
)
