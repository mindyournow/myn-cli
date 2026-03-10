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

type inboxLoadedMsg struct{ tasks []api.UnifiedTask }
type inboxErrMsg struct{ err error }

// InboxScreen shows unprocessed inbox items (tasks with no priority).
type InboxScreen struct {
	app    *app.App
	width  int
	height int
	cursor int

	tasks   []api.UnifiedTask
	loading bool
	err     error

	showDetail bool
	detail     *api.UnifiedTask
}

// NewInboxScreen creates the Inbox screen model.
func NewInboxScreen(application *app.App) InboxScreen {
	return InboxScreen{app: application, loading: true}
}

// Init implements tea.Model.
func (s InboxScreen) Init() tea.Cmd {
	return s.loadData()
}

func (s InboxScreen) loadData() tea.Cmd {
	return func() tea.Msg {
		if s.app == nil {
			return inboxLoadedMsg{}
		}
		ctx := context.Background()
		// Inbox = tasks with no priority
		tasks, err := s.app.Client.ListTasks(ctx, api.TaskListParams{Type: "TASK"})
		if err != nil {
			return inboxErrMsg{err}
		}
		// Filter to only null-priority (inbox) items
		var inbox []api.UnifiedTask
		for _, t := range tasks {
			if t.PriorityString() == "" && !t.IsCompleted {
				inbox = append(inbox, t)
			}
		}
		return inboxLoadedMsg{inbox}
	}
}

// Update implements tea.Model.
func (s InboxScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case inboxLoadedMsg:
		s.tasks = msg.tasks
		s.loading = false

	case inboxErrMsg:
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
		switch msg.String() {
		case "j", "down":
			if s.cursor < len(s.tasks)-1 {
				s.cursor++
			}
		case "k", "up":
			if s.cursor > 0 {
				s.cursor--
			}
		case "enter":
			if s.cursor < len(s.tasks) {
				t := s.tasks[s.cursor]
				s.detail = &t
				s.showDetail = true
			}
		case "g":
			s.loading = true
			return s, s.loadData()
		}
	}
	return s, nil
}

// View implements tea.Model.
func (s InboxScreen) View() string {
	if s.showDetail && s.detail != nil {
		return renderTaskDetail(s.detail, s.width)
	}

	title := titleStyle.Render(fmt.Sprintf("INBOX (%d items)", len(s.tasks)))

	if s.loading {
		return lipgloss.JoinVertical(lipgloss.Left, title, dimStyle.Render("  Loading..."))
	}
	if s.err != nil {
		return lipgloss.JoinVertical(lipgloss.Left, title, errorStyle.Render(fmt.Sprintf("  Error: %v", s.err)))
	}

	var rows []string
	rows = append(rows, title)

	if len(s.tasks) == 0 {
		rows = append(rows, successStyle.Render("  Inbox is empty!"))
	} else {
		for i, t := range s.tasks {
			rows = append(rows, renderInboxRow(t, i == s.cursor))
		}
	}

	rows = append(rows, "")
	rows = append(rows, dimStyle.Render("  j/k navigate  Enter=detail  c=critical  o=opportunity  d=delete  n=add"))

	return lipgloss.JoinVertical(lipgloss.Left, rows...)
}

func renderInboxRow(t api.UnifiedTask, selected bool) string {
	age := formatAge(t.CreatedDate)
	row := fmt.Sprintf("%-42s  %s", truncate(t.Title, 42), age)
	if selected {
		return selectedItemStyle.Render("► " + row)
	}
	return itemStyle.Render("  " + row)
}

func formatAge(dateStr string) string {
	if dateStr == "" {
		return ""
	}
	// Try a few common formats
	formats := []string{
		time.RFC3339,
		"2006-01-02T15:04:05",
		"2006-01-02",
	}
	var t time.Time
	var err error
	for _, f := range formats {
		t, err = time.Parse(f, dateStr)
		if err == nil {
			break
		}
	}
	if err != nil {
		return dateStr
	}

	diff := time.Since(t)
	hours := int(diff.Hours())
	switch {
	case hours < 1:
		return "added just now"
	case hours < 24:
		return fmt.Sprintf("added %dh ago", hours)
	case hours < 48:
		return "added yesterday"
	default:
		return "added " + t.Format("Jan 2")
	}
}
