package screens

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/mindyournow/myn-cli/internal/api"
	"github.com/mindyournow/myn-cli/internal/util"
)

var (
	titleStyle        = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#7c3aed")).Padding(1, 2)
	sectionStyle      = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#64748b")).MarginLeft(2).MarginTop(1)
	itemStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("#e2e8f0")).PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#7c3aed")).PaddingLeft(2).Bold(true)
	dimStyle          = lipgloss.NewStyle().Foreground(lipgloss.Color("#475569"))
	successStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("#22c55e"))
	errorStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("#ef4444"))
)

// renderTaskDetail renders a task detail overlay panel.
func renderTaskDetail(t *api.UnifiedTask, width int) string {
	if t == nil {
		return ""
	}

	boxWidth := 50
	if width > 0 && width < 54 {
		boxWidth = width - 4
	}
	if boxWidth < 20 {
		boxWidth = 20
	}

	inner := boxWidth - 2

	priorityStr := priorityLabel(t.PriorityString())
	durStr := "—"
	if t.Duration > 0 {
		durStr = util.FormatDuration(t.Duration * 60)
	}
	proj := t.ProjectName
	if proj == "" {
		proj = "—"
	}
	date := t.StartDate
	if date == "" {
		date = t.DueDate
	}
	if date == "" {
		date = "—"
	}

	rows := []string{
		fmt.Sprintf("%-10s %s", "Title:", t.Title),
		fmt.Sprintf("%-10s %s", "Priority:", priorityStr),
		fmt.Sprintf("%-10s %s", "Date:", date),
		fmt.Sprintf("%-10s %s", "Duration:", durStr),
		fmt.Sprintf("%-10s %s", "Project:", proj),
		fmt.Sprintf("%-10s %s", "Type:", t.TaskType),
	}

	top := "╔" + strings.Repeat("═", inner) + "╗"
	headerText := padRight("TASK DETAIL", inner-6) + "[Esc]"
	header := "║" + headerText + "║"
	sep1 := "╠" + strings.Repeat("═", inner) + "╣"
	bottom := "╚" + strings.Repeat("═", inner) + "╝"

	lines := []string{top, header, sep1}
	for _, row := range rows {
		lines = append(lines, "║ "+padRight(truncate(row, inner-2), inner-2)+" ║")
	}

	if t.Description != "" {
		lines = append(lines, sep1)
		lines = append(lines, "║ "+padRight("Description:", inner-2)+" ║")
		desc := truncate(t.Description, inner-4)
		lines = append(lines, "║   "+padRight(desc, inner-4)+" ║")
	}

	lines = append(lines, bottom)
	lines = append(lines, "")
	lines = append(lines, dimStyle.Render("  Esc=close"))

	return strings.Join(lines, "\n")
}

func priorityLabel(p string) string {
	switch p {
	case "CRITICAL":
		return "● Critical Now"
	case "OPPORTUNITY_NOW":
		return "○ Opportunity Now"
	case "OVER_THE_HORIZON":
		return "◌ Over The Horizon"
	case "PARKING_LOT":
		return "⊡ Parking Lot"
	default:
		return "— (inbox)"
	}
}

func padRight(s string, width int) string {
	runes := []rune(s)
	if len(runes) >= width {
		return string(runes[:width])
	}
	return s + strings.Repeat(" ", width-len(runes))
}
