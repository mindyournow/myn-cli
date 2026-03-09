package main

import (
	"github.com/mindyournow/myn-cli/internal/app"
	"github.com/spf13/cobra"
)

func newTaskCmd(a *app.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "task",
		Short: "Manage tasks",
	}

	cmd.AddCommand(
		newTaskListCmd(a),
		newTaskAddCmd(a),
		newTaskShowCmd(a),
		newTaskEditCmd(a),
		newTaskDoneCmd(a),
		newTaskDeleteCmd(a),
		newTaskRestoreCmd(a),
		newTaskArchiveCmd(a),
		newTaskSnoozeCmd(a),
	)

	return cmd
}

func newTaskListCmd(a *app.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List tasks",
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return a.TaskList()
		},
	}

	cmd.Flags().String("priority", "", "Filter by priority zone (critical, opportunity, horizon, parking)")
	cmd.Flags().String("type", "", "Filter by type (task, habit, chore)")
	cmd.Flags().String("project", "", "Filter by project")
	cmd.Flags().Bool("completed", false, "Include completed tasks")
	cmd.Flags().Bool("archived", false, "Show archived tasks")
	cmd.Flags().Bool("today", false, "Only tasks for today")
	cmd.Flags().Bool("overdue", false, "Only overdue tasks")
	cmd.Flags().Bool("household", false, "Include household tasks")
	cmd.Flags().String("sort", "priority", "Sort by: priority, date, title, created")
	cmd.Flags().Bool("reverse", false, "Reverse sort order")
	cmd.Flags().Int("page", 0, "Page number")
	cmd.Flags().Int("limit", 50, "Page size")

	return cmd
}

func newTaskAddCmd(a *app.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add <title>",
		Short: "Add a new task",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			title := args[0]
			return a.TaskAdd(title)
		},
	}

	cmd.Flags().String("priority", "", "Priority zone (critical, opportunity, horizon, parking)")
	cmd.Flags().String("project", "", "Project ID or name")
	cmd.Flags().String("due", "", "Due date (natural language)")
	cmd.Flags().StringArray("tags", nil, "Tags to add")
	cmd.Flags().String("note", "", "Additional notes")

	return cmd
}

func newTaskShowCmd(a *app.App) *cobra.Command {
	return &cobra.Command{
		Use:   "show <id>",
		Short: "Show task details",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return a.TaskShow(args[0])
		},
	}
}

func newTaskEditCmd(a *app.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "edit <id>",
		Short: "Edit a task",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return a.TaskEdit(args[0])
		},
	}

	cmd.Flags().String("title", "", "New title")
	cmd.Flags().String("priority", "", "New priority")
	cmd.Flags().String("due", "", "New due date")
	cmd.Flags().String("note", "", "New note")

	return cmd
}

func newTaskDoneCmd(a *app.App) *cobra.Command {
	return &cobra.Command{
		Use:   "done <id>",
		Short: "Mark task as completed",
		Aliases: []string{"complete", "finish"},
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return a.TaskDone(args[0])
		},
	}
}

func newTaskDeleteCmd(a *app.App) *cobra.Command {
	return &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a task",
		Aliases: []string{"rm", "remove"},
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return a.TaskDelete(args[0])
		},
	}
}

func newTaskRestoreCmd(a *app.App) *cobra.Command {
	return &cobra.Command{
		Use:   "restore <id>",
		Short: "Restore a deleted task",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return a.TaskRestore(args[0])
		},
	}
}

func newTaskArchiveCmd(a *app.App) *cobra.Command {
	return &cobra.Command{
		Use:   "archive <id>",
		Short: "Archive a task",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return a.TaskArchive(args[0])
		},
	}
}

func newTaskSnoozeCmd(a *app.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "snooze <id> [duration]",
		Short: "Snooze a task",
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			duration := "1d"
			if len(args) > 1 {
				duration = args[1]
			}
			return a.TaskSnooze(args[0], duration)
		},
	}

	return cmd
}
