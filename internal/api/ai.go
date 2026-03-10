package api

import (
	"context"
	"fmt"
)

// AIConversation represents an AI chat conversation.
type AIConversation struct {
	ID                 string `json:"conversationId"`
	Title              string `json:"title,omitempty"`
	CreatedAt          string `json:"createdAt,omitempty"`
	UpdatedAt          string `json:"updatedAt,omitempty"`
	MessageCount       int    `json:"messageCount,omitempty"`
	IsVoice            bool   `json:"isVoice"`
	LastMessagePreview string `json:"lastMessagePreview,omitempty"`
	IsArchived         bool   `json:"isArchived"`
}

// AIChatRequest is the body for an AI chat request.
type AIChatRequest struct {
	CurrentMessage          string  `json:"currentMessage"`
	ConversationID          string  `json:"conversationId,omitempty"`
	IsMobile                bool    `json:"isMobile"`
	IsVoice                 bool    `json:"isVoice"`
	VoiceSessionID          *string `json:"voiceSessionId"`
	Context                 *string `json:"context"`
	AdditionalSystemContext *string `json:"additionalSystemContext"`
	Regenerate              bool    `json:"regenerate"`
}

// ListAIConversations fetches AI chat conversations.
func (c *Client) ListAIConversations(ctx context.Context) ([]AIConversation, error) {
	resp, err := c.Get(ctx, "/api/v1/ai/conversations", nil)
	if err != nil {
		return nil, err
	}
	var conversations []AIConversation
	if err := resp.DecodeJSON(&conversations); err != nil {
		return nil, fmt.Errorf("failed to parse conversations: %w", err)
	}
	return conversations, nil
}

// GetAIConversation fetches a single conversation with messages.
func (c *Client) GetAIConversation(ctx context.Context, id string) (map[string]interface{}, error) {
	resp, err := c.Get(ctx, "/api/v1/ai/conversations/"+id+"/messages", nil)
	if err != nil {
		return nil, err
	}
	var conv map[string]interface{}
	if err := resp.DecodeJSON(&conv); err != nil {
		return nil, fmt.Errorf("failed to parse conversation: %w", err)
	}
	return conv, nil
}

// DeleteAIConversation deletes a conversation.
func (c *Client) DeleteAIConversation(ctx context.Context, id string) error {
	_, err := c.Delete(ctx, "/api/v1/ai/conversations/"+id)
	return err
}

// AIChatStream sends a chat message and streams the response via SSE.
func (c *Client) AIChatStream(ctx context.Context, req AIChatRequest, handler SSEHandler) error {
	return c.StreamPost(ctx, "/api/ai/chat/stream", req, handler)
}

// CreateAIConversation creates a new conversation.
func (c *Client) CreateAIConversation(ctx context.Context, title string) (*AIConversation, error) {
	body := map[string]string{}
	if title != "" {
		body["title"] = title
	}
	resp, err := c.Post(ctx, "/api/v1/ai/conversations", body)
	if err != nil {
		return nil, err
	}
	var conv AIConversation
	if err := resp.DecodeJSON(&conv); err != nil {
		return nil, fmt.Errorf("failed to parse conversation: %w", err)
	}
	return &conv, nil
}

// ArchiveAIConversation archives a conversation.
func (c *Client) ArchiveAIConversation(ctx context.Context, id string) error {
	_, err := c.Patch(ctx, "/api/v1/ai/conversations/"+id+"/status",
		map[string]bool{"isArchived": true})
	return err
}

// FavoriteAIConversation favorites a conversation.
func (c *Client) FavoriteAIConversation(ctx context.Context, id string) error {
	_, err := c.Patch(ctx, "/api/v1/ai/conversations/"+id+"/status",
		map[string]bool{"favorited": true})
	return err
}

// SearchAIConversations searches conversations by query.
func (c *Client) SearchAIConversations(ctx context.Context, query string) ([]AIConversation, error) {
	resp, err := c.Get(ctx, "/api/v1/ai/conversations/search",
		map[string]string{"q": query})
	if err != nil {
		return nil, err
	}
	var convs []AIConversation
	if err := resp.DecodeJSON(&convs); err != nil {
		return nil, fmt.Errorf("failed to parse conversations: %w", err)
	}
	return convs, nil
}

// GetAIConversationStats fetches AI conversation usage stats.
func (c *Client) GetAIConversationStats(ctx context.Context) (map[string]interface{}, error) {
	resp, err := c.Get(ctx, "/api/v1/ai/conversations/stats", nil)
	if err != nil {
		return nil, err
	}
	var stats map[string]interface{}
	if err := resp.DecodeJSON(&stats); err != nil {
		return nil, fmt.Errorf("failed to parse stats: %w", err)
	}
	return stats, nil
}

// ContinueAIConversation continues a conversation with a new message.
func (c *Client) ContinueAIConversation(ctx context.Context, id string, req AIChatRequest, handler SSEHandler) error {
	return c.StreamPost(ctx, "/api/v1/ai/conversations/"+id+"/continue", req, handler)
}
