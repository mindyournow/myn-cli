package api

import (
	"context"
	"fmt"
)

// ChoreInstance represents a scheduled chore instance.
type ChoreInstance struct {
	ID           string `json:"id"`
	ChoreID      string `json:"choreId"`
	Title        string `json:"title"`
	AssignedTo   string `json:"assignedTo,omitempty"`
	ScheduledDate string `json:"scheduledDate,omitempty"`
	IsCompleted  bool   `json:"isCompleted"`
}

// ListTodayChores fetches today's chores.
func (c *Client) ListTodayChores(ctx context.Context, date, timezone, householdID string) ([]ChoreInstance, error) {
	params := map[string]string{}
	if date != "" {
		params["date"] = date
	}
	if timezone != "" {
		params["timezone"] = timezone
	}
	if householdID != "" {
		params["householdId"] = householdID
	}
	resp, err := c.Get(ctx, "/api/v2/chores/today", params)
	if err != nil {
		return nil, err
	}
	var chores []ChoreInstance
	if err := resp.DecodeJSON(&chores); err != nil {
		return nil, fmt.Errorf("failed to parse chores: %w", err)
	}
	return chores, nil
}

// GetChoreSchedule fetches the chore schedule for a date.
func (c *Client) GetChoreSchedule(ctx context.Context, date string) ([]ChoreInstance, error) {
	resp, err := c.Get(ctx, "/api/v2/chores/schedule/"+date, nil)
	if err != nil {
		return nil, err
	}
	var schedule []ChoreInstance
	if err := resp.DecodeJSON(&schedule); err != nil {
		return nil, fmt.Errorf("failed to parse chore schedule: %w", err)
	}
	return schedule, nil
}

// CompleteChoreInstance marks a chore instance as complete.
func (c *Client) CompleteChoreInstance(ctx context.Context, instanceID, note string) error {
	body := map[string]string{}
	if note != "" {
		body["note"] = note
	}
	_, err := c.Post(ctx, "/api/v2/chores/instances/"+instanceID+"/complete", body)
	return err
}

// GetChoreStatistics fetches chore completion statistics.
func (c *Client) GetChoreStatistics(ctx context.Context) (map[string]interface{}, error) {
	resp, err := c.Get(ctx, "/api/v2/chores/statistics", nil)
	if err != nil {
		return nil, err
	}
	var stats map[string]interface{}
	if err := resp.DecodeJSON(&stats); err != nil {
		return nil, fmt.Errorf("failed to parse chore stats: %w", err)
	}
	return stats, nil
}

// GetChoreRotationStatus fetches the rotation status for a chore.
func (c *Client) GetChoreRotationStatus(ctx context.Context, choreID string) (map[string]interface{}, error) {
	resp, err := c.Get(ctx, "/api/v2/unified-tasks/"+choreID+"/rotation/status", nil)
	if err != nil {
		return nil, err
	}
	var status map[string]interface{}
	if err := resp.DecodeJSON(&status); err != nil {
		return nil, fmt.Errorf("failed to parse rotation status: %w", err)
	}
	return status, nil
}

// AdvanceChoreRotation moves the rotation to the next member.
func (c *Client) AdvanceChoreRotation(ctx context.Context, choreID string) error {
	_, err := c.Post(ctx, "/api/v2/unified-tasks/"+choreID+"/rotation/advance", nil)
	return err
}

// ResetChoreRotation resets the rotation to the first member.
func (c *Client) ResetChoreRotation(ctx context.Context, choreID string) error {
	_, err := c.Post(ctx, "/api/v2/unified-tasks/"+choreID+"/rotation/reset", nil)
	return err
}
