package screens

import (
	"context"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/mindyournow/myn-cli/internal/api"
	"github.com/mindyournow/myn-cli/internal/app"
	"github.com/mindyournow/myn-cli/internal/util"
)

type tasksLoadedMsg struct{ tasks []api.UnifiedTask }
type tasksErrMsg struct{ err error }

// TasksScreen shows all tasks grouped by priority zone.
type TasksScreen struct {
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

// NewTasksScreen creates the Tasks screen model.
func NewTasksScreen(application *app.App) TasksScreen {
	return TasksScreen{app: application, loading: true}
}

// Init implements tea.Model.
func (s TasksScreen) Init() tea.Cmd {
	return s.loadData()
}

func (s TasksScreen) loadData() tea.Cmd {
	return func() tea.Msg {
		if s.app == nil {
			return tasksLoadedMsg{}
		}
		ctx := context.Background()
		tasks, err := s.app.Client.ListTasks(ctx, api.TaskListParams{})
		if err != nil {
			return tasksErrMsg{err}
		}
		return tasksLoadedMsg{tasks}
	}
}

// Update implements tea.Model.
func (s TasksScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tasksLoadedMsg:
		s.tasks = msg.tasks
		s.loading = false
		nav := s.navigableItems()
		if s.cursor >= len(nav) {
			s.cursor = max(0, len(nav)-1)
		}

	case tasksErrMsg:
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
		navigable := s.navigableItems()
		switch msg.String() {
		case "j", "down":
			if s.cursor < len(navigable)-1 {
				s.cursor++
			}
		case "k", "up":
			if s.cursor > 0 {
				s.cursor--
			}
		case "enter":
			if s.cursor < len(navigable) {
				t := navigable[s.cursor]
				s.detail = t
				s.showDetail = true
			}
		case "g":
			s.loading = true
			return s, s.loadData()
		}
	}
	return s, nil
}

// navigableItems returns the non-completed tasks in display order.
func (s TasksScreen) navigableItems() []*api.UnifiedTask {
	priority := []string{"CRITICAL", "OPPORTUNITY_NOW", "OVER_THE_HORIZON", "PARKING_LOT", ""}
	var result []*api.UnifiedTask
	for _, p := range priority {
		for i := range s.tasks {
			t := &s.tasks[i]
			if t.IsCompleted {
				continue
			}
			if t.PriorityString() == p {
				result = append(result, t)
			}
		}
	}
	return result
}

// View implements tea.Model.
func (s TasksScreen) View() string {
	if s.showDetail && s.detail != nil {
		return renderTaskDetail(s.detail, s.width)
	}

	title := titleStyle.Render("TASKS")

	if s.loading {
		return lipgloss.JoinVertical(lipgloss.Left, title, dimStyle.Render("  Loading..."))
	}
	if s.err != nil {
		return lipgloss.JoinVertical(lipgloss.Left, title, formatError(s.err))
	}

	type zone struct {
		label    string
		priority string
		sym      string
	}
	zones := []zone{
		{"CRITICAL NOW (●)", "CRITICAL", "●"},
		{"OPPORTUNITY NOW (○)", "OPPORTUNITY_NOW", "○"},
		{"OVER THE HORIZON (◌)", "OVER_THE_HORIZON", "◌"},
		{"PARKING LOT (⊡)", "PARKING_LOT", "⊡"},
		{"INBOX", "", " "},
	}

	var rows []string
	rows = append(rows, title)

	globalIdx := 0
	for _, z := range zones {
		var group []api.UnifiedTask
		for _, t := range s.tasks {
			if t.IsCompleted {
				continue
			}
			if t.PriorityString() == z.priority {
				group = append(group, t)
			}
		}
		if len(group) == 0 {
			continue
		}
		rows = append(rows, sectionStyle.Render(z.label))
		for _, t := range group {
			rows = append(rows, renderFullTaskRow(t, z.sym, globalIdx == s.cursor))
			globalIdx++
		}
	}

	if globalIdx == 0 {
		rows = append(rows, dimStyle.Render("  No tasks found."))
	}

	rows = append(rows, "")
	rows = append(rows, dimStyle.Render("  j/k navigate  Enter=detail  d=done  s=snooze  n=new  g=refresh"))

	return lipgloss.JoinVertical(lipgloss.Left, rows...)
}

func renderFullTaskRow(t api.UnifiedTask, sym string, selected bool) string {
	dur := ""
	if t.Duration > 0 {
		dur = util.FormatDuration(t.Duration * 60)
	}
	proj := t.ProjectName
	if proj == "" {
		proj = "—"
	}
	row := fmt.Sprintf("%s %-38s  %-5s  %s", sym, truncate(t.Title, 38), dur, truncate(proj, 20))
	if selected {
		return selectedItemStyle.Render("► " + row)
	}
	return itemStyle.Render("  " + row)
}
