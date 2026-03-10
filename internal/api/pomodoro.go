package api

import (
	"context"
	"fmt"
)

// PomodoroSession represents a Pomodoro session.
type PomodoroSession struct {
	ID        string `json:"id"`
	Status    string `json:"status"`
	Phase     string `json:"phase,omitempty"`
	TaskID    string `json:"taskId,omitempty"`
	StartedAt string `json:"startedAt,omitempty"`
	Duration  int    `json:"duration,omitempty"`
}

// PomodoroSettings holds Pomodoro timer settings.
type PomodoroSettings struct {
	WorkDuration       int  `json:"workDuration"`
	ShortBreakDuration int  `json:"shortBreakDuration"`
	LongBreakDuration  int  `json:"longBreakDuration"`
	SessionsBeforeLong int  `json:"sessionsBeforeLong"`
	AutoStartBreaks    bool `json:"autoStartBreaks"`
	AutoStartWork      bool `json:"autoStartWork"`
}

// StartPomodoroRequest is the body for starting a Pomodoro.
type StartPomodoroRequest struct {
	WorkDuration int    `json:"workDuration,omitempty"`
	TaskID       string `json:"taskId,omitempty"`
	Label        string `json:"label,omitempty"`
}

// StartPomodoro starts a Pomodoro session.
func (c *Client) StartPomodoro(ctx context.Context, req StartPomodoroRequest) (*PomodoroSession, error) {
	resp, err := c.Post(ctx, "/api/v1/pomodoro/start", req)
	if err != nil {
		return nil, err
	}
	var session PomodoroSession
	if err := resp.DecodeJSON(&session); err != nil {
		return nil, fmt.Errorf("failed to parse pomodoro session: %w", err)
	}
	return &session, nil
}

// StartSmartPomodoro starts a smart Pomodoro with AI suggestions.
func (c *Client) StartSmartPomodoro(ctx context.Context, availableMinutes int) (map[string]interface{}, error) {
	resp, err := c.Post(ctx, "/api/v1/pomodoro/smart-start",
		map[string]int{"availableMinutes": availableMinutes})
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	if err := resp.DecodeJSON(&result); err != nil {
		return nil, fmt.Errorf("failed to parse smart pomodoro: %w", err)
	}
	return result, nil
}

// GetCurrentPomodoro fetches the current active Pomodoro session.
func (c *Client) GetCurrentPomodoro(ctx context.Context) (*PomodoroSession, error) {
	resp, err := c.Get(ctx, "/api/v1/pomodoro/current", nil)
	if err != nil {
		return nil, err
	}
	var session PomodoroSession
	if err := resp.DecodeJSON(&session); err != nil {
		return nil, fmt.Errorf("failed to parse current pomodoro: %w", err)
	}
	return &session, nil
}

// PausePomodoro pauses the current session.
func (c *Client) PausePomodoro(ctx context.Context) (*PomodoroSession, error) {
	resp, err := c.Post(ctx, "/api/v1/pomodoro/pause", nil)
	if err != nil {
		return nil, err
	}
	var session PomodoroSession
	if err := resp.DecodeJSON(&session); err != nil {
		return nil, fmt.Errorf("failed to parse paused pomodoro: %w", err)
	}
	return &session, nil
}

// ResumePomodoro resumes a paused session.
func (c *Client) ResumePomodoro(ctx context.Context) (*PomodoroSession, error) {
	resp, err := c.Post(ctx, "/api/v1/pomodoro/resume", nil)
	if err != nil {
		return nil, err
	}
	var session PomodoroSession
	if err := resp.DecodeJSON(&session); err != nil {
		return nil, fmt.Errorf("failed to parse resumed pomodoro: %w", err)
	}
	return &session, nil
}

