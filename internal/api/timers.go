package api

import (
	"context"
	"fmt"
)

// Timer represents a timer from the backend.
type Timer struct {
	ID        string `json:"id"`
	Type      string `json:"type"`
	Duration  int    `json:"duration,omitempty"`
	AlarmTime string `json:"alarmTime,omitempty"`
	Status    string `json:"status"`
	CreatedAt string `json:"createdAt,omitempty"`
	Label     string `json:"label,omitempty"`
}

// ListTimers fetches all timers.
func (c *Client) ListTimers(ctx context.Context, includeCompleted bool) ([]Timer, error) {
	params := map[string]string{}
	if includeCompleted {
		params["includeCompleted"] = "true"
	}
	resp, err := c.Get(ctx, "/api/v2/timers", params)
	if err != nil {
		return nil, err
	}
	var timers []Timer
	if err := resp.DecodeJSON(&timers); err != nil {
		return nil, fmt.Errorf("failed to parse timers: %w", err)
	}
	return timers, nil
}

// StartCountdownTimer creates a countdown timer.
func (c *Client) StartCountdownTimer(ctx context.Context, durationSeconds int, label string) (*Timer, error) {
	body := map[string]interface{}{
		"duration": durationSeconds,
	}
	if label != "" {
		body["label"] = label
	}
	resp, err := c.Post(ctx, "/api/v2/timers/countdown", body)
	if err != nil {
		return nil, err
	}
	var timer Timer
	if err := resp.DecodeJSON(&timer); err != nil {
		return nil, fmt.Errorf("failed to parse timer: %w", err)
	}
	return &timer, nil
}

// StartAlarmTimer creates an alarm timer for a specific time.
func (c *Client) StartAlarmTimer(ctx context.Context, alarmTime string, label string) (*Timer, error) {
	body := map[string]interface{}{
		"alarmTime": alarmTime,
	}
	if label != "" {
		body["label"] = label
	}
	resp, err := c.Post(ctx, "/api/v2/timers/alarm", body)
	if err != nil {
		return nil, err
	}
	var timer Timer
	if err := resp.DecodeJSON(&timer); err != nil {
		return nil, fmt.Errorf("failed to parse timer: %w", err)
	}
	return &timer, nil
}

// PauseTimer pauses a running timer.
func (c *Client) PauseTimer(ctx context.Context, id string) (*Timer, error) {
	resp, err := c.Post(ctx, "/api/v2/timers/"+id+"/pause", nil)
	if err != nil {
		return nil, err
	}
	var timer Timer
	if err := resp.DecodeJSON(&timer); err != nil {
		return nil, fmt.Errorf("failed to parse timer: %w", err)
	}
	return &timer, nil
}

// ResumeTimer resumes a paused timer.
func (c *Client) ResumeTimer(ctx context.Context, id string) (*Timer, error) {
	resp, err := c.Post(ctx, "/api/v2/timers/"+id+"/resume", nil)
	if err != nil {
		return nil, err
	}
	var timer Timer
	if err := resp.DecodeJSON(&timer); err != nil {
		return nil, fmt.Errorf("failed to parse timer: %w", err)
	}
	return &timer, nil
}

// CompleteTimer marks a timer as complete.
func (c *Client) CompleteTimer(ctx context.Context, id string) (*Timer, error) {
	resp, err := c.Post(ctx, "/api/v2/timers/"+id+"/complete", nil)
	if err != nil {
		return nil, err
	}
	var timer Timer
	if err := resp.DecodeJSON(&timer); err != nil {
		return nil, fmt.Errorf("failed to parse timer: %w", err)
	}
	return &timer, nil
}

// DismissCompletedTimers deletes all completed timers.
func (c *Client) DismissCompletedTimers(ctx context.Context) error {
	_, err := c.Delete(ctx, "/api/v2/timers/completed")
	return err
}

// CancelTimer cancels an active timer.
func (c *Client) CancelTimer(ctx context.Context, id string) (*Timer, error) {
	resp, err := c.Post(ctx, "/api/v2/timers/"+id+"/cancel", nil)
	if err != nil {
		return nil, err
	}
	var timer Timer
	if err := resp.DecodeJSON(&timer); err != nil {
		return nil, fmt.Errorf("failed to parse cancelled timer: %w", err)
	}
	return &timer, nil
}

// SnoozeTimer snoozes a completed timer for additional minutes.
func (c *Client) SnoozeTimer(ctx context.Context, id string, snoozeMinutes int) (*Timer, error) {
	resp, err := c.Post(ctx, "/api/v2/timers/"+id+"/snooze",
		map[string]int{"snoozeMinutes": snoozeMinutes})
	if err != nil {
		return nil, err
	}
	var timer Timer
	if err := resp.DecodeJSON(&timer); err != nil {
		return nil, fmt.Errorf("failed to parse snoozed timer: %w", err)
	}
	return &timer, nil
}

// GetTimerCount returns the number of active timers.
func (c *Client) GetTimerCount(ctx context.Context) (int, error) {
	resp, err := c.Get(ctx, "/api/v2/timers/count", nil)
	if err != nil {
		return 0, err
	}
	var result struct {
		Count int `json:"count"`
	}
	if err := resp.DecodeJSON(&result); err != nil {
		return 0, fmt.Errorf("failed to parse timer count: %w", err)
	}
	return result.Count, nil
}
