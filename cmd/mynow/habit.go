package main

import (
	"github.com/mindyournow/myn-cli/internal/app"
	"github.com/spf13/cobra"
)

func newHabitCmd(a *app.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "habit",
		Short: "Manage habits",
	}

	cmd.AddCommand(
		newHabitListCmd(a),
		newHabitDoneCmd(a),
		newHabitSkipCmd(a),
		newHabitStreakCmd(a),
		newHabitChainsCmd(a),
		newHabitScheduleCmd(a),
		newHabitRemindersCmd(a),
	)

	return cmd
}

func newHabitListCmd(a *app.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List habits",
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return a.HabitList()
		},
	}
	cmd.Flags().Bool("today", false, "Show only today's habits")
	cmd.Flags().Int("days", 7, "Number of days to show")
	return cmd
}

func newHabitDoneCmd(a *app.App) *cobra.Command {
	return &cobra.Command{
		Use:   "done <id>",
		Short: "Mark a habit as done",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return a.HabitDone(args[0])
		},
	}
}

func newHabitSkipCmd(a *app.App) *cobra.Command {
	return &cobra.Command{
		Use:   "skip <id>",
		Short: "Skip a habit for today",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return a.HabitSkip(args[0])
		},
	}
}

func newHabitStreakCmd(a *app.App) *cobra.Command {
	return &cobra.Command{
		Use:   "streak [id]",
		Short: "Show habit streaks",
		RunE: func(cmd *cobra.Command, args []string) error {
			id := ""
			if len(args) > 0 {
				id = args[0]
			}
			return a.HabitStreak(id)
		},
	}
}

func newHabitChainsCmd(a *app.App) *cobra.Command {
	return &cobra.Command{
		Use:   "chains",
		Short: "Manage habit chains",
		RunE: func(cmd *cobra.Command, args []string) error {
			return a.HabitChains()
		},
	}
}

func newHabitScheduleCmd(a *app.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "schedule <id>",
		Short: "View or edit habit schedule",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return a.HabitSchedule(args[0])
		},
	}
	cmd.Flags().StringArray("days", nil, "Days of week (mon,tue,wed,thu,fri,sat,sun)")
	return cmd
}

func newHabitRemindersCmd(a *app.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reminders <id>",
		Short: "Manage habit reminders",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return a.HabitReminders(args[0])
		},
	}
	cmd.Flags().StringArray("time", nil, "Reminder times (e.g., 09:00, 18:00)")
	return cmd
}
