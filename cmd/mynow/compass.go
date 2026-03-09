package main

import (
	"github.com/mindyournow/myn-cli/internal/app"
	"github.com/spf13/cobra"
)

func newCompassCmd(a *app.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "compass",
		Short: "Daily Compass briefing",
	}

	cmd.AddCommand(
		newCompassShowCmd(a),
		newCompassGenerateCmd(a),
		newCompassCorrectCmd(a),
		newCompassCompleteCmd(a),
		newCompassStatusCmd(a),
		newCompassHistoryCmd(a),
	)

	return cmd
}

func newCompassShowCmd(a *app.App) *cobra.Command {
	return &cobra.Command{
		Use:     "show",
		Short:   "Show today's Compass briefing",
		Aliases: []string{"display", "view"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return a.CompassShow()
		},
	}
}

func newCompassGenerateCmd(a *app.App) *cobra.Command {
	return &cobra.Command{
		Use:   "generate",
		Short: "Generate a new Compass briefing",
		RunE: func(cmd *cobra.Command, args []string) error {
			return a.CompassGenerate()
		},
	}
}

func newCompassCorrectCmd(a *app.App) *cobra.Command {
	return &cobra.Command{
		Use:   "correct",
		Short: "Correct yesterday's Compass results",
		RunE: func(cmd *cobra.Command, args []string) error {
			return a.CompassCorrect()
		},
	}
}

func newCompassCompleteCmd(a *app.App) *cobra.Command {
	return &cobra.Command{
		Use:   "complete <item-id>",
		Short: "Mark a Compass item as completed",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return a.CompassComplete(args[0])
		},
	}
}

func newCompassStatusCmd(a *app.App) *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show Compass generation status",
		RunE: func(cmd *cobra.Command, args []string) error {
			return a.CompassStatus()
		},
	}
}

func newCompassHistoryCmd(a *app.App) *cobra.Command {
	return &cobra.Command{
		Use:   "history",
		Short: "Show Compass history",
		RunE: func(cmd *cobra.Command, args []string) error {
			return a.CompassHistory()
		},
	}
}
