package components

import (
	"fmt"
	"math"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	progressLabelStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#e2e8f0")).
				Background(lipgloss.Color("#0f172a"))

	progressPctStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#475569")).
				Background(lipgloss.Color("#0f172a"))

	progressFilledPurple = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#7c3aed")).
				Background(lipgloss.Color("#0f172a"))

	progressFilledGreen = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#22c55e")).
				Background(lipgloss.Color("#0f172a"))

	progressEmptyStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#475569")).
				Background(lipgloss.Color("#0f172a"))

	progressBgStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#0f172a"))
)

// ProgressBar renders a horizontal labeled progress bar.
type ProgressBar struct {
	label string
	value float64
	max   float64
	width int
}

// NewProgressBar creates a new ProgressBar.
func NewProgressBar(label string, value, max float64) ProgressBar {
	return ProgressBar{label: label, value: value, max: max}
}

// SetSize sets the available width.
func (p *ProgressBar) SetSize(w int) {
	p.width = w
}

// View renders the progress bar.
func (p ProgressBar) View() string {
	w := p.width
	if w < 1 {
		w = 40
	}

	pct := 0.0
	if p.max > 0 {
		pct = p.value / p.max
		if pct > 1.0 {
			pct = 1.0
		}
		if pct < 0.0 {
			pct = 0.0
		}
	}

	pctText := fmt.Sprintf("%3.0f%%", math.Round(pct*100))

	// Fixed widths for label and percentage
	labelRendered := progressLabelStyle.Render(p.label)
	pctRendered := progressPctStyle.Render(pctText)

	labelW := lipgloss.Width(labelRendered)
	pctW := lipgloss.Width(pctRendered)

	// Bar width = total - label - spaces - brackets - percentage
	// Format: "label  [████░░░░░]  42%"
	// Overhead: 2 spaces + 2 brackets + 2 spaces = 6
	barWidth := w - labelW - pctW - 6
	if barWidth < 1 {
		barWidth = 1
	}

	filled := int(math.Round(float64(barWidth) * pct))
	if filled > barWidth {
		filled = barWidth
	}
	empty := barWidth - filled

	var barStyle *lipgloss.Style
	if pct >= 0.80 {
		barStyle = &progressFilledGreen
	} else {
		barStyle = &progressFilledPurple
	}

	filledStr := barStyle.Render(strings.Repeat("█", filled))
	emptyStr := progressEmptyStyle.Render(strings.Repeat("░", empty))

	bracket := progressEmptyStyle.Render

	bar := bracket("[") + filledStr + emptyStr + bracket("]")

	result := labelRendered + progressBgStyle.Render("  ") + bar + progressBgStyle.Render("  ") + pctRendered
	return progressBgStyle.Render(result)
}
