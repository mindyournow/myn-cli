package api

import (
	"context"
	"fmt"
	"time"
)

// UnifiedTask represents a task from the MYN backend.
type UnifiedTask struct {
	ID             string      `json:"id"`
	Title          string      `json:"title"`
	Description    string      `json:"description,omitempty"`
	TaskType       string      `json:"taskType"`
	Priority       interface{} `json:"priority"` // null or string
	StartDate      string      `json:"startDate,omitempty"`
	DueDate        string      `json:"dueDate,omitempty"`
	Duration       int         `json:"duration,omitempty"`
	IsCompleted    bool        `json:"isCompleted"`
	IsArchived     bool        `json:"isArchived"`
	RecurrenceRule string      `json:"recurrenceRule,omitempty"`
	CreatedDate    string      `json:"createdDate,omitempty"`
	LastUpdated    string      `json:"lastUpdated,omitempty"`
	StreakCount    interface{} `json:"streakCount,omitempty"`
	CommentCount   int         `json:"commentCount,omitempty"`
	ProjectName    string      `json:"projectName,omitempty"`
	ProjectID      string      `json:"projectId,omitempty"`
}

// PriorityString returns the priority as a string (handles null).
func (t *UnifiedTask) PriorityString() string {
	switch v := t.Priority.(type) {
	case string:
		return v
	default:
		return ""
	}
}

// TaskListParams are query parameters for listing tasks.
type TaskListParams struct {
	Type             string
	Priority         string
	IsCompleted      bool
	IncludeHousehold bool
	Archived         bool
	Today            bool
	Overdue          bool
	ProjectID        string
	Sort             string
	Page             int
	Limit            int
}

// ListTasks fetches tasks from the backend.
func (c *Client) ListTasks(ctx context.Context, p TaskListParams) ([]UnifiedTask, error) {
	path := "/api/v2/unified-tasks"
	if p.Archived {
		path = "/api/v2/unified-tasks/archived"
	}
	params := map[string]string{}
	if p.Type != "" {
		params["type"] = p.Type
	}
	if p.IsCompleted {
		params["isCompleted"] = "true"
	}
	if p.IncludeHousehold {
		params["includeHousehold"] = "true"
	}
	if p.Sort != "" {
		params["sort"] = p.Sort
	}
	if p.Limit > 0 {
		params["size"] = fmt.Sprintf("%d", p.Limit)
	}
	if p.Page > 0 {
		params["page"] = fmt.Sprintf("%d", p.Page)
	}
	if p.Priority != "" {
		params["priority"] = p.Priority
	}
	if p.Today {
		params["date"] = time.Now().Format("2006-01-02")
	}
	if p.Overdue {
		params["overdue"] = "true"
	}

	resp, err := c.Get(ctx, path, params)
	if err != nil {
		return nil, err
	}
	var tasks []UnifiedTask
	if err := resp.DecodeJSON(&tasks); err != nil {
		return nil, fmt.Errorf("failed to parse task list: %w", err)
	}
	return tasks, nil
}

// GetTask fetches a single task by ID.
func (c *Client) GetTask(ctx context.Context, id string) (*UnifiedTask, error) {
	resp, err := c.Get(ctx, "/api/v2/unified-tasks/"+id, nil)
	if err != nil {
		return nil, err
	}
	var task UnifiedTask
	if err := resp.DecodeJSON(&task); err != nil {
		return nil, fmt.Errorf("failed to parse task: %w", err)
	}
	return &task, nil
}

// CreateTaskRequest is the body for creating a task.
type CreateTaskRequest struct {
	ID             string `json:"id"`
	Title          string `json:"title"`
	Description    string `json:"description,omitempty"`
	TaskType       string `json:"taskType,omitempty"`
	Priority       string `json:"priority,omitempty"`
	StartDate      string `json:"startDate,omitempty"`
	DueDate        string `json:"dueDate,omitempty"`
	Duration       int    `json:"duration,omitempty"`
	RecurrenceRule string `json:"recurrenceRule,omitempty"`
	ProjectID      string `json:"projectId,omitempty"`
	IsAutoScheduled bool  `json:"isAutoScheduled,omitempty"`
}

// CreateTask creates a new task.
func (c *Client) CreateTask(ctx context.Context, req CreateTaskRequest) (*UnifiedTask, error) {
	resp, err := c.Post(ctx, "/api/v2/unified-tasks", req)
	if err != nil {
		return nil, err
	}
	var task UnifiedTask
	if err := resp.DecodeJSON(&task); err != nil {
		return nil, fmt.Errorf("failed to parse created task: %w", err)
	}
	return &task, nil
}

