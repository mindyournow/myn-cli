package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// TestTokenPersistence_FullLifecycle tests the exact bug scenario:
// 1. OAuth login issues tokens
// 2. Refresh token is saved to disk via KeyStore
// 3. A NEW TokenCache (simulating a new process) loads the saved token
// 4. The new cache can successfully refresh the access token
func TestTokenPersistence_FullLifecycle(t *testing.T) {
	// Set up a temp directory to act as config dir
	tmpDir := t.TempDir()

	// Create a mock OAuth server that issues and accepts refresh tokens
	const (
		initialAccessToken  = "access-token-v1"
		initialRefreshToken = "refresh-token-v1"
		refreshedAccessTok  = "access-token-v2"
		rotatedRefreshTok   = "refresh-token-v2"
	)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/mcp/oauth/token" && r.Method == http.MethodPost {
			if err := r.ParseForm(); err != nil {
				http.Error(w, "bad form", http.StatusBadRequest)
				return
			}
			grantType := r.FormValue("grant_type")
			rt := r.FormValue("refresh_token")

			if grantType == "refresh_token" && rt == initialRefreshToken {
				// First refresh: return new tokens with a rotated refresh token
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(TokenResponse{
					AccessToken:  refreshedAccessTok,
					RefreshToken: rotatedRefreshTok,
					TokenType:    "Bearer",
					ExpiresIn:    3600,
				})
				return
			}
			if grantType == "refresh_token" && rt == rotatedRefreshTok {
				// Second refresh with rotated token should also work
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(TokenResponse{
					AccessToken:  "access-token-v3",
					RefreshToken: rotatedRefreshTok, // Same refresh token (no rotation this time)
					TokenType:    "Bearer",
					ExpiresIn:    3600,
				})
				return
			}
			http.Error(w, `{"error":"invalid_grant"}`, http.StatusBadRequest)
			return
		}
		http.Error(w, "not found", http.StatusNotFound)
	}))
	defer server.Close()

	// --- STEP 1: Simulate login saving the refresh token ---
	fileKeyring := NewKeyring(tmpDir)
	keyStore := NewKeyStore(fileKeyring, "file") // Force file backend (no OS keyring in tests)

	// Save the refresh token (this is what app.Login() now does)
	if err := keyStore.SaveRefreshToken(initialRefreshToken); err != nil {
		t.Fatalf("SaveRefreshToken() error = %v", err)
	}

	// Verify the file was created on disk
	credPath := filepath.Join(tmpDir, "credentials", "refresh_token.enc")
	if _, err := os.Stat(credPath); os.IsNotExist(err) {
		t.Fatal("refresh_token.enc was not created on disk")
	}

	// --- STEP 2: Simulate a NEW process creating a fresh TokenCache ---
	// This is what NewWithConfig() does when a new `mynow tui` process starts
	fileKeyring2 := NewKeyring(tmpDir)
	keyStore2 := NewKeyStore(fileKeyring2, "file")
	oauthClient2 := NewOAuthClient(server.URL, keyStore2)

	// Set client ID (normally this would come from registration)
	oauthClient2.ClientID = "test-client-id"

	tokenCache2 := NewTokenCache(keyStore2, oauthClient2)

	// The cache starts empty — no in-memory access token
	if !tokenCache2.IsExpired() {
		t.Fatal("New token cache should start expired (no in-memory token)")
	}

	// --- STEP 3: GetAccessToken should trigger a refresh using the stored token ---
	ctx := context.Background()
	accessToken, err := tokenCache2.GetAccessToken(ctx)
	if err != nil {
		t.Fatalf("GetAccessToken() error = %v", err)
	}

	if accessToken != refreshedAccessTok {
		t.Errorf("GetAccessToken() = %q, want %q", accessToken, refreshedAccessTok)
	}

	// Verify the in-memory cache is now populated
	if tokenCache2.IsExpired() {
		t.Error("Token cache should not be expired after successful refresh")
	}

	// --- STEP 4: Verify the rotated refresh token was saved ---
	loadedToken, err := keyStore2.LoadRefreshToken()
	if err != nil {
		t.Fatalf("LoadRefreshToken() after rotation error = %v", err)
	}
	if loadedToken != rotatedRefreshTok {
		t.Errorf("Rotated refresh token not saved: got %q, want %q", loadedToken, rotatedRefreshTok)
	}

	// --- STEP 5: Simulate ANOTHER new process — should use disk-cached access token ---
	fileKeyring3 := NewKeyring(tmpDir)
	keyStore3 := NewKeyStore(fileKeyring3, "file")
	oauthClient3 := NewOAuthClient(server.URL, keyStore3)
	oauthClient3.ClientID = "test-client-id"
	tokenCache3 := NewTokenCache(keyStore3, oauthClient3)

	accessToken3, err := tokenCache3.GetAccessToken(ctx)
	if err != nil {
		t.Fatalf("GetAccessToken() on third process error = %v", err)
	}
	// Should return the disk-cached access token from step 3 (still valid, not expired)
	if accessToken3 != refreshedAccessTok {
		t.Errorf("Third process GetAccessToken() = %q, want %q (disk-cached)", accessToken3, refreshedAccessTok)
	}
}

