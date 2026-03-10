package main

import (
	"fmt"
	"os"
	"strings"

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
		newCompassCmd(&application),
		newCalendarCmd(&application),
		newTimerCmd(&application),
		newGroceryCmd(&application),
		newProjectCmd(&application),
		newSearchCmd(&application),
		newProfileCmd(&application),
		newMemoryCmd(&application),
		newHouseholdCmd(&application),
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

	// task list
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List tasks",
		RunE: func(cmd *cobra.Command, args []string) error {
			priority, _ := cmd.Flags().GetString("priority")
			taskType, _ := cmd.Flags().GetString("type")
			completed, _ := cmd.Flags().GetBool("completed")
			archived, _ := cmd.Flags().GetBool("archived")
			today, _ := cmd.Flags().GetBool("today")
			overdue, _ := cmd.Flags().GetBool("overdue")
			household, _ := cmd.Flags().GetBool("household")
			sortBy, _ := cmd.Flags().GetString("sort")
			page, _ := cmd.Flags().GetInt("page")
			limit, _ := cmd.Flags().GetInt("limit")
			return (*a).TaskListFull(cmd.Context(), app.TaskListOptions{
				Priority: priority, Type: taskType, Completed: completed,
				Archived: archived, Today: today, Overdue: overdue,
				Household: household, Sort: sortBy, Page: page, Limit: limit,
			})
		},
	}
	listCmd.Flags().String("priority", "", "Filter by priority (critical, opportunity, horizon, parking)")
	listCmd.Flags().String("type", "", "Filter by type (task, habit, chore)")
	listCmd.Flags().Bool("completed", false, "Include completed tasks")
	listCmd.Flags().Bool("archived", false, "Show archived tasks")
	listCmd.Flags().Bool("today", false, "Only tasks for today")
	listCmd.Flags().Bool("overdue", false, "Only overdue tasks")
	listCmd.Flags().Bool("household", false, "Include household tasks")
	listCmd.Flags().String("sort", "", "Sort by (priority, date, title, created)")
	listCmd.Flags().Int("page", 0, "Page number")
	listCmd.Flags().Int("limit", 50, "Page size")
	cmd.AddCommand(listCmd)

	// task add
	addCmd := &cobra.Command{
		Use:   "add <title>",
		Short: "Create a new task",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			priority, _ := cmd.Flags().GetString("priority")
			date, _ := cmd.Flags().GetString("date")
			duration, _ := cmd.Flags().GetString("duration")
			desc, _ := cmd.Flags().GetString("description")
			taskType, _ := cmd.Flags().GetString("type")
			recurrence, _ := cmd.Flags().GetString("recurrence")
			project, _ := cmd.Flags().GetString("project")
			return (*a).TaskAdd(cmd.Context(), args[0], app.TaskAddOptions{
				Priority: priority, Date: date, Duration: duration,
				Description: desc, Type: taskType, Recurrence: recurrence,
				ProjectID: project,
			})
		},
	}
	addCmd.Flags().String("priority", "", "Priority zone")
	addCmd.Flags().String("date", "", "Start date")
	addCmd.Flags().String("duration", "", "Duration (30m, 1h, etc.)")
	addCmd.Flags().String("description", "", "Description")
	addCmd.Flags().String("type", "task", "Task type (task, habit, chore)")
	addCmd.Flags().String("recurrence", "", "Recurrence (daily, weekly, RRULE)")
	addCmd.Flags().String("project", "", "Project ID")
	cmd.AddCommand(addCmd)

	// task show
	cmd.AddCommand(&cobra.Command{
		Use:   "show <id>",
		Short: "Show task details",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return (*a).TaskShow(cmd.Context(), args[0])
		},
	})

	// task edit
	editCmd := &cobra.Command{
		Use:   "edit <id>",
		Short: "Edit a task",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			title, _ := cmd.Flags().GetString("title")
			priority, _ := cmd.Flags().GetString("priority")
			date, _ := cmd.Flags().GetString("date")
			duration, _ := cmd.Flags().GetString("duration")
			desc, _ := cmd.Flags().GetString("description")
			project, _ := cmd.Flags().GetString("project")
			return (*a).TaskEdit(cmd.Context(), args[0], app.TaskEditOptions{
				Title: title, Priority: priority, Date: date,
				Duration: duration, Description: desc, ProjectID: project,
			})
		},
	}
	editCmd.Flags().String("title", "", "New title")
	editCmd.Flags().String("priority", "", "New priority")
	editCmd.Flags().String("date", "", "New start date")
	editCmd.Flags().String("duration", "", "New duration")
	editCmd.Flags().String("description", "", "New description")
	editCmd.Flags().String("project", "", "New project ID")
	cmd.AddCommand(editCmd)

	// task done
	cmd.AddCommand(&cobra.Command{
		Use:   "done <id>",
		Short: "Mark a task as done",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return (*a).TaskComplete(cmd.Context(), args[0])
		},
	})

	// task uncomplete
	cmd.AddCommand(&cobra.Command{
		Use:   "uncomplete <id>",
		Short: "Mark a completed task as not done",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return (*a).TaskUncomplete(cmd.Context(), args[0])
		},
	})

	// task delete
	deleteCmd := &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a task",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			permanent, _ := cmd.Flags().GetBool("permanent")
			return (*a).TaskDelete(cmd.Context(), args[0], permanent)
		},
	}
	deleteCmd.Flags().Bool("permanent", false, "Permanently delete (cannot be restored)")
	cmd.AddCommand(deleteCmd)

	// task restore
	cmd.AddCommand(&cobra.Command{
		Use:   "restore <id>",
		Short: "Restore a deleted task",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return (*a).TaskRestore(cmd.Context(), args[0])
		},
	})

	// task snooze
	snoozeCmd := &cobra.Command{
		Use:   "snooze <id>",
		Short: "Snooze a task to a later date",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			date, _ := cmd.Flags().GetString("date")
			days, _ := cmd.Flags().GetInt("days")
			return (*a).TaskSnoozeTask(cmd.Context(), args[0], app.TaskSnoozeOpt{Date: date, Days: days})
		},
	}
	snoozeCmd.Flags().String("date", "", "Target date (default: tomorrow)")
	snoozeCmd.Flags().Int("days", 0, "Snooze by N days")
	cmd.AddCommand(snoozeCmd)

	// task batch
	batchCmd := &cobra.Command{
		Use:   "batch",
		Short: "Apply updates to multiple tasks",
		RunE: func(cmd *cobra.Command, args []string) error {
			idsStr, _ := cmd.Flags().GetString("ids")
			priority, _ := cmd.Flags().GetString("priority")
			project, _ := cmd.Flags().GetString("project")
			date, _ := cmd.Flags().GetString("date")
			if idsStr == "" {
				return fmt.Errorf("--ids is required")
			}
			ids := splitCSV(idsStr)
			return (*a).TaskBatch(cmd.Context(), app.TaskBatchOptions{
				IDs: ids, Priority: priority, ProjectID: project, Date: date,
			})
		},
	}
	batchCmd.Flags().String("ids", "", "Comma-separated task IDs (required)")
	batchCmd.Flags().String("priority", "", "Set priority for all")
	batchCmd.Flags().String("project", "", "Move all to project ID")
	batchCmd.Flags().String("date", "", "Set start date for all")
	cmd.AddCommand(batchCmd)

	// task move
	cmd.AddCommand(&cobra.Command{
		Use:   "move <task-id> <project-id>",
		Short: "Move a task to a project",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return (*a).TaskMove(cmd.Context(), args[0], args[1])
		},
	})

	return cmd
}

