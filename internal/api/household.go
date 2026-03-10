package api

import (
	"context"
	"fmt"
)

// Household represents household info.
type Household struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// HouseholdMember represents a household member.
type HouseholdMember struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

// GetHouseholdInfo fetches household info from the customer profile.
func (c *Client) GetHouseholdInfo(ctx context.Context) ([]Household, error) {
	resp, err := c.Get(ctx, "/api/v1/customers", nil)
	if err != nil {
		return nil, err
	}
	var profile struct {
		Households []Household `json:"households"`
	}
	if err := resp.DecodeJSON(&profile); err != nil {
		return nil, fmt.Errorf("failed to parse household info: %w", err)
	}
	return profile.Households, nil
}

// ListHouseholdMembers fetches members of a household.
func (c *Client) ListHouseholdMembers(ctx context.Context, householdID string) ([]HouseholdMember, error) {
	resp, err := c.Get(ctx, "/api/v1/households/"+householdID+"/members", nil)
	if err != nil {
		return nil, err
	}
	var members []HouseholdMember
	if err := resp.DecodeJSON(&members); err != nil {
		return nil, fmt.Errorf("failed to parse household members: %w", err)
	}
	return members, nil
}

// InviteToHousehold sends an invitation to join a household.
func (c *Client) InviteToHousehold(ctx context.Context, householdID, email, role string) error {
	body := map[string]string{
		"email":       email,
		"householdId": householdID,
		"role":        role,
	}
	_, err := c.Post(ctx, "/api/v1/households/invites", body)
	return err
}

// GetHouseholdLeaderboard fetches the household leaderboard.
func (c *Client) GetHouseholdLeaderboard(ctx context.Context, householdID, period string) (map[string]interface{}, error) {
	params := map[string]string{}
	if period != "" {
		params["period"] = period
	}
	resp, err := c.Get(ctx, "/api/v1/gamification/households/"+householdID+"/leaderboard", params)
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	if err := resp.DecodeJSON(&result); err != nil {
		return nil, fmt.Errorf("failed to parse leaderboard: %w", err)
	}
	return result, nil
}

// GetHouseholdChallenges fetches active household challenges.
func (c *Client) GetHouseholdChallenges(ctx context.Context, householdID string) ([]map[string]interface{}, error) {
	resp, err := c.Get(ctx, "/api/v1/gamification/households/"+householdID+"/challenges", nil)
	if err != nil {
		return nil, err
	}
	var challenges []map[string]interface{}
	if err := resp.DecodeJSON(&challenges); err != nil {
		return nil, fmt.Errorf("failed to parse challenges: %w", err)
	}
	return challenges, nil
}
