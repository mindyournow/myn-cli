package app

import (
	"context"
	"encoding/json"
	"fmt"
)

// StatsStreaks shows gamification streaks.
func (a *App) StatsStreaks(ctx context.Context) error {
	if err := a.ensureAuth(ctx); err != nil {
		return err
	}
	streaks, err := a.Client.GetStreaks(ctx)
	if err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to get streaks: %v", err))
		return err
	}
	if a.Formatter.JSON {
		return a.Formatter.Print(streaks)
	}
	if len(streaks) == 0 {
		return a.Formatter.Println("No active streaks.")
	}
	tbl := a.Formatter.NewTable("TITLE", "COUNT", "TYPE")
	for _, s := range streaks {
		tbl.AddRow(s.Title, fmt.Sprintf("%d", s.Count), s.Type)
	}
	tbl.Render()
	return nil
}

// StatsAchievements shows achievements.
func (a *App) StatsAchievements(ctx context.Context, available bool) error {
	if err := a.ensureAuth(ctx); err != nil {
		return err
	}
	achievements, err := a.Client.GetAchievements(ctx, available)
	if err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to get achievements: %v", err))
		return err
	}
	if a.Formatter.JSON {
		return a.Formatter.Print(achievements)
	}
	if len(achievements) == 0 {
		return a.Formatter.Println("No achievements.")
	}
	tbl := a.Formatter.NewTable("TITLE", "UNLOCKED", "PROGRESS")
	for _, ach := range achievements {
		unlocked := "✗"
		if ach.IsUnlocked {
			unlocked = "✓"
		}
		progress := ""
		if ach.Total > 0 {
			progress = fmt.Sprintf("%d/%d", ach.Progress, ach.Total)
		}
		tbl.AddRow(ach.Title, unlocked, progress)
	}
	tbl.Render()
	return nil
}

// StatsPoints shows total points.
func (a *App) StatsPoints(ctx context.Context) error {
	if err := a.ensureAuth(ctx); err != nil {
		return err
	}
	points, err := a.Client.GetPoints(ctx)
	if err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to get points: %v", err))
		return err
	}
	if a.Formatter.JSON {
		return a.Formatter.Print(map[string]int{"points": points})
	}
	return a.Formatter.Println(fmt.Sprintf("Points: %d", points))
}

// StatsProductivity shows productivity stats.
func (a *App) StatsProductivity(ctx context.Context) error {
	if err := a.ensureAuth(ctx); err != nil {
		return err
	}
	stats, err := a.Client.GetProductivityStats(ctx)
	if err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to get stats: %v", err))
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

// StatsPomodoroStats shows Pomodoro statistics.
func (a *App) StatsPomodoroStats(ctx context.Context) error {
	if err := a.ensureAuth(ctx); err != nil {
		return err
	}
	stats, err := a.Client.GetPomodoroStats(ctx)
	if err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to get pomodoro stats: %v", err))
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
