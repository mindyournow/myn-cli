package config

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

const (
	DefaultBaseURL = "https://api.mindyournow.com"
	AppName        = "mynow"
	ConfigFileName = "config.yaml"
)

// Config is the resolved, fully-loaded application configuration.
// Fields at the top level are convenience accessors maintained for backward
// compatibility; the structured sub-fields mirror the YAML file layout.
type Config struct {
	// Resolved top-level fields (backward compat)
	BaseURL    string
	ConfigDir  string
	PluginDir  string
	ConfigFile string
	APIKey     string
	Debug      bool

	// Structured YAML config sections
	API      APIConfig      `yaml:"api"`
	Auth     AuthConfig     `yaml:"auth"`
	Display  DisplayConfig  `yaml:"display"`
	TUI      TUIConfig      `yaml:"tui"`
	Defaults DefaultsConfig `yaml:"defaults"`
}

// APIConfig holds backend connection settings.
type APIConfig struct {
	URL     string `yaml:"url"`
	Timeout string `yaml:"timeout"` // e.g. "30s"
	Retries int    `yaml:"retries"`
}

// TimeoutDuration parses the Timeout string; falls back to 30s on error.
func (a APIConfig) TimeoutDuration() time.Duration {
	d, err := time.ParseDuration(a.Timeout)
	if err != nil || d <= 0 {
		return 30 * time.Second
	}
	return d
}

// AuthConfig holds authentication preferences.
type AuthConfig struct {
	Method  string `yaml:"method"`  // api-key | oauth | device
	Keyring string `yaml:"keyring"` // auto | gnome | kde | pass | file
}

// DisplayConfig holds output formatting preferences.
type DisplayConfig struct {
	Color         string `yaml:"color"`          // auto | always | never
	DateFormat    string `yaml:"date_format"`    // relative | iso | short | long
	TimeFormat    string `yaml:"time_format"`    // 12h | 24h
	DefaultOutput string `yaml:"default_output"` // text | json | table
}

// TUIConfig holds TUI-specific settings.
type TUIConfig struct {
	Theme           string `yaml:"theme"`            // dark | light | auto
	RefreshInterval string `yaml:"refresh_interval"` // e.g. "30s"
	VimKeys         bool   `yaml:"vim_keys"`
	Mouse           bool   `yaml:"mouse"`
	Animations      bool   `yaml:"animations"`
}

// RefreshIntervalDuration parses RefreshInterval; falls back to 30s on error.
func (t TUIConfig) RefreshIntervalDuration() time.Duration {
	d, err := time.ParseDuration(t.RefreshInterval)
	if err != nil || d <= 0 {
		return 30 * time.Second
	}
	return d
}

// DefaultsConfig holds default values for new items.
type DefaultsConfig struct {
	Priority          string `yaml:"priority"`
	TaskType          string `yaml:"task_type"`
	CalendarDays      int    `yaml:"calendar_days"`
	HabitScheduleDays int    `yaml:"habit_schedule_days"`
}

// defaults returns the baseline Config with all sensible defaults filled in.
func defaults() Config {
	return Config{
		API: APIConfig{
			URL:     DefaultBaseURL,
			Timeout: "30s",
			Retries: 3,
		},
		Auth: AuthConfig{
			Method:  "oauth",
			Keyring: "auto",
		},
		Display: DisplayConfig{
			Color:         "auto",
			DateFormat:    "relative",
			TimeFormat:    "12h",
			DefaultOutput: "text",
		},
		TUI: TUIConfig{
			Theme:           "dark",
			RefreshInterval: "30s",
			VimKeys:         true,
			Mouse:           false,
			Animations:      true,
		},
		Defaults: DefaultsConfig{
			Priority:          "OPPORTUNITY_NOW",
			TaskType:          "TASK",
			CalendarDays:      7,
			HabitScheduleDays: 7,
		},
	}
}

// Load loads the configuration from the YAML file and environment variables.
func Load() (*Config, error) {
	return LoadWithOverrides("")
}

