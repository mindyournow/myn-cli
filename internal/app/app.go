package app

import (
	"fmt"

	"github.com/mindyournow/myn-cli/internal/api"
	"github.com/mindyournow/myn-cli/internal/config"
)

// App is the central application struct shared by CLI and TUI.
type App struct {
	Config *config.Config
	Client *api.Client
}

func New() *App {
	cfg := config.Load()
	client := api.NewClient(cfg.BaseURL)
	return &App{
		Config: cfg,
		Client: client,
	}
}

func (a *App) Login(device bool) error {
	if device {
		fmt.Println("Device authorization flow not yet implemented.")
		return nil
	}
	fmt.Println("Browser-based login not yet implemented.")
	return nil
}

func (a *App) Logout() error {
	fmt.Println("Logout not yet implemented.")
	return nil
}

func (a *App) InboxAdd(title string) error {
	fmt.Printf("Inbox add not yet implemented: %s\n", title)
	return nil
}

func (a *App) InboxList() error {
	fmt.Println("Inbox list not yet implemented.")
	return nil
}

func (a *App) NowList() error {
	fmt.Println("Now list not yet implemented.")
	return nil
}

func (a *App) NowFocus() error {
	fmt.Println("Now focus not yet implemented.")
	return nil
}

func (a *App) TaskDone(id string) error {
	fmt.Printf("Task done not yet implemented: %s\n", id)
	return nil
}

func (a *App) TaskSnooze(id string) error {
	fmt.Printf("Task snooze not yet implemented: %s\n", id)
	return nil
}

func (a *App) ReviewDaily() error {
	fmt.Println("Daily review not yet implemented.")
	return nil
}

func (a *App) RunTUI() error {
	fmt.Println("TUI not yet implemented.")
	return nil
}

func (a *App) PluginList() error {
	fmt.Println("Plugin list not yet implemented.")
	return nil
}

func (a *App) PluginEnable(name string) error {
	fmt.Printf("Plugin enable not yet implemented: %s\n", name)
	return nil
}
