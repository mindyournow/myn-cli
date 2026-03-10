package api

import (
	"context"
	"fmt"
)

// CalendarEvent represents a calendar event from the backend.
type CalendarEvent struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	StartTime   string `json:"startTime,omitempty"`
	EndTime     string `json:"endTime,omitempty"`
	Date        string `json:"date,omitempty"`
	IsAllDay    bool   `json:"isAllDay"`
	Location    string `json:"location,omitempty"`
	EventType   string `json:"eventType,omitempty"`
}

// ListCalendarEvents fetches calendar events for a date range.
func (c *Client) ListCalendarEvents(ctx context.Context, date string, days int) ([]CalendarEvent, error) {
	params := map[string]string{}
	if date != "" {
		params["date"] = date
	}
	if days > 0 {
		params["days"] = fmt.Sprintf("%d", days)
	}
	resp, err := c.Get(ctx, "/api/v2/calendar/events", params)
	if err != nil {
		return nil, err
	}
	var events []CalendarEvent
	if err := resp.DecodeJSON(&events); err != nil {
		return nil, fmt.Errorf("failed to parse calendar events: %w", err)
	}
	return events, nil
}

// CreateCalendarEventRequest is the body for creating a calendar event.
type CreateCalendarEventRequest struct {
	Title          string   `json:"title"`
	Description    string   `json:"description,omitempty"`
	StartTime      string   `json:"startTime,omitempty"`
	EndTime        string   `json:"endTime,omitempty"`
	Date           string   `json:"date,omitempty"`
	IsAllDay       bool     `json:"isAllDay,omitempty"`
	Location       string   `json:"location,omitempty"`
	Attendees      []string `json:"attendees,omitempty"`
	RecurrenceRule string   `json:"recurrenceRule,omitempty"`
}

// CreateCalendarEvent creates a new standalone calendar event.
func (c *Client) CreateCalendarEvent(ctx context.Context, req CreateCalendarEventRequest) (*CalendarEvent, error) {
	resp, err := c.Post(ctx, "/api/v2/calendar/standalone-events", req)
	if err != nil {
		return nil, err
	}
	var event CalendarEvent
	if err := resp.DecodeJSON(&event); err != nil {
		return nil, fmt.Errorf("failed to parse created event: %w", err)
	}
	return &event, nil
}

// DeleteCalendarEvent deletes a calendar event.
func (c *Client) DeleteCalendarEvent(ctx context.Context, id string) error {
	_, err := c.Delete(ctx, "/api/v2/calendar/events/"+id)
	return err
}

// DeclineCalendarMeeting declines a meeting invitation.
func (c *Client) DeclineCalendarMeeting(ctx context.Context, id string) error {
	_, err := c.Post(ctx, "/api/v2/calendar/meetings/"+id+"/decline", nil)
	return err
}

// SkipCalendarMeeting skips a recurring meeting occurrence.
func (c *Client) SkipCalendarMeeting(ctx context.Context, id string) error {
	_, err := c.Post(ctx, "/api/v2/calendar/meetings/"+id+"/skip", nil)
	return err
}
