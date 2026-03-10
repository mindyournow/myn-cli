package auth

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"
)

const (
	minTokenTTL = 60    // 1 minute
	maxTokenTTL = 86400 // 24 hours
)

// clampTTL returns a safe TTL in seconds.
// The second return value is true if clamping occurred.
func clampTTL(ttl int64) (time.Duration, bool) {
	if ttl < minTokenTTL {
		return time.Duration(minTokenTTL) * time.Second, true
	}
	if ttl > maxTokenTTL {
		return time.Duration(maxTokenTTL) * time.Second, true
	}
	return time.Duration(ttl) * time.Second, false
}

// AccessTokenStore is an optional interface for persisting access tokens to disk.
type AccessTokenStore interface {
	SaveAccessToken(token string, expiresAt time.Time) error
	LoadAccessToken() (string, time.Time, error)
}

// TokenCache holds the in-memory access token with its expiry time.
// Thread-safe; multiple goroutines may call GetAccessToken concurrently.
type TokenCache struct {
	mu          sync.RWMutex
	accessToken string
	expiresAt   time.Time
	authMethod  AuthMethod

	// For OAuth: used to refresh when the access token expires
	store       CredentialStore
	oauthClient *OAuthClient
}

// NewTokenCache creates a new token cache.
func NewTokenCache(store CredentialStore, oauthClient *OAuthClient) *TokenCache {
	return &TokenCache{
		store:       store,
		oauthClient: oauthClient,
	}
}

// SetAccessToken stores an access token with the given TTL.
// Also persists to disk if the store supports it.
func (tc *TokenCache) SetAccessToken(token string, expiresIn time.Duration) {
	tc.mu.Lock()
	defer tc.mu.Unlock()
	tc.accessToken = token
	// Subtract 60s buffer so we refresh before actual expiry
	tc.expiresAt = time.Now().Add(expiresIn - 60*time.Second)

	// Persist to disk for cross-process access (skip empty tokens)
	if token != "" && tc.store != nil {
		if ats, ok := tc.store.(AccessTokenStore); ok {
			if err := ats.SaveAccessToken(token, tc.expiresAt); err != nil {
				fmt.Fprintf(os.Stderr, "Warning: failed to cache access token: %v\n", err)
			}
		}
	}
}

// SetAuthMethod records which auth method is active.
func (tc *TokenCache) SetAuthMethod(m AuthMethod) {
	tc.mu.Lock()
	defer tc.mu.Unlock()
	tc.authMethod = m
}

// AuthMethod returns the active auth method.
func (tc *TokenCache) GetAuthMethod() AuthMethod {
	tc.mu.RLock()
	defer tc.mu.RUnlock()
	return tc.authMethod
}

// IsExpired returns true if the cached access token is absent or expired.
func (tc *TokenCache) IsExpired() bool {
	tc.mu.RLock()
	defer tc.mu.RUnlock()
	return tc.accessToken == "" || time.Now().After(tc.expiresAt)
}

// GetAccessToken returns the cached access token, refreshing if needed.
// For API-key auth this always returns "" (the API key is passed separately).
func (tc *TokenCache) GetAccessToken(ctx context.Context) (string, error) {
	if tc.GetAuthMethod() == AuthMethodAPIKey {
		return "", nil
	}

	if !tc.IsExpired() {
		tc.mu.RLock()
		defer tc.mu.RUnlock()
		return tc.accessToken, nil
	}

	// Try loading a cached access token from disk (cross-process persistence)
	if tc.store != nil {
		if ats, ok := tc.store.(AccessTokenStore); ok {
			token, expiresAt, err := ats.LoadAccessToken()
			if err == nil && token != "" && time.Now().Before(expiresAt) {
				tc.mu.Lock()
				tc.accessToken = token
				tc.expiresAt = expiresAt
				tc.mu.Unlock()
				return token, nil
			}
		}
	}

	// Access token expired — refresh using the stored refresh token
	return tc.Refresh(ctx)
}

// Refresh forces a new access token using the stored refresh token.
func (tc *TokenCache) Refresh(ctx context.Context) (string, error) {
	if tc.oauthClient == nil || tc.store == nil {
		return "", fmt.Errorf("no OAuth client configured for token refresh")
	}

	refreshToken, err := tc.store.LoadRefreshToken()
	if err != nil {
		return "", fmt.Errorf("no refresh token available (run 'mynow login'): %w", err)
	}

	resp, err := tc.oauthClient.RefreshToken(ctx, refreshToken)
	if err != nil {
		return "", fmt.Errorf("token refresh failed: %w", err)
	}

	// Save rotated refresh token if the server issued a new one
	if resp.RefreshToken != "" && resp.RefreshToken != refreshToken {
		if err := tc.store.SaveRefreshToken(resp.RefreshToken); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to save rotated refresh token: %v\n", err)
		}
	}

	expiresIn := int64(resp.ExpiresIn)
	if expiresIn <= 0 {
		expiresIn = 3600
	}
	ttl, clamped := clampTTL(expiresIn)
	if clamped {
		fmt.Fprintf(os.Stderr, "Warning: token expires_in %d clamped to %v\n", expiresIn, ttl)
	}
	tc.SetAccessToken(resp.AccessToken, ttl)
	return resp.AccessToken, nil
}

// ExpiresAt returns the expiry time of the current access token.
func (tc *TokenCache) ExpiresAt() time.Time {
	tc.mu.RLock()
	defer tc.mu.RUnlock()
	return tc.expiresAt
}
