package auth

import (
	"context"
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
	// Backend returns id as a number (Java Long), so use a numeric JSON payload
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
		w.Write([]byte(`{"id":1234,"email":"alice@example.com","username":"alice","firstName":"Alice"}`))
	}))
	defer srv.Close()

	client := NewAPIKeyClient(srv.URL, &mockCredStore{})
	got, err := client.Validate(context.Background(), "test-api-key")
	if err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
	if got.Email != "alice@example.com" {
		t.Errorf("Validate() email = %q, want %q", got.Email, "alice@example.com")
	}
	if got.ID.String() != "1234" {
		t.Errorf("Validate() id = %q, want %q", got.ID.String(), "1234")
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
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"id":5678,"email":"bob@example.com"}`))
	}))
	defer srv.Close()

	store := &mockCredStore{}
	client := NewAPIKeyClient(srv.URL, store)

	const key = "login-api-key"
	got, err := client.Login(context.Background(), key)
	if err != nil {
		t.Fatalf("Login() error = %v", err)
	}
	if got.Email != "bob@example.com" {
		t.Errorf("Login() email = %q, want %q", got.Email, "bob@example.com")
	}
	if store.apiKey != key {
		t.Errorf("Login() stored key = %q, want %q", store.apiKey, key)
	}
}
