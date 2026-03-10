package app

import (
	"context"
	"fmt"

	"github.com/mindyournow/myn-cli/internal/util"
)

// TimerList lists all timers.
func (a *App) TimerList(ctx context.Context, includeCompleted bool) error {
	if err := a.ensureAuth(ctx); err != nil {
		return err
	}
	timers, err := a.Client.ListTimers(ctx, includeCompleted)
	if err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to list timers: %v", err))
		return err
	}
	if a.Formatter.JSON {
		return a.Formatter.Print(timers)
	}
	if len(timers) == 0 {
		return a.Formatter.Println("No active timers.")
	}
	tbl := a.Formatter.NewTable("ID", "TYPE", "STATUS", "DURATION/ALARM")
	for _, t := range timers {
		val := ""
		if t.Duration > 0 {
			val = util.FormatDuration(t.Duration)
		} else if t.AlarmTime != "" {
			val = t.AlarmTime
		}
		tbl.AddRow(t.ID, t.Type, t.Status, val)
	}
	tbl.Render()
	return nil
}

// TimerStart starts a countdown timer.
func (a *App) TimerStart(ctx context.Context, duration, label string) error {
	if err := a.ensureAuth(ctx); err != nil {
		return err
	}
	secs, err := util.ParseDuration(duration)
	if err != nil {
		return fmt.Errorf("invalid duration %q: %w", duration, err)
	}
	timer, err := a.Client.StartCountdownTimer(ctx, secs, label)
	if err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to start timer: %v", err))
		return err
	}
	if a.Formatter.JSON {
		return a.Formatter.Print(timer)
	}
	return a.Formatter.Success(fmt.Sprintf("Started %s countdown timer (%s)", util.FormatDuration(secs), timer.ID))
}

// TimerAlarm starts an alarm timer.
func (a *App) TimerAlarm(ctx context.Context, alarmTime, label string) error {
	if err := a.ensureAuth(ctx); err != nil {
		return err
	}
	timer, err := a.Client.StartAlarmTimer(ctx, alarmTime, label)
	if err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to start alarm: %v", err))
		return err
	}
	if a.Formatter.JSON {
		return a.Formatter.Print(timer)
	}
	return a.Formatter.Success(fmt.Sprintf("Set alarm for %s (%s)", alarmTime, timer.ID))
}

// TimerPause pauses a timer.
func (a *App) TimerPause(ctx context.Context, id string) error {
	if err := a.ensureAuth(ctx); err != nil {
		return err
	}
	timer, err := a.Client.PauseTimer(ctx, id)
	if err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to pause timer: %v", err))
		return err
	}
	if a.Formatter.JSON {
		return a.Formatter.Print(timer)
	}
	return a.Formatter.Success(fmt.Sprintf("Paused timer %s.", id))
}

// TimerResume resumes a paused timer.
func (a *App) TimerResume(ctx context.Context, id string) error {
	if err := a.ensureAuth(ctx); err != nil {
		return err
	}
	timer, err := a.Client.ResumeTimer(ctx, id)
	if err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to resume timer: %v", err))
		return err
	}
	if a.Formatter.JSON {
		return a.Formatter.Print(timer)
	}
	return a.Formatter.Success(fmt.Sprintf("Resumed timer %s.", id))
}

// TimerComplete marks a timer as complete.
func (a *App) TimerComplete(ctx context.Context, id string) error {
	if err := a.ensureAuth(ctx); err != nil {
		return err
	}
	timer, err := a.Client.CompleteTimer(ctx, id)
	if err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to complete timer: %v", err))
		return err
	}
	if a.Formatter.JSON {
		return a.Formatter.Print(timer)
	}
	return a.Formatter.Success(fmt.Sprintf("Completed timer %s.", id))
}

// TimerDismiss dismisses all completed timers.
func (a *App) TimerDismiss(ctx context.Context) error {
	if err := a.ensureAuth(ctx); err != nil {
		return err
	}
	if err := a.Client.DismissCompletedTimers(ctx); err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to dismiss timers: %v", err))
		return err
	}
	return a.Formatter.Success("All completed timers dismissed.")
}

// TimerCancel cancels a timer.
func (a *App) TimerCancel(ctx context.Context, id string) error {
	if err := a.ensureAuth(ctx); err != nil {
		return err
	}
	timer, err := a.Client.CancelTimer(ctx, id)
	if err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to cancel timer: %v", err))
		return err
	}
	if a.Formatter.JSON {
		return a.Formatter.Print(timer)
	}
	return a.Formatter.Success(fmt.Sprintf("Timer cancelled: %s", timer.ID))
}

// TimerSnooze snoozes a completed timer.
func (a *App) TimerSnooze(ctx context.Context, id string, minutes int) error {
	if err := a.ensureAuth(ctx); err != nil {
		return err
	}
	timer, err := a.Client.SnoozeTimer(ctx, id, minutes)
	if err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to snooze timer: %v", err))
		return err
	}
	if a.Formatter.JSON {
		return a.Formatter.Print(timer)
	}
	return a.Formatter.Success(fmt.Sprintf("Timer snoozed %d min: %s", minutes, timer.ID))
}

// TimerCount shows the number of active timers.
func (a *App) TimerCount(ctx context.Context) error {
	if err := a.ensureAuth(ctx); err != nil {
		return err
	}
	count, err := a.Client.GetTimerCount(ctx)
	if err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to get timer count: %v", err))
		return err
	}
	if a.Formatter.JSON {
		return a.Formatter.Print(map[string]int{"count": count})
	}
	return a.Formatter.Println(fmt.Sprintf("%d active timer(s).", count))
}
