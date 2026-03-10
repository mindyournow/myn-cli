package api

import (
	"context"
	"fmt"
)

// HabitStreak holds streak info for a habit.
type HabitStreak struct {
	HabitID     string `json:"habitId"`
	CurrentStreak int  `json:"currentStreak"`
	LongestStreak int  `json:"longestStreak"`
	History     []string `json:"history,omitempty"`
}

// HabitChain represents a chain of habits.
type HabitChain struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Habits []struct {
		ID    string `json:"id"`
		Title string `json:"title"`
	} `json:"habits,omitempty"`
}

// SkipHabit skips a habit for the current period.
func (c *Client) SkipHabit(ctx context.Context, id, reason, date string) (*UnifiedTask, error) {
	body := map[string]string{}
	if reason != "" {
		body["reason"] = reason
	}
	if date != "" {
		body["date"] = date
	}
	resp, err := c.Post(ctx, "/api/v2/unified-tasks/"+id+"/skip", body)
	if err != nil {
		return nil, err
	}
	var task UnifiedTask
	if err := resp.DecodeJSON(&task); err != nil {
		return nil, fmt.Errorf("failed to parse habit skip: %w", err)
	}
	return &task, nil
}

// GetHabitStreak fetches streak info for a habit.
func (c *Client) GetHabitStreak(ctx context.Context, id string) (*HabitStreak, error) {
	resp, err := c.Get(ctx, "/api/v2/unified-tasks/"+id+"/streak", nil)
	if err != nil {
		return nil, err
	}
	var streak HabitStreak
	if err := resp.DecodeJSON(&streak); err != nil {
		return nil, fmt.Errorf("failed to parse streak: %w", err)
	}
	return &streak, nil
}

// ListHabitChains fetches all habit chains.
func (c *Client) ListHabitChains(ctx context.Context) ([]HabitChain, error) {
	resp, err := c.Get(ctx, "/api/habits/chains", nil)
	if err != nil {
		return nil, err
	}
	var chains []HabitChain
	if err := resp.DecodeJSON(&chains); err != nil {
		return nil, fmt.Errorf("failed to parse habit chains: %w", err)
	}
	return chains, nil
}

// CreateHabitChain creates a new habit chain.
func (c *Client) CreateHabitChain(ctx context.Context, name string) (*HabitChain, error) {
	resp, err := c.Post(ctx, "/api/habits/chains", map[string]string{"name": name})
	if err != nil {
		return nil, err
	}
	var chain HabitChain
	if err := resp.DecodeJSON(&chain); err != nil {
		return nil, fmt.Errorf("failed to parse habit chain: %w", err)
	}
	return &chain, nil
}

// AddHabitToChain adds a habit to a chain.
func (c *Client) AddHabitToChain(ctx context.Context, chainID, habitID string) error {
	_, err := c.Post(ctx, "/api/habits/chains/"+chainID+"/habits",
		map[string]string{"habitId": habitID})
	return err
}

// RemoveHabitFromChain removes a habit from a chain.
func (c *Client) RemoveHabitFromChain(ctx context.Context, chainID, habitID string) error {
	_, err := c.Delete(ctx, "/api/habits/chains/"+chainID+"/habits/"+habitID)
	return err
}

// GetHabitChainStatus fetches the completion status of a habit chain.
func (c *Client) GetHabitChainStatus(ctx context.Context, chainID string) (map[string]interface{}, error) {
	resp, err := c.Get(ctx, "/api/habits/chains/"+chainID+"/status", nil)
	if err != nil {
		return nil, err
	}
	var status map[string]interface{}
	if err := resp.DecodeJSON(&status); err != nil {
		return nil, fmt.Errorf("failed to parse chain status: %w", err)
	}
	return status, nil
}

// BatchCompleteHabitChain completes all habits in a chain.
func (c *Client) BatchCompleteHabitChain(ctx context.Context, chainID string) error {
	_, err := c.Post(ctx, "/api/habits/chains/"+chainID+"/batch-complete", nil)
	return err
}

// ScheduleHabits triggers AI habit scheduling.
func (c *Client) ScheduleHabits(ctx context.Context, days int) (map[string]interface{}, error) {
	body := map[string]interface{}{}
	if days > 0 {
		body["numberOfDays"] = days
	}
	resp, err := c.Post(ctx, "/api/v2/scheduling/habits/schedule", body)
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	if err := resp.DecodeJSON(&result); err != nil {
		return nil, fmt.Errorf("failed to parse schedule result: %w", err)
	}
	return result, nil
}

// ListHabitReminders fetches habit reminders.
func (c *Client) ListHabitReminders(ctx context.Context) ([]map[string]interface{}, error) {
	resp, err := c.Get(ctx, "/api/habits/reminders", nil)
	if err != nil {
		return nil, err
	}
	var reminders []map[string]interface{}
	if err := resp.DecodeJSON(&reminders); err != nil {
		return nil, fmt.Errorf("failed to parse reminders: %w", err)
	}
	return reminders, nil
}

// CalculateSmartTime calculates the optimal reminder time for a habit.
func (c *Client) CalculateSmartTime(ctx context.Context, habitID string) (map[string]interface{}, error) {
	resp, err := c.Post(ctx, "/api/habits/reminders/"+habitID+"/calculate-smart-time", nil)
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	if err := resp.DecodeJSON(&result); err != nil {
		return nil, fmt.Errorf("failed to parse smart time: %w", err)
	}
	return result, nil
}
