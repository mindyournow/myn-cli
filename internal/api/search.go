package api

import (
	"context"
	"fmt"
	"strings"
)

// SearchResult represents a single item from a unified search.
type SearchResult struct {
	ID       string      `json:"id"`
	Type     string      `json:"type"`
	Title    string      `json:"title"`
	Snippet  string      `json:"snippet,omitempty"`
	Priority interface{} `json:"priority,omitempty"`
	Date     string      `json:"date,omitempty"`
}

// SearchParams are query parameters for unified search.
type SearchParams struct {
	Query           string
	Types           []string
	Priorities      []string
	IncludeArchived bool
	Limit           int
	Offset          int
}

// Search performs a unified search across all entity types.
func (c *Client) Search(ctx context.Context, p SearchParams) ([]SearchResult, error) {
	params := map[string]string{
		"q": p.Query,
	}
	if len(p.Types) > 0 {
		params["types"] = strings.Join(p.Types, ",")
	}
	if len(p.Priorities) > 0 {
		params["priorities"] = strings.Join(p.Priorities, ",")
	}
	if p.IncludeArchived {
		params["includeArchived"] = "true"
	}
	if p.Limit > 0 {
		params["limit"] = fmt.Sprintf("%d", p.Limit)
	}
	if p.Offset > 0 {
		params["offset"] = fmt.Sprintf("%d", p.Offset)
	}

	resp, err := c.Get(ctx, "/api/v2/search", params)
	if err != nil {
		return nil, err
	}
	var results []SearchResult
	if err := resp.DecodeJSON(&results); err != nil {
		return nil, fmt.Errorf("failed to parse search results: %w", err)
	}
	return results, nil
}
