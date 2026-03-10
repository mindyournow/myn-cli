package app

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mindyournow/myn-cli/internal/api"
)

// AIConversationList lists AI chat conversations.
func (a *App) AIConversationList(ctx context.Context) error {
	if err := a.ensureAuth(ctx); err != nil {
		return err
	}
	convs, err := a.Client.ListAIConversations(ctx)
	if err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to list conversations: %v", err))
		return err
	}
	if a.Formatter.JSON {
		return a.Formatter.Print(convs)
	}
	if len(convs) == 0 {
		return a.Formatter.Println("No AI conversations.")
	}
	tbl := a.Formatter.NewTable("ID", "TITLE", "MESSAGES", "UPDATED")
	for _, c := range convs {
		tbl.AddRow(c.ID, c.Title, fmt.Sprintf("%d", c.MessageCount), c.UpdatedAt)
	}
	tbl.Render()
	return nil
}

// AIConversationDelete deletes an AI conversation.
func (a *App) AIConversationDelete(ctx context.Context, id string) error {
	if err := a.ensureAuth(ctx); err != nil {
		return err
	}
	if err := a.Client.DeleteAIConversation(ctx, id); err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to delete conversation: %v", err))
		return err
	}
	return a.Formatter.Success(fmt.Sprintf("Deleted conversation %s.", id))
}

// AIConversationCreate creates a new AI conversation.
func (a *App) AIConversationCreate(ctx context.Context, title string) error {
	if err := a.ensureAuth(ctx); err != nil {
		return err
	}
	conv, err := a.Client.CreateAIConversation(ctx, title)
	if err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to create conversation: %v", err))
		return err
	}
	if a.Formatter.JSON {
		return a.Formatter.Print(conv)
	}
	return a.Formatter.Success(fmt.Sprintf("Created conversation: %s (%s)", conv.Title, conv.ID))
}

// AIConversationArchive archives a conversation.
func (a *App) AIConversationArchive(ctx context.Context, id string) error {
	if err := a.ensureAuth(ctx); err != nil {
		return err
	}
	if err := a.Client.ArchiveAIConversation(ctx, id); err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to archive conversation: %v", err))
		return err
	}
	return a.Formatter.Success(fmt.Sprintf("Conversation archived: %s", id))
}

// AIConversationSearch searches conversations.
func (a *App) AIConversationSearch(ctx context.Context, query string) error {
	if err := a.ensureAuth(ctx); err != nil {
		return err
	}
	convs, err := a.Client.SearchAIConversations(ctx, query)
	if err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to search conversations: %v", err))
		return err
	}
	if a.Formatter.JSON {
		return a.Formatter.Print(convs)
	}
	if len(convs) == 0 {
		return a.Formatter.Println("No conversations found.")
	}
	tbl := a.Formatter.NewTable("ID", "TITLE", "MESSAGES")
	for _, c := range convs {
		tbl.AddRow(c.ID, c.Title, fmt.Sprintf("%d", c.MessageCount))
	}
	tbl.Render()
	return nil
}

// AIConversationStats shows conversation stats.
func (a *App) AIConversationStats(ctx context.Context) error {
	if err := a.ensureAuth(ctx); err != nil {
		return err
	}
	stats, err := a.Client.GetAIConversationStats(ctx)
	if err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to get stats: %v", err))
		return err
	}
	if a.Formatter.JSON {
		return a.Formatter.Print(stats)
	}
	b, err := json.MarshalIndent(stats, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to format output: %w", err)
	}
	return a.Formatter.Println(string(b))
}

// HabitCalculateSmartTime calculates optimal reminder time for a habit.
func (a *App) HabitCalculateSmartTime(ctx context.Context, habitID string) error {
	if err := a.ensureAuth(ctx); err != nil {
		return err
	}
	result, err := a.Client.CalculateSmartTime(ctx, habitID)
	if err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to calculate smart time: %v", err))
		return err
	}
	if a.Formatter.JSON {
		return a.Formatter.Print(result)
	}
	b, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to format output: %w", err)
	}
	return a.Formatter.Println(string(b))
}

// AIChat sends a message and streams the AI response.
func (a *App) AIChat(ctx context.Context, message, conversationID string) error {
	if err := a.ensureAuth(ctx); err != nil {
		return err
	}
	req := api.AIChatRequest{
		CurrentMessage: message,
		ConversationID: conversationID,
	}
	err := a.Client.AIChatStream(ctx, req, func(event api.SSEEvent) error {
		if event.Data == "[DONE]" {
			return nil
		}
		_, werr := fmt.Fprint(a.Formatter.Writer(), event.Data)
		return werr
	})
	if err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("chat error: %v", err))
		return err
	}
	// final newline after streamed content
	return a.Formatter.Println("")
}