// splitCSV splits a comma-separated string into a slice of trimmed strings.
func splitCSV(s string) []string {
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
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

func newCompassCmd(a **app.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "compass",
		Short: "Daily compass briefing",
		RunE: func(cmd *cobra.Command, args []string) error {
			return (*a).CompassShow(cmd.Context())
		},
	}

	generateCmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate a new compass briefing",
		RunE: func(cmd *cobra.Command, args []string) error {
			briefingType, _ := cmd.Flags().GetString("type")
			async, _ := cmd.Flags().GetBool("async")
			return (*a).CompassGenerate(cmd.Context(), briefingType, async)
		},
	}
	generateCmd.Flags().String("type", "on-demand", "Briefing type (daily, evening, weekly, on-demand)")
	generateCmd.Flags().Bool("async", false, "Don't wait for result")
	cmd.AddCommand(generateCmd)

	correctCmd := &cobra.Command{
		Use:   "correct",
		Short: "Apply a correction to the compass briefing",
		RunE: func(cmd *cobra.Command, args []string) error {
			summaryID, _ := cmd.Flags().GetString("summary-id")
			taskID, _ := cmd.Flags().GetString("task")
			decision, _ := cmd.Flags().GetString("decision")
			newDate, _ := cmd.Flags().GetString("new-date")
			reason, _ := cmd.Flags().GetString("reason")
			return (*a).CompassCorrect(cmd.Context(), app.CompassCorrectOptions{
				SummaryID: summaryID, TaskID: taskID, Decision: decision,
				NewDate: newDate, Reason: reason,
			})
		},
	}
	correctCmd.Flags().String("summary-id", "", "Compass summary ID")
	correctCmd.Flags().String("task", "", "Task ID")
	correctCmd.Flags().String("decision", "", "Decision (accepted, rejected, modified, completed, archived)")
	correctCmd.Flags().String("new-date", "", "New date (for modified decision)")
	correctCmd.Flags().String("reason", "", "Reason for correction")
	cmd.AddCommand(correctCmd)

	completeCmd := &cobra.Command{
		Use:   "complete",
		Short: "Mark the compass session as complete",
		RunE: func(cmd *cobra.Command, args []string) error {
			summary, _ := cmd.Flags().GetString("summary")
			decisionsStr, _ := cmd.Flags().GetString("decisions")
			var decisions []string
			if decisionsStr != "" {
				decisions = splitCSV(decisionsStr)
			}
			return (*a).CompassComplete(cmd.Context(), summary, decisions)
		},
	}
	completeCmd.Flags().String("summary", "", "Session summary")
	completeCmd.Flags().String("decisions", "", "Comma-separated key decisions")
	cmd.AddCommand(completeCmd)

	cmd.AddCommand(&cobra.Command{
		Use:   "status",
		Short: "Show compass session status",
		RunE: func(cmd *cobra.Command, args []string) error {
			return (*a).CompassStatus(cmd.Context())
		},
	})

	historyCmd := &cobra.Command{
		Use:   "history",
		Short: "Show past compass briefings",
		RunE: func(cmd *cobra.Command, args []string) error {
			limit, _ := cmd.Flags().GetInt("limit")
			return (*a).CompassHistory(cmd.Context(), limit)
		},
	}
	historyCmd.Flags().Int("limit", 10, "Number of entries to show")
	cmd.AddCommand(historyCmd)

	return cmd
}

