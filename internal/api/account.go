package api

import (
	"context"
	"fmt"
)

// APIKey represents an API key.
type APIKey struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description,omitempty"`
	Prefix      string   `json:"prefix,omitempty"`
	Scopes      []string `json:"scopes,omitempty"`
	ExpiresAt   string   `json:"expiresAt,omitempty"`
	IsEnabled   bool     `json:"isEnabled"`
	CreatedAt   string   `json:"createdAt,omitempty"`
	LastUsedAt  string   `json:"lastUsedAt,omitempty"`
}

// CreateAPIKeyRequest is the body for creating an API key.
type CreateAPIKeyRequest struct {
	Name        string   `json:"name"`
	Description string   `json:"description,omitempty"`
	Scopes      []string `json:"scopes,omitempty"`
	ExpiresAt   string   `json:"expiresAt,omitempty"`
}

// CreateAPIKeyResponse includes the one-time secret value.
type CreateAPIKeyResponse struct {
	APIKey
	Secret string `json:"secret,omitempty"`
}

// ListAPIKeys fetches all API keys.
func (c *Client) ListAPIKeys(ctx context.Context) ([]APIKey, error) {
	resp, err := c.Get(ctx, "/api/v1/api-keys", nil)
	if err != nil {
		return nil, err
	}
	var keys []APIKey
	if err := resp.DecodeJSON(&keys); err != nil {
		return nil, fmt.Errorf("failed to parse API keys: %w", err)
	}
	return keys, nil
}

// CreateAPIKey creates a new API key.
func (c *Client) CreateAPIKey(ctx context.Context, req CreateAPIKeyRequest) (*CreateAPIKeyResponse, error) {
	resp, err := c.Post(ctx, "/api/v1/api-keys", req)
	if err != nil {
		return nil, err
	}
	var key CreateAPIKeyResponse
	if err := resp.DecodeJSON(&key); err != nil {
		return nil, fmt.Errorf("failed to parse created API key: %w", err)
	}
	return &key, nil
}

// GetAPIKey fetches a single API key.
func (c *Client) GetAPIKey(ctx context.Context, id string) (*APIKey, error) {
	resp, err := c.Get(ctx, "/api/v1/api-keys/"+id, nil)
	if err != nil {
		return nil, err
	}
	var key APIKey
	if err := resp.DecodeJSON(&key); err != nil {
		return nil, fmt.Errorf("failed to parse API key: %w", err)
	}
	return &key, nil
}

// RevokeAPIKey deletes an API key.
func (c *Client) RevokeAPIKey(ctx context.Context, id string) error {
	_, err := c.Delete(ctx, "/api/v1/api-keys/"+id)
	return err
}

// ExportRequest is the body for requesting a data export.
type ExportRequest struct {
	Format   string   `json:"format,omitempty"`
	Includes []string `json:"includes,omitempty"`
}

// Export represents a data export job.
type Export struct {
	ID        string `json:"id"`
	Status    string `json:"status"`
	Format    string `json:"format,omitempty"`
	CreatedAt string `json:"createdAt,omitempty"`
	ExpiresAt string `json:"expiresAt,omitempty"`
}

// RequestExport creates a new data export job.
func (c *Client) RequestExport(ctx context.Context, req ExportRequest) (*Export, error) {
	resp, err := c.Post(ctx, "/api/v1/customers/request-export", req)
	if err != nil {
		return nil, err
	}
	var export Export
	if err := resp.DecodeJSON(&export); err != nil {
		return nil, fmt.Errorf("failed to parse export: %w", err)
	}
	return &export, nil
}

// ListExports fetches all export jobs.
func (c *Client) ListExports(ctx context.Context) ([]Export, error) {
	resp, err := c.Get(ctx, "/api/v1/customers/exports", nil)
	if err != nil {
		return nil, err
	}
	var exports []Export
	if err := resp.DecodeJSON(&exports); err != nil {
		return nil, fmt.Errorf("failed to parse exports: %w", err)
	}
	return exports, nil
}

// DownloadExport fetches the raw export data.
func (c *Client) DownloadExport(ctx context.Context, id string) ([]byte, error) {
	resp, err := c.Get(ctx, "/api/v1/customers/exports/"+id+"/download", nil)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

// DeleteExport deletes an export job.
func (c *Client) DeleteExport(ctx context.Context, id string) error {
	_, err := c.Delete(ctx, "/api/v1/customers/exports/"+id)
	return err
}

// GetAccountUsage fetches today's account usage.
func (c *Client) GetAccountUsage(ctx context.Context) (map[string]interface{}, error) {
	resp, err := c.Get(ctx, "/api/v1/usage/today", nil)
	if err != nil {
		return nil, err
	}
	var usage map[string]interface{}
	if err := resp.DecodeJSON(&usage); err != nil {
		return nil, fmt.Errorf("failed to parse usage: %w", err)
	}
	return usage, nil
}

// GetAccountLimits fetches subscription limits.
func (c *Client) GetAccountLimits(ctx context.Context) (map[string]interface{}, error) {
	resp, err := c.Get(ctx, "/api/v1/usage/limits", nil)
	if err != nil {
		return nil, err
	}
	var limits map[string]interface{}
	if err := resp.DecodeJSON(&limits); err != nil {
		return nil, fmt.Errorf("failed to parse limits: %w", err)
	}
	return limits, nil
}
