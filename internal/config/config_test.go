package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// --- Load / defaults ---

func TestLoad_Defaults(t *testing.T) {
	t.Setenv("MYN_API_URL", "")
	t.Setenv("XDG_CONFIG_HOME", "")
	t.Setenv("MYNOW_CONFIG", "")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.BaseURL != DefaultBaseURL {
		t.Errorf("BaseURL = %q, want %q", cfg.BaseURL, DefaultBaseURL)
	}
	if cfg.API.URL != DefaultBaseURL {
		t.Errorf("API.URL = %q, want %q", cfg.API.URL, DefaultBaseURL)
	}
	if cfg.ConfigDir == "" {
		t.Error("ConfigDir should not be empty")
	}
	// Verify defaults
	if cfg.API.Retries != 3 {
		t.Errorf("API.Retries = %d, want 3", cfg.API.Retries)
	}
	if cfg.Display.Color != "auto" {
		t.Errorf("Display.Color = %q, want auto", cfg.Display.Color)
	}
	if cfg.TUI.VimKeys != true {
		t.Error("TUI.VimKeys should default to true")
	}
	if cfg.Defaults.CalendarDays != 7 {
		t.Errorf("Defaults.CalendarDays = %d, want 7", cfg.Defaults.CalendarDays)
	}
}

func TestLoad_CustomBaseURL(t *testing.T) {
	customURL := "https://custom.example.com"
	t.Setenv("MYN_API_URL", customURL)
	t.Setenv("XDG_CONFIG_HOME", "")
	t.Setenv("MYNOW_CONFIG", "")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.BaseURL != customURL {
		t.Errorf("BaseURL = %q, want %q", cfg.BaseURL, customURL)
	}
	if cfg.API.URL != customURL {
		t.Errorf("API.URL = %q, want %q", cfg.API.URL, customURL)
	}
}

func TestLoad_XDGConfigHome(t *testing.T) {
	xdgDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", xdgDir)
	t.Setenv("MYN_API_URL", "")
	t.Setenv("MYNOW_CONFIG", "")

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
	xdgDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", xdgDir)
	t.Setenv("MYN_API_URL", "")
	t.Setenv("MYNOW_CONFIG", "")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	expectedPluginDir := filepath.Join(xdgDir, AppName, "plugins")
	if cfg.PluginDir != expectedPluginDir {
		t.Errorf("PluginDir = %q, want %q", cfg.PluginDir, expectedPluginDir)
	}
}

// --- YAML config file loading ---

func TestLoad_YAMLFile(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", dir)
	t.Setenv("MYN_API_URL", "")
	t.Setenv("MYNOW_CONFIG", "")

	cfgFile := filepath.Join(dir, AppName, ConfigFileName)
	if err := os.MkdirAll(filepath.Dir(cfgFile), 0700); err != nil {
		t.Fatal(err)
	}
	yaml := `
api:
  url: https://my-server.example.com
  timeout: 60s
  retries: 5
display:
  color: always
  date_format: iso
tui:
  theme: light
  vim_keys: true
defaults:
  calendar_days: 14
`
	if err := os.WriteFile(cfgFile, []byte(yaml), 0600); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.API.URL != "https://my-server.example.com" {
		t.Errorf("API.URL = %q, want https://my-server.example.com", cfg.API.URL)
	}
	if cfg.API.Timeout != "60s" {
		t.Errorf("API.Timeout = %q, want 60s", cfg.API.Timeout)
	}
	if cfg.API.Retries != 5 {
		t.Errorf("API.Retries = %d, want 5", cfg.API.Retries)
	}
	if cfg.Display.Color != "always" {
		t.Errorf("Display.Color = %q, want always", cfg.Display.Color)
	}
	if cfg.Display.DateFormat != "iso" {
		t.Errorf("Display.DateFormat = %q, want iso", cfg.Display.DateFormat)
	}
	if cfg.TUI.Theme != "light" {
		t.Errorf("TUI.Theme = %q, want light", cfg.TUI.Theme)
	}
	if cfg.Defaults.CalendarDays != 14 {
		t.Errorf("Defaults.CalendarDays = %d, want 14", cfg.Defaults.CalendarDays)
	}
	// Unspecified fields should keep defaults
	if cfg.Display.TimeFormat != "12h" {
		t.Errorf("Display.TimeFormat = %q, want 12h (default)", cfg.Display.TimeFormat)
	}
}

