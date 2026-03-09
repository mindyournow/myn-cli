package app

import (
	"github.com/mindyournow/myn-cli/internal/api"
	"github.com/mindyournow/myn-cli/internal/auth"
	"github.com/mindyournow/myn-cli/internal/config"
)

// App is the central application struct shared by CLI and TUI.
type App struct {
	Config *config.Config
	Client *api.Client
	Auth   *auth.Manager
}

// New creates a new App instance.
func New() *App {
	cfg := config.Load()
	client := api.NewClient(cfg.BaseURL())
	authManager := auth.NewManager(cfg.BaseURL())

	return &App{
		Config: cfg,
		Client: client,
		Auth:   authManager,
	}
}

// Login performs authentication based on flags.
func (a *App) Login(device bool, apiKey string) error {
	if apiKey != "" {
		return a.Auth.LoginWithAPIKey(apiKey)
	}
	if device {
		return a.loginDevice()
	}
	return a.Auth.LoginWithOAuth()
}

// loginDevice performs device authorization flow (not yet implemented).
func (a *App) loginDevice() error {
	return ErrNotImplemented
}

// Logout clears stored credentials.
func (a *App) Logout() error {
	return a.Auth.Logout()
}

// InboxAdd adds an item to the inbox.
func (a *App) InboxAdd(title string) error {
	return ErrNotImplemented
}

// InboxList lists inbox items.
func (a *App) InboxList() error {
	return ErrNotImplemented
}

// NowList shows current focus/now items.
func (a *App) NowList() error {
	return ErrNotImplemented
}

// NowFocus enters focus mode.
func (a *App) NowFocus() error {
	return ErrNotImplemented
}


// ReviewDaily runs the daily review.
func (a *App) ReviewDaily() error {
	return ErrNotImplemented
}

// RunTUI launches the TUI.
func (a *App) RunTUI() error {
	return ErrNotImplemented
}

// PluginList lists available plugins.
func (a *App) PluginList() error {
	return ErrNotImplemented
}

// PluginEnable enables a plugin.
func (a *App) PluginEnable(name string) error {
	return ErrNotImplemented
}
