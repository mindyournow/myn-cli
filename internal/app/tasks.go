package app

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/mindyournow/myn-cli/internal/api"
	"github.com/mindyournow/myn-cli/internal/output"
	"github.com/mindyournow/myn-cli/internal/util"
)

// TaskListOptions are filters for task list.
type TaskListOptions struct {
	Priority         string
	Type             string
	Project          string
	Completed        bool
	Archived         bool
	Today            bool
	Overdue          bool
	Household        bool
	Sort             string
	Page             int
	Limit            int
}

// TaskListFull lists tasks with full options.
func (a *App) TaskListFull(ctx context.Context, opts TaskListOptions) error {
	if err := a.ensureAuth(ctx); err != nil {
		return err
	}

	params := api.TaskListParams{
		Type:             opts.Type,
		IsCompleted:      opts.Completed,
		IncludeHousehold: opts.Household,
		Archived:         opts.Archived,
		Sort:             opts.Sort,
		Page:             opts.Page,
		Limit:            opts.Limit,
	}

	tasks, err := a.Client.ListTasks(ctx, params)
	if err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to list tasks: %v", err))
		return err
	}

	// Client-side filters
	var filtered []api.UnifiedTask
	for _, t := range tasks {
		if opts.Priority != "" && t.PriorityString() != priorityToAPI(opts.Priority) {
			continue
		}
		if opts.Today && t.StartDate != time.Now().Format("2006-01-02") {
			continue
		}
		if opts.Overdue {
			if t.StartDate == "" || t.StartDate >= time.Now().Format("2006-01-02") {
				continue
			}
		}
		filtered = append(filtered, t)
	}

	if a.Formatter.JSON {
		type result struct {
			Tasks []api.UnifiedTask `json:"tasks"`
			Count int               `json:"count"`
		}
		return a.Formatter.Print(result{Tasks: filtered, Count: len(filtered)})
	}

	if len(filtered) == 0 {
		return a.Formatter.Println("No tasks found.")
	}

	tbl := a.Formatter.NewTable("", "TITLE", "PRIORITY", "DATE", "TYPE")
	for _, t := range filtered {
		priority := output.PriorityColored(t.PriorityString(), a.Formatter.NoColor)
		symbol := output.PriorityColored(t.PriorityString(), a.Formatter.NoColor)
		date := formatTaskDate(t.StartDate)
		tbl.AddRow(symbol, t.Title, priority, date, strings.ToLower(t.TaskType))
	}
	tbl.Render()
	return nil
}

// TaskAdd adds a new task.
func (a *App) TaskAdd(ctx context.Context, title string, opts TaskAddOptions) error {
	if err := a.ensureAuth(ctx); err != nil {
		return err
	}

	req := api.CreateTaskRequest{
		ID:    uuid.New().String(),
		Title: title,
	}
	if opts.Priority != "" {
		req.Priority = priorityToAPI(opts.Priority)
	}
	if opts.Date != "" {
		d, err := util.ParseDate(opts.Date)
		if err != nil {
			return err
		}
		req.StartDate = d.Format("2006-01-02")
	}
	if opts.Duration != "" {
		secs, err := util.ParseDuration(opts.Duration)
		if err != nil {
			return err
		}
		req.Duration = secs / 60 // API expects minutes
	}
	if opts.Description != "" {
		req.Description = opts.Description
	}
	if opts.Type != "" {
		req.TaskType = strings.ToUpper(opts.Type)
	}
	if opts.Recurrence != "" {
		req.RecurrenceRule = util.ParseRecurrence(opts.Recurrence)
	}
	if opts.ProjectID != "" {
		req.ProjectID = opts.ProjectID
	}

	task, err := a.Client.CreateTask(ctx, req)
	if err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to create task: %v", err))
		return err
	}

	if a.Formatter.JSON {
		return a.Formatter.Print(task)
	}
	symbol := output.PriorityColored(task.PriorityString(), a.Formatter.NoColor)
	return a.Formatter.Success(fmt.Sprintf("%s Created: %s (%s)", symbol, task.Title, task.ID))
}

// TaskAddOptions are flags for task add.
type TaskAddOptions struct {
	Priority    string
	Date        string
	Duration    string
	Description string
	Type        string
	Recurrence  string
	ProjectID   string
}

