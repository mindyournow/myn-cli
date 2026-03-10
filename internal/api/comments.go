package api

import (
	"context"
	"fmt"
)

// TaskComment represents a comment on a task.
type TaskComment struct {
	ID        string `json:"id"`
	TaskID    string `json:"taskId"`
	Content   string `json:"content"`
	AuthorID  string `json:"authorId,omitempty"`
	Author    string `json:"author,omitempty"`
	CreatedAt string `json:"createdAt,omitempty"`
	UpdatedAt string `json:"updatedAt,omitempty"`
}

// ListTaskComments fetches comments for a task.
func (c *Client) ListTaskComments(ctx context.Context, taskID string) ([]TaskComment, error) {
	resp, err := c.Get(ctx, "/api/v2/unified-tasks/"+taskID+"/comments", nil)
	if err != nil {
		return nil, err
	}
	var comments []TaskComment
	if err := resp.DecodeJSON(&comments); err != nil {
		return nil, fmt.Errorf("failed to parse comments: %w", err)
	}
	return comments, nil
}

// AddTaskComment adds a comment to a task.
func (c *Client) AddTaskComment(ctx context.Context, taskID, content string) (*TaskComment, error) {
	resp, err := c.Post(ctx, "/api/v2/unified-tasks/"+taskID+"/comments",
		map[string]string{"content": content})
	if err != nil {
		return nil, err
	}
	var comment TaskComment
	if err := resp.DecodeJSON(&comment); err != nil {
		return nil, fmt.Errorf("failed to parse comment: %w", err)
	}
	return &comment, nil
}

// EditTaskComment edits an existing comment.
func (c *Client) EditTaskComment(ctx context.Context, taskID, commentID, content string) (*TaskComment, error) {
	resp, err := c.Put(ctx, "/api/v2/unified-tasks/"+taskID+"/comments/"+commentID,
		map[string]string{"content": content})
	if err != nil {
		return nil, err
	}
	var comment TaskComment
	if err := resp.DecodeJSON(&comment); err != nil {
		return nil, fmt.Errorf("failed to parse comment: %w", err)
	}
	return &comment, nil
}

// DeleteTaskComment deletes a comment.
func (c *Client) DeleteTaskComment(ctx context.Context, taskID, commentID string) error {
	_, err := c.Delete(ctx, "/api/v2/unified-tasks/"+taskID+"/comments/"+commentID)
	return err
}

// GetTaskCommentCount fetches the comment count for a task.
func (c *Client) GetTaskCommentCount(ctx context.Context, taskID string) (int, error) {
	resp, err := c.Get(ctx, "/api/v2/unified-tasks/"+taskID+"/comments/count", nil)
	if err != nil {
		return 0, err
	}
	var result struct {
		Count int `json:"count"`
	}
	if err := resp.DecodeJSON(&result); err != nil {
		return 0, fmt.Errorf("failed to parse comment count: %w", err)
	}
	return result.Count, nil
}

// ShareTaskRequest is the body for sharing a task.
type ShareTaskRequest struct {
	MemberID string `json:"memberId"`
	Type     string `json:"type,omitempty"` // view, edit, delegate
	Message  string `json:"message,omitempty"`
}

// ShareTask shares a task with a household member.
func (c *Client) ShareTask(ctx context.Context, taskID string, req ShareTaskRequest) error {
	_, err := c.Post(ctx, "/api/v2/unified-tasks/"+taskID+"/share", req)
	return err
}

// RespondToShare accepts or declines a task share.
func (c *Client) RespondToShare(ctx context.Context, taskID, decision, note string) error {
	body := map[string]string{"decision": decision}
	if note != "" {
		body["note"] = note
	}
	_, err := c.Post(ctx, "/api/v2/unified-tasks/"+taskID+"/share/respond", body)
	return err
}

// RevokeTaskShare revokes a task share.
func (c *Client) RevokeTaskShare(ctx context.Context, taskID, memberID string) error {
	_, err := c.Delete(ctx, "/api/v2/unified-tasks/"+taskID+"/share/"+memberID)
	return err
}

// ListTaskShares lists shares for a task.
func (c *Client) ListTaskShares(ctx context.Context, taskID string) ([]map[string]interface{}, error) {
	resp, err := c.Get(ctx, "/api/v2/unified-tasks/"+taskID+"/shares", nil)
	if err != nil {
		return nil, err
	}
	var shares []map[string]interface{}
	if err := resp.DecodeJSON(&shares); err != nil {
		return nil, fmt.Errorf("failed to parse shares: %w", err)
	}
	return shares, nil
}

// GetSharedInbox fetches tasks shared with the current user.
func (c *Client) GetSharedInbox(ctx context.Context) ([]UnifiedTask, error) {
	resp, err := c.Get(ctx, "/api/v2/unified-tasks/shared-with-me", nil)
	if err != nil {
		return nil, err
	}
	var tasks []UnifiedTask
	if err := resp.DecodeJSON(&tasks); err != nil {
		return nil, fmt.Errorf("failed to parse shared inbox: %w", err)
	}
	return tasks, nil
}
