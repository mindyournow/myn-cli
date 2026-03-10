package components

import (
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mindyournow/myn-cli/internal/api"
)

var (
	taskListStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#0f172a"))

	taskListEmptyStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#64748b")).
				Background(lipgloss.Color("#0f172a")).
				Padding(1, 2)
)

// TaskList is a filterable, sortable, cursor-navigable list of tasks.
type TaskList struct {
	tasks    []api.UnifiedTask // original full set
	visible  []api.UnifiedTask // filtered + sorted view
	filter   string
	sortBy   string // "priority" | "title" | "date"
	cursor   int
	Width    int
	Height   int
}

// NewTaskList creates a TaskList pre-loaded with tasks.
func NewTaskList(tasks []api.UnifiedTask) TaskList {
	tl := TaskList{
		tasks:  tasks,
		sortBy: "priority",
	}
	tl.rebuild()
	return tl
}

// SetFilter updates the filter string and rebuilds the visible list.
func (tl *TaskList) SetFilter(q string) {
	tl.filter = q
	tl.rebuild()
}

// SetSort updates the sort field and rebuilds the visible list.
// Accepted values: "priority", "title", "date".
func (tl *TaskList) SetSort(field string) {
	tl.sortBy = field
	tl.rebuild()
}

// SelectedTask returns the currently highlighted task, or nil if list is empty.
func (tl *TaskList) SelectedTask() *api.UnifiedTask {
	if len(tl.visible) == 0 {
		return nil
	}
	t := tl.visible[tl.cursor]
	return &t
}

// rebuild filters and sorts the task slice into tl.visible.
func (tl *TaskList) rebuild() {
	q := strings.ToLower(tl.filter)
	filtered := make([]api.UnifiedTask, 0, len(tl.tasks))
	for _, t := range tl.tasks {
		if q == "" || strings.Contains(strings.ToLower(t.Title), q) {
			filtered = append(filtered, t)
		}
	}

	priorityOrder := map[string]int{
		"CRITICAL":         0,
		"OPPORTUNITY_NOW":  1,
		"OVER_THE_HORIZON": 2,
		"PARKING_LOT":      3,
		"":                 4,
	}

	switch tl.sortBy {
	case "title":
		sort.SliceStable(filtered, func(i, j int) bool {
			return strings.ToLower(filtered[i].Title) < strings.ToLower(filtered[j].Title)
		})
	case "date":
		sort.SliceStable(filtered, func(i, j int) bool {
			di := filtered[i].DueDate
			if di == "" {
				di = filtered[i].StartDate
			}
			dj := filtered[j].DueDate
			if dj == "" {
				dj = filtered[j].StartDate
			}
			if di == dj {
				return filtered[i].Title < filtered[j].Title
			}
			if di == "" {
				return false
			}
			if dj == "" {
				return true
			}
			return di < dj
		})
	default: // "priority"
		sort.SliceStable(filtered, func(i, j int) bool {
			pi := priorityOrder[filtered[i].PriorityString()]
			pj := priorityOrder[filtered[j].PriorityString()]
			if pi != pj {
				return pi < pj
			}
			return strings.ToLower(filtered[i].Title) < strings.ToLower(filtered[j].Title)
		})
	}

	tl.visible = filtered
	if tl.cursor >= len(tl.visible) {
		tl.cursor = max(0, len(tl.visible)-1)
	}
}

// Update handles keyboard navigation.
func (tl TaskList) Update(msg tea.Msg) (TaskList, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "j", "down":
			if tl.cursor < len(tl.visible)-1 {
				tl.cursor++
			}
		case "k", "up":
			if tl.cursor > 0 {
				tl.cursor--
			}
		}
	}
	return tl, nil
}

// View renders the visible task rows.
func (tl TaskList) View() string {
	if len(tl.visible) == 0 {
		msg := "No tasks"
		if tl.filter != "" {
			msg = "No tasks match \"" + tl.filter + "\""
		}
		return taskListEmptyStyle.Render(msg)
	}

	rowWidth := tl.Width
	if rowWidth < 20 {
		rowWidth = 20
	}

	maxRows := tl.Height
	if maxRows <= 0 {
		maxRows = len(tl.visible)
	}

	// Scroll window: keep cursor visible.
	start := 0
	if tl.cursor >= maxRows {
		start = tl.cursor - maxRows + 1
	}
	end := start + maxRows
	if end > len(tl.visible) {
		end = len(tl.visible)
	}

	rows := make([]string, 0, end-start)
	for i := start; i < end; i++ {
		rows = append(rows, RenderTaskRow(tl.visible[i], i == tl.cursor, rowWidth))
	}

	return taskListStyle.Render(strings.Join(rows, "\n"))
}
