package main

import (
	"github.com/mindyournow/myn-cli/internal/app"
	"github.com/spf13/cobra"
)

func newCalendarCmd(a *app.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "calendar",
		Short: "Manage calendar events",
		Aliases: []string{"cal"},
	}

	cmd.AddCommand(
		newCalendarListCmd(a),
		newCalendarAddCmd(a),
		newCalendarDeleteCmd(a),
		newCalendarDeclineCmd(a),
		newCalendarSkipCmd(a),
	)

	return cmd
}

func newCalendarListCmd(a *app.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List calendar events",
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return a.CalendarList()
		},
	}
	cmd.Flags().Int("days", 7, "Number of days to show")
	return cmd
}

func newCalendarAddCmd(a *app.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add <title>",
		Short: "Add a calendar event",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return a.CalendarAdd(args[0])
		},
	}
	cmd.Flags().String("start", "", "Start time")
	cmd.Flags().String("end", "", "End time")
	cmd.Flags().Bool("all-day", false, "All-day event")
	return cmd
}

func newCalendarDeleteCmd(a *app.App) *cobra.Command {
	return &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a calendar event",
		Aliases: []string{"rm", "remove"},
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return a.CalendarDelete(args[0])
		},
	}
}

func newCalendarDeclineCmd(a *app.App) *cobra.Command {
	return &cobra.Command{
		Use:   "decline <id>",
		Short: "Decline a calendar invitation",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return a.CalendarDecline(args[0])
		},
	}
}

func newCalendarSkipCmd(a *app.App) *cobra.Command {
	return &cobra.Command{
		Use:   "skip <id>",
		Short: "Skip a recurring calendar event",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return a.CalendarSkip(args[0])
		},
	}
}