// TaskShow displays detailed info for a single task.
func (a *App) TaskShow(ctx context.Context, id string) error {
	if err := a.ensureAuth(ctx); err != nil {
		return err
	}
	task, err := a.Client.GetTask(ctx, id)
	if err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to get task: %v", err))
		return err
	}
	if a.Formatter.JSON {
		return a.Formatter.Print(task)
	}

	symbol := output.PriorityColored(task.PriorityString(), a.Formatter.NoColor)
	_ = a.Formatter.Println(fmt.Sprintf("%s %s", symbol, output.Bold(task.Title, a.Formatter.NoColor)))
	_ = a.Formatter.Println(fmt.Sprintf("  ID:       %s", task.ID))
	_ = a.Formatter.Println(fmt.Sprintf("  Type:     %s", strings.ToLower(task.TaskType)))
	_ = a.Formatter.Println(fmt.Sprintf("  Priority: %s", task.PriorityString()))
	if task.StartDate != "" {
		_ = a.Formatter.Println(fmt.Sprintf("  Date:     %s", task.StartDate))
	}
	if task.Duration > 0 {
		_ = a.Formatter.Println(fmt.Sprintf("  Duration: %s", util.FormatDuration(task.Duration*60)))
	}
	if task.Description != "" {
		_ = a.Formatter.Println("")
		return a.Formatter.PrintMarkdown(task.Description)
	}
	return nil
}

// TaskEdit edits a task.
func (a *App) TaskEdit(ctx context.Context, id string, opts TaskEditOptions) error {
	if err := a.ensureAuth(ctx); err != nil {
		return err
	}
	req := api.UpdateTaskRequest{}
	if opts.Title != "" {
		req.Title = opts.Title
	}
	if opts.Priority != "" {
		req.Priority = priorityToAPI(opts.Priority)
	}
	if opts.Date != "" {
		d, err := util.ParseDate(opts.Date)
		if err != nil {
			return err
		}
		req.StartDate = d.Format("2006-01-02")
	}
	if opts.Duration != "" {
		secs, err := util.ParseDuration(opts.Duration)
		if err != nil {
			return err
		}
		req.Duration = secs / 60
	}
	if opts.Description != "" {
		req.Description = opts.Description
	}
	if opts.ProjectID != "" {
		req.ProjectID = opts.ProjectID
	}

	task, err := a.Client.UpdateTask(ctx, id, req)
	if err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to update task: %v", err))
		return err
	}
	if a.Formatter.JSON {
		return a.Formatter.Print(task)
	}
	return a.Formatter.Success(fmt.Sprintf("Updated: %s", task.Title))
}

// TaskEditOptions are flags for task edit.
type TaskEditOptions struct {
	Title       string
	Priority    string
	Date        string
	Duration    string
	Description string
	ProjectID   string
}

// TaskComplete marks a task as done.
func (a *App) TaskComplete(ctx context.Context, id string) error {
	if err := a.ensureAuth(ctx); err != nil {
		return err
	}
	task, err := a.Client.CompleteTask(ctx, id)
	if err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to complete task: %v", err))
		return err
	}
	if a.Formatter.JSON {
		return a.Formatter.Print(task)
	}
	return a.Formatter.Success(fmt.Sprintf("%s Done: %s", output.SymbolDone, task.Title))
}

// TaskUncomplete marks a task as not done.
func (a *App) TaskUncomplete(ctx context.Context, id string) error {
	if err := a.ensureAuth(ctx); err != nil {
		return err
	}
	task, err := a.Client.UncompleteTask(ctx, id)
	if err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to uncomplete task: %v", err))
		return err
	}
	if a.Formatter.JSON {
		return a.Formatter.Print(task)
	}
	return a.Formatter.Success(fmt.Sprintf("Uncompleted: %s", task.Title))
}

// TaskDelete deletes a task.
func (a *App) TaskDelete(ctx context.Context, id string, permanent bool) error {
	if err := a.ensureAuth(ctx); err != nil {
		return err
	}
	if err := a.Client.DeleteTask(ctx, id, permanent); err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to delete task: %v", err))
		return err
	}
	if permanent {
		return a.Formatter.Success(fmt.Sprintf("Permanently deleted task %s", id))
	}
	return a.Formatter.Success(fmt.Sprintf("Deleted task %s (restorable with 'task restore')", id))
}

// TaskArchive archives a task.
func (a *App) TaskArchive(ctx context.Context, id string) error {
	if err := a.ensureAuth(ctx); err != nil {
		return err
	}
	task, err := a.Client.ArchiveTask(ctx, id)
	if err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to archive task: %v", err))
		return err
	}
	if a.Formatter.JSON {
		return a.Formatter.Print(task)
	}
	return a.Formatter.Success(fmt.Sprintf("Task archived: %s", task.Title))
}