func TestLoad_MYNOWConfigEnvVar(t *testing.T) {
	dir := t.TempDir()
	cfgFile := filepath.Join(dir, "custom-config.yaml")
	yaml := "api:\n  url: https://custom-env.example.com\n"
	if err := os.WriteFile(cfgFile, []byte(yaml), 0600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("MYNOW_CONFIG", cfgFile)
	t.Setenv("MYN_API_URL", "")
	t.Setenv("XDG_CONFIG_HOME", "")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.API.URL != "https://custom-env.example.com" {
		t.Errorf("API.URL = %q", cfg.API.URL)
	}
	if cfg.ConfigFile != cfgFile {
		t.Errorf("ConfigFile = %q, want %q", cfg.ConfigFile, cfgFile)
	}
}

func TestLoad_EnvVarAPIKey(t *testing.T) {
	t.Setenv("MYN_API_KEY", "myapikey123")
	t.Setenv("MYN_API_URL", "")
	t.Setenv("XDG_CONFIG_HOME", "")
	t.Setenv("MYNOW_CONFIG", "")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.APIKey != "myapikey123" {
		t.Errorf("APIKey = %q, want myapikey123", cfg.APIKey)
	}
}

func TestLoad_NoColor(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	t.Setenv("MYN_API_URL", "")
	t.Setenv("XDG_CONFIG_HOME", "")
	t.Setenv("MYNOW_CONFIG", "")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.Display.Color != "never" {
		t.Errorf("Display.Color = %q, want never (NO_COLOR set)", cfg.Display.Color)
	}
}

func TestLoad_MYNOW_DEBUG(t *testing.T) {
	t.Setenv("MYNOW_DEBUG", "1")
	t.Setenv("MYN_API_URL", "")
	t.Setenv("XDG_CONFIG_HOME", "")
	t.Setenv("MYNOW_CONFIG", "")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if !cfg.Debug {
		t.Error("Debug should be true when MYNOW_DEBUG is set")
	}
}

// --- validateAPIURL / isLocalhost ---

func TestValidateAPIURL(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{"https production", "https://api.mindyournow.com", false},
		{"https with path", "https://api.example.com/v1", false},
		{"http localhost", "http://localhost:8080", false},
		{"http 127.0.0.1", "http://127.0.0.1:9000", false},
		{"http ::1", "http://[::1]:8080", false},
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
		{"LOCALHOST", true},
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
	t.Setenv("MYNOW_CONFIG", "")

	_, err := Load()
	if err == nil {
		t.Fatal("Load() should fail for http non-localhost URL")
	}
}

func TestLoadWithOverrides_ValidOverride(t *testing.T) {
	t.Setenv("MYN_API_URL", "")
	t.Setenv("XDG_CONFIG_HOME", "")
	t.Setenv("MYNOW_CONFIG", "")

	cfg, err := LoadWithOverrides("https://custom.example.com")
	if err != nil {
		t.Fatalf("LoadWithOverrides() error = %v", err)
	}
	if cfg.BaseURL != "https://custom.example.com" {
		t.Errorf("BaseURL = %q, want https://custom.example.com", cfg.BaseURL)
	}
}

func TestLoadWithOverrides_InvalidOverride(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", "")
	t.Setenv("MYNOW_CONFIG", "")

	_, err := LoadWithOverrides("http://evil.com")
	if err == nil {
		t.Fatal("LoadWithOverrides() should reject non-localhost http URL")
	}
}

// --- GetValue / SetValue ---

func TestGetValue(t *testing.T) {
	cfg := defaults()

	val, err := GetValue(&cfg, "api.url")
	if err != nil {
		t.Fatalf("GetValue api.url: %v", err)
	}
	if val != DefaultBaseURL {
		t.Errorf("api.url = %q, want %q", val, DefaultBaseURL)
	}

	val, err = GetValue(&cfg, "api.retries")
	if err != nil {
		t.Fatalf("GetValue api.retries: %v", err)
	}
	if val != "3" {
		t.Errorf("api.retries = %q, want 3", val)
	}

	val, err = GetValue(&cfg, "tui.vim_keys")
	if err != nil {
		t.Fatalf("GetValue tui.vim_keys: %v", err)
	}
	if val != "true" {
		t.Errorf("tui.vim_keys = %q, want true", val)
	}

	_, err = GetValue(&cfg, "nonexistent.key")
	if err == nil {
		t.Error("GetValue should return error for unknown key")
	}
}

func TestSetValue(t *testing.T) {
	cfg := defaults()

	if err := SetValue(&cfg, "api.url", "https://new.example.com"); err != nil {
		t.Fatalf("SetValue api.url: %v", err)
	}
	if cfg.API.URL != "https://new.example.com" {
		t.Errorf("API.URL = %q after set", cfg.API.URL)
	}
	// BaseURL should also be updated
	if cfg.BaseURL != "https://new.example.com" {
		t.Errorf("BaseURL = %q after set api.url", cfg.BaseURL)
	}

	if err := SetValue(&cfg, "api.retries", "5"); err != nil {
		t.Fatalf("SetValue api.retries: %v", err)
	}
	if cfg.API.Retries != 5 {
		t.Errorf("API.Retries = %d after set", cfg.API.Retries)
	}

	if err := SetValue(&cfg, "display.color", "always"); err != nil {
		t.Fatalf("SetValue display.color: %v", err)
	}
	if cfg.Display.Color != "always" {
		t.Errorf("Display.Color = %q after set", cfg.Display.Color)
	}

	if err := SetValue(&cfg, "tui.vim_keys", "false"); err != nil {
		t.Fatalf("SetValue tui.vim_keys: %v", err)
	}
	if cfg.TUI.VimKeys != false {
		t.Error("TUI.VimKeys should be false after set")
	}
}

func TestSetValue_Validation(t *testing.T) {
	cfg := defaults()

	tests := []struct {
		key   string
		value string
	}{
		{"api.url", "http://evil.com"},      // non-localhost http rejected
		{"api.retries", "-1"},               // negative rejected
		{"api.timeout", "notaduration"},     // invalid duration
		{"display.color", "rainbow"},        // invalid enum
		{"auth.method", "magic"},            // invalid enum
		{"auth.keyring", "windows"},         // invalid enum
		{"tui.vim_keys", "yes"},             // not a bool
		{"nonexistent.key", "value"},        // unknown key
	}

	for _, tt := range tests {
		t.Run(tt.key+"="+tt.value, func(t *testing.T) {
			if err := SetValue(&cfg, tt.key, tt.value); err == nil {
				t.Errorf("SetValue(%q, %q) should return error", tt.key, tt.value)
			}
		})
	}
}

// --- Save / Reset ---

func TestSave_And_Reload(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", dir)
	t.Setenv("MYN_API_URL", "")
	t.Setenv("MYNOW_CONFIG", "")

	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}

	cfg.API.URL = "https://saved.example.com"
	cfg.Display.DateFormat = "long"

	if err := Save(cfg); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	// Reload
	cfg2, err := Load()
	if err != nil {
		t.Fatalf("Load() after Save() error = %v", err)
	}
	if cfg2.API.URL != "https://saved.example.com" {
		t.Errorf("API.URL after reload = %q", cfg2.API.URL)
	}
	if cfg2.Display.DateFormat != "long" {
		t.Errorf("Display.DateFormat after reload = %q", cfg2.Display.DateFormat)
	}
}

