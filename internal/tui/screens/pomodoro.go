package screens

import (
	"context"
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/mindyournow/myn-cli/internal/api"
	"github.com/mindyournow/myn-cli/internal/app"
	"github.com/mindyournow/myn-cli/internal/tui/components"
)

type pomodoroSettingsMsg struct{ settings api.PomodoroSettings }
type pomodoroTickMsg struct{}
type pomodoroErrMsg struct{ err error }

// PomodoroScreen shows the Pomodoro focus mode timer.
type PomodoroScreen struct {
	app       *app.App
	width     int
	height    int
	ring      components.PomodoroRing
	phase     string
	paused    bool
	remaining time.Duration
	total     time.Duration
	sessionID string
	loading   bool
	err       error
	toast     components.Toast
}

// NewPomodoroScreen creates the Pomodoro screen model.
func NewPomodoroScreen(application *app.App) PomodoroScreen {
	total := 25 * time.Minute
	return PomodoroScreen{
		app:       application,
		ring:      components.NewPomodoroRing(),
		phase:     "work",
		total:     total,
		remaining: total,
		toast:     components.NewToast(),
	}
}

// Init implements tea.Model.
func (s PomodoroScreen) Init() tea.Cmd {
	return tea.Batch(s.tickCmd(), s.loadSettings())
}

func (s PomodoroScreen) loadSettings() tea.Cmd {
	return func() tea.Msg {
		if s.app == nil {
			return pomodoroSettingsMsg{settings: api.PomodoroSettings{
				WorkDuration:       25,
				ShortBreakDuration: 5,
				LongBreakDuration:  15,
				SessionsBeforeLong: 4,
			}}
		}
		ctx := context.Background()
		settings, err := s.app.Client.GetPomodoroSettings(ctx)
		if err != nil {
			return pomodoroErrMsg{err}
		}
		return pomodoroSettingsMsg{settings: *settings}
	}
}

func (s PomodoroScreen) tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return pomodoroTickMsg{}
	})
}

// Update implements tea.Model.
func (s PomodoroScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case pomodoroSettingsMsg:
		if msg.settings.WorkDuration > 0 {
			s.total = time.Duration(msg.settings.WorkDuration) * time.Minute
			s.remaining = s.total
		}

	case pomodoroErrMsg:
		s.err = msg.err

	case pomodoroTickMsg:
		if !s.paused && s.remaining > 0 {
			s.remaining -= time.Second
		}
		s.updateRing()
		return s, s.tickCmd()

	case tea.WindowSizeMsg:
		s.width = msg.Width
		s.height = msg.Height
		s.ring.SetSize(s.width)

	case tea.KeyMsg:
		switch msg.String() {
		case "p":
			s.paused = true
			s.toast.Show("Paused", "info")
			return s, s.toast.Tick()
		case "r":
			s.paused = false
			s.toast.Show("Resumed", "success")
			return s, tea.Batch(s.tickCmd(), s.toast.Tick())
		case "s":
			s.remaining = 0
			s.sessionID = ""
			s.toast.Show("Session stopped", "info")
			return s, s.toast.Tick()
		case "n":
			s.phase = "work"
			s.remaining = s.total
			s.paused = false
			return s, tea.Batch(s.tickCmd(), s.loadSettings())
		case "g":
			return s, s.loadSettings()
		}
	}

	// Always pass through to toast so it can handle its internal tick.
	var toastCmd tea.Cmd
	s.toast, toastCmd = s.toast.Update(msg)
	s.updateRing()
	return s, toastCmd
}

func (s *PomodoroScreen) updateRing() {
	var fraction float64
	if s.total > 0 {
		fraction = 1.0 - float64(s.remaining)/float64(s.total)
	}
	s.ring.SetProgress(fraction)
	s.ring.SetPhase(s.phase)
}

// View implements tea.Model.
func (s PomodoroScreen) View() string {
	title := titleStyle.Render("POMODORO FOCUS")

	var phaseLabel string
	switch s.phase {
	case "short_break":
		phaseLabel = "Short Break"
	case "long_break":
		phaseLabel = "Long Break"
	default:
		phaseLabel = "Work Session"
	}

	mins := int(s.remaining.Minutes())
	secs := int(s.remaining.Seconds()) % 60
	timeStr := fmt.Sprintf("%02d:%02d", mins, secs)

	centeredRing := lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(s.ring.View())
	centeredPhase := lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Bold(true).Render(phaseLabel)
	centeredTime := lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Bold(true).
		Foreground(lipgloss.Color("#e2e8f0")).Render(timeStr)

	var rows []string
	rows = append(rows, title)
	rows = append(rows, "")
	rows = append(rows, centeredRing)
	rows = append(rows, "")
	rows = append(rows, centeredPhase)
	rows = append(rows, centeredTime)
	rows = append(rows, "")
	if s.paused {
		rows = append(rows, dimStyle.Render("  [PAUSED]"))
	}
	rows = append(rows, "")
	rows = append(rows, dimStyle.Render("  p=pause  r=resume  s=stop  n=new  g=refresh"))

	tv := s.toast.View()
	if tv != "" {
		rows = append(rows, tv)
	}

	return lipgloss.JoinVertical(lipgloss.Left, rows...)
}
