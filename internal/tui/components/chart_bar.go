package components

import (
	"fmt"
	"math"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	chartTitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#e2e8f0")).
			Background(lipgloss.Color("#0f172a")).
			Bold(true)

	chartBarStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7c3aed")).
			Background(lipgloss.Color("#0f172a"))

	chartLabelStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#e2e8f0")).
			Background(lipgloss.Color("#0f172a"))

	chartValueStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#475569")).
			Background(lipgloss.Color("#0f172a"))

	chartBgStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#0f172a"))
)

// BarEntry is a single bar in the chart.
type BarEntry struct {
	Label string
	Value float64
}

// BarChart renders an ASCII horizontal bar chart.
type BarChart struct {
	title  string
	bars   []BarEntry
	width  int
	height int
}

// NewBarChart creates a new BarChart.
func NewBarChart(title string, bars []BarEntry) BarChart {
	return BarChart{title: title, bars: bars}
}

// SetSize sets the available terminal dimensions.
func (b *BarChart) SetSize(w, h int) {
	b.width = w
	b.height = h
}

// View renders the bar chart.
func (b BarChart) View() string {
	w := b.width
	if w < 1 {
		w = 60
	}

	var rows []string

	if b.title != "" {
		rows = append(rows, chartTitleStyle.Render(b.title))
		rows = append(rows, "")
	}

	if len(b.bars) == 0 {
		rows = append(rows, chartValueStyle.Render("No data."))
		return chartBgStyle.Width(w).Render(strings.Join(rows, "\n"))
	}

	// Find max value and max label width
	maxVal := 0.0
	maxLabelW := 0
	for _, bar := range b.bars {
		if bar.Value > maxVal {
			maxVal = bar.Value
		}
		if len(bar.Label) > maxLabelW {
			maxLabelW = len(bar.Label)
		}
	}

	// Value display width (e.g., "1234")
	valDisplayW := len(fmt.Sprintf("%.0f", maxVal))
	if valDisplayW < 1 {
		valDisplayW = 1
	}

	// Bar area = total width - label col - 2 spaces - 2 spaces - value col
	barAreaW := w - maxLabelW - 2 - 2 - valDisplayW
	if barAreaW < 1 {
		barAreaW = 1
	}

	for _, bar := range b.bars {
		barLen := 0
		if maxVal > 0 {
			barLen = int(math.Round(float64(barAreaW) * bar.Value / maxVal))
		}
		if barLen > barAreaW {
			barLen = barAreaW
		}

		label := fmt.Sprintf("%-*s", maxLabelW, bar.Label)
		barStr := strings.Repeat("█", barLen)
		valStr := fmt.Sprintf("%*.0f", valDisplayW, bar.Value)

		row := chartLabelStyle.Render(label) +
			chartBgStyle.Render("  ") +
			chartBarStyle.Render(barStr) +
			chartBgStyle.Render(strings.Repeat(" ", barAreaW-barLen+2)) +
			chartValueStyle.Render(valStr)

		rows = append(rows, row)
	}

	return chartBgStyle.Width(w).Render(strings.Join(rows, "\n"))
}
