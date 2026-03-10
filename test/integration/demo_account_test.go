package integration

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"
)

// TestDemoAccountFlow validates the demo account lifecycle that integration
// tests depend on: create a demo account, get a token, use it to hit the API.
func TestDemoAccountFlow(t *testing.T) {
	EnsureBackend(t)

	// Read demo API key from environment — never hardcode it (MED-5 fix: consistent with docker-compose.yml)
	demoAPIKey := os.Getenv("MYN_TEST_DEMO_KEY")
	if demoAPIKey == "" {
		t.Fatal("MYN_TEST_DEMO_KEY environment variable must be set")
	}

	// Use a client with a timeout instead of http.DefaultClient (MED-5 fix)
	httpClient := &http.Client{Timeout: 30 * time.Second}

	baseURL := BackendURL()

	// Step 1: Create demo account
	req, err := http.NewRequest("POST", baseURL+"/api/v1/admin/demo/recreate-account", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("X-Demo-API-Key", demoAPIKey)

	resp, err := httpClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to create demo account: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected 200, got %d", resp.StatusCode)
	}

	var demoResp struct {
		Token      string `json:"token"`
		CustomerID string `json:"customerId"`
		Email      string `json:"email"`
		Username   string `json:"username"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&demoResp); err != nil {
		t.Fatalf("Failed to decode demo response: %v", err)
	}

	if demoResp.Token == "" {
		t.Fatal("Demo account token is empty")
	}
	t.Logf("Demo account created: %s (%s)", demoResp.Username, demoResp.Email)

	// Step 2: Use the token to list tasks (validates auth works end-to-end)
	taskReq, err := http.NewRequest("GET", baseURL+"/api/v2/unified-tasks", nil)
	if err != nil {
		t.Fatalf("Failed to create task request: %v", err)
	}
	taskReq.Header.Set("Authorization", "Bearer "+demoResp.Token)

	taskResp, err := httpClient.Do(taskReq)
	if err != nil {
		t.Fatalf("Failed to list tasks: %v", err)
	}
	defer taskResp.Body.Close()

	if taskResp.StatusCode != http.StatusOK {
		t.Fatalf("Expected 200 for task list, got %d", taskResp.StatusCode)
	}

	t.Log("Successfully authenticated and listed tasks via demo account.")

	// Step 3: Create a task via API (validates write operations)
	taskBody := `{"title":"Integration test task","priority":"CRITICAL_NOW"}`
	createReq, err := http.NewRequest("POST", baseURL+"/api/v2/unified-tasks",
		strings.NewReader(taskBody))
	if err != nil {
		t.Fatalf("Failed to create task request: %v", err)
	}
	createReq.Header.Set("Authorization", "Bearer "+demoResp.Token)
	createReq.Header.Set("Content-Type", "application/json")

	createResp, err := httpClient.Do(createReq)
	if err != nil {
		t.Fatalf("Failed to create task: %v", err)
	}
	defer createResp.Body.Close()

	if createResp.StatusCode != http.StatusOK && createResp.StatusCode != http.StatusCreated {
		t.Fatalf("Expected 200/201 for task create, got %d", createResp.StatusCode)
	}

	t.Log("Successfully created a task via API.")
}
