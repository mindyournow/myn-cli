package api

import (
	"context"
	"fmt"
)

// Memory represents a memory item from the backend.
type Memory struct {
	ID        string   `json:"id"`
	Content   string   `json:"content"`
	Tags      []string `json:"tags,omitempty"`
	CreatedAt string   `json:"createdAt,omitempty"`
	UpdatedAt string   `json:"updatedAt,omitempty"`
}

// ListMemories fetches all memories for the current user.
func (c *Client) ListMemories(ctx context.Context) ([]Memory, error) {
	resp, err := c.Get(ctx, "/api/v1/customers/memories", nil)
	if err != nil {
		return nil, err
	}
	var memories []Memory
	if err := resp.DecodeJSON(&memories); err != nil {
		return nil, fmt.Errorf("failed to parse memories: %w", err)
	}
	return memories, nil
}

// GetMemory fetches a single memory by ID.
func (c *Client) GetMemory(ctx context.Context, id string) (*Memory, error) {
	resp, err := c.Get(ctx, "/api/v1/customers/memories/"+id, nil)
	if err != nil {
		return nil, err
	}
	var memory Memory
	if err := resp.DecodeJSON(&memory); err != nil {
		return nil, fmt.Errorf("failed to parse memory: %w", err)
	}
	return &memory, nil
}

// CreateMemoryRequest is the body for creating a memory.
type CreateMemoryRequest struct {
	Content string   `json:"content"`
	Tags    []string `json:"tags,omitempty"`
}

// AddMemory creates a new memory.
func (c *Client) AddMemory(ctx context.Context, req CreateMemoryRequest) (*Memory, error) {
	resp, err := c.Post(ctx, "/api/v1/customers/memories", req)
	if err != nil {
		return nil, err
	}
	var memory Memory
	if err := resp.DecodeJSON(&memory); err != nil {
		return nil, fmt.Errorf("failed to parse created memory: %w", err)
	}
	return &memory, nil
}

// UpdateMemory updates an existing memory.
func (c *Client) UpdateMemory(ctx context.Context, id string, req CreateMemoryRequest) (*Memory, error) {
	resp, err := c.Put(ctx, "/api/v1/customers/memories/"+id, req)
	if err != nil {
		return nil, err
	}
	var memory Memory
	if err := resp.DecodeJSON(&memory); err != nil {
		return nil, fmt.Errorf("failed to parse updated memory: %w", err)
	}
	return &memory, nil
}

// SearchMemories searches memories by query.
func (c *Client) SearchMemories(ctx context.Context, query string) ([]Memory, error) {
	resp, err := c.Get(ctx, "/api/v1/customers/memories/search",
		map[string]string{"q": query})
	if err != nil {
		return nil, err
	}
	var memories []Memory
	if err := resp.DecodeJSON(&memories); err != nil {
		return nil, fmt.Errorf("failed to parse memory search results: %w", err)
	}
	return memories, nil
}

// DeleteMemory deletes a memory by ID.
func (c *Client) DeleteMemory(ctx context.Context, id string) error {
	_, err := c.Delete(ctx, "/api/v1/customers/memories/"+id)
	return err
}

// DeleteAllMemories deletes all memories for the current user.
func (c *Client) DeleteAllMemories(ctx context.Context) error {
	_, err := c.Delete(ctx, "/api/v1/customers/memories")
	return err
}

// ExportMemories exports all memories.
func (c *Client) ExportMemories(ctx context.Context) ([]byte, error) {
	resp, err := c.Get(ctx, "/api/v1/customers/memories/export", nil)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}