func newCalendarCmd(a **app.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "calendar",
		Short: "Manage calendar events",
		RunE: func(cmd *cobra.Command, args []string) error {
			date, _ := cmd.Flags().GetString("date")
			days, _ := cmd.Flags().GetInt("days")
			return (*a).CalendarList(cmd.Context(), date, days)
		},
	}
	cmd.Flags().String("date", "", "Date (default: today)")
	cmd.Flags().Int("days", 1, "Number of days to show")

	addCmd := &cobra.Command{
		Use:   "add <title>",
		Short: "Add a calendar event",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			date, _ := cmd.Flags().GetString("date")
			start, _ := cmd.Flags().GetString("start")
			end, _ := cmd.Flags().GetString("end")
			allDay, _ := cmd.Flags().GetBool("all-day")
			location, _ := cmd.Flags().GetString("location")
			attendees, _ := cmd.Flags().GetString("attendees")
			desc, _ := cmd.Flags().GetString("description")
			recurrence, _ := cmd.Flags().GetString("recurrence")
			return (*a).CalendarAdd(cmd.Context(), args[0], app.CalendarAddOptions{
				Date: date, Start: start, End: end, AllDay: allDay,
				Location: location, Attendees: attendees,
				Description: desc, Recurrence: recurrence,
			})
		},
	}
	addCmd.Flags().String("date", "", "Date for all-day event")
	addCmd.Flags().String("start", "", "Start time")
	addCmd.Flags().String("end", "", "End time")
	addCmd.Flags().Bool("all-day", false, "All-day event")
	addCmd.Flags().String("location", "", "Location")
	addCmd.Flags().String("attendees", "", "Comma-separated attendee emails")
	addCmd.Flags().String("description", "", "Description")
	addCmd.Flags().String("recurrence", "", "Recurrence rule")
	cmd.AddCommand(addCmd)

	cmd.AddCommand(&cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a calendar event",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return (*a).CalendarDelete(cmd.Context(), args[0])
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "decline <id>",
		Short: "Decline a meeting invitation",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return (*a).CalendarDecline(cmd.Context(), args[0])
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "skip <id>",
		Short: "Skip a recurring meeting occurrence",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return (*a).CalendarSkip(cmd.Context(), args[0])
		},
	})

	return cmd
}

