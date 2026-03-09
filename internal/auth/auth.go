package auth

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

// Manager handles all authentication operations.
type Manager struct {
	Keyring     *Keyring
	OAuth       *OAuthClient
	baseURL     string
}

// NewManager creates a new auth manager.
func NewManager(baseURL string) *Manager {
	keyring := NewKeyring()
	return &Manager{
		Keyring: keyring,
		OAuth:   NewOAuthClient(baseURL, keyring),
		baseURL: baseURL,
	}
}

// LoginWithAPIKey authenticates using an API key.
func (m *Manager) LoginWithAPIKey(apiKey string) error {
	// Validate the API key by making a test request
	if err := m.validateAPIKey(apiKey); err != nil {
		return fmt.Errorf("invalid API key: %w", err)
	}

	if err := m.Keyring.SaveAPIKey(apiKey); err != nil {
		return fmt.Errorf("storing API key: %w", err)
	}

	return nil
}

// LoginWithOAuth performs OAuth PKCE authentication.
func (m *Manager) LoginWithOAuth() error {
	tokens, err := m.OAuth.Login()
	if err != nil {
		return err
	}

	// The OAuth client already stored the refresh token
	fmt.Printf("Authenticated successfully. Token expires in %d seconds.\n", tokens.ExpiresIn)
	return nil
}

// Logout clears all stored credentials.
func (m *Manager) Logout() error {
	if err := m.Keyring.Clear(); err != nil {
		return fmt.Errorf("clearing credentials: %w", err)
	}
	return nil
}

// GetAPIKey returns the stored API key if available.
func (m *Manager) GetAPIKey() (string, error) {
	return m.Keyring.LoadAPIKey()
}

// HasCredentials returns true if any credentials are stored.
func (m *Manager) HasCredentials() bool {
	if _, err := m.Keyring.LoadAPIKey(); err == nil {
		return true
	}
	if _, err := m.Keyring.LoadRefreshToken(); err == nil {
		return true
	}
	return false
}

// GetAuthMethod returns the current authentication method.
func (m *Manager) GetAuthMethod() string {
	if _, err := m.Keyring.LoadAPIKey(); err == nil {
		return "api-key"
	}
	if _, err := m.Keyring.LoadRefreshToken(); err == nil {
		return "oauth"
	}
	return "none"
}

// validateAPIKey validates an API key by making a test request.
func (m *Manager) validateAPIKey(apiKey string) error {
	// This will be implemented by the api package
	// For now, just do basic format validation
	if !strings.HasPrefix(apiKey, "myn_") {
		return fmt.Errorf("API key must start with 'myn_'")
	}
	return nil
}

// openBrowser opens the default browser to the given URL.
// This is called by the OAuth client but defined here for reusability.
func openBrowser(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "linux":
		cmd = "xdg-open"
	case "darwin":
		cmd = "open"
	case "windows":
		cmd = "rundll32"
		args = []string{"url.dll,FileProtocolHandler"}
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}