// TestTokenPersistence_NoRefreshToken verifies that GetAccessToken fails
// gracefully when no refresh token is stored (user never logged in).
func TestTokenPersistence_NoRefreshToken(t *testing.T) {
	tmpDir := t.TempDir()
	fileKeyring := NewKeyring(tmpDir)
	keyStore := NewKeyStore(fileKeyring, "file")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("Server should not be called when no refresh token exists")
	}))
	defer server.Close()

	oauthClient := NewOAuthClient(server.URL, keyStore)
	oauthClient.ClientID = "test-client-id"
	tokenCache := NewTokenCache(keyStore, oauthClient)

	_, err := tokenCache.GetAccessToken(context.Background())
	if err == nil {
		t.Fatal("GetAccessToken() should fail when no refresh token is stored")
	}
}

// TestTokenPersistence_ExpiredThenRefresh verifies that an expired in-memory
// token triggers a refresh using the stored refresh token.
func TestTokenPersistence_ExpiredThenRefresh(t *testing.T) {
	tmpDir := t.TempDir()
	fileKeyring := NewKeyring(tmpDir)
	keyStore := NewKeyStore(fileKeyring, "file")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/mcp/oauth/token" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(TokenResponse{
				AccessToken:  "fresh-access-token",
				RefreshToken: "same-refresh-token",
				TokenType:    "Bearer",
				ExpiresIn:    3600,
			})
			return
		}
	}))
	defer server.Close()

	oauthClient := NewOAuthClient(server.URL, keyStore)
	oauthClient.ClientID = "test-client-id"
	tokenCache := NewTokenCache(keyStore, oauthClient)

	// Save a refresh token
	if err := keyStore.SaveRefreshToken("same-refresh-token"); err != nil {
		t.Fatal(err)
	}

	// Set an already-expired access token in the cache
	tokenCache.SetAccessToken("old-expired-token", -1*time.Second)

	if !tokenCache.IsExpired() {
		t.Fatal("Token should be expired")
	}

	// GetAccessToken should refresh
	token, err := tokenCache.GetAccessToken(context.Background())
	if err != nil {
		t.Fatalf("GetAccessToken() error = %v", err)
	}
	if token != "fresh-access-token" {
		t.Errorf("Got %q, want %q", token, "fresh-access-token")
	}
}

// TestTokenPersistence_APIKeyBypass verifies that GetAccessToken returns ""
// for API key auth, never hitting the OAuth refresh flow.
func TestTokenPersistence_APIKeyBypass(t *testing.T) {
	tmpDir := t.TempDir()
	fileKeyring := NewKeyring(tmpDir)
	keyStore := NewKeyStore(fileKeyring, "file")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("Server should not be called for API key auth")
	}))
	defer server.Close()

	oauthClient := NewOAuthClient(server.URL, keyStore)
	tokenCache := NewTokenCache(keyStore, oauthClient)
	tokenCache.SetAuthMethod(AuthMethodAPIKey)

	token, err := tokenCache.GetAccessToken(context.Background())
	if err != nil {
		t.Fatalf("GetAccessToken() error = %v", err)
	}
	if token != "" {
		t.Errorf("Got %q, want empty string for API key auth", token)
	}
}

// TestTokenPersistence_DoubleRefreshSave verifies that both the oauth.go RefreshToken()
// AND the tokens.go Refresh() paths save the new refresh token (preventing double-save
// conflicts or missed saves).
func TestTokenPersistence_DoubleRefreshSave(t *testing.T) {
	tmpDir := t.TempDir()
	fileKeyring := NewKeyring(tmpDir)
	keyStore := NewKeyStore(fileKeyring, "file")

	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/mcp/oauth/token" {
			callCount++
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(TokenResponse{
				AccessToken:  fmt.Sprintf("access-%d", callCount),
				RefreshToken: fmt.Sprintf("refresh-%d", callCount),
				TokenType:    "Bearer",
				ExpiresIn:    3600,
			})
			return
		}
	}))
	defer server.Close()

	// Save initial refresh token
	if err := keyStore.SaveRefreshToken("refresh-0"); err != nil {
		t.Fatal(err)
	}

	oauthClient := NewOAuthClient(server.URL, keyStore)
	oauthClient.ClientID = "test-client-id"
	tokenCache := NewTokenCache(keyStore, oauthClient)

	// First refresh via TokenCache.Refresh()
	token, err := tokenCache.Refresh(context.Background())
	if err != nil {
		t.Fatalf("First Refresh() error = %v", err)
	}
	if token != "access-1" {
		t.Errorf("First Refresh() = %q, want %q", token, "access-1")
	}

	// Verify refresh-1 was saved (by both oauth.go and tokens.go)
	stored, err := keyStore.LoadRefreshToken()
	if err != nil {
		t.Fatalf("LoadRefreshToken() error = %v", err)
	}
	if stored != "refresh-1" {
		t.Errorf("Stored token = %q, want %q", stored, "refresh-1")
	}

	// The access token is now disk-cached. A second GetAccessToken call
	// should return the cached "access-1" (still valid), not trigger a new refresh.
	token2, err := tokenCache.GetAccessToken(context.Background())
	if err != nil {
		t.Fatalf("Second GetAccessToken() error = %v", err)
	}
	if token2 != "access-1" {
		t.Errorf("Second GetAccessToken() = %q, want %q (cached)", token2, "access-1")
	}
}