func newTimerCmd(a **app.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "timer",
		Short: "Manage timers",
		RunE: func(cmd *cobra.Command, args []string) error {
			completed, _ := cmd.Flags().GetBool("completed")
			return (*a).TimerList(cmd.Context(), completed)
		},
	}
	cmd.Flags().Bool("completed", false, "Include completed timers")

	startCmd := &cobra.Command{
		Use:   "start <duration>",
		Short: "Start a countdown timer (e.g. 25m, 1h)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			label, _ := cmd.Flags().GetString("label")
			return (*a).TimerStart(cmd.Context(), args[0], label)
		},
	}
	startCmd.Flags().String("label", "", "Timer label")
	cmd.AddCommand(startCmd)

	alarmCmd := &cobra.Command{
		Use:   "alarm <time>",
		Short: "Set an alarm for a specific time",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			label, _ := cmd.Flags().GetString("label")
			return (*a).TimerAlarm(cmd.Context(), args[0], label)
		},
	}
	alarmCmd.Flags().String("label", "", "Alarm label")
	cmd.AddCommand(alarmCmd)

	cmd.AddCommand(&cobra.Command{
		Use:   "pause <id>",
		Short: "Pause a running timer",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return (*a).TimerPause(cmd.Context(), args[0])
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "resume <id>",
		Short: "Resume a paused timer",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return (*a).TimerResume(cmd.Context(), args[0])
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "complete <id>",
		Short: "Mark a timer as complete",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return (*a).TimerComplete(cmd.Context(), args[0])
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "dismiss",
		Short: "Dismiss all completed timers",
		RunE: func(cmd *cobra.Command, args []string) error {
			return (*a).TimerDismiss(cmd.Context())
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "count",
		Short: "Show number of active timers",
		RunE: func(cmd *cobra.Command, args []string) error {
			return (*a).TimerCount(cmd.Context())
		},
	})

	return cmd
}

