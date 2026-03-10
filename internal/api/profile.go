package api

import (
	"context"
	"fmt"
)

// CustomerGoals represents a customer's goals.
type CustomerGoals struct {
	Goals []string `json:"goals,omitempty"`
}

// CustomerPreferences represents a customer's preferences.
type CustomerPreferences struct {
	Timezone          string `json:"timezone,omitempty"`
	Language          string `json:"language,omitempty"`
	CoachingIntensity string `json:"coachingIntensity,omitempty"`
}

// GetCustomerProfile fetches the current user's profile.
func (c *Client) GetCustomerProfile(ctx context.Context) (map[string]interface{}, error) {
	resp, err := c.Get(ctx, "/api/v1/customers", nil)
	if err != nil {
		return nil, err
	}
	var profile map[string]interface{}
	if err := resp.DecodeJSON(&profile); err != nil {
		return nil, fmt.Errorf("failed to parse profile: %w", err)
	}
	return profile, nil
}

// GetCustomerGoals fetches the current user's goals.
func (c *Client) GetCustomerGoals(ctx context.Context) (*CustomerGoals, error) {
	resp, err := c.Get(ctx, "/api/v1/customers/goals", nil)
	if err != nil {
		return nil, err
	}
	var goals CustomerGoals
	if err := resp.DecodeJSON(&goals); err != nil {
		return nil, fmt.Errorf("failed to parse goals: %w", err)
	}
	return &goals, nil
}

// SetCustomerGoals updates the customer's goals.
func (c *Client) SetCustomerGoals(ctx context.Context, goals CustomerGoals) (*CustomerGoals, error) {
	resp, err := c.Put(ctx, "/api/v1/customers/goals", goals)
	if err != nil {
		return nil, err
	}
	var updated CustomerGoals
	if err := resp.DecodeJSON(&updated); err != nil {
		return nil, fmt.Errorf("failed to parse updated goals: %w", err)
	}
	return &updated, nil
}

// GetCustomerPreferences fetches the customer's preferences.
func (c *Client) GetCustomerPreferences(ctx context.Context) (*CustomerPreferences, error) {
	resp, err := c.Get(ctx, "/api/v1/customers/preferences", nil)
	if err != nil {
		return nil, err
	}
	var prefs CustomerPreferences
	if err := resp.DecodeJSON(&prefs); err != nil {
		return nil, fmt.Errorf("failed to parse preferences: %w", err)
	}
	return &prefs, nil
}

// SetCustomerPreferences updates the customer's preferences.
func (c *Client) SetCustomerPreferences(ctx context.Context, prefs CustomerPreferences) (*CustomerPreferences, error) {
	resp, err := c.Put(ctx, "/api/v1/customers/preferences", prefs)
	if err != nil {
		return nil, err
	}
	var updated CustomerPreferences
	if err := resp.DecodeJSON(&updated); err != nil {
		return nil, fmt.Errorf("failed to parse updated preferences: %w", err)
	}
	return &updated, nil
}

// GetCoachingIntensity fetches the current coaching intensity.
func (c *Client) GetCoachingIntensity(ctx context.Context) (string, error) {
	resp, err := c.Get(ctx, "/api/v1/customers/coaching-intensity", nil)
	if err != nil {
		return "", err
	}
	var result struct {
		Intensity string `json:"intensity"`
	}
	if err := resp.DecodeJSON(&result); err != nil {
		return "", fmt.Errorf("failed to parse coaching intensity: %w", err)
	}
	return result.Intensity, nil
}

// SetCoachingIntensity updates the coaching intensity.
func (c *Client) SetCoachingIntensity(ctx context.Context, intensity string) error {
	_, err := c.Put(ctx, "/api/v1/customers/coaching-intensity",
		map[string]string{"intensity": intensity})
	return err
}

// GetNotificationPreferences fetches notification preferences.
func (c *Client) GetNotificationPreferences(ctx context.Context) (map[string]interface{}, error) {
	resp, err := c.Get(ctx, "/api/v1/customers/notification-preferences", nil)
	if err != nil {
		return nil, err
	}
	var prefs map[string]interface{}
	if err := resp.DecodeJSON(&prefs); err != nil {
		return nil, fmt.Errorf("failed to parse notification preferences: %w", err)
	}
	return prefs, nil
}

// GetTimerPreferences fetches timer preferences.
func (c *Client) GetTimerPreferences(ctx context.Context) (map[string]interface{}, error) {
	resp, err := c.Get(ctx, "/api/v2/customers/me/timer-preferences", nil)
	if err != nil {
		return nil, err
	}
	var prefs map[string]interface{}
	if err := resp.DecodeJSON(&prefs); err != nil {
		return nil, fmt.Errorf("failed to parse timer preferences: %w", err)
	}
	return prefs, nil
}