// LoadWithOverrides loads config, optionally overriding the API URL.
// An empty overrideURL means use YAML config / MYN_API_URL / default.
func LoadWithOverrides(overrideURL string) (*Config, error) {
	configDir, err := defaultConfigDir()
	if err != nil {
		return nil, fmt.Errorf("failed to determine config directory: %w", err)
	}

	// Determine config file path (MYNOW_CONFIG env var or default)
	configFile := os.Getenv("MYNOW_CONFIG")
	if configFile == "" {
		configFile = filepath.Join(configDir, ConfigFileName)
	}

	// Start with defaults
	cfg := defaults()
	cfg.ConfigDir = configDir
	cfg.PluginDir = filepath.Join(configDir, "plugins")
	cfg.ConfigFile = configFile

	// Load YAML file if it exists (missing file is not an error)
	if err := loadYAMLFile(configFile, &cfg); err != nil {
		return nil, fmt.Errorf("failed to load config file %s: %w", configFile, err)
	}

	// Apply environment variable overrides
	if v := os.Getenv("MYN_API_URL"); v != "" {
		cfg.API.URL = v
	}
	if v := os.Getenv("MYN_API_KEY"); v != "" {
		cfg.APIKey = v
	}
	if v := os.Getenv("MYNOW_KEYRING"); v != "" {
		cfg.Auth.Keyring = v
	}
	if os.Getenv("NO_COLOR") != "" {
		cfg.Display.Color = "never"
	}
	if os.Getenv("MYNOW_DEBUG") != "" {
		cfg.Debug = true
	}

	// Apply --api-url flag override (highest priority)
	if overrideURL != "" {
		cfg.API.URL = overrideURL
	}

	// Validate the final API URL
	errLabel := "MYN_API_URL"
	if overrideURL != "" {
		errLabel = "--api-url"
	}
	if err := validateAPIURL(cfg.API.URL); err != nil {
		return nil, fmt.Errorf("invalid %s: %w", errLabel, err)
	}

	// Set resolved top-level field for backward compat
	cfg.BaseURL = cfg.API.URL

	return &cfg, nil
}

// loadYAMLFile reads the YAML config file into cfg. Missing file is silently ignored.
func loadYAMLFile(path string, cfg *Config) error {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // no config file is fine
		}
		return err
	}
	// Unmarshal into a temporary struct so we only override fields present in the file
	var fileCfg Config
	if err := yaml.Unmarshal(data, &fileCfg); err != nil {
		return fmt.Errorf("invalid YAML: %w", err)
	}
	// Merge: only override non-zero values from file
	mergeConfig(cfg, &fileCfg)
	return nil
}

// mergeConfig copies non-zero values from src into dst.
func mergeConfig(dst, src *Config) {
	if src.API.URL != "" {
		dst.API.URL = src.API.URL
	}
	if src.API.Timeout != "" {
		dst.API.Timeout = src.API.Timeout
	}
	if src.API.Retries != 0 {
		dst.API.Retries = src.API.Retries
	}
	if src.Auth.Method != "" {
		dst.Auth.Method = src.Auth.Method
	}
	if src.Auth.Keyring != "" {
		dst.Auth.Keyring = src.Auth.Keyring
	}
	if src.Display.Color != "" {
		dst.Display.Color = src.Display.Color
	}
	if src.Display.DateFormat != "" {
		dst.Display.DateFormat = src.Display.DateFormat
	}
	if src.Display.TimeFormat != "" {
		dst.Display.TimeFormat = src.Display.TimeFormat
	}
	if src.Display.DefaultOutput != "" {
		dst.Display.DefaultOutput = src.Display.DefaultOutput
	}
	if src.TUI.Theme != "" {
		dst.TUI.Theme = src.TUI.Theme
	}
	if src.TUI.RefreshInterval != "" {
		dst.TUI.RefreshInterval = src.TUI.RefreshInterval
	}
	// booleans: use src value (yaml.Unmarshal sets false for absent bools, so we
	// can't distinguish "false" from "not set" without pointer types; for now
	// we overwrite only when the whole TUI section was present in the file).
	if src.TUI.VimKeys {
		dst.TUI.VimKeys = true
	}
	if src.TUI.Mouse {
		dst.TUI.Mouse = true
	}
	if src.TUI.Animations {
		dst.TUI.Animations = true
	}
	if src.Defaults.Priority != "" {
		dst.Defaults.Priority = src.Defaults.Priority
	}
	if src.Defaults.TaskType != "" {
		dst.Defaults.TaskType = src.Defaults.TaskType
	}
	if src.Defaults.CalendarDays != 0 {
		dst.Defaults.CalendarDays = src.Defaults.CalendarDays
	}
	if src.Defaults.HabitScheduleDays != 0 {
		dst.Defaults.HabitScheduleDays = src.Defaults.HabitScheduleDays
	}
}