// StopPomodoro cancels the current session.
func (c *Client) StopPomodoro(ctx context.Context) (*PomodoroSession, error) {
	resp, err := c.Post(ctx, "/api/v1/pomodoro/stop", nil)
	if err != nil {
		return nil, err
	}
	var session PomodoroSession
	if err := resp.DecodeJSON(&session); err != nil {
		return nil, fmt.Errorf("failed to parse stopped pomodoro: %w", err)
	}
	return &session, nil
}

// CompletePomodoro marks the current session as complete.
func (c *Client) CompletePomodoro(ctx context.Context, note string) (*PomodoroSession, error) {
	body := map[string]string{}
	if note != "" {
		body["note"] = note
	}
	resp, err := c.Post(ctx, "/api/v1/pomodoro/complete", body)
	if err != nil {
		return nil, err
	}
	var session PomodoroSession
	if err := resp.DecodeJSON(&session); err != nil {
		return nil, fmt.Errorf("failed to parse completed pomodoro: %w", err)
	}
	return &session, nil
}

// InterruptPomodoro records an interruption.
func (c *Client) InterruptPomodoro(ctx context.Context, sessionID, reason string) (*PomodoroSession, error) {
	resp, err := c.Put(ctx, "/api/v1/pomodoro/sessions/"+sessionID,
		map[string]string{"interruptReason": reason})
	if err != nil {
		return nil, err
	}
	var session PomodoroSession
	if err := resp.DecodeJSON(&session); err != nil {
		return nil, fmt.Errorf("failed to parse interrupted pomodoro: %w", err)
	}
	return &session, nil
}

// GetPomodoroSuggestions fetches task suggestions for a Pomodoro session.
func (c *Client) GetPomodoroSuggestions(ctx context.Context, availableMinutes, maxSuggestions int) ([]map[string]interface{}, error) {
	params := map[string]string{}
	if availableMinutes > 0 {
		params["availableMinutes"] = fmt.Sprintf("%d", availableMinutes)
	}
	if maxSuggestions > 0 {
		params["maxSuggestions"] = fmt.Sprintf("%d", maxSuggestions)
	}
	resp, err := c.Get(ctx, "/api/v1/pomodoro/suggestions", params)
	if err != nil {
		return nil, err
	}
	var suggestions []map[string]interface{}
	if err := resp.DecodeJSON(&suggestions); err != nil {
		return nil, fmt.Errorf("failed to parse suggestions: %w", err)
	}
	return suggestions, nil
}

// ListPomodoroSessions fetches Pomodoro session history.
func (c *Client) ListPomodoroSessions(ctx context.Context, params map[string]string) ([]PomodoroSession, error) {
	resp, err := c.Get(ctx, "/api/v1/pomodoro/sessions", params)
	if err != nil {
		return nil, err
	}
	var sessions []PomodoroSession
	if err := resp.DecodeJSON(&sessions); err != nil {
		return nil, fmt.Errorf("failed to parse pomodoro sessions: %w", err)
	}
	return sessions, nil
}

// GetPomodoroSettings fetches Pomodoro settings.
func (c *Client) GetPomodoroSettings(ctx context.Context) (*PomodoroSettings, error) {
	resp, err := c.Get(ctx, "/api/v1/pomodoro/settings", nil)
	if err != nil {
		return nil, err
	}
	var settings PomodoroSettings
	if err := resp.DecodeJSON(&settings); err != nil {
		return nil, fmt.Errorf("failed to parse pomodoro settings: %w", err)
	}
	return &settings, nil
}

// UpdatePomodoroSettings updates Pomodoro settings.
func (c *Client) UpdatePomodoroSettings(ctx context.Context, settings map[string]interface{}) (*PomodoroSettings, error) {
	resp, err := c.Put(ctx, "/api/v1/pomodoro/settings", settings)
	if err != nil {
		return nil, err
	}
	var updated PomodoroSettings
	if err := resp.DecodeJSON(&updated); err != nil {
		return nil, fmt.Errorf("failed to parse updated pomodoro settings: %w", err)
	}
	return &updated, nil
}
