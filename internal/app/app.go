package app

import (
	"context"
	"fmt"
	"time"

	"github.com/mindyournow/myn-cli/internal/api"
	"github.com/mindyournow/myn-cli/internal/auth"
	"github.com/mindyournow/myn-cli/internal/config"
	"github.com/mindyournow/myn-cli/internal/output"
)

// App is the central application struct shared by CLI and TUI.
type App struct {
	Config     *config.Config
	Client     *api.Client
	Keyring    *auth.Keyring
	KeyStore   *auth.KeyStore
	TokenCache *auth.TokenCache
	Formatter  *output.Formatter
}

// New creates a new App instance using environment variables and defaults.
// Returns an error if configuration cannot be loaded.
func New() (*App, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}
	return NewWithConfig(cfg), nil
}

// NewWithConfig creates an App from an already-loaded Config.
// Used when flags (e.g., --api-url) override the loaded configuration (BUG-3 fix).
func NewWithConfig(cfg *config.Config) *App {
	fileKeyring := auth.NewKeyring(cfg.ConfigDir)
	keyStore := auth.NewKeyStore(fileKeyring, cfg.Auth.Keyring)
	oauthClient := auth.NewOAuthClient(cfg.BaseURL, keyStore)
	tokenCache := auth.NewTokenCache(keyStore, oauthClient)

	return &App{
		Config:     cfg,
		Client:     api.NewClient(cfg.BaseURL),
		Keyring:    fileKeyring,
		KeyStore:   keyStore,
		TokenCache: tokenCache,
		Formatter:  output.NewFormatter(false, false, false),
	}
}

// SetFormatter sets the output formatter for the app.
func (a *App) SetFormatter(f *output.Formatter) {
	a.Formatter = f
}

// Login performs OAuth PKCE authentication.
func (a *App) Login(ctx context.Context, device bool) error {
	if device {
		d := auth.NewDeviceClient(a.Config.BaseURL)
		if err := d.Authorize(ctx); err != nil {
			_ = a.Formatter.Error(err.Error())
			return err
		}
		return nil
	}

	oauthClient := auth.NewOAuthClient(a.Config.BaseURL, a.KeyStore)
	tokens, err := oauthClient.Authenticate(ctx)
	if err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("authentication failed: %v", err))
		return err
	}

	ttl := time.Duration(tokens.ExpiresIn) * time.Second
	if ttl <= 0 {
		ttl = 3600 * time.Second
	}
	a.TokenCache.SetAccessToken(tokens.AccessToken, ttl)
	a.TokenCache.SetAuthMethod(auth.AuthMethodOAuth)
	a.Client.SetToken(tokens.AccessToken)
	return a.Formatter.Success("Successfully authenticated!")
}

// LoginAPIKey authenticates using an API key.
func (a *App) LoginAPIKey(ctx context.Context, apiKey string) error {
	client := auth.NewAPIKeyClient(a.Config.BaseURL, a.KeyStore)
	profile, err := client.Login(ctx, apiKey)
	if err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("API key authentication failed: %v", err))
		return err
	}
	a.Client.SetAPIKey(apiKey)
	a.TokenCache.SetAuthMethod(auth.AuthMethodAPIKey)
	return a.Formatter.Success(fmt.Sprintf("Authenticated as %s (%s) using API key.", profile.Name, profile.Email))
}

// Logout clears stored credentials.
func (a *App) Logout(ctx context.Context) error {
	if err := a.KeyStore.Clear(); err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to clear credentials: %v", err))
		return err
	}
	a.Client.SetToken("")
	a.Client.SetAPIKey("")
	return a.Formatter.Success("Logged out successfully.")
}

