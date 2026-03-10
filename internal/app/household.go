package app

import (
	"context"
	"encoding/json"
	"fmt"
)

// HouseholdInfo shows household information.
func (a *App) HouseholdInfo(ctx context.Context) error {
	if err := a.ensureAuth(ctx); err != nil {
		return err
	}
	households, err := a.Client.GetHouseholdInfo(ctx)
	if err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to get household info: %v", err))
		return err
	}
	if a.Formatter.JSON {
		return a.Formatter.Print(households)
	}
	if len(households) == 0 {
		return a.Formatter.Println("No households found.")
	}
	for _, h := range households {
		_ = a.Formatter.Println(fmt.Sprintf("ID:   %s", h.ID))
		_ = a.Formatter.Println(fmt.Sprintf("Name: %s", h.Name))
	}
	return nil
}

// HouseholdMembers lists household members.
func (a *App) HouseholdMembers(ctx context.Context) error {
	householdID, err := a.ensureHousehold(ctx)
	if err != nil {
		return err
	}
	members, err := a.Client.ListHouseholdMembers(ctx, householdID)
	if err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to list members: %v", err))
		return err
	}
	if a.Formatter.JSON {
		return a.Formatter.Print(members)
	}
	tbl := a.Formatter.NewTable("NAME", "EMAIL", "ROLE")
	for _, m := range members {
		tbl.AddRow(m.Name, m.Email, m.Role)
	}
	tbl.Render()
	return nil
}

// HouseholdInvite invites someone to the household.
func (a *App) HouseholdInvite(ctx context.Context, email, role string) error {
	householdID, err := a.ensureHousehold(ctx)
	if err != nil {
		return err
	}
	if role == "" {
		role = "MEMBER"
	}
	if err := a.Client.InviteToHousehold(ctx, householdID, email, role); err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to send invite: %v", err))
		return err
	}
	return a.Formatter.Success(fmt.Sprintf("Invitation sent to %s.", email))
}

// HouseholdLeaderboard shows the household leaderboard.
func (a *App) HouseholdLeaderboard(ctx context.Context, period string) error {
	householdID, err := a.ensureHousehold(ctx)
	if err != nil {
		return err
	}
	if period == "" {
		period = "WEEKLY"
	}
	lb, err := a.Client.GetHouseholdLeaderboard(ctx, householdID, period)
	if err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to get leaderboard: %v", err))
		return err
	}
	if a.Formatter.JSON {
		return a.Formatter.Print(lb)
	}
	b, err := json.MarshalIndent(lb, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to format output: %w", err)
	}
	return a.Formatter.Println(string(b))
}

// HouseholdChallenges shows active household challenges.
func (a *App) HouseholdChallenges(ctx context.Context) error {
	householdID, err := a.ensureHousehold(ctx)
	if err != nil {
		return err
	}
	challenges, err := a.Client.GetHouseholdChallenges(ctx, householdID)
	if err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to get challenges: %v", err))
		return err
	}
	if a.Formatter.JSON {
		return a.Formatter.Print(challenges)
	}
	if len(challenges) == 0 {
		return a.Formatter.Println("No active challenges.")
	}
	b, err := json.MarshalIndent(challenges, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to format output: %w", err)
	}
	return a.Formatter.Println(string(b))
}
