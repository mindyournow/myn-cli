package app

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mindyournow/myn-cli/internal/api"
)

// PomodoroStart starts a Pomodoro session.
func (a *App) PomodoroStart(ctx context.Context, taskID, label string, duration int) error {
	if err := a.ensureAuth(ctx); err != nil {
		return err
	}
	req := api.StartPomodoroRequest{
		TaskID:       taskID,
		Label:        label,
		WorkDuration: duration,
	}
	session, err := a.Client.StartPomodoro(ctx, req)
	if err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to start pomodoro: %v", err))
		return err
	}
	if a.Formatter.JSON {
		return a.Formatter.Print(session)
	}
	return a.Formatter.Success(fmt.Sprintf("🍅 Pomodoro started: %s (status: %s)", session.ID, session.Status))
}

// PomodoroSmartStart starts a smart Pomodoro with AI suggestions.
func (a *App) PomodoroSmartStart(ctx context.Context, availableMinutes int) error {
	if err := a.ensureAuth(ctx); err != nil {
		return err
	}
	result, err := a.Client.StartSmartPomodoro(ctx, availableMinutes)
	if err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to start smart pomodoro: %v", err))
		return err
	}
	if a.Formatter.JSON {
		return a.Formatter.Print(result)
	}
	b, _ := json.MarshalIndent(result, "", "  ")
	return a.Formatter.Println(string(b))
}

// PomodoroStatus shows the current Pomodoro session.
func (a *App) PomodoroStatus(ctx context.Context) error {
	if err := a.ensureAuth(ctx); err != nil {
		return err
	}
	session, err := a.Client.GetCurrentPomodoro(ctx)
	if err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to get pomodoro: %v", err))
		return err
	}
	if a.Formatter.JSON {
		return a.Formatter.Print(session)
	}
	_ = a.Formatter.Println(fmt.Sprintf("ID:      %s", session.ID))
	_ = a.Formatter.Println(fmt.Sprintf("Status:  %s", session.Status))
	if session.Phase != "" {
		_ = a.Formatter.Println(fmt.Sprintf("Phase:   %s", session.Phase))
	}
	if session.TaskID != "" {
		_ = a.Formatter.Println(fmt.Sprintf("Task:    %s", session.TaskID))
	}
	return nil
}

// PomodoroPause pauses the current session.
func (a *App) PomodoroPause(ctx context.Context) error {
	if err := a.ensureAuth(ctx); err != nil {
		return err
	}
	session, err := a.Client.PausePomodoro(ctx)
	if err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to pause pomodoro: %v", err))
		return err
	}
	if a.Formatter.JSON {
		return a.Formatter.Print(session)
	}
	return a.Formatter.Success(fmt.Sprintf("Pomodoro paused: %s", session.ID))
}

// PomodoroResume resumes a paused session.
func (a *App) PomodoroResume(ctx context.Context) error {
	if err := a.ensureAuth(ctx); err != nil {
		return err
	}
	session, err := a.Client.ResumePomodoro(ctx)
	if err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to resume pomodoro: %v", err))
		return err
	}
	if a.Formatter.JSON {
		return a.Formatter.Print(session)
	}
	return a.Formatter.Success(fmt.Sprintf("Pomodoro resumed: %s", session.ID))
}

// PomodoroStop cancels the current session.
func (a *App) PomodoroStop(ctx context.Context) error {
	if err := a.ensureAuth(ctx); err != nil {
		return err
	}
	session, err := a.Client.StopPomodoro(ctx)
	if err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to stop pomodoro: %v", err))
		return err
	}
	if a.Formatter.JSON {
		return a.Formatter.Print(session)
	}
	return a.Formatter.Success(fmt.Sprintf("Pomodoro stopped: %s", session.ID))
}

// PomodoroComplete marks the current session as complete.
func (a *App) PomodoroComplete(ctx context.Context, note string) error {
	if err := a.ensureAuth(ctx); err != nil {
		return err
	}
	session, err := a.Client.CompletePomodoro(ctx, note)
	if err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to complete pomodoro: %v", err))
		return err
	}
	if a.Formatter.JSON {
		return a.Formatter.Print(session)
	}
	return a.Formatter.Success(fmt.Sprintf("🍅 Pomodoro complete: %s", session.ID))
}

// PomodoroInterrupt records an interruption.
func (a *App) PomodoroInterrupt(ctx context.Context, sessionID, reason string) error {
	if err := a.ensureAuth(ctx); err != nil {
		return err
	}
	session, err := a.Client.InterruptPomodoro(ctx, sessionID, reason)
	if err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to record interruption: %v", err))
		return err
	}
	if a.Formatter.JSON {
		return a.Formatter.Print(session)
	}
	return a.Formatter.Success(fmt.Sprintf("Interruption recorded for session %s.", session.ID))
}

// PomodoroSuggestions shows task suggestions for a Pomodoro.
func (a *App) PomodoroSuggestions(ctx context.Context, availableMinutes, maxSuggestions int) error {
	if err := a.ensureAuth(ctx); err != nil {
		return err
	}
	suggestions, err := a.Client.GetPomodoroSuggestions(ctx, availableMinutes, maxSuggestions)
	if err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to get suggestions: %v", err))
		return err
	}
	if a.Formatter.JSON {
		return a.Formatter.Print(suggestions)
	}
	b, _ := json.MarshalIndent(suggestions, "", "  ")
	return a.Formatter.Println(string(b))
}

// PomodoroHistory shows session history.
func (a *App) PomodoroHistory(ctx context.Context, params map[string]string) error {
	if err := a.ensureAuth(ctx); err != nil {
		return err
	}
	sessions, err := a.Client.ListPomodoroSessions(ctx, params)
	if err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to get history: %v", err))
		return err
	}
	if a.Formatter.JSON {
		return a.Formatter.Print(sessions)
	}
	if len(sessions) == 0 {
		return a.Formatter.Println("No pomodoro sessions.")
	}
	tbl := a.Formatter.NewTable("ID", "STATUS", "PHASE", "STARTED")
	for _, s := range sessions {
		tbl.AddRow(s.ID, s.Status, s.Phase, s.StartedAt)
	}
	tbl.Render()
	return nil
}

// PomodoroSettings shows Pomodoro settings.
func (a *App) PomodoroSettings(ctx context.Context) error {
	if err := a.ensureAuth(ctx); err != nil {
		return err
	}
	settings, err := a.Client.GetPomodoroSettings(ctx)
	if err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to get settings: %v", err))
		return err
	}
	if a.Formatter.JSON {
		return a.Formatter.Print(settings)
	}
	_ = a.Formatter.Println(fmt.Sprintf("Work:        %d min", settings.WorkDuration))
	_ = a.Formatter.Println(fmt.Sprintf("Short break: %d min", settings.ShortBreakDuration))
	_ = a.Formatter.Println(fmt.Sprintf("Long break:  %d min", settings.LongBreakDuration))
	_ = a.Formatter.Println(fmt.Sprintf("Sessions before long: %d", settings.SessionsBeforeLong))
	return nil
}

// PomodoroSettingsUpdate updates Pomodoro settings.
func (a *App) PomodoroSettingsUpdate(ctx context.Context, settings map[string]interface{}) error {
	if err := a.ensureAuth(ctx); err != nil {
		return err
	}
	updated, err := a.Client.UpdatePomodoroSettings(ctx, settings)
	if err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to update settings: %v", err))
		return err
	}
	if a.Formatter.JSON {
		return a.Formatter.Print(updated)
	}
	return a.Formatter.Success("Pomodoro settings updated.")
}
