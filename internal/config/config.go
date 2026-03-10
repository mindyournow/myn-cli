package config

import (
	"fmt"
	"os"
	"path/filepath"
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
// Returns an error if the config directory cannot be determined.
func Load() (*Config, error) {
	configDir, err := defaultConfigDir()
	if err != nil {
		return nil, fmt.Errorf("failed to determine config directory: %w", err)
	}

	baseURL := os.Getenv("MYN_API_URL")
	if baseURL == "" {
		baseURL = DefaultBaseURL
	}

	return &Config{
		BaseURL:   baseURL,
		ConfigDir: configDir,
		PluginDir: filepath.Join(configDir, "plugins"),
	}, nil
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