// Save writes the current config (excluding resolved top-level fields) to the config file.
func Save(cfg *Config) error {
	if err := os.MkdirAll(filepath.Dir(cfg.ConfigFile), 0700); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}
	// Build a YAML-friendly struct without the resolved top-level fields
	type yamlConfig struct {
		API      APIConfig      `yaml:"api"`
		Auth     AuthConfig     `yaml:"auth"`
		Display  DisplayConfig  `yaml:"display"`
		TUI      TUIConfig      `yaml:"tui"`
		Defaults DefaultsConfig `yaml:"defaults"`
	}
	yc := yamlConfig{
		API:      cfg.API,
		Auth:     cfg.Auth,
		Display:  cfg.Display,
		TUI:      cfg.TUI,
		Defaults: cfg.Defaults,
	}
	data, err := yaml.Marshal(yc)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}
	header := "# MYN CLI configuration\n"
	if err := os.WriteFile(cfg.ConfigFile, append([]byte(header), data...), 0600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}
	return nil
}

// GetValue returns the string value of a dot-notation config key (e.g. "api.url").
func GetValue(cfg *Config, key string) (string, error) {
	switch key {
	case "api.url":
		return cfg.API.URL, nil
	case "api.timeout":
		return cfg.API.Timeout, nil
	case "api.retries":
		return strconv.Itoa(cfg.API.Retries), nil
	case "auth.method":
		return cfg.Auth.Method, nil
	case "auth.keyring":
		return cfg.Auth.Keyring, nil
	case "display.color":
		return cfg.Display.Color, nil
	case "display.date_format":
		return cfg.Display.DateFormat, nil
	case "display.time_format":
		return cfg.Display.TimeFormat, nil
	case "display.default_output":
		return cfg.Display.DefaultOutput, nil
	case "tui.theme":
		return cfg.TUI.Theme, nil
	case "tui.refresh_interval":
		return cfg.TUI.RefreshInterval, nil
	case "tui.vim_keys":
		return strconv.FormatBool(cfg.TUI.VimKeys), nil
	case "tui.mouse":
		return strconv.FormatBool(cfg.TUI.Mouse), nil
	case "tui.animations":
		return strconv.FormatBool(cfg.TUI.Animations), nil
	case "defaults.priority":
		return cfg.Defaults.Priority, nil
	case "defaults.task_type":
		return cfg.Defaults.TaskType, nil
	case "defaults.calendar_days":
		return strconv.Itoa(cfg.Defaults.CalendarDays), nil
	case "defaults.habit_schedule_days":
		return strconv.Itoa(cfg.Defaults.HabitScheduleDays), nil
	default:
		return "", fmt.Errorf("unknown config key: %q (use 'mynow config show' to list valid keys)", key)
	}
}