func TestReset(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", dir)
	t.Setenv("MYN_API_URL", "")
	t.Setenv("MYNOW_CONFIG", "")

	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	cfg.Display.Color = "always"
	if err := Save(cfg); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	if err := Reset(cfg); err != nil {
		t.Fatalf("Reset() error = %v", err)
	}

	// File should be gone
	if _, err := os.Stat(cfg.ConfigFile); !os.IsNotExist(err) {
		t.Error("Config file should be removed after Reset()")
	}

	// Reload gets defaults
	cfg2, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if cfg2.Display.Color != "auto" {
		t.Errorf("Display.Color = %q after reset, want auto", cfg2.Display.Color)
	}
}

// --- Duration helpers ---

func TestAPIConfig_TimeoutDuration(t *testing.T) {
	a := APIConfig{Timeout: "45s"}
	if a.TimeoutDuration().Seconds() != 45 {
		t.Errorf("TimeoutDuration = %v, want 45s", a.TimeoutDuration())
	}
	// Invalid falls back to 30s
	a.Timeout = "invalid"
	if a.TimeoutDuration().Seconds() != 30 {
		t.Errorf("TimeoutDuration with invalid = %v, want 30s", a.TimeoutDuration())
	}
}

// --- defaultConfigDir ---

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

func TestLoad_InvalidYAML(t *testing.T) {
	dir := t.TempDir()
	cfgFile := filepath.Join(dir, "bad.yaml")
	if err := os.WriteFile(cfgFile, []byte("not: valid: yaml: :"), 0600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("MYNOW_CONFIG", cfgFile)
	t.Setenv("MYN_API_URL", "")
	t.Setenv("XDG_CONFIG_HOME", "")

	_, err := Load()
	if err == nil {
		t.Error("Load() should fail on invalid YAML")
	}
	if !strings.Contains(err.Error(), "failed to load config file") {
		t.Errorf("error = %q, want 'failed to load config file'", err.Error())
	}
}
