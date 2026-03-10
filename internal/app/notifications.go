package app

import (
	"context"
	"fmt"
)

// NotificationsList lists notifications.
func (a *App) NotificationsList(ctx context.Context, unread bool, limit int) error {
	if err := a.ensureAuth(ctx); err != nil {
		return err
	}
	notifications, err := a.Client.ListNotifications(ctx, 0, limit)
	if err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to list notifications: %v", err))
		return err
	}
	if unread {
		filtered := notifications[:0]
		for _, n := range notifications {
			if !n.IsRead {
				filtered = append(filtered, n)
			}
		}
		notifications = filtered
	}
	if a.Formatter.JSON {
		return a.Formatter.Print(notifications)
	}
	if len(notifications) == 0 {
		return a.Formatter.Println("No notifications.")
	}
	tbl := a.Formatter.NewTable("READ", "TITLE", "DATE")
	for _, n := range notifications {
		read := "✓"
		if !n.IsRead {
			read = "•"
		}
		tbl.AddRow(read, n.Title, n.CreatedAt)
	}
	tbl.Render()
	return nil
}

// NotificationsUnread shows unread notification count.
func (a *App) NotificationsUnread(ctx context.Context) error {
	if err := a.ensureAuth(ctx); err != nil {
		return err
	}
	count, err := a.Client.GetUnreadNotificationCount(ctx)
	if err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to get unread count: %v", err))
		return err
	}
	if a.Formatter.JSON {
		return a.Formatter.Print(map[string]int{"unread": count})
	}
	return a.Formatter.Println(fmt.Sprintf("%d unread notification(s).", count))
}

// NotificationsRead marks a notification as read.
func (a *App) NotificationsRead(ctx context.Context, id string) error {
	if err := a.ensureAuth(ctx); err != nil {
		return err
	}
	if err := a.Client.MarkNotificationRead(ctx, id); err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to mark notification: %v", err))
		return err
	}
	return a.Formatter.Success(fmt.Sprintf("Marked read: %s", id))
}

// NotificationsReadAll marks all notifications as read.
func (a *App) NotificationsReadAll(ctx context.Context) error {
	if err := a.ensureAuth(ctx); err != nil {
		return err
	}
	if err := a.Client.MarkAllNotificationsRead(ctx); err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to mark all notifications: %v", err))
		return err
	}
	return a.Formatter.Success("All notifications marked as read.")
}

// NotificationsDelete deletes a notification.
func (a *App) NotificationsDelete(ctx context.Context, id string) error {
	if err := a.ensureAuth(ctx); err != nil {
		return err
	}
	if err := a.Client.DeleteNotification(ctx, id); err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to delete notification: %v", err))
		return err
	}
	return a.Formatter.Success(fmt.Sprintf("Deleted notification %s.", id))
}
