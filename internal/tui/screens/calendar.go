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

type calendarLoadedMsg struct {
	events []api.CalendarEvent
	date   time.Time
}
type calendarErrMsg struct{ err error }

// CalendarScreen shows calendar events for a selected day.
type CalendarScreen struct {
	app    *app.App
	width  int
	height int
	cursor int

	events  []api.CalendarEvent
	date    time.Time
	loading bool
	err     error
}

// NewCalendarScreen creates the Calendar screen model.
func NewCalendarScreen(application *app.App) CalendarScreen {
	return CalendarScreen{app: application, loading: true, date: time.Now()}
}

// Init implements tea.Model.
func (s CalendarScreen) Init() tea.Cmd {
	return s.loadData(s.date)
}

func (s CalendarScreen) loadData(date time.Time) tea.Cmd {
	return func() tea.Msg {
		if s.app == nil {
			return calendarLoadedMsg{date: date}
		}
		ctx := context.Background()
		dateStr := date.Format("2006-01-02")
		events, err := s.app.Client.ListCalendarEvents(ctx, dateStr, 1)
		if err != nil {
			return calendarErrMsg{err}
		}
		return calendarLoadedMsg{events: events, date: date}
	}
}

// Update implements tea.Model.
func (s CalendarScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case calendarLoadedMsg:
		s.events = msg.events
		s.date = msg.date
		s.loading = false
		s.cursor = 0

	case calendarErrMsg:
		s.err = msg.err
		s.loading = false

	case tea.WindowSizeMsg:
		s.width = msg.Width
		s.height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "j", "down":
			if s.cursor < len(s.events)-1 {
				s.cursor++
			}
		case "k", "up":
			if s.cursor > 0 {
				s.cursor--
			}
		case "]":
			next := s.date.AddDate(0, 0, 1)
			s.loading = true
			s.date = next
			return s, s.loadData(next)
		case "[":
			prev := s.date.AddDate(0, 0, -1)
			s.loading = true
			s.date = prev
			return s, s.loadData(prev)
		case "g":
			s.loading = true
			return s, s.loadData(s.date)
		}
	}
	return s, nil
}

// View implements tea.Model.
func (s CalendarScreen) View() string {
	dayStr := s.date.Format("Monday, January 2")
	title := titleStyle.Render(fmt.Sprintf("CALENDAR — %s", dayStr))

	if s.loading {
		return lipgloss.JoinVertical(lipgloss.Left, title, dimStyle.Render("  Loading..."))
	}
	if s.err != nil {
		return lipgloss.JoinVertical(lipgloss.Left, title, errorStyle.Render(fmt.Sprintf("  Error: %v", s.err)))
	}

	var rows []string
	rows = append(rows, title)

	if len(s.events) == 0 {
		rows = append(rows, dimStyle.Render("  No events scheduled."))
	} else {
		for i, ev := range s.events {
			rows = append(rows, renderCalEventRow(ev, i == s.cursor))
		}
	}

	rows = append(rows, "")
	rows = append(rows, dimStyle.Render("  j/k navigate  [=prev day  ]=next day  d=delete  g=refresh"))

	return lipgloss.JoinVertical(lipgloss.Left, rows...)
}

func renderCalEventRow(ev api.CalendarEvent, selected bool) string {
	startStr := ""
	endStr := ""
	if ev.StartTime != "" && len(ev.StartTime) >= 5 {
		startStr = ev.StartTime[:5]
	}
	if ev.EndTime != "" && len(ev.EndTime) >= 5 {
		endStr = ev.EndTime[:5]
	}

	timeRange := fmt.Sprintf("%s - %s", startStr, endStr)
	if startStr == "" && endStr == "" {
		timeRange = "all-day"
	}

	loc := ev.Location
	if loc == "" {
		loc = "—"
	}

	row := fmt.Sprintf("  %-15s  %-32s  %s",
		timeRange,
		truncate(ev.Title, 32),
		truncate(loc, 20),
	)

	if selected {
		return selectedItemStyle.Render("►" + row)
	}
	return itemStyle.Render(" " + row)
}
