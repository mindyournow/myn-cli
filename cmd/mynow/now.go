package main

import (
	"github.com/mindyournow/myn-cli/internal/app"
	"github.com/spf13/cobra"
)

func newNowCmd(a *app.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "now",
		Short: "Now/Focus mode",
	}

	cmd.AddCommand(
		newNowListCmd(a),
		newNowFocusCmd(a),
		newNowCompleteCmd(a),
		newNowSnoozeCmd(a),
	)

	return cmd
}

func newNowListCmd(a *app.App) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "Show current now/focus items",
		Aliases: []string{"ls", "show"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return a.NowList()
		},
	}
}

func newNowFocusCmd(a *app.App) *cobra.Command {
	return &cobra.Command{
		Use:   "focus [id]",
		Short: "Enter focus mode on a task",
		RunE: func(cmd *cobra.Command, args []string) error {
			id := ""
			if len(args) > 0 {
				id = args[0]
			}
			return a.NowFocus(id)
		},
	}
}

func newNowCompleteCmd(a *app.App) *cobra.Command {
	return &cobra.Command{
		Use:   "complete [id]",
		Short: "Complete the current focus task",
		RunE: func(cmd *cobra.Command, args []string) error {
			id := ""
			if len(args) > 0 {
				id = args[0]
			}
			return a.NowComplete(id)
		},
	}
}

func newNowSnoozeCmd(a *app.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "snooze [duration]",
		Short: "Snooze the current focus task",
		RunE: func(cmd *cobra.Command, args []string) error {
			duration := "1h"
			if len(args) > 0 {
				duration = args[0]
			}
			return a.NowSnooze(duration)
		},
	}
	cmd.Flags().String("id", "", "Task ID (instead of current focus)")
	return cmd
}
