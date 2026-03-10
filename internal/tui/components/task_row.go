package components

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/mindyournow/myn-cli/internal/api"
)

var (
	taskRowSelectedStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#ffffff")).
				Background(lipgloss.Color("#4c1d95"))

	taskRowNormalStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#e2e8f0")).
				Background(lipgloss.Color("#0f172a"))

	taskRowDimStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#64748b")).
			Background(lipgloss.Color("#0f172a"))

	taskRowSelectedDimStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#a78bfa")).
				Background(lipgloss.Color("#4c1d95"))
)

// RenderTaskRow renders a single task as a list row.
// selected highlights the row. width is the total available width.
func RenderTaskRow(t api.UnifiedTask, selected bool, width int) string {
	badge := PriorityBadge(t.PriorityString())

	// Build the metadata suffix (duration + project)
	meta := ""
	if t.Duration > 0 {
		meta += fmt.Sprintf(" %dm", t.Duration)
	}
	if t.ProjectName != "" {
		meta += "  " + t.ProjectName
	}

	// Reserve space: 1 badge + 1 space + meta + 2 padding
	reserved := 1 + 1 + len(meta) + 2
	titleWidth := width - reserved
	if titleWidth < 1 {
		titleWidth = 1
	}

	title := t.Title
	if len([]rune(title)) > titleWidth {
		title = string([]rune(title)[:titleWidth-1]) + "…"
	}

	if selected {
		row := badge + " " +
			taskRowSelectedStyle.Render(title) +
			taskRowSelectedDimStyle.Render(meta)
		return taskRowSelectedStyle.Width(width).Render(row)
	}

	row := badge + " " +
		taskRowNormalStyle.Render(title) +
		taskRowDimStyle.Render(meta)
	return taskRowNormalStyle.Width(width).Render(row)
}
