package api

import (
	"context"
	"fmt"
)

// CompassBriefing represents a compass/briefing from the backend.
type CompassBriefing struct {
	ID        string `json:"id"`
	Type      string `json:"type"`
	Summary   string `json:"summary"`
	CreatedAt string `json:"createdAt"`
	Status    string `json:"status"`
}

// CompassStatus holds the current compass session status.
type CompassStatus struct {
	SessionActive     bool   `json:"sessionActive"`
	BriefingID        string `json:"briefingId,omitempty"`
	PendingCorrections int   `json:"pendingCorrections,omitempty"`
	LastBriefing      string `json:"lastBriefing,omitempty"`
}

// GetCurrentCompass fetches the current compass briefing.
func (c *Client) GetCurrentCompass(ctx context.Context) (*CompassBriefing, error) {
	resp, err := c.Get(ctx, "/api/v2/compass/current", nil)
	if err != nil {
		return nil, err
	}
	var briefing CompassBriefing
	if err := resp.DecodeJSON(&briefing); err != nil {
		return nil, fmt.Errorf("failed to parse compass: %w", err)
	}
	return &briefing, nil
}

// GenerateCompassRequest is the body for generating a compass briefing.
type GenerateCompassRequest struct {
	Type string `json:"type"`
	Sync bool   `json:"sync"`
}

// GenerateCompass triggers a new compass briefing.
func (c *Client) GenerateCompass(ctx context.Context, req GenerateCompassRequest) (*CompassBriefing, error) {
	resp, err := c.Post(ctx, "/api/v2/compass/generate", req)
	if err != nil {
		return nil, err
	}
	var briefing CompassBriefing
	if err := resp.DecodeJSON(&briefing); err != nil {
		return nil, fmt.Errorf("failed to parse generated compass: %w", err)
	}
	return &briefing, nil
}

// CompassCorrectionRequest is the body for applying a compass correction.
type CompassCorrectionRequest struct {
	SummaryID string `json:"summaryId,omitempty"`
	TaskID    string `json:"taskId,omitempty"`
	Decision  string `json:"decision"`
	NewDate   string `json:"newDate,omitempty"`
	Reason    string `json:"reason,omitempty"`
}

// ApplyCompassCorrection applies a correction to a compass briefing.
func (c *Client) ApplyCompassCorrection(ctx context.Context, req CompassCorrectionRequest) (*CompassBriefing, error) {
	resp, err := c.Post(ctx, "/api/v2/compass/corrections/apply", req)
	if err != nil {
		return nil, err
	}
	var briefing CompassBriefing
	if err := resp.DecodeJSON(&briefing); err != nil {
		return nil, fmt.Errorf("failed to parse compass correction: %w", err)
	}
	return &briefing, nil
}

// CompleteCompassRequest is the body for completing a compass session.
type CompleteCompassRequest struct {
	Summary   string   `json:"summary,omitempty"`
	Decisions []string `json:"decisions,omitempty"`
}

// CompleteCompass marks the current compass session as complete.
func (c *Client) CompleteCompass(ctx context.Context, req CompleteCompassRequest) (*CompassBriefing, error) {
	resp, err := c.Post(ctx, "/api/v2/compass/complete", req)
	if err != nil {
		return nil, err
	}
	var briefing CompassBriefing
	if err := resp.DecodeJSON(&briefing); err != nil {
		return nil, fmt.Errorf("failed to parse completed compass: %w", err)
	}
	return &briefing, nil
}

// GetCompassStatus fetches the current compass session status.
func (c *Client) GetCompassStatus(ctx context.Context) (*CompassStatus, error) {
	resp, err := c.Get(ctx, "/api/v2/compass/status", nil)
	if err != nil {
		return nil, err
	}
	var status CompassStatus
	if err := resp.DecodeJSON(&status); err != nil {
		return nil, fmt.Errorf("failed to parse compass status: %w", err)
	}
	return &status, nil
}

// GetCompassHistory fetches past compass briefings.
func (c *Client) GetCompassHistory(ctx context.Context, limit int) ([]CompassBriefing, error) {
	params := map[string]string{}
	if limit > 0 {
		params["limit"] = fmt.Sprintf("%d", limit)
	}
	resp, err := c.Get(ctx, "/api/v2/compass/history", params)
	if err != nil {
		return nil, err
	}
	var history []CompassBriefing
	if err := resp.DecodeJSON(&history); err != nil {
		return nil, fmt.Errorf("failed to parse compass history: %w", err)
	}
	return history, nil
}
