package main

import (
	"fmt"
	"os"

	"github.com/mindyournow/myn-cli/internal/app"
	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	application := app.New()

	rootCmd := &cobra.Command{
		Use:   "mynow",
		Short: "Mind Your Now — CLI & TUI for MYN",
		Long:  "A fast, scriptable, Linux-native terminal client for Mind Your Now.",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Default action: launch TUI
			return application.RunTUI()
		},
	}

	rootCmd.AddCommand(
		newVersionCmd(),
		newLoginCmd(application),
		newLogoutCmd(application),
		newInboxCmd(application),
		newNowCmd(application),
		newTaskCmd(application),
		newCompassCmd(application),
		newHabitCmd(application),
		newCalendarCmd(application),
		newTimerCmd(application),
		newGroceryCmd(application),
		newReviewCmd(application),
		newTUICmd(application),
		newPluginCmd(application),
	)

	// Global flags
	rootCmd.PersistentFlags().Bool("json", false, "Output in JSON format")
	rootCmd.PersistentFlags().Bool("quiet", false, "Suppress non-essential output")
	rootCmd.PersistentFlags().Bool("no-color", false, "Disable color output")

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
	_ = application
}

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("mynow %s (commit: %s, built: %s)\n", version, commit, date)
		},
	}
}

func newLoginCmd(a *app.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login",
		Short: "Authenticate with MYN backend",
		RunE: func(cmd *cobra.Command, args []string) error {
			device, _ := cmd.Flags().GetBool("device")
			apiKey, _ := cmd.Flags().GetString("api-key")
			return a.Login(device, apiKey)
		},
	}
	cmd.Flags().Bool("device", false, "Use device authorization flow (headless environments)")
	cmd.Flags().String("api-key", "", "Authenticate with an API key")
	return cmd
}

func newLogoutCmd(a *app.App) *cobra.Command {
	return &cobra.Command{
		Use:   "logout",
		Short: "Clear stored credentials",
		RunE: func(cmd *cobra.Command, args []string) error {
			return a.Logout()
		},
	}
}

func newReviewCmd(a *app.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "review",
		Short: "Run review workflows",
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "daily",
		Short: "Run daily review",
		RunE: func(cmd *cobra.Command, args []string) error {
			return a.ReviewDaily()
		},
	})

	return cmd
}

func newTUICmd(a *app.App) *cobra.Command {
	return &cobra.Command{
		Use:   "tui",
		Short: "Launch interactive TUI",
		RunE: func(cmd *cobra.Command, args []string) error {
			return a.RunTUI()
		},
	}
}

func newPluginCmd(a *app.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "plugin",
		Short: "Manage plugins",
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "List installed plugins",
		RunE: func(cmd *cobra.Command, args []string) error {
			return a.PluginList()
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "enable [name]",
		Short: "Enable a plugin",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return a.PluginEnable(args[0])
		},
	})

	return cmd
}
