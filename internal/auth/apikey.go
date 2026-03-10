package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// APIKeyClient validates and stores API keys.
type APIKeyClient struct {
	BaseURL    string
	HTTPClient *http.Client
	Store      CredentialStore
}

// NewAPIKeyClient creates a new API key auth client.
func NewAPIKeyClient(baseURL string, store CredentialStore) *APIKeyClient {
	return &APIKeyClient{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		Store: store,
	}
}

// CustomerProfile holds the user profile returned by GET /api/v1/customers.
type CustomerProfile struct {
	ID       json.Number `json:"id"`
	Email    string      `json:"email"`
	Username string      `json:"username"`
	Name     string      `json:"firstName"`
}

// Validate checks the API key against GET /api/v1/customers and returns the profile.
func (c *APIKeyClient) Validate(ctx context.Context, apiKey string) (*CustomerProfile, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseURL+"/api/v1/customers", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("X-API-KEY", apiKey)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, fmt.Errorf("invalid API key")
	}
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return nil, fmt.Errorf("unexpected response %s: %s", resp.Status, string(body))
	}

	var profile CustomerProfile
	if err := json.NewDecoder(resp.Body).Decode(&profile); err != nil {
		return nil, fmt.Errorf("failed to decode profile: %w", err)
	}
	return &profile, nil
}

// Login validates the API key and stores it in the credential store.
func (c *APIKeyClient) Login(ctx context.Context, apiKey string) (*CustomerProfile, error) {
	profile, err := c.Validate(ctx, apiKey)
	if err != nil {
		return nil, err
	}
	if err := c.Store.SaveAPIKey(apiKey); err != nil {
		return nil, fmt.Errorf("failed to store API key: %w", err)
	}
	return profile, nil
}
