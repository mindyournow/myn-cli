package screens

import (
	"context"
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/mindyournow/myn-cli/internal/api"
	"github.com/mindyournow/myn-cli/internal/app"
)

type timersLoadedMsg struct{ timers []api.Timer }
type timersErrMsg struct{ err error }

// TimersScreen shows active and completed timers.
type TimersScreen struct {
	app    *app.App
	width  int
	height int
	cursor int

	timers  []api.Timer
	loading bool
	err     error
}

// NewTimersScreen creates the Timers screen model.
func NewTimersScreen(application *app.App) TimersScreen {
	return TimersScreen{app: application, loading: true}
}

// Init implements tea.Model.
func (s TimersScreen) Init() tea.Cmd {
	return s.loadData()
}

func (s TimersScreen) loadData() tea.Cmd {
	return func() tea.Msg {
		if s.app == nil {
			return timersLoadedMsg{}
		}
		ctx := context.Background()
		timers, err := s.app.Client.ListTimers(ctx, true)
		if err != nil {
			return timersErrMsg{err}
		}
		return timersLoadedMsg{timers}
	}
}

// Update implements tea.Model.
func (s TimersScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case timersLoadedMsg:
		s.timers = msg.timers
		s.loading = false
		active := s.activeTimers()
		if s.cursor >= len(active) {
			s.cursor = max(0, len(active)-1)
		}

	case timersErrMsg:
		s.err = msg.err
		s.loading = false

	case tea.WindowSizeMsg:
		s.width = msg.Width
		s.height = msg.Height

	case tea.KeyMsg:
		active := s.activeTimers()
		switch msg.String() {
		case "j", "down":
			if s.cursor < len(active)-1 {
				s.cursor++
			}
		case "k", "up":
			if s.cursor > 0 {
				s.cursor--
			}
		case "g":
			s.loading = true
			return s, s.loadData()
		}
	}
	return s, nil
}

func (s TimersScreen) activeTimers() []api.Timer {
	var active []api.Timer
	for _, t := range s.timers {
		if t.Status == "RUNNING" || t.Status == "PAUSED" {
			active = append(active, t)
		}
	}
	return active
}

func (s TimersScreen) completedTimers() []api.Timer {
	var done []api.Timer
	for _, t := range s.timers {
		if t.Status == "COMPLETED" || t.Status == "DONE" {
			done = append(done, t)
		}
	}
	return done
}

// View implements tea.Model.
func (s TimersScreen) View() string {
	title := titleStyle.Render("TIMERS")

	if s.loading {
		return lipgloss.JoinVertical(lipgloss.Left, title, dimStyle.Render("  Loading..."))
	}
	if s.err != nil {
		return lipgloss.JoinVertical(lipgloss.Left, title, errorStyle.Render(fmt.Sprintf("  Error: %v", s.err)))
	}

	active := s.activeTimers()
	completed := s.completedTimers()

	var rows []string
	rows = append(rows, title)

	if len(active) > 0 {
		rows = append(rows, sectionStyle.Render("ACTIVE"))
		for i, t := range active {
			rows = append(rows, renderTimerRow(t, i == s.cursor))
		}
	}

	if len(completed) > 0 {
		rows = append(rows, sectionStyle.Render("COMPLETED"))
		for _, t := range completed {
			rows = append(rows, renderTimerRow(t, false))
		}
	}

	if len(s.timers) == 0 {
		rows = append(rows, dimStyle.Render("  No timers found."))
	}

	rows = append(rows, "")
	rows = append(rows, dimStyle.Render("  j/k navigate  p=pause  r=resume  c=complete  x=dismiss  g=refresh"))

	return lipgloss.JoinVertical(lipgloss.Left, rows...)
}

func renderTimerRow(t api.Timer, selected bool) string {
	label := t.Label
	if label == "" {
		label = t.Type
	}

	remaining := formatTimerRemaining(t)
	status := fmt.Sprintf("[%s]", t.Status)

	row := fmt.Sprintf("%-30s  %-18s  %s",
		truncate(label, 30),
		remaining,
		status,
	)

	if t.Status == "COMPLETED" || t.Status == "DONE" {
		return dimStyle.Render("    " + row)
	}
	if selected {
		return selectedItemStyle.Render("► " + row)
	}
	return itemStyle.Render("  " + row)
}

func formatTimerRemaining(t api.Timer) string {
	if t.Status == "COMPLETED" || t.Status == "DONE" {
		return "Done"
	}

	if t.AlarmTime != "" {
		alarmTime, err := time.Parse(time.RFC3339, t.AlarmTime)
		if err == nil {
			remaining := time.Until(alarmTime)
			if remaining > 0 {
				m := int(remaining.Minutes())
				s := int(remaining.Seconds()) % 60
				return fmt.Sprintf("%02d:%02d remaining", m, s)
			}
			return "elapsed"
		}
	}

	if t.Duration > 0 {
		m := t.Duration / 60
		s := t.Duration % 60
		return fmt.Sprintf("%02d:%02d remaining", m, s)
	}

	return "running"
}
