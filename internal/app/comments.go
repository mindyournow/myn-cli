package app

import (
	"context"
	"fmt"

	"github.com/mindyournow/myn-cli/internal/api"
)

// CommentList lists comments for a task.
func (a *App) CommentList(ctx context.Context, taskID string) error {
	if err := a.ensureAuth(ctx); err != nil {
		return err
	}
	comments, err := a.Client.ListTaskComments(ctx, taskID)
	if err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to list comments: %v", err))
		return err
	}
	if a.Formatter.JSON {
		return a.Formatter.Print(comments)
	}
	if len(comments) == 0 {
		return a.Formatter.Println("No comments.")
	}
	tbl := a.Formatter.NewTable("AUTHOR", "DATE", "COMMENT")
	for _, c := range comments {
		tbl.AddRow(c.Author, c.CreatedAt, c.Content)
	}
	tbl.Render()
	return nil
}

// CommentAdd adds a comment to a task.
func (a *App) CommentAdd(ctx context.Context, taskID, content string) error {
	if err := a.ensureAuth(ctx); err != nil {
		return err
	}
	comment, err := a.Client.AddTaskComment(ctx, taskID, content)
	if err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to add comment: %v", err))
		return err
	}
	if a.Formatter.JSON {
		return a.Formatter.Print(comment)
	}
	return a.Formatter.Success(fmt.Sprintf("Comment added: %s", comment.ID))
}

// CommentEdit edits a comment.
func (a *App) CommentEdit(ctx context.Context, taskID, commentID, content string) error {
	if err := a.ensureAuth(ctx); err != nil {
		return err
	}
	comment, err := a.Client.EditTaskComment(ctx, taskID, commentID, content)
	if err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to edit comment: %v", err))
		return err
	}
	if a.Formatter.JSON {
		return a.Formatter.Print(comment)
	}
	return a.Formatter.Success(fmt.Sprintf("Comment updated: %s", commentID))
}

// CommentDelete deletes a comment.
func (a *App) CommentDelete(ctx context.Context, taskID, commentID string) error {
	if err := a.ensureAuth(ctx); err != nil {
		return err
	}
	if err := a.Client.DeleteTaskComment(ctx, taskID, commentID); err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to delete comment: %v", err))
		return err
	}
	return a.Formatter.Success(fmt.Sprintf("Deleted comment %s.", commentID))
}

// TaskShare shares a task with a household member.
func (a *App) TaskShare(ctx context.Context, taskID, memberID, shareType, message string) error {
	if err := a.ensureAuth(ctx); err != nil {
		return err
	}
	req := api.ShareTaskRequest{
		MemberID: memberID,
		Type:     shareType,
		Message:  message,
	}
	if err := a.Client.ShareTask(ctx, taskID, req); err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to share task: %v", err))
		return err
	}
	return a.Formatter.Success(fmt.Sprintf("Task %s shared with %s.", taskID, memberID))
}

// TaskShareRespond accepts or declines a shared task.
func (a *App) TaskShareRespond(ctx context.Context, taskID, decision, note string) error {
	if err := a.ensureAuth(ctx); err != nil {
		return err
	}
	if err := a.Client.RespondToShare(ctx, taskID, decision, note); err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to respond to share: %v", err))
		return err
	}
	return a.Formatter.Success(fmt.Sprintf("Share %s: %s.", decision, taskID))
}

// TaskShareRevoke revokes a task share.
func (a *App) TaskShareRevoke(ctx context.Context, taskID, memberID string) error {
	if err := a.ensureAuth(ctx); err != nil {
		return err
	}
	if err := a.Client.RevokeTaskShare(ctx, taskID, memberID); err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to revoke share: %v", err))
		return err
	}
	return a.Formatter.Success(fmt.Sprintf("Share revoked for member %s.", memberID))
}

// SharedInbox shows tasks shared with the current user.
func (a *App) SharedInbox(ctx context.Context) error {
	if err := a.ensureAuth(ctx); err != nil {
		return err
	}
	tasks, err := a.Client.GetSharedInbox(ctx)
	if err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to get shared inbox: %v", err))
		return err
	}
	return a.printTaskList(ctx, tasks, "No shared tasks.")
}
