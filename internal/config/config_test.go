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

func TestValidateAPIURL(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		// Valid URLs
		{"https production", "https://api.mindyournow.com", false},
		{"https with path", "https://api.example.com/v1", false},
		{"http localhost", "http://localhost:8080", false},
		{"http 127.0.0.1", "http://127.0.0.1:9000", false},
		{"http ::1", "http://[::1]:8080", false},
		// Invalid URLs — must be rejected
		{"http non-localhost", "http://attacker.example.com", true},
		{"ftp scheme", "ftp://example.com", true},
		{"no scheme", "example.com", true},
		{"empty string", "", true},
		{"scheme only", "https://", true},
		{"file scheme", "file:///etc/passwd", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateAPIURL(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateAPIURL(%q) error = %v, wantErr %v", tt.url, err, tt.wantErr)
			}
		})
	}
}

func TestIsLocalhost(t *testing.T) {
	tests := []struct {
		host string
		want bool
	}{
		{"localhost", true},
		{"127.0.0.1", true},
		{"::1", true},
		{"LOCALHOST", true}, // case-insensitive
		{"192.168.1.1", false},
		{"example.com", false},
		{"notlocalhost", false},
		// Port-suffixed forms must NOT match — url.Hostname() already strips port
		{"localhost:8080", false},
		{"127.0.0.1:9000", false},
	}

	for _, tt := range tests {
		t.Run(tt.host, func(t *testing.T) {
			got := isLocalhost(tt.host)
			if got != tt.want {
				t.Errorf("isLocalhost(%q) = %v, want %v", tt.host, got, tt.want)
			}
		})
	}
}

func TestLoad_InvalidAPIURL(t *testing.T) {
	t.Setenv("MYN_API_URL", "http://attacker.example.com")
	t.Setenv("XDG_CONFIG_HOME", "")

	_, err := Load()
	if err == nil {
		t.Fatal("Load() should fail for http non-localhost URL")
	}
}

func TestLoadWithOverrides_ValidOverride(t *testing.T) {
	t.Setenv("MYN_API_URL", "")
	t.Setenv("XDG_CONFIG_HOME", "")

	cfg, err := LoadWithOverrides("https://custom.example.com")
	if err != nil {
		t.Fatalf("LoadWithOverrides() error = %v", err)
	}
	if cfg.BaseURL != "https://custom.example.com" {
		t.Errorf("BaseURL = %q, want %q", cfg.BaseURL, "https://custom.example.com")
	}
}

func TestLoadWithOverrides_InvalidOverride(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", "")

	_, err := LoadWithOverrides("http://evil.com")
	if err == nil {
		t.Fatal("LoadWithOverrides() should reject non-localhost http URL")
	}
}
