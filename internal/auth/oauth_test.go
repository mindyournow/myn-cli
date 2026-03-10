package auth

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestGenerateCodeVerifier(t *testing.T) {
	verifier, err := generateCodeVerifier()
	if err != nil {
		t.Fatalf("generateCodeVerifier() error = %v", err)
	}

	// Should be URL-safe base64 encoded
	if strings.Contains(verifier, "+") || strings.Contains(verifier, "/") {
		t.Error("Code verifier should be URL-safe base64")
	}

	// Decode and verify length
	decoded, err := base64.RawURLEncoding.DecodeString(verifier)
	if err != nil {
		t.Fatalf("Failed to decode verifier: %v", err)
	}

	if len(decoded) != codeVerifierLength {
		t.Errorf("Code verifier length = %d, want %d", len(decoded), codeVerifierLength)
	}
}

func TestGenerateState(t *testing.T) {
	state, err := generateState()
	if err != nil {
		t.Fatalf("generateState() error = %v", err)
	}

	// Should be URL-safe base64 encoded
	if strings.Contains(state, "+") || strings.Contains(state, "/") {
		t.Error("State should be URL-safe base64")
	}

	// Decode and verify length
	decoded, err := base64.RawURLEncoding.DecodeString(state)
	if err != nil {
		t.Fatalf("Failed to decode state: %v", err)
	}

	if len(decoded) != stateLength {
		t.Errorf("State length = %d, want %d", len(decoded), stateLength)
	}
}

func TestGenerateCodeChallenge(t *testing.T) {
	verifier := "test-verifier-string"
	challenge := generateCodeChallenge(verifier)

	// Expected: base64url(sha256(verifier))
	hash := sha256.Sum256([]byte(verifier))
	expected := base64.RawURLEncoding.EncodeToString(hash[:])

	if challenge != expected {
		t.Errorf("Code challenge = %q, want %q", challenge, expected)
	}
}

func TestGenerateCodeVerifier_UsesCryptoRand(t *testing.T) {
	// Verify crypto/rand is used by checking that we get different values
	verifier1, err := generateCodeVerifier()
	if err != nil {
		t.Fatalf("generateCodeVerifier() error = %v", err)
	}

	verifier2, err := generateCodeVerifier()
	if err != nil {
		t.Fatalf("generateCodeVerifier() error = %v", err)
	}

	if verifier1 == verifier2 {
		t.Error("Code verifiers should be randomly generated and unique")
	}
}

func TestGenerateState_UsesCryptoRand(t *testing.T) {
	state1, err := generateState()
	if err != nil {
		t.Fatalf("generateState() error = %v", err)
	}

	state2, err := generateState()
	if err != nil {
		t.Fatalf("generateState() error = %v", err)
	}

	if state1 == state2 {
		t.Error("State values should be randomly generated and unique")
	}
}

func TestOAuthClient_NewOAuthClient(t *testing.T) {
	client := NewOAuthClient("https://api.example.com", nil)

	if client.BaseURL != "https://api.example.com" {
		t.Errorf("BaseURL = %q, want %q", client.BaseURL, "https://api.example.com")
	}

	if client.HTTPClient == nil {
		t.Error("HTTPClient should not be nil")
	}

	if client.TokenStore != nil {
		t.Error("TokenStore should be nil when passed nil")
	}
}

func TestOAuthClient_buildAuthURL(t *testing.T) {
	client := NewOAuthClient("https://api.example.com", nil)
	client.ClientID = "test-client-id"

	codeVerifier := "test-verifier"
	state := "test-state"
	redirectURI := "http://localhost:8080/callback"

	url := client.buildAuthURL(codeVerifier, state, redirectURI)

	// Verify URL contains all required parameters
	if !strings.Contains(url, "response_type=code") {
		t.Error("URL should contain response_type=code")
	}
	if !strings.Contains(url, "client_id=test-client-id") {
		t.Error("URL should contain client_id")
	}
	if !strings.Contains(url, "redirect_uri=") {
		t.Error("URL should contain redirect_uri")
	}
	if !strings.Contains(url, "code_challenge=") {
		t.Error("URL should contain code_challenge")
	}
	if !strings.Contains(url, "code_challenge_method=S256") {
		t.Error("URL should contain code_challenge_method=S256")
	}
	if !strings.Contains(url, "state=test-state") {
		t.Error("URL should contain state")
	}
}

func TestOAuthClient_RegisterClient(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST, got %s", r.Method)
		}
		if r.URL.Path != registerPath {
			t.Errorf("Expected path %s, got %s", registerPath, r.URL.Path)
		}

		var req map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Failed to decode request: %v", err)
		}

		// Verify required fields
		if req["client_name"] != "MYN CLI" {
			t.Errorf("client_name = %v, want 'MYN CLI'", req["client_name"])
		}
		if req["token_endpoint_auth_method"] != "none" {
			t.Errorf("auth_method should be 'none' for public client")
		}

		// Return client registration response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{
			"client_id": "registered-client-id",
		})
	}))
	defer server.Close()

	client := NewOAuthClient(server.URL, nil)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := client.registerClient(ctx)
	if err != nil {
		t.Fatalf("registerClient() error = %v", err)
	}

	if client.ClientID != "registered-client-id" {
		t.Errorf("ClientID = %q, want %q", client.ClientID, "registered-client-id")
	}
}