func newGroceryCmd(a **app.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "grocery",
		Short: "Manage household grocery list",
		RunE: func(cmd *cobra.Command, args []string) error {
			return (*a).GroceryList(cmd.Context())
		},
	}

	addCmd := &cobra.Command{
		Use:   "add <name>",
		Short: "Add an item to the grocery list",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			unit, _ := cmd.Flags().GetString("unit")
			qty, _ := cmd.Flags().GetFloat64("qty")
			cat, _ := cmd.Flags().GetString("category")
			return (*a).GroceryAdd(cmd.Context(), args[0], unit, qty, cat)
		},
	}
	addCmd.Flags().String("unit", "", "Unit (e.g. kg, L, pcs)")
	addCmd.Flags().Float64("qty", 0, "Quantity")
	addCmd.Flags().String("category", "", "Category")
	cmd.AddCommand(addCmd)

	cmd.AddCommand(&cobra.Command{
		Use:   "add-bulk",
		Short: "Add multiple items (one per line via stdin or --items flag)",
		RunE: func(cmd *cobra.Command, args []string) error {
			items, _ := cmd.Flags().GetString("items")
			return (*a).GroceryAddBulk(cmd.Context(), items)
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "check <id>",
		Short: "Mark a grocery item as checked",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return (*a).GroceryCheck(cmd.Context(), args[0], true)
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "uncheck <id>",
		Short: "Uncheck a grocery item",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return (*a).GroceryCheck(cmd.Context(), args[0], false)
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "delete <id>",
		Short: "Remove a grocery item",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return (*a).GroceryDelete(cmd.Context(), args[0])
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "clear",
		Short: "Remove all checked items",
		RunE: func(cmd *cobra.Command, args []string) error {
			return (*a).GroceryClear(cmd.Context())
		},
	})

	convertCmd := &cobra.Command{
		Use:   "convert",
		Short: "Convert grocery items to tasks",
		RunE: func(cmd *cobra.Command, args []string) error {
			idsStr, _ := cmd.Flags().GetString("ids")
			return (*a).GroceryConvert(cmd.Context(), splitCSV(idsStr))
		},
	}
	convertCmd.Flags().String("ids", "", "Comma-separated item IDs to convert")
	cmd.AddCommand(convertCmd)

	return cmd
}

func newProjectCmd(a **app.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "project",
		Short: "Manage projects",
		RunE: func(cmd *cobra.Command, args []string) error {
			return (*a).ProjectList(cmd.Context())
		},
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "List all projects",
		RunE: func(cmd *cobra.Command, args []string) error {
			return (*a).ProjectList(cmd.Context())
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "show <id>",
		Short: "Show project details",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return (*a).ProjectShow(cmd.Context(), args[0])
		},
	})

	createCmd := &cobra.Command{
		Use:   "create <title>",
		Short: "Create a new project",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			desc, _ := cmd.Flags().GetString("description")
			return (*a).ProjectCreate(cmd.Context(), args[0], desc)
		},
	}
	createCmd.Flags().String("description", "", "Project description")
	cmd.AddCommand(createCmd)

	return cmd
}

func newSearchCmd(a **app.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "search <query>",
		Short: "Search across all items",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			typesStr, _ := cmd.Flags().GetString("type")
			archivd, _ := cmd.Flags().GetBool("archived")
			limit, _ := cmd.Flags().GetInt("limit")
			var types []string
			if typesStr != "" {
				types = splitCSV(typesStr)
			}
			return (*a).SearchAll(cmd.Context(), args[0], app.SearchOptions{
				Types: types, IncludeArchived: archivd, Limit: limit,
			})
		},
	}
	cmd.Flags().String("type", "", "Filter by type (task, habit, etc.)")
	cmd.Flags().Bool("archived", false, "Include archived items")
	cmd.Flags().Int("limit", 20, "Max results")
	return cmd
}

func newProfileCmd(a **app.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "profile",
		Short: "Manage user profile and preferences",
		RunE: func(cmd *cobra.Command, args []string) error {
			return (*a).ProfileShow(cmd.Context())
		},
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "goals",
		Short: "Show goals",
		RunE: func(cmd *cobra.Command, args []string) error {
			return (*a).ProfileGoals(cmd.Context())
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "prefs",
		Short: "Show preferences",
		RunE: func(cmd *cobra.Command, args []string) error {
			return (*a).ProfilePrefs(cmd.Context())
		},
	})

	coachingCmd := &cobra.Command{
		Use:   "coaching [level]",
		Short: "Show or set coaching intensity (off, gentle, proactive)",
		RunE: func(cmd *cobra.Command, args []string) error {
			intensity := ""
			if len(args) > 0 {
				intensity = args[0]
			}
			return (*a).ProfileCoaching(cmd.Context(), intensity)
		},
	}
	cmd.AddCommand(coachingCmd)

	cmd.AddCommand(&cobra.Command{
		Use:   "notifications",
		Short: "Show notification preferences",
		RunE: func(cmd *cobra.Command, args []string) error {
			return (*a).ProfileNotifications(cmd.Context())
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "timers",
		Short: "Show timer preferences",
		RunE: func(cmd *cobra.Command, args []string) error {
			return (*a).ProfileTimers(cmd.Context())
		},
	})

	return cmd
}

