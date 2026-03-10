package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad_Defaults(t *testing.T) {
	// Clear environment variables
	t.Setenv("MYN_API_URL", "")
	t.Setenv("XDG_CONFIG_HOME", "")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.BaseURL != DefaultBaseURL {
		t.Errorf("BaseURL = %q, want %q", cfg.BaseURL, DefaultBaseURL)
	}

	if cfg.ConfigDir == "" {
		t.Error("ConfigDir should not be empty")
	}
}

func TestLoad_CustomBaseURL(t *testing.T) {
	customURL := "https://custom.example.com"
	t.Setenv("MYN_API_URL", customURL)
	t.Setenv("XDG_CONFIG_HOME", "")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.BaseURL != customURL {
		t.Errorf("BaseURL = %q, want %q", cfg.BaseURL, customURL)
	}
}

func TestLoad_XDGConfigHome(t *testing.T) {
	xdgDir := "/tmp/test-xdg-config"
	t.Setenv("XDG_CONFIG_HOME", xdgDir)
	t.Setenv("MYN_API_URL", "")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	expectedDir := filepath.Join(xdgDir, AppName)
	if cfg.ConfigDir != expectedDir {
		t.Errorf("ConfigDir = %q, want %q", cfg.ConfigDir, expectedDir)
	}
}

func TestLoad_PluginDir(t *testing.T) {
	xdgDir := "/tmp/test-plugins"
	t.Setenv("XDG_CONFIG_HOME", xdgDir)
	t.Setenv("MYN_API_URL", "")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	expectedPluginDir := filepath.Join(xdgDir, AppName, "plugins")
	if cfg.PluginDir != expectedPluginDir {
		t.Errorf("PluginDir = %q, want %q", cfg.PluginDir, expectedPluginDir)
	}
}

func TestDefaultConfigDir_NoXDG(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", "")

	home, err := os.UserHomeDir()
	if err != nil {
		t.Skipf("Cannot get user home dir: %v", err)
	}

	dir, err := defaultConfigDir()
	if err != nil {
		t.Fatalf("defaultConfigDir() error = %v", err)
	}

	expected := filepath.Join(home, ".config", AppName)
	if dir != expected {
		t.Errorf("defaultConfigDir() = %q, want %q", dir, expected)
	}
}

func TestDefaultConfigDir_WithXDG(t *testing.T) {
	xdgDir := "/custom/xdg/path"
	t.Setenv("XDG_CONFIG_HOME", xdgDir)

	dir, err := defaultConfigDir()
	if err != nil {
		t.Fatalf("defaultConfigDir() error = %v", err)
	}

	expected := filepath.Join(xdgDir, AppName)
	if dir != expected {
		t.Errorf("defaultConfigDir() = %q, want %q", dir, expected)
	}
}