// Whoami displays the current authenticated user's profile.
func (a *App) Whoami(ctx context.Context) error {
	// Try API key first
	apiKey, _ := a.KeyStore.LoadAPIKey()
	if apiKey != "" {
		client := auth.NewAPIKeyClient(a.Config.BaseURL, a.KeyStore)
		profile, err := client.Validate(ctx, apiKey)
		if err != nil {
			_ = a.Formatter.Error(fmt.Sprintf("failed to get profile: %v", err))
			return err
		}
		if a.Formatter.JSON {
			return a.Formatter.Print(profile)
		}
		_ = a.Formatter.Println(fmt.Sprintf("Name:     %s", profile.Name))
		_ = a.Formatter.Println(fmt.Sprintf("Email:    %s", profile.Email))
		_ = a.Formatter.Println(fmt.Sprintf("Username: %s", profile.Username))
		return a.Formatter.Println("Auth:     API key")
	}

	// OAuth: ensure we have a valid access token
	accessToken, err := a.TokenCache.GetAccessToken(ctx)
	if err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("not authenticated (run 'mynow login'): %v", err))
		return err
	}
	a.Client.SetToken(accessToken)

	resp, err := a.Client.Get(ctx, "/api/v1/customers", nil)
	if err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to get profile: %v", err))
		return err
	}
	var profile auth.CustomerProfile
	if err := resp.DecodeJSON(&profile); err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to parse profile: %v", err))
		return err
	}
	if a.Formatter.JSON {
		return a.Formatter.Print(profile)
	}
	_ = a.Formatter.Println(fmt.Sprintf("Name:     %s", profile.Name))
	_ = a.Formatter.Println(fmt.Sprintf("Email:    %s", profile.Email))
	_ = a.Formatter.Println(fmt.Sprintf("Username: %s", profile.Username))
	return a.Formatter.Println("Auth:     OAuth")
}

// AuthStatus shows the current authentication status.
func (a *App) AuthStatus(ctx context.Context) error {
	method := a.TokenCache.GetAuthMethod()
	if method == "" {
		// Try to detect from stored credentials
		if apiKey, err := a.KeyStore.LoadAPIKey(); err == nil && apiKey != "" {
			method = auth.AuthMethodAPIKey
		} else if _, err := a.KeyStore.LoadRefreshToken(); err == nil {
			method = auth.AuthMethodOAuth
		} else {
			return a.Formatter.Println("Not authenticated. Run 'mynow login' to authenticate.")
		}
	}
	if a.Formatter.JSON {
		return a.Formatter.Print(map[string]interface{}{
			"method":    string(method),
			"expiresAt": a.TokenCache.ExpiresAt().Format(time.RFC3339),
		})
	}
	_ = a.Formatter.Println(fmt.Sprintf("Auth method:   %s", method))
	if method == auth.AuthMethodOAuth && !a.TokenCache.ExpiresAt().IsZero() {
		_ = a.Formatter.Println(fmt.Sprintf("Token expires: %s", a.TokenCache.ExpiresAt().Format("2006-01-02 15:04:05")))
	}
	_ = a.Formatter.Println(fmt.Sprintf("Keyring:       %s", a.Config.Auth.Keyring))
	return nil
}

// AuthRefresh forces a token refresh.
func (a *App) AuthRefresh(ctx context.Context) error {
	if a.TokenCache.GetAuthMethod() == auth.AuthMethodAPIKey {
		return a.Formatter.Info("API key auth does not require token refresh.")
	}
	_, err := a.TokenCache.Refresh(ctx)
	if err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("token refresh failed: %v", err))
		return err
	}
	return a.Formatter.Success("Token refreshed successfully.")
}

// InboxAdd adds an item to the inbox.
func (a *App) InboxAdd(ctx context.Context, title string) error {
	return a.Formatter.Printf("Inbox add not yet implemented: %s", title)
}

// InboxList lists inbox items.
func (a *App) InboxList(ctx context.Context) error {
	return a.Formatter.Println("Inbox list not yet implemented.")
}

// NowList lists current focus items.
func (a *App) NowList(ctx context.Context) error {
	return a.Formatter.Println("Now list not yet implemented.")
}

// NowFocus shows or sets current focus.
func (a *App) NowFocus(ctx context.Context) error {
	return a.Formatter.Println("Now focus not yet implemented.")
}

// TaskDone marks a task as done.
func (a *App) TaskDone(ctx context.Context, id string) error {
	return a.Formatter.Printf("Task done not yet implemented: %s", id)
}

// TaskSnooze snoozes a task.
func (a *App) TaskSnooze(ctx context.Context, id string) error {
	return a.Formatter.Printf("Task snooze not yet implemented: %s", id)
}

// ReviewDaily runs the daily review.
func (a *App) ReviewDaily(ctx context.Context) error {
	return a.Formatter.Println("Daily review not yet implemented.")
}

// RunTUI launches the interactive TUI.
func (a *App) RunTUI(ctx context.Context) error {
	return a.Formatter.Println("TUI not yet implemented.")
}

// PluginList lists installed plugins.
func (a *App) PluginList(ctx context.Context) error {
	return a.Formatter.Println("Plugin list not yet implemented.")
}

