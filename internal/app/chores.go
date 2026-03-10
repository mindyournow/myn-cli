package app

import (
	"context"
	"encoding/json"
	"fmt"
)

// ChoreList lists today's chores.
func (a *App) ChoreList(ctx context.Context, date string) error {
	if err := a.ensureAuth(ctx); err != nil {
		return err
	}
	if date == "" {
		date = nowDate()
	}
	chores, err := a.Client.ListTodayChores(ctx, date, "", "")
	if err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to list chores: %v", err))
		return err
	}
	if a.Formatter.JSON {
		return a.Formatter.Print(chores)
	}
	if len(chores) == 0 {
		return a.Formatter.Println("No chores scheduled.")
	}
	tbl := a.Formatter.NewTable("ID", "TITLE", "ASSIGNED TO", "DATE", "DONE")
	for _, ch := range chores {
		done := " "
		if ch.IsCompleted {
			done = "✓"
		}
		tbl.AddRow(ch.ID, ch.Title, ch.AssignedTo, ch.ScheduledDate, done)
	}
	tbl.Render()
	return nil
}

// ChoreDone marks a chore instance as complete.
func (a *App) ChoreDone(ctx context.Context, instanceID, note string) error {
	if err := a.ensureAuth(ctx); err != nil {
		return err
	}
	if err := a.Client.CompleteChoreInstance(ctx, instanceID, note); err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to complete chore: %v", err))
		return err
	}
	return a.Formatter.Success(fmt.Sprintf("✓ Chore completed: %s", instanceID))
}

// ChoreSchedule shows the chore schedule for a date.
func (a *App) ChoreSchedule(ctx context.Context, date string) error {
	if err := a.ensureAuth(ctx); err != nil {
		return err
	}
	if date == "" {
		date = nowDate()
	}
	schedule, err := a.Client.GetChoreSchedule(ctx, date)
	if err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to get chore schedule: %v", err))
		return err
	}
	if a.Formatter.JSON {
		return a.Formatter.Print(schedule)
	}
	if len(schedule) == 0 {
		return a.Formatter.Println("No chores scheduled.")
	}
	tbl := a.Formatter.NewTable("TITLE", "ASSIGNED TO", "DATE")
	for _, ch := range schedule {
		tbl.AddRow(ch.Title, ch.AssignedTo, ch.ScheduledDate)
	}
	tbl.Render()
	return nil
}

// ChoreStats shows chore statistics.
func (a *App) ChoreStats(ctx context.Context) error {
	if err := a.ensureAuth(ctx); err != nil {
		return err
	}
	stats, err := a.Client.GetChoreStatistics(ctx)
	if err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to get chore stats: %v", err))
		return err
	}
	if a.Formatter.JSON {
		return a.Formatter.Print(stats)
	}
	b, err := json.MarshalIndent(stats, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to format output: %w", err)
	}
	return a.Formatter.Println(string(b))
}

// ChoreRotation shows the rotation status for a chore.
func (a *App) ChoreRotation(ctx context.Context, choreID string) error {
	if err := a.ensureAuth(ctx); err != nil {
		return err
	}
	status, err := a.Client.GetChoreRotationStatus(ctx, choreID)
	if err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to get rotation: %v", err))
		return err
	}
	if a.Formatter.JSON {
		return a.Formatter.Print(status)
	}
	b, err := json.MarshalIndent(status, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to format output: %w", err)
	}
	return a.Formatter.Println(string(b))
}

// ChoreRotationAdvance moves rotation to the next member.
func (a *App) ChoreRotationAdvance(ctx context.Context, choreID string) error {
	if err := a.ensureAuth(ctx); err != nil {
		return err
	}
	if err := a.Client.AdvanceChoreRotation(ctx, choreID); err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to advance rotation: %v", err))
		return err
	}
	return a.Formatter.Success("Rotation advanced to next member.")
}

// ChoreRotationReset resets the rotation.
func (a *App) ChoreRotationReset(ctx context.Context, choreID string) error {
	if err := a.ensureAuth(ctx); err != nil {
		return err
	}
	if err := a.Client.ResetChoreRotation(ctx, choreID); err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to reset rotation: %v", err))
		return err
	}
	return a.Formatter.Success("Rotation reset to first member.")
}
