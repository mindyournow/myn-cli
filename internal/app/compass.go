package app

import (
	"context"
	"fmt"
	"strings"

	"github.com/mindyournow/myn-cli/internal/api"
)

// CompassShow shows the current compass briefing.
func (a *App) CompassShow(ctx context.Context) error {
	if err := a.ensureAuth(ctx); err != nil {
		return err
	}
	briefing, err := a.Client.GetCurrentCompass(ctx)
	if err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to get compass: %v", err))
		return err
	}
	if a.Formatter.JSON {
		return a.Formatter.Print(briefing)
	}
	return a.Formatter.PrintMarkdown(briefing.Summary)
}

// CompassGenerate generates a new compass briefing.
func (a *App) CompassGenerate(ctx context.Context, briefingType string, async bool) error {
	if err := a.ensureAuth(ctx); err != nil {
		return err
	}
	t := "ON_DEMAND"
	if briefingType != "" {
		t = strings.ToUpper(strings.ReplaceAll(briefingType, "-", "_"))
	}
	briefing, err := a.Client.GenerateCompass(ctx, api.GenerateCompassRequest{
		Type: t,
		Sync: !async,
	})
	if err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to generate compass: %v", err))
		return err
	}
	if a.Formatter.JSON {
		return a.Formatter.Print(briefing)
	}
	return a.Formatter.PrintMarkdown(briefing.Summary)
}

// CompassCorrect applies a correction to a compass briefing.
func (a *App) CompassCorrect(ctx context.Context, opts CompassCorrectOptions) error {
	if err := a.ensureAuth(ctx); err != nil {
		return err
	}
	req := api.CompassCorrectionRequest{
		SummaryID: opts.SummaryID,
		TaskID:    opts.TaskID,
		Decision:  opts.Decision,
		NewDate:   opts.NewDate,
		Reason:    opts.Reason,
	}
	briefing, err := a.Client.ApplyCompassCorrection(ctx, req)
	if err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to apply correction: %v", err))
		return err
	}
	if a.Formatter.JSON {
		return a.Formatter.Print(briefing)
	}
	return a.Formatter.Success("Correction applied.")
}

// CompassCorrectOptions are options for compass correct.
type CompassCorrectOptions struct {
	SummaryID string
	TaskID    string
	Decision  string
	NewDate   string
	Reason    string
}

// CompassComplete marks the current compass session as complete.
func (a *App) CompassComplete(ctx context.Context, summary string, decisions []string) error {
	if err := a.ensureAuth(ctx); err != nil {
		return err
	}
	briefing, err := a.Client.CompleteCompass(ctx, api.CompleteCompassRequest{
		Summary:   summary,
		Decisions: decisions,
	})
	if err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to complete compass: %v", err))
		return err
	}
	if a.Formatter.JSON {
		return a.Formatter.Print(briefing)
	}
	return a.Formatter.Success("Compass session completed.")
}

// CompassStatus shows the current compass session status.
func (a *App) CompassStatus(ctx context.Context) error {
	if err := a.ensureAuth(ctx); err != nil {
		return err
	}
	status, err := a.Client.GetCompassStatus(ctx)
	if err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to get compass status: %v", err))
		return err
	}
	if a.Formatter.JSON {
		return a.Formatter.Print(status)
	}
	activeStr := "inactive"
	if status.SessionActive {
		activeStr = "active"
	}
	_ = a.Formatter.Println(fmt.Sprintf("Session:             %s", activeStr))
	if status.BriefingID != "" {
		_ = a.Formatter.Println(fmt.Sprintf("Briefing ID:         %s", status.BriefingID))
	}
	_ = a.Formatter.Println(fmt.Sprintf("Pending corrections: %d", status.PendingCorrections))
	return nil
}

// CompassHistory shows past compass briefings.
func (a *App) CompassHistory(ctx context.Context, limit int) error {
	if err := a.ensureAuth(ctx); err != nil {
		return err
	}
	history, err := a.Client.GetCompassHistory(ctx, limit)
	if err != nil {
		_ = a.Formatter.Error(fmt.Sprintf("failed to get compass history: %v", err))
		return err
	}
	if a.Formatter.JSON {
		return a.Formatter.Print(history)
	}
	if len(history) == 0 {
		return a.Formatter.Println("No compass history.")
	}
	tbl := a.Formatter.NewTable("DATE", "TYPE", "STATUS")
	for _, h := range history {
		tbl.AddRow(h.CreatedAt, h.Type, h.Status)
	}
	tbl.Render()
	return nil
}