// PluginEnable enables a plugin.
func (a *App) PluginEnable(ctx context.Context, name string) error {
	return a.Formatter.Printf("Plugin enable not yet implemented: %s", name)
}

// ConfigShow prints the resolved configuration (secrets redacted).
func (a *App) ConfigShow(ctx context.Context) error {
	cfg := a.Config
	if a.Formatter.JSON {
		type jsonCfg struct {
			ConfigFile string               `json:"config_file"`
			API        config.APIConfig     `json:"api"`
			Auth       config.AuthConfig    `json:"auth"`
			Display    config.DisplayConfig `json:"display"`
			TUI        config.TUIConfig     `json:"tui"`
			Defaults   config.DefaultsConfig `json:"defaults"`
			APIKey     string               `json:"api_key,omitempty"`
		}
		out := jsonCfg{
			ConfigFile: cfg.ConfigFile,
			API:        cfg.API,
			Auth:       cfg.Auth,
			Display:    cfg.Display,
			TUI:        cfg.TUI,
			Defaults:   cfg.Defaults,
		}
		if cfg.APIKey != "" {
			out.APIKey = "***redacted***"
		}
		return a.Formatter.Print(out)
	}
	lines := []string{
		fmt.Sprintf("config file:                   %s", cfg.ConfigFile),
		fmt.Sprintf("api.url:                       %s", cfg.API.URL),
		fmt.Sprintf("api.timeout:                   %s", cfg.API.Timeout),
		fmt.Sprintf("api.retries:                   %d", cfg.API.Retries),
		fmt.Sprintf("auth.method:                   %s", cfg.Auth.Method),
		fmt.Sprintf("auth.keyring:                  %s", cfg.Auth.Keyring),
		fmt.Sprintf("display.color:                 %s", cfg.Display.Color),
		fmt.Sprintf("display.date_format:           %s", cfg.Display.DateFormat),
		fmt.Sprintf("display.time_format:           %s", cfg.Display.TimeFormat),
		fmt.Sprintf("display.default_output:        %s", cfg.Display.DefaultOutput),
		fmt.Sprintf("tui.theme:                     %s", cfg.TUI.Theme),
		fmt.Sprintf("tui.refresh_interval:          %s", cfg.TUI.RefreshInterval),
		fmt.Sprintf("tui.vim_keys:                  %v", cfg.TUI.VimKeys),
		fmt.Sprintf("tui.mouse:                     %v", cfg.TUI.Mouse),
		fmt.Sprintf("tui.animations:                %v", cfg.TUI.Animations),
		fmt.Sprintf("defaults.priority:             %s", cfg.Defaults.Priority),
		fmt.Sprintf("defaults.task_type:            %s", cfg.Defaults.TaskType),
		fmt.Sprintf("defaults.calendar_days:        %d", cfg.Defaults.CalendarDays),
		fmt.Sprintf("defaults.habit_schedule_days:  %d", cfg.Defaults.HabitScheduleDays),
	}
	if cfg.APIKey != "" {
		lines = append(lines, "api_key:                       ***redacted***")
	}
	for _, line := range lines {
		if err := a.Formatter.Println(line); err != nil {
			return err
		}
	}
	return nil
}

// ConfigGet prints the value of a single config key.
func (a *App) ConfigGet(ctx context.Context, key string) error {
	val, err := config.GetValue(a.Config, key)
	if err != nil {
		_ = a.Formatter.Error(err.Error())
		return err
	}
	return a.Formatter.Println(val)
}

// ConfigSet sets a config key and persists it to the config file.
func (a *App) ConfigSet(ctx context.Context, key, value string) error {
	if err := config.SetValue(a.Config, key, value); err != nil {
		_ = a.Formatter.Error(err.Error())
		return err
	}
	if err := config.Save(a.Config); err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to save config: %v", err))
		return err
	}
	return a.Formatter.Success(fmt.Sprintf("Set %s = %s", key, value))
}

// ConfigReset removes the config file, reverting to defaults.
func (a *App) ConfigReset(ctx context.Context) error {
	if err := config.Reset(a.Config); err != nil {
		_ = a.Formatter.Error(err.Error())
		return err
	}
	return a.Formatter.Success("Configuration reset to defaults.")
}

// ConfigPath prints the path to the config file.
func (a *App) ConfigPath(ctx context.Context) error {
	return a.Formatter.Println(a.Config.ConfigFile)
}
