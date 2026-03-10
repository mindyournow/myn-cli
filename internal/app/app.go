package app

import (
	"context"
	"fmt"

	"github.com/mindyournow/myn-cli/internal/api"
	"github.com/mindyournow/myn-cli/internal/auth"
	"github.com/mindyournow/myn-cli/internal/config"
	"github.com/mindyournow/myn-cli/internal/output"
)

// App is the central application struct shared by CLI and TUI.
type App struct {
	Config    *config.Config
	Client    *api.Client
	Keyring   *auth.Keyring
	Formatter *output.Formatter
}

// New creates a new App instance.
// Returns an error if configuration cannot be loaded.
func New() (*App, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	client := api.NewClient(cfg.BaseURL)
	keyring := auth.NewKeyring(cfg.ConfigDir)

	return &App{
		Config:    cfg,
		Client:    client,
		Keyring:   keyring,
		Formatter: output.NewFormatter(false, false, false),
	}, nil
}

// SetFormatter sets the output formatter for the app.
func (a *App) SetFormatter(f *output.Formatter) {
	a.Formatter = f
}

// Login performs authentication with the MYN backend.
func (a *App) Login(ctx context.Context, device bool) error {
	if device {
		return a.Formatter.Info("Device authorization flow not yet implemented.")
	}

	oauthClient := auth.NewOAuthClient(a.Config.BaseURL, a.Keyring)
	tokens, err := oauthClient.Authenticate(ctx)
	if err != nil {
		return a.Formatter.Error(fmt.Sprintf("authentication failed: %v", err))
	}

	a.Client.SetToken(tokens.AccessToken)
	return a.Formatter.Success("Successfully authenticated!")
}

// Logout clears stored credentials.
func (a *App) Logout(ctx context.Context) error {
	if err := a.Keyring.Clear(); err != nil {
		return a.Formatter.Error(fmt.Sprintf("failed to clear credentials: %v", err))
	}
	a.Client.SetToken("")
	return a.Formatter.Success("Logged out successfully.")
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