func TestOAuthClient_ExchangeCode(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST, got %s", r.Method)
		}
		if r.URL.Path != tokenPath {
			t.Errorf("Expected path %s, got %s", tokenPath, r.URL.Path)
		}

		// Verify content type
		contentType := r.Header.Get("Content-Type")
		if contentType != "application/x-www-form-urlencoded" {
			t.Errorf("Content-Type = %q, want application/x-www-form-urlencoded", contentType)
		}

		// Parse form data
		if err := r.ParseForm(); err != nil {
			t.Fatalf("Failed to parse form: %v", err)
		}

		// Verify required parameters (B3 fix - client_id should be present)
		if r.Form.Get("grant_type") != "authorization_code" {
			t.Error("grant_type should be authorization_code")
		}
		if r.Form.Get("code") != "test-code" {
			t.Error("code parameter mismatch")
		}
		if r.Form.Get("redirect_uri") != "http://localhost:8080/callback" {
			t.Error("redirect_uri parameter mismatch")
		}
		if r.Form.Get("code_verifier") != "test-verifier" {
			t.Error("code_verifier parameter mismatch")
		}
		if r.Form.Get("client_id") != "test-client" {
			t.Error("client_id should be present for public clients")
		}

		// Return token response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(TokenResponse{
			AccessToken:  "access-token-123",
			RefreshToken: "refresh-token-456",
			TokenType:    "Bearer",
			ExpiresIn:    3600,
		})
	}))
	defer server.Close()

	tokenStore := &mockTokenStore{}
	client := NewOAuthClient(server.URL, tokenStore)
	client.ClientID = "test-client"

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tokens, err := client.exchangeCode(ctx, "test-code", "test-verifier", "http://localhost:8080/callback")
	if err != nil {
		t.Fatalf("exchangeCode() error = %v", err)
	}

	if tokens.AccessToken != "access-token-123" {
		t.Errorf("AccessToken = %q, want %q", tokens.AccessToken, "access-token-123")
	}

	// Verify refresh token was saved (B10 fix - error not ignored)
	if tokenStore.savedToken != "refresh-token-456" {
		t.Error("Refresh token should have been saved")
	}
}

func TestOAuthClient_RefreshToken(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Fatalf("Failed to parse form: %v", err)
		}

		if r.Form.Get("grant_type") != "refresh_token" {
			t.Error("grant_type should be refresh_token")
		}
		if r.Form.Get("refresh_token") != "old-refresh-token" {
			t.Error("refresh_token parameter mismatch")
		}
		if r.Form.Get("client_id") != "test-client" {
			t.Error("client_id should be present")
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(TokenResponse{
			AccessToken:  "new-access-token",
			RefreshToken: "new-refresh-token",
			TokenType:    "Bearer",
			ExpiresIn:    3600,
		})
	}))
	defer server.Close()

	tokenStore := &mockTokenStore{}
	client := NewOAuthClient(server.URL, tokenStore)
	client.ClientID = "test-client"

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tokens, err := client.RefreshToken(ctx, "old-refresh-token")
	if err != nil {
		t.Fatalf("RefreshToken() error = %v", err)
	}

	if tokens.AccessToken != "new-access-token" {
		t.Errorf("AccessToken = %q, want %q", tokens.AccessToken, "new-access-token")
	}

	if tokenStore.savedToken != "new-refresh-token" {
		t.Error("New refresh token should have been saved")
	}
}

// mockTokenStore is a test double for TokenStore
type mockTokenStore struct {
	savedToken string
	saveErr    error
	loadToken  string
	loadErr    error
	clearErr   error
}

func (m *mockTokenStore) SaveRefreshToken(token string) error {
	if m.saveErr != nil {
		return m.saveErr
	}
	m.savedToken = token
	return nil
}

func (m *mockTokenStore) LoadRefreshToken() (string, error) {
	if m.loadErr != nil {
		return "", m.loadErr
	}
	return m.loadToken, nil
}

func (m *mockTokenStore) Clear() error {
	return m.clearErr
}

func TestOAuthClient_RefreshToken_SaveError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(TokenResponse{
			AccessToken:  "new-access-token",
			RefreshToken: "new-refresh-token",
			TokenType:    "Bearer",
			ExpiresIn:    3600,
		})
	}))
	defer server.Close()

	tokenStore := &mockTokenStore{saveErr: errTest}
	client := NewOAuthClient(server.URL, tokenStore)
	client.ClientID = "test-client"

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := client.RefreshToken(ctx, "old-token")
	if err == nil {
		t.Error("RefreshToken() should return error when save fails (B10 fix)")
	}
}

var errTest = context.DeadlineExceeded
