package screens

import (
	"context"
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/mindyournow/myn-cli/internal/api"
	"github.com/mindyournow/myn-cli/internal/app"
	"github.com/mindyournow/myn-cli/internal/util"
)

// nowDataMsg carries the loaded Now screen data.
type nowDataMsg struct {
	tasks  []api.UnifiedTask
	events []api.CalendarEvent
}

// nowErrMsg carries a load error for the Now screen.
type nowErrMsg struct{ err error }

// NowScreen shows the current focus tasks (Critical Now + calendar + habits due today).
type NowScreen struct {
	app    *app.App
	width  int
	height int
	cursor int

	tasks  []api.UnifiedTask
	events []api.CalendarEvent

	loading bool
	err     error

	showDetail bool
	detail     *api.UnifiedTask
}

// NewNowScreen creates the Now screen model.
func NewNowScreen(application *app.App) NowScreen {
	return NowScreen{app: application, loading: true}
}

// Init implements tea.Model.
func (s NowScreen) Init() tea.Cmd {
	return s.loadData()
}

func (s NowScreen) loadData() tea.Cmd {
	return func() tea.Msg {
		if s.app == nil {
			return nowDataMsg{}
		}
		ctx := context.Background()
		tasks, err := s.app.Client.ListTasks(ctx, api.TaskListParams{Today: true})
		if err != nil {
			return nowErrMsg{err}
		}
		today := time.Now().Format("2006-01-02")
		events, err := s.app.Client.ListCalendarEvents(ctx, today, 1)
		if err != nil {
			// non-fatal — show tasks without events
			events = nil
		}
		return nowDataMsg{tasks: tasks, events: events}
	}
}

// Update implements tea.Model.
func (s NowScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case nowDataMsg:
		s.tasks = msg.tasks
		s.events = msg.events
		s.loading = false

	case nowErrMsg:
		s.err = msg.err
		s.loading = false

	case tea.WindowSizeMsg:
		s.width = msg.Width
		s.height = msg.Height

	case tea.KeyMsg:
		if s.showDetail {
			if msg.String() == "esc" || msg.String() == "q" {
				s.showDetail = false
			}
			return s, nil
		}
		allItems := s.allNavigableItems()
		switch msg.String() {
		case "j", "down":
			if s.cursor < len(allItems)-1 {
				s.cursor++
			}
		case "k", "up":
			if s.cursor > 0 {
				s.cursor--
			}
		case "enter":
			if s.cursor < len(allItems) {
				s.detail = allItems[s.cursor]
				if s.detail != nil {
					s.showDetail = true
				}
			}
		case "g":
			s.loading = true
			return s, s.loadData()
		}
	}
	return s, nil
}

// allNavigableItems returns tasks as navigable pointers (nil for non-task items).
func (s NowScreen) allNavigableItems() []*api.UnifiedTask {
	var items []*api.UnifiedTask
	for i := range s.tasks {
		t := &s.tasks[i]
		items = append(items, t)
	}
	return items
}

// View implements tea.Model.
func (s NowScreen) View() string {
	if s.showDetail && s.detail != nil {
		return renderTaskDetail(s.detail, s.width)
	}

	now := time.Now()
	dayStr := now.Format("Monday, January 2")
	title := titleStyle.Render(fmt.Sprintf("NOW — %s", dayStr))

	if s.loading {
		return lipgloss.JoinVertical(lipgloss.Left, title, dimStyle.Render("  Loading..."))
	}
	if s.err != nil {
		return lipgloss.JoinVertical(lipgloss.Left, title, formatError(s.err))
	}

	var sections []string
	sections = append(sections, title)

	// Critical tasks
	var critical []api.UnifiedTask
	var habitsDue []api.UnifiedTask
	for _, t := range s.tasks {
		if t.IsCompleted {
			continue
		}
		switch t.TaskType {
		case "HABIT":
			habitsDue = append(habitsDue, t)
		default:
			p := t.PriorityString()
			if p == "CRITICAL" {
				critical = append(critical, t)
			}
		}
	}

	globalIdx := 0

	if len(critical) > 0 {
		sections = append(sections, sectionStyle.Render("CRITICAL NOW"))
		for _, t := range critical {
			sections = append(sections, renderTaskRow(t, globalIdx == s.cursor))
			globalIdx++
		}
	}

	if len(s.events) > 0 {
		sections = append(sections, sectionStyle.Render("UPCOMING TODAY"))
		for _, ev := range s.events {
			sections = append(sections, renderEventRow(ev))
		}
	}

	if len(habitsDue) > 0 {
		sections = append(sections, sectionStyle.Render("HABITS DUE"))
		for _, t := range habitsDue {
			sections = append(sections, renderHabitRow(t, globalIdx == s.cursor))
			globalIdx++
		}
	}

	if len(critical) == 0 && len(s.events) == 0 && len(habitsDue) == 0 {
		sections = append(sections, dimStyle.Render("  Nothing scheduled for today."))
	}

	sections = append(sections, "")
	sections = append(sections, dimStyle.Render("  j/k navigate  Enter=detail  d=done  g=refresh  n=new task"))

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

func renderTaskRow(t api.UnifiedTask, selected bool) string {
	var sym string
	switch t.PriorityString() {
	case "CRITICAL":
		sym = "●"
	case "OPPORTUNITY_NOW":
		sym = "○"
	case "OVER_THE_HORIZON":
		sym = "◌"
	default:
		sym = " "
	}

	dur := ""
	if t.Duration > 0 {
		dur = util.FormatDuration(t.Duration * 60)
	}

	proj := t.ProjectName
	if proj == "" {
		proj = "—"
	}

	row := fmt.Sprintf("%s %-36s  %-5s  %s", sym, truncate(t.Title, 36), dur, proj)

	if selected {
		return selectedItemStyle.Render("► " + row)
	}
	return itemStyle.Render("  " + row)
}

func renderEventRow(ev api.CalendarEvent) string {
	timeStr := ""
	if ev.StartTime != "" {
		timeStr = ev.StartTime
		if len(timeStr) > 5 {
			timeStr = timeStr[:5]
		}
	}
	loc := ev.Location
	if loc == "" {
		loc = "—"
	}
	row := fmt.Sprintf("  %5s  %-32s  %s", timeStr, truncate(ev.Title, 32), loc)
	return itemStyle.Render(row)
}

func renderHabitRow(t api.UnifiedTask, selected bool) string {
	dur := ""
	if t.Duration > 0 {
		dur = util.FormatDuration(t.Duration * 60)
	}

	streak := ""
	switch v := t.StreakCount.(type) {
	case float64:
		if v > 0 {
			streak = fmt.Sprintf("  🔥 %d", int(v))
		}
	case int:
		if v > 0 {
			streak = fmt.Sprintf("  🔥 %d", v)
		}
	}

	row := fmt.Sprintf("◆ %-36s  %-5s%s", truncate(t.Title, 36), dur, streak)

	if selected {
		return selectedItemStyle.Render("► " + row)
	}
	return itemStyle.Render("  " + row)
}

func truncate(s string, max int) string {
	runes := []rune(s)
	if len(runes) <= max {
		return s
	}
	return string(runes[:max-1]) + "…"
}