// UpdateTaskRequest is the body for updating a task.
type UpdateTaskRequest struct {
	Title          string `json:"title,omitempty"`
	Description    string `json:"description,omitempty"`
	Priority       string `json:"priority,omitempty"`
	StartDate      string `json:"startDate,omitempty"`
	DueDate        string `json:"dueDate,omitempty"`
	Duration       int    `json:"duration,omitempty"`
	RecurrenceRule string `json:"recurrenceRule,omitempty"`
	ProjectID      string `json:"projectId,omitempty"`
}

// UpdateTask partially updates a task.
func (c *Client) UpdateTask(ctx context.Context, id string, req UpdateTaskRequest) (*UnifiedTask, error) {
	resp, err := c.Patch(ctx, "/api/v2/unified-tasks/"+id, req)
	if err != nil {
		return nil, err
	}
	var task UnifiedTask
	if err := resp.DecodeJSON(&task); err != nil {
		return nil, fmt.Errorf("failed to parse updated task: %w", err)
	}
	return &task, nil
}

// CompleteTask marks a task as done.
func (c *Client) CompleteTask(ctx context.Context, id string) (*UnifiedTask, error) {
	resp, err := c.Post(ctx, "/api/v2/unified-tasks/"+id+"/complete", nil)
	if err != nil {
		return nil, err
	}
	var task UnifiedTask
	if err := resp.DecodeJSON(&task); err != nil {
		return nil, fmt.Errorf("failed to parse completed task: %w", err)
	}
	return &task, nil
}

// UncompleteTask marks a task as not done.
func (c *Client) UncompleteTask(ctx context.Context, id string) (*UnifiedTask, error) {
	resp, err := c.Post(ctx, "/api/v2/unified-tasks/"+id+"/uncomplete", nil)
	if err != nil {
		return nil, err
	}
	var task UnifiedTask
	if err := resp.DecodeJSON(&task); err != nil {
		return nil, fmt.Errorf("failed to parse task: %w", err)
	}
	return &task, nil
}

// DeleteTask soft-deletes a task.
func (c *Client) DeleteTask(ctx context.Context, id string, permanent bool) error {
	path := "/api/v2/unified-tasks/" + id
	if permanent {
		path += "/permanent"
	}
	_, err := c.Delete(ctx, path)
	return err
}

// ArchiveTask archives a task (soft archive, not deletion).
func (c *Client) ArchiveTask(ctx context.Context, id string) (*UnifiedTask, error) {
	resp, err := c.Post(ctx, "/api/v2/unified-tasks/"+id+"/archive", nil)
	if err != nil {
		return nil, err
	}
	var task UnifiedTask
	if err := resp.DecodeJSON(&task); err != nil {
		return nil, fmt.Errorf("failed to parse archived task: %w", err)
	}
	return &task, nil
}

// RestoreTask restores a soft-deleted task.
func (c *Client) RestoreTask(ctx context.Context, id string) (*UnifiedTask, error) {
	resp, err := c.Post(ctx, "/api/v2/unified-tasks/"+id+"/restore", nil)
	if err != nil {
		return nil, err
	}
	var task UnifiedTask
	if err := resp.DecodeJSON(&task); err != nil {
		return nil, fmt.Errorf("failed to parse restored task: %w", err)
	}
	return &task, nil
}

// BatchUpdateRequest is the body for batch updating tasks.
type BatchUpdateRequest struct {
	IDs     []string           `json:"ids"`
	Updates UpdateTaskRequest  `json:"updates"`
}

// BatchUpdateTasks applies the same updates to multiple tasks.
func (c *Client) BatchUpdateTasks(ctx context.Context, req BatchUpdateRequest) ([]UnifiedTask, error) {
	resp, err := c.Patch(ctx, "/api/v2/unified-tasks/batch", req)
	if err != nil {
		return nil, err
	}
	var tasks []UnifiedTask
	if err := resp.DecodeJSON(&tasks); err != nil {
		return nil, fmt.Errorf("failed to parse batch result: %w", err)
	}
	return tasks, nil
}

// MoveTask moves a task to a project.
func (c *Client) MoveTask(ctx context.Context, taskID, projectID string) (*UnifiedTask, error) {
	resp, err := c.Put(ctx, "/api/project/"+projectID+"/moveTaskToProject/"+taskID, nil)
	if err != nil {
		return nil, err
	}
	var task UnifiedTask
	if err := resp.DecodeJSON(&task); err != nil {
		return nil, fmt.Errorf("failed to parse moved task: %w", err)
	}
	return &task, nil
}
