package integration

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"
)

func TestSearchTask(t *testing.T) {
	EnsureBackend(t)

	token := getDemoToken(t)
	client := &http.Client{Timeout: 30 * time.Second}
	base := BackendURL()

	// Step 1: Create a task with a distinctive title
	uniqueTitle := "SearchTest-CLI33-" + fmt.Sprintf("%d", time.Now().UnixNano())
	taskBody := fmt.Sprintf(`{"title":%q,"priority":"CRITICAL_NOW"}`, uniqueTitle)

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
	t.Logf("Created task %s with title %q", createdTask.ID, uniqueTitle)

	// Step 2: Search for the task by title
	searchURL := fmt.Sprintf("%s/api/v2/unified-tasks?q=%s", base, url.QueryEscape(uniqueTitle))
	searchReq, err := http.NewRequest("GET", searchURL, nil)
	if err != nil {
		t.Fatalf("failed to build search request: %v", err)
	}
	searchReq.Header.Set("Authorization", "Bearer "+token)

	searchResp, err := client.Do(searchReq)
	if err != nil {
		t.Fatalf("search request failed: %v", err)
	}
	defer searchResp.Body.Close()

	if searchResp.StatusCode != http.StatusOK {
		t.Fatalf("search: expected 200, got %d", searchResp.StatusCode)
	}

	var results []struct {
		ID    string `json:"id"`
		Title string `json:"title"`
	}
	if err := json.NewDecoder(searchResp.Body).Decode(&results); err != nil {
		t.Fatalf("failed to decode search results: %v", err)
	}

	found := false
	for _, r := range results {
		if r.ID == createdTask.ID {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("task %s not found in search results (%d results)", createdTask.ID, len(results))
	}
	t.Logf("Search returned %d result(s); task found: %v", len(results), found)

	// Cleanup: delete the task
	deleteReq, err := http.NewRequest("DELETE", fmt.Sprintf("%s/api/v2/unified-tasks/%s", base, createdTask.ID), nil)
	if err != nil {
		t.Fatalf("failed to build delete request: %v", err)
	}
	deleteReq.Header.Set("Authorization", "Bearer "+token)
	deleteResp, err := client.Do(deleteReq)
	if err != nil {
		t.Logf("cleanup delete failed: %v", err)
		return
	}
	defer deleteResp.Body.Close()
	t.Logf("Cleanup delete status: %d", deleteResp.StatusCode)
}
