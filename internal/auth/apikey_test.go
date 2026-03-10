package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// mockCredStore is a minimal in-memory CredentialStore for testing.
type mockCredStore struct {
	refreshToken string
	apiKey       string
}

func (m *mockCredStore) SaveRefreshToken(token string) error {
	m.refreshToken = token
	return nil
}

func (m *mockCredStore) LoadRefreshToken() (string, error) {
	if m.refreshToken == "" {
		return "", fmt.Errorf("not found")
	}
	return m.refreshToken, nil
}

func (m *mockCredStore) Clear() error {
	m.refreshToken = ""
	m.apiKey = ""
	return nil
}

func (m *mockCredStore) SaveAPIKey(key string) error {
	m.apiKey = key
	return nil
}

func (m *mockCredStore) LoadAPIKey() (string, error) {
	if m.apiKey == "" {
		return "", fmt.Errorf("not found")
	}
	return m.apiKey, nil
}

func TestAPIKeyClient_Validate_Success(t *testing.T) {
	want := CustomerProfile{
		ID:       "user-1",
		Email:    "alice@example.com",
		Username: "alice",
		Name:     "Alice",
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/customers" {
			http.NotFound(w, r)
			return
		}
		if r.Header.Get("X-API-KEY") == "" {
			http.Error(w, "missing key", http.StatusUnauthorized)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(want)
	}))
	defer srv.Close()

	client := NewAPIKeyClient(srv.URL, &mockCredStore{})
	got, err := client.Validate(context.Background(), "test-api-key")
	if err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
	if got.Email != want.Email {
		t.Errorf("Validate() email = %q, want %q", got.Email, want.Email)
	}
	if got.ID != want.ID {
		t.Errorf("Validate() id = %q, want %q", got.ID, want.ID)
	}
}

func TestAPIKeyClient_Validate_Unauthorized(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
	}))
	defer srv.Close()

	client := NewAPIKeyClient(srv.URL, &mockCredStore{})
	_, err := client.Validate(context.Background(), "bad-key")
	if err == nil {
		t.Fatal("Validate() should return error on 401")
	}
	if err.Error() != "invalid API key" {
		t.Errorf("Validate() error = %q, want %q", err.Error(), "invalid API key")
	}
}

func TestAPIKeyClient_Validate_Forbidden(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "forbidden", http.StatusForbidden)
	}))
	defer srv.Close()

	client := NewAPIKeyClient(srv.URL, &mockCredStore{})
	_, err := client.Validate(context.Background(), "bad-key")
	if err == nil {
		t.Fatal("Validate() should return error on 403")
	}
	if err.Error() != "invalid API key" {
		t.Errorf("Validate() error = %q, want %q", err.Error(), "invalid API key")
	}
}

func TestAPIKeyClient_Validate_NetworkError(t *testing.T) {
	// Start and immediately close a server so the port is unreachable.
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	srv.Close()

	client := NewAPIKeyClient(srv.URL, &mockCredStore{})
	_, err := client.Validate(context.Background(), "any-key")
	if err == nil {
		t.Fatal("Validate() should return error when server is unreachable")
	}
}

func TestAPIKeyClient_Login_SavesKey(t *testing.T) {
	profile := CustomerProfile{
		ID:    "user-2",
		Email: "bob@example.com",
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(profile)
	}))
	defer srv.Close()

	store := &mockCredStore{}
	client := NewAPIKeyClient(srv.URL, store)

	const key = "login-api-key"
	got, err := client.Login(context.Background(), key)
	if err != nil {
		t.Fatalf("Login() error = %v", err)
	}
	if got.Email != profile.Email {
		t.Errorf("Login() email = %q, want %q", got.Email, profile.Email)
	}
	if store.apiKey != key {
		t.Errorf("Login() stored key = %q, want %q", store.apiKey, key)
	}
}
