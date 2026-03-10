package api

import (
	"context"
	"fmt"
)

// Project represents a project from the backend.
type Project struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	Status      string `json:"status,omitempty"`
	CreatedDate string `json:"createdDate,omitempty"`
}

// ListProjects fetches all projects.
func (c *Client) ListProjects(ctx context.Context) ([]Project, error) {
	resp, err := c.Get(ctx, "/api/project", nil)
	if err != nil {
		return nil, err
	}
	var projects []Project
	if err := resp.DecodeJSON(&projects); err != nil {
		return nil, fmt.Errorf("failed to parse projects: %w", err)
	}
	return projects, nil
}

// GetProject fetches a single project by ID.
func (c *Client) GetProject(ctx context.Context, id string) (*Project, error) {
	resp, err := c.Get(ctx, "/api/project/"+id, nil)
	if err != nil {
		return nil, err
	}
	var project Project
	if err := resp.DecodeJSON(&project); err != nil {
		return nil, fmt.Errorf("failed to parse project: %w", err)
	}
	return &project, nil
}

// CreateProjectRequest is the body for creating a project.
type CreateProjectRequest struct {
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
}

// CreateProject creates a new project.
func (c *Client) CreateProject(ctx context.Context, req CreateProjectRequest) (*Project, error) {
	resp, err := c.Post(ctx, "/api/project/create", req)
	if err != nil {
		return nil, err
	}
	var project Project
	if err := resp.DecodeJSON(&project); err != nil {
		return nil, fmt.Errorf("failed to parse created project: %w", err)
	}
	return &project, nil
}
