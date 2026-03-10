package api

import (
	"context"
	"fmt"
)

// Streak represents a gamification streak.
type Streak struct {
	HabitID string `json:"habitId,omitempty"`
	Title   string `json:"title,omitempty"`
	Count   int    `json:"count"`
	Type    string `json:"type,omitempty"`
}

// Achievement represents a gamification achievement.
type Achievement struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	IsUnlocked  bool   `json:"isUnlocked"`
	Progress    int    `json:"progress,omitempty"`
	Total       int    `json:"total,omitempty"`
}

// GetStreaks fetches active streaks.
func (c *Client) GetStreaks(ctx context.Context) ([]Streak, error) {
	resp, err := c.Get(ctx, "/api/v1/gamification/streaks", nil)
	if err != nil {
		return nil, err
	}
	var streaks []Streak
	if err := resp.DecodeJSON(&streaks); err != nil {
		return nil, fmt.Errorf("failed to parse streaks: %w", err)
	}
	return streaks, nil
}

// GetAchievements fetches achievements.
func (c *Client) GetAchievements(ctx context.Context, available bool) ([]Achievement, error) {
	path := "/api/v1/gamification/achievements"
	if available {
		path += "/available"
	}
	resp, err := c.Get(ctx, path, nil)
	if err != nil {
		return nil, err
	}
	var achievements []Achievement
	if err := resp.DecodeJSON(&achievements); err != nil {
		return nil, fmt.Errorf("failed to parse achievements: %w", err)
	}
	return achievements, nil
}

// GetPoints fetches the user's total points.
func (c *Client) GetPoints(ctx context.Context) (int, error) {
	resp, err := c.Get(ctx, "/api/v1/gamification/points", nil)
	if err != nil {
		return 0, err
	}
	var result struct {
		Points int `json:"points"`
	}
	if err := resp.DecodeJSON(&result); err != nil {
		return 0, fmt.Errorf("failed to parse points: %w", err)
	}
	return result.Points, nil
}

// GetProductivityStats fetches overall productivity statistics.
func (c *Client) GetProductivityStats(ctx context.Context) (map[string]interface{}, error) {
	resp, err := c.Get(ctx, "/api/v1/usage/today", nil)
	if err != nil {
		return nil, err
	}
	var stats map[string]interface{}
	if err := resp.DecodeJSON(&stats); err != nil {
		return nil, fmt.Errorf("failed to parse stats: %w", err)
	}
	return stats, nil
}

// GetPomodoroStats fetches Pomodoro statistics.
func (c *Client) GetPomodoroStats(ctx context.Context) (map[string]interface{}, error) {
	resp, err := c.Get(ctx, "/api/v1/pomodoro/stats", nil)
	if err != nil {
		return nil, err
	}
	var stats map[string]interface{}
	if err := resp.DecodeJSON(&stats); err != nil {
		return nil, fmt.Errorf("failed to parse pomodoro stats: %w", err)
	}
	return stats, nil
}
