package main

import (
	"github.com/mindyournow/myn-cli/internal/app"
	"github.com/spf13/cobra"
)

func newTimerCmd(a *app.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "timer",
		Short: "Manage timers",
	}

	cmd.AddCommand(
		newTimerStartCmd(a),
		newTimerAlarmCmd(a),
		newTimerPauseCmd(a),
		newTimerResumeCmd(a),
		newTimerCompleteCmd(a),
		newTimerDismissCmd(a),
		newTimerCountCmd(a),
	)

	return cmd
}

func newTimerStartCmd(a *app.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start <duration> [name]",
		Short: "Start a countdown timer",
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := ""
			if len(args) > 1 {
				name = args[1]
			}
			return a.TimerStart(args[0], name)
		},
	}
	return cmd
}

func newTimerAlarmCmd(a *app.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "alarm <time> [name]",
		Short: "Set an alarm for a specific time",
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := ""
			if len(args) > 1 {
				name = args[1]
			}
			return a.TimerAlarm(args[0], name)
		},
	}
	return cmd
}

func newTimerPauseCmd(a *app.App) *cobra.Command {
	return &cobra.Command{
		Use:   "pause [id]",
		Short: "Pause a timer",
		RunE: func(cmd *cobra.Command, args []string) error {
			id := ""
			if len(args) > 0 {
				id = args[0]
			}
			return a.TimerPause(id)
		},
	}
}

func newTimerResumeCmd(a *app.App) *cobra.Command {
	return &cobra.Command{
		Use:   "resume [id]",
		Short: "Resume a paused timer",
		RunE: func(cmd *cobra.Command, args []string) error {
			id := ""
			if len(args) > 0 {
				id = args[0]
			}
			return a.TimerResume(id)
		},
	}
}

func newTimerCompleteCmd(a *app.App) *cobra.Command {
	return &cobra.Command{
		Use:   "complete [id]",
		Short: "Complete a timer early",
		RunE: func(cmd *cobra.Command, args []string) error {
			id := ""
			if len(args) > 0 {
				id = args[0]
			}
			return a.TimerComplete(id)
		},
	}
}

func newTimerDismissCmd(a *app.App) *cobra.Command {
	return &cobra.Command{
		Use:   "dismiss [id]",
		Short: "Dismiss an alarm",
		RunE: func(cmd *cobra.Command, args []string) error {
			id := ""
			if len(args) > 0 {
				id = args[0]
			}
			return a.TimerDismiss(id)
		},
	}
}

func newTimerCountCmd(a *app.App) *cobra.Command {
	return &cobra.Command{
		Use:   "count",
		Short: "Show active timer count",
		RunE: func(cmd *cobra.Command, args []string) error {
			return a.TimerCount()
		},
	}
}
