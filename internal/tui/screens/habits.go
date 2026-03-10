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

type habitsLoadedMsg struct{ tasks []api.UnifiedTask }
type habitsErrMsg struct{ err error }

// HabitsScreen shows habits due today and completed habits.
type HabitsScreen struct {
	app    *app.App
	width  int
	height int
	cursor int

	tasks   []api.UnifiedTask
	loading bool
	err     error
}

// NewHabitsScreen creates the Habits screen model.
func NewHabitsScreen(application *app.App) HabitsScreen {
	return HabitsScreen{app: application, loading: true}
}

// Init implements tea.Model.
func (s HabitsScreen) Init() tea.Cmd {
	return s.loadData()
}

func (s HabitsScreen) loadData() tea.Cmd {
	return func() tea.Msg {
		if s.app == nil {
			return habitsLoadedMsg{}
		}
		ctx := context.Background()
		tasks, err := s.app.Client.ListTasks(ctx, api.TaskListParams{Type: "HABIT", Today: true})
		if err != nil {
			return habitsErrMsg{err}
		}
		return habitsLoadedMsg{tasks}
	}
}

// Update implements tea.Model.
func (s HabitsScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case habitsLoadedMsg:
		s.tasks = msg.tasks
		s.loading = false
		if s.cursor >= len(s.dueHabits()) {
			s.cursor = max(0, len(s.dueHabits())-1)
		}

	case habitsErrMsg:
		s.err = msg.err
		s.loading = false

	case tea.WindowSizeMsg:
		s.width = msg.Width
		s.height = msg.Height

	case tea.KeyMsg:
		due := s.dueHabits()
		switch msg.String() {
		case "j", "down":
			if s.cursor < len(due)-1 {
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

func (s HabitsScreen) dueHabits() []api.UnifiedTask {
	var due []api.UnifiedTask
	for _, t := range s.tasks {
		if !t.IsCompleted {
			due = append(due, t)
		}
	}
	return due
}

func (s HabitsScreen) completedHabits() []api.UnifiedTask {
	var done []api.UnifiedTask
	for _, t := range s.tasks {
		if t.IsCompleted {
			done = append(done, t)
		}
	}
	return done
}

// View implements tea.Model.
func (s HabitsScreen) View() string {
	title := titleStyle.Render("HABITS")

	if s.loading {
		return lipgloss.JoinVertical(lipgloss.Left, title, dimStyle.Render("  Loading..."))
	}
	if s.err != nil {
		return lipgloss.JoinVertical(lipgloss.Left, title, errorStyle.Render(fmt.Sprintf("  Error: %v", s.err)))
	}

	due := s.dueHabits()
	completed := s.completedHabits()

	var rows []string
	rows = append(rows, title)

	if len(due) > 0 {
		rows = append(rows, sectionStyle.Render("DUE TODAY"))
		for i, t := range due {
			rows = append(rows, renderHabitListRow(t, i == s.cursor, false))
		}
	}

	if len(completed) > 0 {
		rows = append(rows, sectionStyle.Render("COMPLETED"))
		for _, t := range completed {
			rows = append(rows, renderHabitListRow(t, false, true))
		}
	}

	if len(s.tasks) == 0 {
		rows = append(rows, dimStyle.Render("  No habits scheduled for today."))
	}

	rows = append(rows, "")
	rows = append(rows, dimStyle.Render("  j/k navigate  Enter=done  s=skip  g=refresh"))

	return lipgloss.JoinVertical(lipgloss.Left, rows...)
}

func renderHabitListRow(t api.UnifiedTask, selected bool, done bool) string {
	sym := "◆"
	if done {
		sym = "✓"
	}

	dur := ""
	if t.Duration > 0 {
		dur = util.FormatDuration(t.Duration * 60)
	}

	streak := ""
	switch v := t.StreakCount.(type) {
	case float64:
		if v > 0 {
			streak = fmt.Sprintf("  🔥 %d-day streak", int(v))
		}
	case int:
		if v > 0 {
			streak = fmt.Sprintf("  🔥 %d-day streak", v)
		}
	}

	row := fmt.Sprintf("%s %-36s  %-5s%s", sym, truncate(t.Title, 36), dur, streak)

	if done {
		return dimStyle.Render("    " + row)
	}
	if selected {
		return selectedItemStyle.Render("► " + row)
	}
	return itemStyle.Render("  " + row)
}
