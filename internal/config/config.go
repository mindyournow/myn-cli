package config

import (
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

func Load() *Config {
	configDir := defaultConfigDir()
	baseURL := os.Getenv("MYN_API_URL")
	if baseURL == "" {
		baseURL = DefaultBaseURL
	}

	return &Config{
		BaseURL:   baseURL,
		ConfigDir: configDir,
		PluginDir: filepath.Join(configDir, "plugins"),
	}
}

func defaultConfigDir() string {
	if dir := os.Getenv("XDG_CONFIG_HOME"); dir != "" {
		return filepath.Join(dir, AppName)
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", AppName)
}
