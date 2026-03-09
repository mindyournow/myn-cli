package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

const (
	DefaultBaseURL = "https://api.mindyournow.com"
	AppName        = "mynow"
)

// APIConfig holds API-related configuration.
type APIConfig struct {
	URL     string        `yaml:"url"`
	Timeout time.Duration `yaml:"timeout"`
	Retries int           `yaml:"retries"`
}

// AuthConfig holds authentication configuration.
type AuthConfig struct {
	Method  string `yaml:"method"`
	Keyring string `yaml:"keyring"`
}

// DisplayConfig holds display/output configuration.
type DisplayConfig struct {
	Color          string `yaml:"color"`
	DateFormat     string `yaml:"date_format"`
	TimeFormat     string `yaml:"time_format"`
	DefaultOutput  string `yaml:"default_output"`
}

// TUIConfig holds TUI-specific configuration.
type TUIConfig struct {
	Theme            string        `yaml:"theme"`
	RefreshInterval  time.Duration `yaml:"refresh_interval"`
	VimKeys          bool          `yaml:"vim_keys"`
	Mouse            bool          `yaml:"mouse"`
	Animations       bool          `yaml:"animations"`
}

// DefaultsConfig holds default values for new items.
type DefaultsConfig struct {
	Priority            string `yaml:"priority"`
	TaskType            string `yaml:"task_type"`
	CalendarDays        int    `yaml:"calendar_days"`
	HabitScheduleDays   int    `yaml:"habit_schedule_days"`
}

// Config is the complete application configuration.
type Config struct {
	API       APIConfig      `yaml:"api"`
	Auth      AuthConfig     `yaml:"auth"`
	Display   DisplayConfig  `yaml:"display"`
	TUI       TUIConfig      `yaml:"tui"`
	Defaults  DefaultsConfig `yaml:"defaults"`

	// Internal fields (not persisted)
	ConfigDir string `yaml:"-"`
	PluginDir string `yaml:"-"`
	DataDir   string `yaml:"-"`
	CacheDir  string `yaml:"-"`
}

// DefaultConfig returns a Config with default values.
func DefaultConfig() *Config {
	return &Config{
		API: APIConfig{
			URL:     DefaultBaseURL,
			Timeout: 30 * time.Second,
			Retries: 3,
		},
		Auth: AuthConfig{
			Method:  "api-key",
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
			RefreshInterval: 30 * time.Second,
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

// Load loads configuration from file and environment variables.
func Load() *Config {
	cfg := DefaultConfig()

	// Set up directory paths
	cfg.ConfigDir = defaultConfigDir()
	cfg.PluginDir = filepath.Join(cfg.ConfigDir, "plugins")
	cfg.DataDir = defaultDataDir()
	cfg.CacheDir = defaultCacheDir()

	// Ensure config directory exists
	_ = os.MkdirAll(cfg.ConfigDir, 0755)

	// Load from config file if it exists
	configFile := cfg.FilePath()
	if data, err := os.ReadFile(configFile); err == nil {
		_ = yaml.Unmarshal(data, cfg)
	}

	// Apply environment variable overrides
	cfg.applyEnvVars()

	return cfg
}

// FilePath returns the path to the config file.
func (c *Config) FilePath() string {
	if path := os.Getenv("MYNOW_CONFIG"); path != "" {
		return path
	}
	return filepath.Join(c.ConfigDir, "config.yaml")
}

// Save writes the configuration to the config file.
func (c *Config) Save() error {
	// Ensure config directory exists
	if err := os.MkdirAll(c.ConfigDir, 0755); err != nil {
		return fmt.Errorf("creating config directory: %w", err)
	}

	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("marshaling config: %w", err)
	}

	// Add header comment
	header := "# MYN CLI configuration\n"
	data = append([]byte(header), data...)

	if err := os.WriteFile(c.FilePath(), data, 0600); err != nil {
		return fmt.Errorf("writing config file: %w", err)
	}

	return nil
}

// applyEnvVars applies environment variable overrides to the config.
func (c *Config) applyEnvVars() {
	if url := os.Getenv("MYN_API_URL"); url != "" {
		c.API.URL = url
	}
	if key := os.Getenv("MYN_API_KEY"); key != "" {
		// API key is handled by auth package, but we note it's set
		c.Auth.Method = "api-key"
	}
	if keyring := os.Getenv("MYNOW_KEYRING"); keyring != "" {
		c.Auth.Keyring = keyring
	}
	if os.Getenv("NO_COLOR") != "" {
		c.Display.Color = "never"
	}
}

// BaseURL returns the API base URL.
func (c *Config) BaseURL() string {
	return c.API.URL
}

// defaultConfigDir returns the default configuration directory.
func defaultConfigDir() string {
	if dir := os.Getenv("XDG_CONFIG_HOME"); dir != "" {
		return filepath.Join(dir, AppName)
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", AppName)
}

// defaultDataDir returns the default data directory.
func defaultDataDir() string {
	if dir := os.Getenv("XDG_DATA_HOME"); dir != "" {
		return filepath.Join(dir, AppName)
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".local", "share", AppName)
}

// defaultCacheDir returns the default cache directory.
func defaultCacheDir() string {
	if dir := os.Getenv("XDG_CACHE_HOME"); dir != "" {
		return filepath.Join(dir, AppName)
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".cache", AppName)
}
