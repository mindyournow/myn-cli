package screens

import (
	"context"
	"errors"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/mindyournow/myn-cli/internal/api"
	"github.com/mindyournow/myn-cli/internal/app"
	"github.com/mindyournow/myn-cli/internal/output"
	"github.com/mindyournow/myn-cli/internal/tui/components"
)

type compassLoadedMsg struct{ briefing *api.CompassBriefing }
type compassErrMsg struct{ err error }
type compassActionMsg struct{ msg string }
type compassGeneratingMsg struct{}

// CompassScreen shows the Compass briefing.
type CompassScreen struct {
	app          *app.App
	width        int
	height       int
	briefing     *api.CompassBriefing
	loading      bool
	err          error
	toast        components.Toast
	scrollOffset int
}

// NewCompassScreen creates the Compass screen model.
func NewCompassScreen(application *app.App) CompassScreen {
	return CompassScreen{
		app:     application,
		loading: true,
		toast:   components.NewToast(),
	}
}

// Init implements tea.Model.
func (s CompassScreen) Init() tea.Cmd {
	return s.loadData()
}

func (s CompassScreen) loadData() tea.Cmd {
	return func() tea.Msg {
		if s.app == nil {
			return compassLoadedMsg{nil}
		}
		ctx := context.Background()
		briefing, err := s.app.Client.GetCurrentCompass(ctx)
		if err != nil {
			return compassErrMsg{err}
		}
		return compassLoadedMsg{briefing}
	}
}

func (s CompassScreen) generateBriefing() tea.Cmd {
	return func() tea.Msg {
		if s.app == nil {
			return compassErrMsg{errors.New("no app")}
		}
		ctx := context.Background()
		briefing, err := s.app.Client.GenerateCompass(ctx, api.GenerateCompassRequest{
			Type: "DAILY",
			Sync: true,
		})
		if err != nil {
			return compassErrMsg{err}
		}
		return compassLoadedMsg{briefing}
	}
}

// Update implements tea.Model.
func (s CompassScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case compassLoadedMsg:
		s.briefing = msg.briefing
		s.loading = false

	case compassErrMsg:
		s.err = msg.err
		s.loading = false

	case compassGeneratingMsg:
		// handled via toast already

	case compassActionMsg:
		s.toast.Show(msg.msg, "success")
		return s, tea.Batch(s.loadData(), s.toast.Tick())

	case tea.WindowSizeMsg:
		s.width = msg.Width
		s.height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "g":
			s.loading = true
			return s, s.loadData()
		case "G":
			s.toast.Show("Generating briefing...", "info")
			return s, tea.Batch(s.generateBriefing(), s.toast.Tick())
		case "c":
			if s.briefing != nil {
				s.toast.Show("Applying correction...", "info")
				return s, tea.Batch(s.correctBriefing(), s.toast.Tick())
			}
			return s, nil
		case "enter":
			if s.briefing != nil {
				s.toast.Show("Completing session...", "info")
				return s, tea.Batch(s.completeBriefing(), s.toast.Tick())
			}
			return s, nil
		case "j", "down":
			s.scrollOffset++
		case "k", "up":
			if s.scrollOffset > 0 {
				s.scrollOffset--
			}
		}
	}

	// Always pass through to toast so it can handle its internal tick.
	var toastCmd tea.Cmd
	s.toast, toastCmd = s.toast.Update(msg)
	return s, toastCmd
}

func (s CompassScreen) correctBriefing() tea.Cmd {
	return func() tea.Msg {
		if s.app == nil || s.briefing == nil {
			return compassErrMsg{errors.New("no briefing to correct")}
		}
		ctx := context.Background()
		_, err := s.app.Client.ApplyCompassCorrection(ctx, api.CompassCorrectionRequest{
			SummaryID: s.briefing.ID,
			Decision:  "correct",
		})
		if err != nil {
			return compassErrMsg{err}
		}
		return compassActionMsg{msg: "Correction applied"}
	}
}

func (s CompassScreen) completeBriefing() tea.Cmd {
	return func() tea.Msg {
		if s.app == nil || s.briefing == nil {
			return compassErrMsg{errors.New("no briefing to complete")}
		}
		ctx := context.Background()
		_, err := s.app.Client.CompleteCompass(ctx, api.CompleteCompassRequest{
			Summary: s.briefing.Summary,
		})
		if err != nil {
			return compassErrMsg{err}
		}
		return compassActionMsg{msg: "Session completed"}
	}
}

// View implements tea.Model.
func (s CompassScreen) View() string {
	title := titleStyle.Render("COMPASS BRIEFING")

	var rows []string
	rows = append(rows, title)

	if s.loading {
		rows = append(rows, dimStyle.Render("  Loading..."))
	} else if s.err != nil {
		rows = append(rows, errorStyle.Render(fmt.Sprintf("  Error: %v", s.err)))
	} else if s.briefing == nil {
		rows = append(rows, dimStyle.Render("  No briefing available. Press G to generate."))
	} else {
		meta := dimStyle.Render(fmt.Sprintf("  %s · %s", s.briefing.Type, s.briefing.CreatedAt))
		rendered := output.RenderMarkdown(s.briefing.Summary, false)
		rows = append(rows, meta)
		rows = append(rows, "")
		rows = append(rows, rendered)
		rows = append(rows, "")
		rows = append(rows, dimStyle.Render("  g=refresh  G=generate  j/k=scroll  c=correct  Enter=complete"))
	}

	tv := s.toast.View()
	if tv != "" {
		rows = append(rows, tv)
	}

	return lipgloss.JoinVertical(lipgloss.Left, rows...)
}
