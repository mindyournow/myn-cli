package main

import (
	"fmt"
	"os"

	"github.com/mindyournow/myn-cli/internal/app"
	"github.com/mindyournow/myn-cli/internal/config"
	mynerrors "github.com/mindyournow/myn-cli/internal/errors"
	"github.com/mindyournow/myn-cli/internal/output"
	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	// application is nil until PersistentPreRunE runs (BUG-3 fix: defer config loading to after flag parsing)
	var application *app.App

	var (
		jsonFlag    bool
		quietFlag   bool
		noColorFlag bool
		apiURLFlag  string
		debugFlag   bool
	)

	rootCmd := &cobra.Command{
		Use:   "mynow",
		Short: "Mind Your Now — CLI & TUI for MYN",
		Long:  "A fast, scriptable, Linux-native terminal client for Mind Your Now.",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// Initialize the app here, after cobra has parsed all flags (BUG-3 fix)
			cfg, err := config.LoadWithOverrides(apiURLFlag)
			if err != nil {
				return fmt.Errorf("configuration error: %w", err)
			}
			application = app.NewWithConfig(cfg)
			application.SetFormatter(output.NewFormatter(jsonFlag, quietFlag, noColorFlag))
			_ = debugFlag // TODO: wire debug flag to logger
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			// Default action: launch TUI
			return application.RunTUI(cmd.Context())
		},
	}

	// Silence usage on errors (B14 fix)
	rootCmd.SilenceUsage = true

	// Global flags (Spec §4.1)
	rootCmd.PersistentFlags().BoolVar(&jsonFlag, "json", false, "Output in JSON format")
	rootCmd.PersistentFlags().BoolVar(&quietFlag, "quiet", false, "Suppress non-essential output")
	rootCmd.PersistentFlags().BoolVar(&noColorFlag, "no-color", false, "Disable color output")
	rootCmd.PersistentFlags().StringVar(&apiURLFlag, "api-url", "", "Override API base URL (default: https://api.mindyournow.com)")
	rootCmd.PersistentFlags().BoolVar(&debugFlag, "debug", false, "Enable debug logging")

	rootCmd.AddCommand(
		newVersionCmd(),
		newLoginCmd(&application),
		newLogoutCmd(&application),
		newWhoamiCmd(&application),
		newAuthCmd(&application),
		newInboxCmd(&application),
		newNowCmd(&application),
		newTaskCmd(&application),
		newReviewCmd(&application),
		newTUICmd(&application),
		newPluginCmd(&application),
		newConfigCmd(&application),
	)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(mynerrors.ExitCode(err))
	}
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

func newLoginCmd(a **app.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login",
		Short: "Authenticate with MYN backend",
		RunE: func(cmd *cobra.Command, args []string) error {
			device, _ := cmd.Flags().GetBool("device")
			apiKeyFlag, _ := cmd.Flags().GetBool("api-key")
			if apiKeyFlag {
				fmt.Fprint(cmd.OutOrStdout(), "Enter your MYN API key: ")
				var key string
				if _, err := fmt.Fscan(cmd.InOrStdin(), &key); err != nil {
					return fmt.Errorf("failed to read API key: %w", err)
				}
				return (*a).LoginAPIKey(cmd.Context(), key)
			}
			return (*a).Login(cmd.Context(), device)
		},
	}
	cmd.Flags().Bool("device", false, "Use device authorization flow (headless environments)")
	cmd.Flags().Bool("api-key", false, "Authenticate using an API key")
	return cmd
}

func newWhoamiCmd(a **app.App) *cobra.Command {
	return &cobra.Command{
		Use:   "whoami",
		Short: "Show current authenticated user",
		RunE: func(cmd *cobra.Command, args []string) error {
			return (*a).Whoami(cmd.Context())
		},
	}
}

func newAuthCmd(a **app.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auth",
		Short: "Authentication management",
	}
	cmd.AddCommand(&cobra.Command{
		Use:   "status",
		Short: "Show current authentication status",
		RunE: func(cmd *cobra.Command, args []string) error {
			return (*a).AuthStatus(cmd.Context())
		},
	})
	cmd.AddCommand(&cobra.Command{
		Use:   "refresh",
		Short: "Force token refresh",
		RunE: func(cmd *cobra.Command, args []string) error {
			return (*a).AuthRefresh(cmd.Context())
		},
	})
	return cmd
}

func newLogoutCmd(a **app.App) *cobra.Command {
	return &cobra.Command{
		Use:   "logout",
		Short: "Clear stored credentials",
		RunE: func(cmd *cobra.Command, args []string) error {
			return (*a).Logout(cmd.Context())
		},
	}
}

func newInboxCmd(a **app.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "inbox",
		Short: "Manage inbox items",
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "add <title>",
		Short: "Add an item to the inbox",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return (*a).InboxAdd(cmd.Context(), args[0])
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "List inbox items",
		RunE: func(cmd *cobra.Command, args []string) error {
			return (*a).InboxList(cmd.Context())
		},
	})

	return cmd
}

func newNowCmd(a **app.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "now",
		Short: "Manage current focus",
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "List current focus items",
		RunE: func(cmd *cobra.Command, args []string) error {
			return (*a).NowList(cmd.Context())
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "focus",
		Short: "Show or set current focus",
		RunE: func(cmd *cobra.Command, args []string) error {
			return (*a).NowFocus(cmd.Context())
		},
	})

	return cmd
}

func newTaskCmd(a **app.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "task",
		Short: "Manage tasks",
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "done <id>",
		Short: "Mark a task as done",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return (*a).TaskDone(cmd.Context(), args[0])
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "snooze <id>",
		Short: "Snooze a task",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return (*a).TaskSnooze(cmd.Context(), args[0])
		},
	})

	return cmd
}

func newReviewCmd(a **app.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "review",
		Short: "Run review workflows",
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "daily",
		Short: "Run daily review",
		RunE: func(cmd *cobra.Command, args []string) error {
			return (*a).ReviewDaily(cmd.Context())
		},
	})

	return cmd
}

func newTUICmd(a **app.App) *cobra.Command {
	return &cobra.Command{
		Use:   "tui",
		Short: "Launch interactive TUI",
		RunE: func(cmd *cobra.Command, args []string) error {
			return (*a).RunTUI(cmd.Context())
		},
	}
}

func newPluginCmd(a **app.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "plugin",
		Short: "Manage plugins",
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "List installed plugins",
		RunE: func(cmd *cobra.Command, args []string) error {
			return (*a).PluginList(cmd.Context())
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "enable <name>",
		Short: "Enable a plugin",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return (*a).PluginEnable(cmd.Context(), args[0])
		},
	})

	return cmd
}

func newConfigCmd(a **app.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage configuration",
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "show",
		Short: "Print resolved configuration (secrets redacted)",
		RunE: func(cmd *cobra.Command, args []string) error {
			return (*a).ConfigShow(cmd.Context())
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "get <key>",
		Short: "Get a configuration value",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return (*a).ConfigGet(cmd.Context(), args[0])
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "set <key> <value>",
		Short: "Set a configuration value",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return (*a).ConfigSet(cmd.Context(), args[0], args[1])
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "reset",
		Short: "Reset configuration to defaults",
		RunE: func(cmd *cobra.Command, args []string) error {
			return (*a).ConfigReset(cmd.Context())
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "path",
		Short: "Print the config file path",
		RunE: func(cmd *cobra.Command, args []string) error {
			return (*a).ConfigPath(cmd.Context())
		},
	})

	return cmd
}
