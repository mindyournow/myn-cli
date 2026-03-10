package api

import (
	"context"
	"fmt"
)

// Notification represents a notification from the backend.
type Notification struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Body      string `json:"body,omitempty"`
	IsRead    bool   `json:"isRead"`
	CreatedAt string `json:"createdAt,omitempty"`
	Type      string `json:"type,omitempty"`
}

// ListNotifications fetches notifications.
func (c *Client) ListNotifications(ctx context.Context, page, size int) ([]Notification, error) {
	params := map[string]string{}
	if page > 0 {
		params["page"] = fmt.Sprintf("%d", page)
	}
	if size > 0 {
		params["size"] = fmt.Sprintf("%d", size)
	}
	resp, err := c.Get(ctx, "/api/v1/notifications", params)
	if err != nil {
		return nil, err
	}
	var notifications []Notification
	if err := resp.DecodeJSON(&notifications); err != nil {
		return nil, fmt.Errorf("failed to parse notifications: %w", err)
	}
	return notifications, nil
}

// GetUnreadNotificationCount fetches the unread notification count.
func (c *Client) GetUnreadNotificationCount(ctx context.Context) (int, error) {
	resp, err := c.Get(ctx, "/api/v1/notifications/unread-count", nil)
	if err != nil {
		return 0, err
	}
	var result struct {
		Count int `json:"count"`
	}
	if err := resp.DecodeJSON(&result); err != nil {
		return 0, fmt.Errorf("failed to parse notification count: %w", err)
	}
	return result.Count, nil
}

// MarkNotificationRead marks a notification as read.
func (c *Client) MarkNotificationRead(ctx context.Context, id string) error {
	_, err := c.Put(ctx, "/api/v1/notifications/"+id+"/read", nil)
	return err
}

// MarkAllNotificationsRead marks all notifications as read.
func (c *Client) MarkAllNotificationsRead(ctx context.Context) error {
	_, err := c.Put(ctx, "/api/v1/notifications/read-all", nil)
	return err
}

// DeleteNotification deletes a notification.
func (c *Client) DeleteNotification(ctx context.Context, id string) error {
	_, err := c.Delete(ctx, "/api/v1/notifications/"+id)
	return err
}
