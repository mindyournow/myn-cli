package auth

import (
	"context"
	"fmt"
	"sync"
	"time"
)

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
func (tc *TokenCache) SetAccessToken(token string, expiresIn time.Duration) {
	tc.mu.Lock()
	defer tc.mu.Unlock()
	tc.accessToken = token
	// Subtract 60s buffer so we refresh before actual expiry
	tc.expiresAt = time.Now().Add(expiresIn - 60*time.Second)
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

	ttl := time.Duration(resp.ExpiresIn) * time.Second
	if ttl <= 0 {
		ttl = 3600 * time.Second
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
