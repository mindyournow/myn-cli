package app

import (
	"context"
	"fmt"

	mynerrors "github.com/mindyournow/myn-cli/internal/errors"
)

// ensureAuth loads the access token into the API client.
// Must be called before any authenticated API request.
func (a *App) ensureAuth(ctx context.Context) error {
	// Try API key first
	if apiKey, err := a.KeyStore.LoadAPIKey(); err == nil && apiKey != "" {
		a.Client.SetAPIKey(apiKey)
		return nil
	}
	// OAuth: get (or refresh) access token
	accessToken, err := a.TokenCache.GetAccessToken(ctx)
	if err != nil {
		return mynerrors.Auth("not authenticated", err).WithHint("Run 'mynow login' to authenticate")
	}
	a.Client.SetToken(accessToken)
	return nil
}

// ensureHousehold returns the first household ID from the customer profile.
// Used by commands that require a household context (grocery, household commands).
func (a *App) ensureHousehold(ctx context.Context) (string, error) {
	if err := a.ensureAuth(ctx); err != nil {
		return "", err
	}
	households, err := a.Client.GetHouseholdInfo(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get household: %w", err)
	}
	if len(households) == 0 {
		return "", fmt.Errorf("no household found — you need to be part of a household to use this command")
	}
	return households[0].ID, nil
}
