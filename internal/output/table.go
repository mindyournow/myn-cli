package output

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

// Table formats data as a column-aligned text table.
type Table struct {
	headers []string
	rows    [][]string
	widths  []int
}

// NewTable creates a new table with the given headers.
func NewTable(headers ...string) *Table {
	widths := make([]int, len(headers))
	for i, h := range headers {
		widths[i] = displayWidth(h)
	}
	return &Table{
		headers: headers,
		widths:  widths,
	}
}

// AddRow adds a row to the table.
func (t *Table) AddRow(cells ...string) *Table {
	row := make([]string, len(t.headers))
	for i, cell := range cells {
		if i < len(row) {
			row[i] = cell
		}
	}

	// Update column widths
	for i, cell := range row {
		if i < len(t.widths) {
			w := displayWidth(cell)
			if w > t.widths[i] {
				t.widths[i] = w
			}
		}
	}

	t.rows = append(t.rows, row)
	return t
}

// String returns the formatted table.
func (t *Table) String() string {
	if len(t.rows) == 0 {
		return ""
	}

	var lines []string

	// Header row
	headerRow := t.formatRow(t.headers, true)
	lines = append(lines, headerRow)

	// Separator
	sepParts := make([]string, len(t.widths))
	for i, w := range t.widths {
		sepParts[i] = strings.Repeat("-", w)
	}
	lines = append(lines, strings.Join(sepParts, "  "))

	// Data rows
	for _, row := range t.rows {
		lines = append(lines, t.formatRow(row, false))
	}

	return strings.Join(lines, "\n")
}

// formatRow formats a single row.
func (t *Table) formatRow(cells []string, isHeader bool) string {
	parts := make([]string, len(t.widths))
	for i, width := range t.widths {
		cell := ""
		if i < len(cells) {
			cell = cells[i]
		}

		// Pad or truncate
		cellWidth := displayWidth(cell)
		if cellWidth < width {
			cell += strings.Repeat(" ", width-cellWidth)
		}

		if isHeader && !ColorEnabled {
			cell = Bold(cell)
		}

		parts[i] = cell
	}

	return strings.Join(parts, "  ")
}

// displayWidth returns the display width of a string (accounting for ANSI codes).
func displayWidth(s string) int {
	// Strip ANSI codes before measuring
	stripped := Strip(s)
	return utf8.RuneCountInString(stripped)
}

// SimpleTable creates and returns a simple table string.
func SimpleTable(headers []string, rows [][]string) string {
	t := NewTable(headers...)
	for _, row := range rows {
		t.AddRow(row...)
	}
	return t.String()
}

// KVTable creates a key-value table.
func KVTable(pairs map[string]string) string {
	// Find max key length
	maxLen := 0
	for k := range pairs {
		if w := displayWidth(k); w > maxLen {
			maxLen = w
		}
	}

	var lines []string
	for k, v := range pairs {
		padding := strings.Repeat(" ", maxLen-displayWidth(k))
		lines = append(lines, fmt.Sprintf("%s%s: %s", k, padding, v))
	}

	return strings.Join(lines, "\n")
}
