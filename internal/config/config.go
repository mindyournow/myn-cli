package config

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

const (
	DefaultBaseURL = "https://api.mindyournow.com"
	AppName        = "mynow"
)

type Config struct {
	BaseURL   string
	ConfigDir string
	PluginDir string
}

// Load loads the configuration from environment variables and defaults.
// Returns an error if the config directory cannot be determined or if
// MYN_API_URL is invalid (HIGH-3 fix).
func Load() (*Config, error) {
	return LoadWithOverrides("")
}

// LoadWithOverrides loads config, optionally overriding the API URL (e.g., from --api-url flag).
// An empty overrideURL means use MYN_API_URL env var or the default.
func LoadWithOverrides(overrideURL string) (*Config, error) {
	configDir, err := defaultConfigDir()
	if err != nil {
		return nil, fmt.Errorf("failed to determine config directory: %w", err)
	}

	baseURL := overrideURL
	if baseURL == "" {
		baseURL = os.Getenv("MYN_API_URL")
	}
	if baseURL == "" {
		baseURL = DefaultBaseURL
	}

	// Validate URL to prevent SSRF/redirect attacks (HIGH-3 fix)
	errLabel := "MYN_API_URL"
	if overrideURL != "" {
		errLabel = "--api-url"
	}
	if err := validateAPIURL(baseURL); err != nil {
		return nil, fmt.Errorf("invalid %s: %w", errLabel, err)
	}

	return &Config{
		BaseURL:   baseURL,
		ConfigDir: configDir,
		PluginDir: filepath.Join(configDir, "plugins"),
	}, nil
}

// validateAPIURL validates the API URL to prevent SSRF attacks.
// Only allows http://localhost* for development and https:// for production.
func validateAPIURL(baseURL string) error {
	u, err := url.Parse(baseURL)
	if err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}

	// Must have a scheme
	if u.Scheme != "https" && u.Scheme != "http" {
		return fmt.Errorf("URL scheme must be http or https, got: %s", u.Scheme)
	}

	// HTTP only allowed for localhost development
	if u.Scheme == "http" {
		host := strings.ToLower(u.Hostname())
		if !isLocalhost(host) {
			return fmt.Errorf("http:// is only allowed for localhost development, use https:// for: %s", u.Host)
		}
	}

	// Must have a host
	if u.Host == "" {
		return fmt.Errorf("URL must have a host")
	}

	return nil
}

// isLocalhost checks if a hostname is localhost (127.0.0.1, ::1, or localhost variants)
func isLocalhost(host string) bool {
	host = strings.ToLower(host)
	return host == "localhost" ||
		host == "127.0.0.1" ||
		host == "::1" ||
		strings.HasPrefix(host, "localhost:") ||
		strings.HasPrefix(host, "127.0.0.1:")
}

// defaultConfigDir returns the configuration directory path.
// It checks XDG_CONFIG_HOME first, then falls back to ~/.config/mynow.
func defaultConfigDir() (string, error) {
	if dir := os.Getenv("XDG_CONFIG_HOME"); dir != "" {
		return filepath.Join(dir, AppName), nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}
	return filepath.Join(home, ".config", AppName), nil
}
