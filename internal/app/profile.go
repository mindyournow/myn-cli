package app

import (
	"context"
	"encoding/json"
	"fmt"
)

// ProfileShow shows the user's full profile.
func (a *App) ProfileShow(ctx context.Context) error {
	if err := a.ensureAuth(ctx); err != nil {
		return err
	}
	profile, err := a.Client.GetCustomerProfile(ctx)
	if err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to get profile: %v", err))
		return err
	}
	if a.Formatter.JSON {
		return a.Formatter.Print(profile)
	}
	printMap := func(m map[string]interface{}, prefix string) {
		for k, v := range m {
			_ = a.Formatter.Println(fmt.Sprintf("%-20s %v", prefix+k+":", v))
		}
	}
	printMap(profile, "")
	return nil
}

// ProfileGoals shows the user's goals.
func (a *App) ProfileGoals(ctx context.Context) error {
	if err := a.ensureAuth(ctx); err != nil {
		return err
	}
	goals, err := a.Client.GetCustomerGoals(ctx)
	if err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to get goals: %v", err))
		return err
	}
	if a.Formatter.JSON {
		return a.Formatter.Print(goals)
	}
	for i, g := range goals.Goals {
		_ = a.Formatter.Println(fmt.Sprintf("%d. %s", i+1, g))
	}
	return nil
}

// ProfilePrefs shows the user's preferences.
func (a *App) ProfilePrefs(ctx context.Context) error {
	if err := a.ensureAuth(ctx); err != nil {
		return err
	}
	prefs, err := a.Client.GetCustomerPreferences(ctx)
	if err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to get preferences: %v", err))
		return err
	}
	if a.Formatter.JSON {
		return a.Formatter.Print(prefs)
	}
	if prefs.Timezone != "" {
		_ = a.Formatter.Println(fmt.Sprintf("Timezone:  %s", prefs.Timezone))
	}
	if prefs.Language != "" {
		_ = a.Formatter.Println(fmt.Sprintf("Language:  %s", prefs.Language))
	}
	return nil
}

// ProfileCoaching shows or sets coaching intensity.
func (a *App) ProfileCoaching(ctx context.Context, intensity string) error {
	if err := a.ensureAuth(ctx); err != nil {
		return err
	}
	if intensity != "" {
		if err := a.Client.SetCoachingIntensity(ctx, intensity); err != nil {
			_ = a.Formatter.Error(fmt.Sprintf("failed to set coaching intensity: %v", err))
			return err
		}
		return a.Formatter.Success(fmt.Sprintf("Coaching intensity set to %s.", intensity))
	}
	current, err := a.Client.GetCoachingIntensity(ctx)
	if err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to get coaching intensity: %v", err))
		return err
	}
	if a.Formatter.JSON {
		return a.Formatter.Print(map[string]string{"intensity": current})
	}
	return a.Formatter.Println(fmt.Sprintf("Coaching intensity: %s", current))
}

// ProfileNotifications shows notification preferences.
func (a *App) ProfileNotifications(ctx context.Context) error {
	if err := a.ensureAuth(ctx); err != nil {
		return err
	}
	prefs, err := a.Client.GetNotificationPreferences(ctx)
	if err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to get notification preferences: %v", err))
		return err
	}
	if a.Formatter.JSON {
		return a.Formatter.Print(prefs)
	}
	b, _ := json.MarshalIndent(prefs, "", "  ")
	return a.Formatter.Println(string(b))
}

// ProfileTimers shows timer preferences.
func (a *App) ProfileTimers(ctx context.Context) error {
	if err := a.ensureAuth(ctx); err != nil {
		return err
	}
	prefs, err := a.Client.GetTimerPreferences(ctx)
	if err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to get timer preferences: %v", err))
		return err
	}
	if a.Formatter.JSON {
		return a.Formatter.Print(prefs)
	}
	b, _ := json.MarshalIndent(prefs, "", "  ")
	return a.Formatter.Println(string(b))
}
