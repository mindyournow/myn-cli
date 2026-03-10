package api

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestArchiveTask verifies ArchiveTask sends POST to the archive endpoint
// and returns a decoded UnifiedTask with the correct ID.
func TestArchiveTask(t *testing.T) {
	const taskID = "task-abc-123"
	response := UnifiedTask{
		ID:         taskID,
		Title:      "Archived Task",
		IsArchived: true,
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if !strings.Contains(r.URL.Path, "archive") {
			t.Errorf("expected path to contain 'archive', got %s", r.URL.Path)
		}
		if !strings.Contains(r.URL.Path, taskID) {
			t.Errorf("expected path to contain task ID %q, got %s", taskID, r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient(server.URL)
	task, err := client.ArchiveTask(context.Background(), taskID)
	if err != nil {
		t.Fatalf("ArchiveTask() error = %v", err)
	}
	if task.ID != taskID {
		t.Errorf("task.ID = %q, want %q", task.ID, taskID)
	}
	if !task.IsArchived {
		t.Error("task.IsArchived should be true")
	}
}

// TestCancelTimer verifies CancelTimer sends POST to the cancel endpoint
// and returns a decoded Timer.
func TestCancelTimer(t *testing.T) {
	const timerID = "timer-xyz-456"
	response := Timer{
		ID:     timerID,
		Status: "cancelled",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if !strings.Contains(r.URL.Path, "cancel") {
			t.Errorf("expected path to contain 'cancel', got %s", r.URL.Path)
		}
		if !strings.Contains(r.URL.Path, timerID) {
			t.Errorf("expected path to contain timer ID %q, got %s", timerID, r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient(server.URL)
	timer, err := client.CancelTimer(context.Background(), timerID)
	if err != nil {
		t.Fatalf("CancelTimer() error = %v", err)
	}
	if timer.ID != timerID {
		t.Errorf("timer.ID = %q, want %q", timer.ID, timerID)
	}
}

// TestSnoozeTimer verifies SnoozeTimer sends POST to the snooze endpoint
// with snoozeMinutes in the request body.
func TestSnoozeTimer(t *testing.T) {
	const timerID = "timer-snooze-789"
	const snoozeMinutes = 5
	response := Timer{
		ID:     timerID,
		Status: "snoozed",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if !strings.Contains(r.URL.Path, "snooze") {
			t.Errorf("expected path to contain 'snooze', got %s", r.URL.Path)
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("failed to read request body: %v", err)
		}
		var bodyMap map[string]interface{}
		if err := json.Unmarshal(body, &bodyMap); err != nil {
			t.Fatalf("failed to parse request body: %v", err)
		}
		snooze, ok := bodyMap["snoozeMinutes"]
		if !ok {
			t.Error("request body should contain 'snoozeMinutes'")
		} else if int(snooze.(float64)) != snoozeMinutes {
			t.Errorf("snoozeMinutes = %v, want %d", snooze, snoozeMinutes)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient(server.URL)
	timer, err := client.SnoozeTimer(context.Background(), timerID, snoozeMinutes)
	if err != nil {
		t.Fatalf("SnoozeTimer() error = %v", err)
	}
	if timer.ID != timerID {
		t.Errorf("timer.ID = %q, want %q", timer.ID, timerID)
	}
}

// TestCreateAIConversation verifies CreateAIConversation sends POST to the
// conversations endpoint and returns an AIConversation with a conversationId.
func TestCreateAIConversation(t *testing.T) {
	const convID = "conv-abc-001"
	response := AIConversation{
		ID:    convID,
		Title: "New Conversation",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/ai/conversations" {
			t.Errorf("expected path /api/v1/ai/conversations, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient(server.URL)
	conv, err := client.CreateAIConversation(context.Background(), "New Conversation")
	if err != nil {
		t.Fatalf("CreateAIConversation() error = %v", err)
	}
	if conv.ID != convID {
		t.Errorf("conv.ID = %q, want %q", conv.ID, convID)
	}
}

// TestArchiveAIConversation verifies ArchiveAIConversation sends PATCH to the
// status endpoint with {"isArchived": true}.
func TestArchiveAIConversation(t *testing.T) {
	const convID = "conv-archive-002"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("expected PATCH, got %s", r.Method)
		}
		expectedPath := "/api/v1/ai/conversations/" + convID + "/status"
		if r.URL.Path != expectedPath {
			t.Errorf("expected path %q, got %q", expectedPath, r.URL.Path)
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("failed to read request body: %v", err)
		}
		var bodyMap map[string]interface{}
		if err := json.Unmarshal(body, &bodyMap); err != nil {
			t.Fatalf("failed to parse request body: %v", err)
		}
		if isArchived, ok := bodyMap["isArchived"]; !ok || isArchived != true {
			t.Errorf("expected body[\"isArchived\"] = true, got %v", bodyMap)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	}))
	defer server.Close()

	client := NewClient(server.URL)
	err := client.ArchiveAIConversation(context.Background(), convID)
	if err != nil {
		t.Fatalf("ArchiveAIConversation() error = %v", err)
	}
}

// TestFavoriteAIConversation verifies FavoriteAIConversation sends PATCH to the
// status endpoint with {"favorited": true}.
func TestFavoriteAIConversation(t *testing.T) {
	const convID = "conv-fav-003"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("expected PATCH, got %s", r.Method)
		}
		expectedPath := "/api/v1/ai/conversations/" + convID + "/status"
		if r.URL.Path != expectedPath {
			t.Errorf("expected path %q, got %q", expectedPath, r.URL.Path)
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("failed to read request body: %v", err)
		}
		var bodyMap map[string]interface{}
		if err := json.Unmarshal(body, &bodyMap); err != nil {
			t.Fatalf("failed to parse request body: %v", err)
		}
		if favorited, ok := bodyMap["favorited"]; !ok || favorited != true {
			t.Errorf("expected body[\"favorited\"] = true, got %v", bodyMap)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	}))
	defer server.Close()

	client := NewClient(server.URL)
	err := client.FavoriteAIConversation(context.Background(), convID)
	if err != nil {
		t.Fatalf("FavoriteAIConversation() error = %v", err)
	}
}

// TestSearchAIConversations verifies SearchAIConversations sends GET with
// the query param q="hello" and returns a slice of AIConversation.
func TestSearchAIConversations(t *testing.T) {
	response := []AIConversation{
		{ID: "conv-search-001", Title: "Hello World"},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/ai/conversations/search" {
			t.Errorf("expected path /api/v1/ai/conversations/search, got %s", r.URL.Path)
		}
		q := r.URL.Query().Get("q")
		if q != "hello" {
			t.Errorf("query param q = %q, want %q", q, "hello")
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient(server.URL)
	convs, err := client.SearchAIConversations(context.Background(), "hello")
	if err != nil {
		t.Fatalf("SearchAIConversations() error = %v", err)
	}
	if len(convs) != 1 {
		t.Fatalf("len(convs) = %d, want 1", len(convs))
	}
	if convs[0].ID != "conv-search-001" {
		t.Errorf("convs[0].ID = %q, want %q", convs[0].ID, "conv-search-001")
	}
}

// TestGetAIConversationStats verifies GetAIConversationStats sends GET to the
// stats endpoint and returns a map.
func TestGetAIConversationStats(t *testing.T) {
	response := map[string]interface{}{
		"totalConversations": float64(42),
		"archivedCount":      float64(5),
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/ai/conversations/stats" {
			t.Errorf("expected path /api/v1/ai/conversations/stats, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient(server.URL)
	stats, err := client.GetAIConversationStats(context.Background())
	if err != nil {
		t.Fatalf("GetAIConversationStats() error = %v", err)
	}
	if stats["totalConversations"] != float64(42) {
		t.Errorf("totalConversations = %v, want 42", stats["totalConversations"])
	}
}

// TestScheduleHabits_IsPost verifies ScheduleHabits uses POST (not GET)
// and sends numberOfDays in the request body.
func TestScheduleHabits_IsPost(t *testing.T) {
	const days = 7
	response := map[string]interface{}{
		"scheduled": float64(3),
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("failed to read request body: %v", err)
		}
		var bodyMap map[string]interface{}
		if err := json.Unmarshal(body, &bodyMap); err != nil {
			t.Fatalf("failed to parse request body: %v", err)
		}
		if numDays, ok := bodyMap["numberOfDays"]; !ok {
			t.Error("request body should contain 'numberOfDays'")
		} else if int(numDays.(float64)) != days {
			t.Errorf("numberOfDays = %v, want %d", numDays, days)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient(server.URL)
	result, err := client.ScheduleHabits(context.Background(), days)
	if err != nil {
		t.Fatalf("ScheduleHabits() error = %v", err)
	}
	if result["scheduled"] != float64(3) {
		t.Errorf("result[\"scheduled\"] = %v, want 3", result["scheduled"])
	}
}

// TestCalculateSmartTime verifies CalculateSmartTime sends POST to the
// calculate-smart-time endpoint for the given habitId and returns a map.
func TestCalculateSmartTime(t *testing.T) {
	const habitID = "habit-smart-001"
	response := map[string]interface{}{
		"suggestedTime": "08:00",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		expectedPath := "/api/habits/reminders/" + habitID + "/calculate-smart-time"
		if r.URL.Path != expectedPath {
			t.Errorf("expected path %q, got %q", expectedPath, r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient(server.URL)
	result, err := client.CalculateSmartTime(context.Background(), habitID)
	if err != nil {
		t.Fatalf("CalculateSmartTime() error = %v", err)
	}
	if result["suggestedTime"] != "08:00" {
		t.Errorf("result[\"suggestedTime\"] = %v, want \"08:00\"", result["suggestedTime"])
	}
}

// TestListTasks_SendsPriorityParam verifies ListTasks forwards the Priority
// field as the "priority" query parameter.
func TestListTasks_SendsPriorityParam(t *testing.T) {
	response := []UnifiedTask{
		{ID: "task-critical-001", Title: "Critical Task"},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		priority := r.URL.Query().Get("priority")
		if priority != "CRITICAL" {
			t.Errorf("query param priority = %q, want %q", priority, "CRITICAL")
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient(server.URL)
	tasks, err := client.ListTasks(context.Background(), TaskListParams{
		Priority: "CRITICAL",
	})
	if err != nil {
		t.Fatalf("ListTasks() error = %v", err)
	}
	if len(tasks) != 1 {
		t.Fatalf("len(tasks) = %d, want 1", len(tasks))
	}
	if tasks[0].ID != "task-critical-001" {
		t.Errorf("tasks[0].ID = %q, want %q", tasks[0].ID, "task-critical-001")
	}
}
