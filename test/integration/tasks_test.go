package integration

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"
)

// getDemoToken creates a fresh demo account and returns the bearer token.
func getDemoToken(t *testing.T) string {
	t.Helper()

	demoAPIKey := os.Getenv("MYN_TEST_DEMO_KEY")
	if demoAPIKey == "" {
		t.Skip("MYN_TEST_DEMO_KEY not set")
	}

	client := &http.Client{Timeout: 30 * time.Second}

	req, err := http.NewRequest("POST", BackendURL()+"/api/v1/admin/demo/recreate-account", nil)
	if err != nil {
		t.Fatalf("failed to create demo request: %v", err)
	}
	req.Header.Set("X-Demo-API-Key", demoAPIKey)

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("demo account request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("demo account: expected 200, got %d", resp.StatusCode)
	}

	var demoResp struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&demoResp); err != nil {
		t.Fatalf("failed to decode demo response: %v", err)
	}
	if demoResp.Token == "" {
		t.Fatal("demo account returned empty token")
	}
	return demoResp.Token
}

func TestTaskCRUD(t *testing.T) {
	EnsureBackend(t)

	token := getDemoToken(t)
	client := &http.Client{Timeout: 30 * time.Second}
	base := BackendURL()

	// Step 1: Create a task
	taskBody := `{"title":"CLI-33 integration test task","priority":"CRITICAL_NOW"}`
	createReq, err := http.NewRequest("POST", base+"/api/v2/unified-tasks", strings.NewReader(taskBody))
	if err != nil {
		t.Fatalf("failed to build create request: %v", err)
	}
	createReq.Header.Set("Authorization", "Bearer "+token)
	createReq.Header.Set("Content-Type", "application/json")

	createResp, err := client.Do(createReq)
	if err != nil {
		t.Fatalf("create task failed: %v", err)
	}
	defer createResp.Body.Close()

	if createResp.StatusCode != http.StatusOK && createResp.StatusCode != http.StatusCreated {
		t.Fatalf("create task: expected 200/201, got %d", createResp.StatusCode)
	}

	var createdTask struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(createResp.Body).Decode(&createdTask); err != nil {
		t.Fatalf("failed to decode created task: %v", err)
	}
	if createdTask.ID == "" {
		t.Fatal("created task has no ID")
	}
	t.Logf("Created task ID: %s", createdTask.ID)

	// Step 2: List tasks — verify the created task is present
	listReq, err := http.NewRequest("GET", base+"/api/v2/unified-tasks", nil)
	if err != nil {
		t.Fatalf("failed to build list request: %v", err)
	}
	listReq.Header.Set("Authorization", "Bearer "+token)

	listResp, err := client.Do(listReq)
	if err != nil {
		t.Fatalf("list tasks failed: %v", err)
	}
	defer listResp.Body.Close()

	if listResp.StatusCode != http.StatusOK {
		t.Fatalf("list tasks: expected 200, got %d", listResp.StatusCode)
	}

	var taskList []struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(listResp.Body).Decode(&taskList); err != nil {
		t.Fatalf("failed to decode task list: %v", err)
	}

	found := false
	for _, task := range taskList {
		if task.ID == createdTask.ID {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("created task %s not found in task list (%d tasks)", createdTask.ID, len(taskList))
	}
	t.Logf("Task list contains %d tasks; created task found: %v", len(taskList), found)

	// Step 3: Get task by ID
	getReq, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v2/unified-tasks/%s", base, createdTask.ID), nil)
	if err != nil {
		t.Fatalf("failed to build get request: %v", err)
	}
	getReq.Header.Set("Authorization", "Bearer "+token)

	getResp, err := client.Do(getReq)
	if err != nil {
		t.Fatalf("get task failed: %v", err)
	}
	defer getResp.Body.Close()

	if getResp.StatusCode != http.StatusOK {
		t.Fatalf("get task: expected 200, got %d", getResp.StatusCode)
	}
	t.Logf("Get task by ID: OK")

	// Step 4: Complete the task
	completeReq, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v2/unified-tasks/%s/complete", base, createdTask.ID), nil)
	if err != nil {
		t.Fatalf("failed to build complete request: %v", err)
	}
	completeReq.Header.Set("Authorization", "Bearer "+token)
	completeReq.Header.Set("Content-Type", "application/json")

	completeResp, err := client.Do(completeReq)
	if err != nil {
		t.Fatalf("complete task failed: %v", err)
	}
	defer completeResp.Body.Close()

	if completeResp.StatusCode != http.StatusOK && completeResp.StatusCode != http.StatusNoContent {
		t.Fatalf("complete task: expected 200/204, got %d", completeResp.StatusCode)
	}
	t.Logf("Complete task: OK")

	// Step 5: Delete the task
	deleteReq, err := http.NewRequest("DELETE", fmt.Sprintf("%s/api/v2/unified-tasks/%s", base, createdTask.ID), nil)
	if err != nil {
		t.Fatalf("failed to build delete request: %v", err)
	}
	deleteReq.Header.Set("Authorization", "Bearer "+token)

	deleteResp, err := client.Do(deleteReq)
	if err != nil {
		t.Fatalf("delete task failed: %v", err)
	}
	defer deleteResp.Body.Close()

	if deleteResp.StatusCode != http.StatusOK && deleteResp.StatusCode != http.StatusNoContent {
		t.Fatalf("delete task: expected 200/204, got %d", deleteResp.StatusCode)
	}
	t.Logf("Delete task: OK")
}