// TaskRestore restores a soft-deleted task.
func (a *App) TaskRestore(ctx context.Context, id string) error {
	if err := a.ensureAuth(ctx); err != nil {
		return err
	}
	task, err := a.Client.RestoreTask(ctx, id)
	if err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to restore task: %v", err))
		return err
	}
	if a.Formatter.JSON {
		return a.Formatter.Print(task)
	}
	return a.Formatter.Success(fmt.Sprintf("Restored: %s", task.Title))
}

// TaskSnoozeOpt are options for snooze.
type TaskSnoozeOpt struct {
	Date string
	Days int
}

// TaskSnoozeTask snoozes a task to a future date.
func (a *App) TaskSnoozeTask(ctx context.Context, id string, opts TaskSnoozeOpt) error {
	if err := a.ensureAuth(ctx); err != nil {
		return err
	}

	var targetDate string
	if opts.Date != "" {
		d, err := util.ParseDate(opts.Date)
		if err != nil {
			return err
		}
		targetDate = d.Format("2006-01-02")
	} else if opts.Days > 0 {
		targetDate = time.Now().AddDate(0, 0, opts.Days).Format("2006-01-02")
	} else {
		targetDate = time.Now().AddDate(0, 0, 1).Format("2006-01-02") // default: tomorrow
	}

	task, err := a.Client.UpdateTask(ctx, id, api.UpdateTaskRequest{StartDate: targetDate})
	if err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to snooze task: %v", err))
		return err
	}
	if a.Formatter.JSON {
		return a.Formatter.Print(task)
	}
	return a.Formatter.Success(fmt.Sprintf("Snoozed: %s → %s", task.Title, targetDate))
}

// TaskBatchOptions are options for batch operations.
type TaskBatchOptions struct {
	IDs       []string
	Priority  string
	ProjectID string
	Date      string
}

// TaskBatch applies updates to multiple tasks at once.
func (a *App) TaskBatch(ctx context.Context, opts TaskBatchOptions) error {
	if err := a.ensureAuth(ctx); err != nil {
		return err
	}

	updates := api.UpdateTaskRequest{}
	if opts.Priority != "" {
		updates.Priority = priorityToAPI(opts.Priority)
	}
	if opts.ProjectID != "" {
		updates.ProjectID = opts.ProjectID
	}
	if opts.Date != "" {
		d, err := util.ParseDate(opts.Date)
		if err != nil {
			return err
		}
		updates.StartDate = d.Format("2006-01-02")
	}

	tasks, err := a.Client.BatchUpdateTasks(ctx, api.BatchUpdateRequest{
		IDs:     opts.IDs,
		Updates: updates,
	})
	if err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("batch update failed: %v", err))
		return err
	}
	if a.Formatter.JSON {
		return a.Formatter.Print(tasks)
	}
	return a.Formatter.Success(fmt.Sprintf("Updated %d tasks.", len(tasks)))
}

// TaskMove moves a task to a project.
func (a *App) TaskMove(ctx context.Context, taskID, projectID string) error {
	if err := a.ensureAuth(ctx); err != nil {
		return err
	}
	task, err := a.Client.MoveTask(ctx, taskID, projectID)
	if err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to move task: %v", err))
		return err
	}
	if a.Formatter.JSON {
		return a.Formatter.Print(task)
	}
	return a.Formatter.Success(fmt.Sprintf("Moved: %s → project %s", task.Title, projectID))
}

// priorityToAPI maps CLI priority names to API values.
func priorityToAPI(p string) string {
	switch strings.ToLower(p) {
	case "critical", "c":
		return "CRITICAL"
	case "opportunity", "o":
		return "OPPORTUNITY_NOW"
	case "horizon", "h":
		return "OVER_THE_HORIZON"
	case "parking", "x":
		return "PARKING_LOT"
	}
	return p // pass through if already an API value
}

// formatTaskDate formats a date string for display.
func formatTaskDate(date string) string {
	if date == "" {
		return "—"
	}
	t, err := time.Parse("2006-01-02", date)
	if err != nil {
		return date
	}
	today := time.Now().Truncate(24 * time.Hour)
	diff := t.Sub(today)
	switch {
	case diff == 0:
		return "today"
	case diff == 24*time.Hour:
		return "tomorrow"
	case diff == -24*time.Hour:
		return "yesterday"
	case diff < 0:
		return output.Red(t.Format("Jan 2"), false) // overdue
	default:
		return t.Format("Jan 2")
	}
}
