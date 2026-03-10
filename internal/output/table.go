package output

import (
	"fmt"
	"io"
	"strings"
	"unicode/utf8"
)

// Table renders column-aligned text output.
type Table struct {
	headers []string
	rows    [][]string
	w       io.Writer
}

// NewTable creates a table that writes to w.
func NewTable(w io.Writer, headers ...string) *Table {
	return &Table{
		headers: headers,
		w:       w,
	}
}

// AddRow appends a row to the table.
func (t *Table) AddRow(cols ...string) {
	t.rows = append(t.rows, cols)
}

// Render writes the table with aligned columns.
func (t *Table) Render() {
	if len(t.headers) == 0 && len(t.rows) == 0 {
		return
	}

	// Compute column widths
	cols := len(t.headers)
	if len(t.rows) > 0 && len(t.rows[0]) > cols {
		cols = len(t.rows[0])
	}
	widths := make([]int, cols)
	for i, h := range t.headers {
		if w := utf8.RuneCountInString(h); w > widths[i] {
			widths[i] = w
		}
	}
	for _, row := range t.rows {
		for i, cell := range row {
			if i < cols {
				if w := utf8.RuneCountInString(stripANSI(cell)); w > widths[i] {
					widths[i] = w
				}
			}
		}
	}

	// Print header
	if len(t.headers) > 0 {
		printRow(t.w, t.headers, widths, true)
		printSeparator(t.w, widths)
	}

	// Print rows
	for _, row := range t.rows {
		printRow(t.w, row, widths, false)
	}
}

func printRow(w io.Writer, cells []string, widths []int, isHeader bool) {
	parts := make([]string, len(widths))
	for i := 0; i < len(widths); i++ {
		var cell string
		if i < len(cells) {
			cell = cells[i]
		}
		visible := utf8.RuneCountInString(stripANSI(cell))
		pad := widths[i] - visible
		if pad < 0 {
			pad = 0
		}
		if isHeader {
			parts[i] = cell + strings.Repeat(" ", pad)
		} else {
			parts[i] = cell + strings.Repeat(" ", pad)
		}
	}
	fmt.Fprintln(w, strings.Join(parts, "  "))
}

func printSeparator(w io.Writer, widths []int) {
	parts := make([]string, len(widths))
	for i, w := range widths {
		parts[i] = strings.Repeat("─", w)
	}
	fmt.Fprintln(w, strings.Join(parts, "  "))
}

// stripANSI removes ANSI escape sequences from s for width calculation.
func stripANSI(s string) string {
	var out strings.Builder
	i := 0
	for i < len(s) {
		if s[i] == '\033' && i+1 < len(s) && s[i+1] == '[' {
			// skip until 'm'
			j := i + 2
			for j < len(s) && (s[j] < 'A' || s[j] > 'Z') && (s[j] < 'a' || s[j] > 'z') {
				j++
			}
			if j < len(s) {
				j++ // skip the final letter
			}
			i = j
		} else {
			out.WriteByte(s[i])
			i++
		}
	}
	return out.String()
}
