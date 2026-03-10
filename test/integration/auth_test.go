package integration

import (
	"encoding/json"
	"net/http"
	"os"
	"testing"
	"time"
)

func TestAPIKeyAuth(t *testing.T) {
	EnsureBackend(t)

	// Read demo API key from environment
	demoAPIKey := os.Getenv("MYN_TEST_DEMO_KEY")
	if demoAPIKey == "" {
		t.Skip("MYN_TEST_DEMO_KEY not set")
	}

	client := &http.Client{Timeout: 30 * time.Second}

	// Create demo account to get a token
	req, err := http.NewRequest("POST", BackendURL()+"/api/v1/admin/demo/recreate-account", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	req.Header.Set("X-Demo-API-Key", demoAPIKey)

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("demo account failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var demoResp struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&demoResp); err != nil {
		t.Fatalf("failed to decode demo response: %v", err)
	}
	if demoResp.Token == "" {
		t.Fatal("no token in response")
	}

	// Validate token against customers endpoint
	profReq, err := http.NewRequest("GET", BackendURL()+"/api/v1/customers", nil)
	if err != nil {
		t.Fatalf("failed to create profile request: %v", err)
	}
	profReq.Header.Set("Authorization", "Bearer "+demoResp.Token)

	profResp, err := client.Do(profReq)
	if err != nil {
		t.Fatalf("profile request failed: %v", err)
	}
	defer profResp.Body.Close()

	if profResp.StatusCode != http.StatusOK {
		t.Fatalf("profile: expected 200, got %d", profResp.StatusCode)
	}

	t.Log("Auth flow: OK")
}
