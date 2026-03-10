package app

import (
	"context"
	"encoding/json"
	"fmt"
)

// HabitList lists habits.
func (a *App) HabitList(ctx context.Context, due bool) error {
	if err := a.ensureAuth(ctx); err != nil {
		return err
	}
	tasks, err := a.Client.ListTasks(ctx, buildParams("HABIT", due))
	if err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to list habits: %v", err))
		return err
	}
	return a.printTaskList(ctx, tasks, "No habits found.")
}

// HabitDone completes a habit for today.
func (a *App) HabitDone(ctx context.Context, id string) error {
	if err := a.ensureAuth(ctx); err != nil {
		return err
	}
	task, err := a.Client.CompleteTask(ctx, id)
	if err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to complete habit: %v", err))
		return err
	}
	if a.Formatter.JSON {
		return a.Formatter.Print(task)
	}
	return a.Formatter.Success(fmt.Sprintf("✓ Habit done: %s", task.Title))
}

// HabitSkip skips a habit for today.
func (a *App) HabitSkip(ctx context.Context, id, reason, date string) error {
	if err := a.ensureAuth(ctx); err != nil {
		return err
	}
	task, err := a.Client.SkipHabit(ctx, id, reason, date)
	if err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to skip habit: %v", err))
		return err
	}
	if a.Formatter.JSON {
		return a.Formatter.Print(task)
	}
	return a.Formatter.Success(fmt.Sprintf("Skipped: %s", task.Title))
}

// HabitStreak shows streak info for a habit.
func (a *App) HabitStreak(ctx context.Context, id string, history bool) error {
	if err := a.ensureAuth(ctx); err != nil {
		return err
	}
	streak, err := a.Client.GetHabitStreak(ctx, id)
	if err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to get streak: %v", err))
		return err
	}
	if a.Formatter.JSON {
		return a.Formatter.Print(streak)
	}
	_ = a.Formatter.Println(fmt.Sprintf("Current streak: 🔥 %d", streak.CurrentStreak))
	_ = a.Formatter.Println(fmt.Sprintf("Longest streak: %d", streak.LongestStreak))
	if history && len(streak.History) > 0 {
		_ = a.Formatter.Println("\nHistory:")
		for _, h := range streak.History {
			_ = a.Formatter.Println(fmt.Sprintf("  %s", h))
		}
	}
	return nil
}

// HabitChainList lists habit chains.
func (a *App) HabitChainList(ctx context.Context) error {
	if err := a.ensureAuth(ctx); err != nil {
		return err
	}
	chains, err := a.Client.ListHabitChains(ctx)
	if err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to list chains: %v", err))
		return err
	}
	if a.Formatter.JSON {
		return a.Formatter.Print(chains)
	}
	if len(chains) == 0 {
		return a.Formatter.Println("No habit chains.")
	}
	tbl := a.Formatter.NewTable("ID", "NAME", "HABITS")
	for _, ch := range chains {
		tbl.AddRow(ch.ID, ch.Name, fmt.Sprintf("%d", len(ch.Habits)))
	}
	tbl.Render()
	return nil
}

// HabitChainCreate creates a new habit chain.
func (a *App) HabitChainCreate(ctx context.Context, name string) error {
	if err := a.ensureAuth(ctx); err != nil {
		return err
	}
	chain, err := a.Client.CreateHabitChain(ctx, name)
	if err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to create chain: %v", err))
		return err
	}
	if a.Formatter.JSON {
		return a.Formatter.Print(chain)
	}
	return a.Formatter.Success(fmt.Sprintf("Created chain: %s (%s)", chain.Name, chain.ID))
}

// HabitChainDone completes all habits in a chain.
func (a *App) HabitChainDone(ctx context.Context, chainID string) error {
	if err := a.ensureAuth(ctx); err != nil {
		return err
	}
	if err := a.Client.BatchCompleteHabitChain(ctx, chainID); err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to complete chain: %v", err))
		return err
	}
	return a.Formatter.Success(fmt.Sprintf("All habits in chain %s completed.", chainID))
}

// HabitSchedule triggers AI habit scheduling.
func (a *App) HabitSchedule(ctx context.Context, days int) error {
	if err := a.ensureAuth(ctx); err != nil {
		return err
	}
	result, err := a.Client.ScheduleHabits(ctx, days)
	if err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to schedule habits: %v", err))
		return err
	}
	if a.Formatter.JSON {
		return a.Formatter.Print(result)
	}
	b, _ := json.MarshalIndent(result, "", "  ")
	return a.Formatter.Println(string(b))
}

// HabitReminders lists habit reminders.
func (a *App) HabitReminders(ctx context.Context) error {
	if err := a.ensureAuth(ctx); err != nil {
		return err
	}
	reminders, err := a.Client.ListHabitReminders(ctx)
	if err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to list reminders: %v", err))
		return err
	}
	if a.Formatter.JSON {
		return a.Formatter.Print(reminders)
	}
	b, _ := json.MarshalIndent(reminders, "", "  ")
	return a.Formatter.Println(string(b))
}
