package app

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/mindyournow/myn-cli/internal/api"
	"github.com/mindyournow/myn-cli/internal/util"
)

// PlanGoal uses the AI chat endpoint to generate a plan for a goal.
func (a *App) PlanGoal(ctx context.Context, goal string, hours int, deadline, priority string) error {
	if err := a.ensureAuth(ctx); err != nil {
		return err
	}

	// Build a structured prompt with task + calendar context
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("I need help planning: %q\n", goal))
	if hours > 0 {
		sb.WriteString(fmt.Sprintf("Available time: %d hours\n", hours))
	}
	if deadline != "" {
		sb.WriteString(fmt.Sprintf("Deadline: %s\n", deadline))
	}
	if priority != "" {
		sb.WriteString(fmt.Sprintf("Priority focus: %s\n", priority))
	}
	sb.WriteString("\nPlease provide an ordered action plan with time estimates.")

	req := api.AIChatRequest{
		CurrentMessage: sb.String(),
	}

	_ = a.Formatter.Println("📋 Generating plan...")
	err := a.Client.AIChatStream(ctx, req, func(event api.SSEEvent) error {
		if event.Data == "[DONE]" {
			return nil
		}
		_, werr := fmt.Fprint(a.Formatter.Writer(), event.Data)
		return werr
	})
	if err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("plan generation failed: %v", err))
		return err
	}
	return a.Formatter.Println("")
}

// AutoSchedule auto-schedules tasks into available calendar slots for a date.
func (a *App) AutoSchedule(ctx context.Context, date string, bufferMinutes int) error {
	if err := a.ensureAuth(ctx); err != nil {
		return err
	}
	if date == "" {
		date = nowDate()
	}

	// Fetch critical tasks
	tasks, err := a.Client.ListTasks(ctx, api.TaskListParams{Priority: "CRITICAL", Today: true})
	if err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to fetch tasks: %v", err))
		return err
	}

	// Fetch calendar events
	events, err := a.Client.ListCalendarEvents(ctx, date, 1)
	if err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to fetch calendar: %v", err))
		return err
	}

	if a.Formatter.JSON {
		return a.Formatter.Print(map[string]interface{}{
			"date":   date,
			"tasks":  tasks,
			"events": events,
		})
	}

	t, _ := time.Parse("2006-01-02", date)
	_ = a.Formatter.Println(fmt.Sprintf("AUTO-SCHEDULED — %s", t.Format("Monday, January 2")))
	_ = a.Formatter.Println("")

	// Simple greedy scheduling: assign tasks starting at 09:00
	currentMinute := 9 * 60
	endOfDay := 18 * 60

	scheduled := 0
	unscheduled := 0

	for _, task := range tasks {
		dur := task.Duration
		if dur <= 0 {
			dur = 30 // default 30 minutes
		}

		if currentMinute+dur > endOfDay {
			unscheduled++
			continue
		}

		startH := currentMinute / 60
		startM := currentMinute % 60
		endMin := currentMinute + dur
		endH := endMin / 60
		endMMin := endMin % 60

		_ = a.Formatter.Println(fmt.Sprintf("  %02d:%02d - %02d:%02d  %-30s %s",
			startH, startM, endH, endMMin, task.Title, task.PriorityString()))

		currentMinute += dur + bufferMinutes
		scheduled++
	}

	if unscheduled > 0 {
		_ = a.Formatter.Println("")
		_ = a.Formatter.Println(fmt.Sprintf("UNSCHEDULED: %d task(s) (not enough time today)", unscheduled))
	}

	return a.Formatter.Println(fmt.Sprintf("\nScheduled %d tasks.", scheduled))
}

// Reschedule moves one or more tasks to a new date.
func (a *App) Reschedule(ctx context.Context, ids []string, date string, spread int, reason string) error {
	if err := a.ensureAuth(ctx); err != nil {
		return err
	}
	if date == "" {
		return fmt.Errorf("--date is required")
	}
	targetDate, err := util.ParseDate(date)
	if err != nil {
		return fmt.Errorf("invalid date %q: %w", date, err)
	}

	rescheduled := 0
	for i, id := range ids {
		d := targetDate
		if spread > 0 && i > 0 {
			d = targetDate.AddDate(0, 0, (i*spread)/len(ids))
		}
		dateStr := d.Format("2006-01-02")
		_, err := a.Client.UpdateTask(ctx, id, api.UpdateTaskRequest{StartDate: dateStr})
		if err != nil {
			_ = a.Formatter.Error(fmt.Sprintf("failed to reschedule %s: %v", id, err))
			continue
		}
		rescheduled++
	}

	if a.Formatter.JSON {
		return a.Formatter.Print(map[string]interface{}{
			"rescheduled": rescheduled,
			"date":        targetDate.Format("2006-01-02"),
			"reason":      reason,
		})
	}
	return a.Formatter.Success(fmt.Sprintf("Rescheduled %d task(s) to %s.", rescheduled, targetDate.Format("Jan 2")))
}