func newMemoryCmd(a **app.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "memory",
		Short: "Manage AI memories",
		RunE: func(cmd *cobra.Command, args []string) error {
			return (*a).MemoryList(cmd.Context())
		},
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "List all memories",
		RunE: func(cmd *cobra.Command, args []string) error {
			return (*a).MemoryList(cmd.Context())
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "show <id>",
		Short: "Show a memory",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return (*a).MemoryShow(cmd.Context(), args[0])
		},
	})

	addCmd := &cobra.Command{
		Use:   "add <content>",
		Short: "Add a memory",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			tagsStr, _ := cmd.Flags().GetString("tags")
			var tags []string
			if tagsStr != "" {
				tags = splitCSV(tagsStr)
			}
			return (*a).MemoryAdd(cmd.Context(), args[0], tags)
		},
	}
	addCmd.Flags().String("tags", "", "Comma-separated tags")
	cmd.AddCommand(addCmd)

	updateCmd := &cobra.Command{
		Use:   "update <id> <content>",
		Short: "Update a memory",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			tagsStr, _ := cmd.Flags().GetString("tags")
			var tags []string
			if tagsStr != "" {
				tags = splitCSV(tagsStr)
			}
			return (*a).MemoryUpdate(cmd.Context(), args[0], args[1], tags)
		},
	}
	updateCmd.Flags().String("tags", "", "Comma-separated tags")
	cmd.AddCommand(updateCmd)

	cmd.AddCommand(&cobra.Command{
		Use:   "search <query>",
		Short: "Search memories",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return (*a).MemorySearch(cmd.Context(), args[0])
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a memory",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return (*a).MemoryDelete(cmd.Context(), args[0])
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "export",
		Short: "Export all memories",
		RunE: func(cmd *cobra.Command, args []string) error {
			return (*a).MemoryExport(cmd.Context())
		},
	})

	return cmd
}

func newHouseholdCmd(a **app.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "household",
		Short: "Manage household",
		RunE: func(cmd *cobra.Command, args []string) error {
			return (*a).HouseholdInfo(cmd.Context())
		},
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "info",
		Short: "Show household information",
		RunE: func(cmd *cobra.Command, args []string) error {
			return (*a).HouseholdInfo(cmd.Context())
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "members",
		Short: "List household members",
		RunE: func(cmd *cobra.Command, args []string) error {
			return (*a).HouseholdMembers(cmd.Context())
		},
	})

	inviteCmd := &cobra.Command{
		Use:   "invite <email>",
		Short: "Invite someone to the household",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			role, _ := cmd.Flags().GetString("role")
			return (*a).HouseholdInvite(cmd.Context(), args[0], role)
		},
	}
	inviteCmd.Flags().String("role", "MEMBER", "Role (MEMBER, ADMIN)")
	cmd.AddCommand(inviteCmd)

	leaderboardCmd := &cobra.Command{
		Use:   "leaderboard",
		Short: "Show household leaderboard",
		RunE: func(cmd *cobra.Command, args []string) error {
			period, _ := cmd.Flags().GetString("period")
			return (*a).HouseholdLeaderboard(cmd.Context(), period)
		},
	}
	leaderboardCmd.Flags().String("period", "WEEKLY", "Period (WEEKLY, MONTHLY)")
	cmd.AddCommand(leaderboardCmd)

	cmd.AddCommand(&cobra.Command{
		Use:   "challenges",
		Short: "Show active household challenges",
		RunE: func(cmd *cobra.Command, args []string) error {
			return (*a).HouseholdChallenges(cmd.Context())
		},
	})

	return cmd
}