// SetValue sets a dot-notation config key to the given string value, validating the value.
func SetValue(cfg *Config, key, value string) error {
	switch key {
	case "api.url":
		if err := validateAPIURL(value); err != nil {
			return fmt.Errorf("invalid api.url: %w", err)
		}
		cfg.API.URL = value
		cfg.BaseURL = value
	case "api.timeout":
		if _, err := time.ParseDuration(value); err != nil {
			return fmt.Errorf("invalid api.timeout %q: must be a duration like '30s' or '1m'", value)
		}
		cfg.API.Timeout = value
	case "api.retries":
		n, err := strconv.Atoi(value)
		if err != nil || n < 0 {
			return fmt.Errorf("invalid api.retries %q: must be a non-negative integer", value)
		}
		cfg.API.Retries = n
	case "auth.method":
		switch value {
		case "api-key", "oauth", "device":
		default:
			return fmt.Errorf("invalid auth.method %q: must be api-key, oauth, or device", value)
		}
		cfg.Auth.Method = value
	case "auth.keyring":
		switch value {
		case "auto", "gnome", "kde", "pass", "file":
		default:
			return fmt.Errorf("invalid auth.keyring %q: must be auto, gnome, kde, pass, or file", value)
		}
		cfg.Auth.Keyring = value
	case "display.color":
		switch value {
		case "auto", "always", "never":
		default:
			return fmt.Errorf("invalid display.color %q: must be auto, always, or never", value)
		}
		cfg.Display.Color = value
	case "display.date_format":
		switch value {
		case "relative", "iso", "short", "long":
		default:
			return fmt.Errorf("invalid display.date_format %q: must be relative, iso, short, or long", value)
		}
		cfg.Display.DateFormat = value
	case "display.time_format":
		switch value {
		case "12h", "24h":
		default:
			return fmt.Errorf("invalid display.time_format %q: must be 12h or 24h", value)
		}
		cfg.Display.TimeFormat = value
	case "display.default_output":
		switch value {
		case "text", "json", "table":
		default:
			return fmt.Errorf("invalid display.default_output %q: must be text, json, or table", value)
		}
		cfg.Display.DefaultOutput = value
	case "tui.theme":
		switch value {
		case "dark", "light", "auto":
		default:
			return fmt.Errorf("invalid tui.theme %q: must be dark, light, or auto", value)
		}
		cfg.TUI.Theme = value
	case "tui.refresh_interval":
		if _, err := time.ParseDuration(value); err != nil {
			return fmt.Errorf("invalid tui.refresh_interval %q: must be a duration like '30s'", value)
		}
		cfg.TUI.RefreshInterval = value
	case "tui.vim_keys":
		b, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("invalid tui.vim_keys %q: must be true or false", value)
		}
		cfg.TUI.VimKeys = b
	case "tui.mouse":
		b, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("invalid tui.mouse %q: must be true or false", value)
		}
		cfg.TUI.Mouse = b
	case "tui.animations":
		b, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("invalid tui.animations %q: must be true or false", value)
		}
		cfg.TUI.Animations = b
	case "defaults.priority":
		cfg.Defaults.Priority = value
	case "defaults.task_type":
		cfg.Defaults.TaskType = value
	case "defaults.calendar_days":
		n, err := strconv.Atoi(value)
		if err != nil || n < 1 {
			return fmt.Errorf("invalid defaults.calendar_days %q: must be a positive integer", value)
		}
		cfg.Defaults.CalendarDays = n
	case "defaults.habit_schedule_days":
		n, err := strconv.Atoi(value)
		if err != nil || n < 1 {
			return fmt.Errorf("invalid defaults.habit_schedule_days %q: must be a positive integer", value)
		}
		cfg.Defaults.HabitScheduleDays = n
	default:
		return fmt.Errorf("unknown config key: %q (use 'mynow config show' to list valid keys)", key)
	}
	return nil
}

// Reset removes the config file, reverting to defaults on next load.
func Reset(cfg *Config) error {
	if err := os.Remove(cfg.ConfigFile); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove config file: %w", err)
	}
	return nil
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

// isLocalhost checks if a hostname is localhost (127.0.0.1, ::1, or "localhost").
// The host parameter must already be stripped of port (e.g., via url.URL.Hostname()).
func isLocalhost(host string) bool {
	host = strings.ToLower(host)
	return host == "localhost" ||
		host == "127.0.0.1" ||
		host == "::1"
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
