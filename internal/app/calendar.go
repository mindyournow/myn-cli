package app

import (
	"context"
	"fmt"
	"strings"

	"github.com/mindyournow/myn-cli/internal/api"
	"github.com/mindyournow/myn-cli/internal/util"
)

// CalendarList shows calendar events for a date.
func (a *App) CalendarList(ctx context.Context, date string, days int) error {
	if err := a.ensureAuth(ctx); err != nil {
		return err
	}
	if date == "" {
		date = nowDate()
	}
	events, err := a.Client.ListCalendarEvents(ctx, date, days)
	if err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to list calendar events: %v", err))
		return err
	}
	if a.Formatter.JSON {
		return a.Formatter.Print(events)
	}
	if len(events) == 0 {
		return a.Formatter.Println("No events.")
	}
	tbl := a.Formatter.NewTable("TIME", "TITLE", "LOCATION")
	for _, e := range events {
		timeStr := ""
		if e.IsAllDay {
			timeStr = "(all day)"
		} else if e.StartTime != "" {
			timeStr = e.StartTime
			if e.EndTime != "" {
				timeStr += " - " + e.EndTime
			}
		}
		tbl.AddRow(timeStr, e.Title, e.Location)
	}
	tbl.Render()
	return nil
}

// CalendarAdd adds a standalone calendar event.
func (a *App) CalendarAdd(ctx context.Context, title string, opts CalendarAddOptions) error {
	if err := a.ensureAuth(ctx); err != nil {
		return err
	}
	req := api.CreateCalendarEventRequest{Title: title}
	if opts.Date != "" {
		d, err := util.ParseDate(opts.Date)
		if err != nil {
			return err
		}
		req.Date = d.Format("2006-01-02")
	}
	req.IsAllDay = opts.AllDay
	req.StartTime = opts.Start
	req.EndTime = opts.End
	req.Location = opts.Location
	req.Description = opts.Description
	if opts.Attendees != "" {
		req.Attendees = strings.Split(opts.Attendees, ",")
	}
	if opts.Recurrence != "" {
		req.RecurrenceRule = util.ParseRecurrence(opts.Recurrence)
	}

	event, err := a.Client.CreateCalendarEvent(ctx, req)
	if err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to create event: %v", err))
		return err
	}
	if a.Formatter.JSON {
		return a.Formatter.Print(event)
	}
	return a.Formatter.Success(fmt.Sprintf("Created event: %s (%s)", event.Title, event.ID))
}

// CalendarAddOptions are options for calendar add.
type CalendarAddOptions struct {
	Date        string
	Start       string
	End         string
	AllDay      bool
	Location    string
	Attendees   string
	Description string
	Recurrence  string
}

// CalendarDelete deletes a calendar event.
func (a *App) CalendarDelete(ctx context.Context, id string) error {
	if err := a.ensureAuth(ctx); err != nil {
		return err
	}
	if err := a.Client.DeleteCalendarEvent(ctx, id); err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to delete event: %v", err))
		return err
	}
	return a.Formatter.Success(fmt.Sprintf("Deleted event %s.", id))
}

// CalendarDecline declines a meeting invitation.
func (a *App) CalendarDecline(ctx context.Context, id string) error {
	if err := a.ensureAuth(ctx); err != nil {
		return err
	}
	if err := a.Client.DeclineCalendarMeeting(ctx, id); err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to decline meeting: %v", err))
		return err
	}
	return a.Formatter.Success(fmt.Sprintf("Declined meeting %s.", id))
}

// CalendarSkip skips a recurring meeting occurrence.
func (a *App) CalendarSkip(ctx context.Context, id string) error {
	if err := a.ensureAuth(ctx); err != nil {
		return err
	}
	if err := a.Client.SkipCalendarMeeting(ctx, id); err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to skip meeting: %v", err))
		return err
	}
	return a.Formatter.Success(fmt.Sprintf("Skipped meeting %s.", id))
}
