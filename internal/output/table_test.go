package output

import (
	"bytes"
	"strings"
	"testing"

	"github.com/mattn/go-runewidth"
)

// ---- stripANSI tests ----

func TestStripANSI_Basic(t *testing.T) {
	input := "\033[32mgreen\033[0m"
	got := stripANSI(input)
	if got != "green" {
		t.Errorf("stripANSI(%q) = %q, want %q", input, got, "green")
	}
}

func TestStripANSI_256Color(t *testing.T) {
	input := "\033[38;5;196m text \033[0m"
	got := stripANSI(input)
	if got != " text " {
		t.Errorf("stripANSI(%q) = %q, want %q", input, got, " text ")
	}
}

func TestStripANSI_TrueColor(t *testing.T) {
	input := "\033[38;2;255;100;0m text \033[0m"
	got := stripANSI(input)
	if got != " text " {
		t.Errorf("stripANSI(%q) = %q, want %q", input, got, " text ")
	}
}

func TestStripANSI_Bold(t *testing.T) {
	input := "\033[1mBold\033[0m"
	got := stripANSI(input)
	if got != "Bold" {
		t.Errorf("stripANSI(%q) = %q, want %q", input, got, "Bold")
	}
}

func TestStripANSI_Mixed(t *testing.T) {
	input := "plain \033[32mgreen\033[0m and \033[1;31mbold-red\033[0m end"
	got := stripANSI(input)
	want := "plain green and bold-red end"
	if got != want {
		t.Errorf("stripANSI(%q) = %q, want %q", input, got, want)
	}
}

func TestStripANSI_NoANSI(t *testing.T) {
	input := "plain text, no escapes"
	got := stripANSI(input)
	if got != input {
		t.Errorf("stripANSI(%q) = %q, want unchanged", input, got)
	}
}

// ---- Table tests ----

func TestTable_AlignASCII(t *testing.T) {
	var buf bytes.Buffer
	tbl := NewTable(&buf, "NAME", "STATUS")
	tbl.AddRow("alpha", "active")
	tbl.AddRow("beta", "inactive")
	tbl.Render()

	output := buf.String()
	lines := strings.Split(strings.TrimRight(output, "\n"), "\n")

	// Expect: header line, separator, two data rows
	if len(lines) < 4 {
		t.Fatalf("expected at least 4 lines, got %d:\n%s", len(lines), output)
	}

	// Header must contain both column names
	if !strings.Contains(lines[0], "NAME") || !strings.Contains(lines[0], "STATUS") {
		t.Errorf("header line %q missing column names", lines[0])
	}

	// "beta" is shorter than "alpha" but should be padded to same column width.
	// The STATUS column of "alpha" row and "beta" row should start at the same
	// position (both padded to the width of "inactive").
	alphaLine := lines[2]
	betaLine := lines[3]
	// Find where "active" / "inactive" start.
	alphaIdx := strings.Index(alphaLine, "active")
	betaIdx := strings.Index(betaLine, "inactive")
	if alphaIdx != betaIdx {
		t.Errorf("columns misaligned: 'active' at %d vs 'inactive' at %d\nalpha: %q\nbeta:  %q",
			alphaIdx, betaIdx, alphaLine, betaLine)
	}
}

func TestTable_AlignEmoji(t *testing.T) {
	var buf bytes.Buffer
	tbl := NewTable(&buf, "ITEM", "DESC")
	tbl.AddRow("🔥 Fire", "hot")
	tbl.AddRow("Water", "cold")
	tbl.Render()

	output := buf.String()
	// The table should render without panicking and contain both values.
	if !strings.Contains(output, "Fire") {
		t.Errorf("output missing 'Fire': %q", output)
	}
	if !strings.Contains(output, "Water") {
		t.Errorf("output missing 'Water': %q", output)
	}

	// Both data rows' second columns ("hot", "cold") should start at the same
	// visual column. Measure the visual width of the prefix before the second token
	// rather than using byte-offset, because emoji occupy 2 terminal columns.
	lines := strings.Split(strings.TrimRight(output, "\n"), "\n")
	if len(lines) < 4 {
		t.Fatalf("expected at least 4 lines, got %d", len(lines))
	}

	// visualPrefixWidth returns the terminal display width of the text before needle.
	visualPrefixWidth := func(line, needle string) int {
		idx := strings.Index(line, needle)
		if idx < 0 {
			return -1
		}
		return runewidth.StringWidth(line[:idx])
	}

	fireVis := visualPrefixWidth(lines[2], "hot")
	waterVis := visualPrefixWidth(lines[3], "cold")
	if fireVis < 0 {
		t.Fatalf("'hot' not found in line: %q", lines[2])
	}
	if waterVis < 0 {
		t.Fatalf("'cold' not found in line: %q", lines[3])
	}
	if fireVis != waterVis {
		t.Errorf("emoji cell misaligned second column: fire row col2 starts at visual col %d, water row at %d\n%s",
			fireVis, waterVis, output)
	}
}

func TestTable_AlignCJK(t *testing.T) {
	var buf bytes.Buffer
	tbl := NewTable(&buf, "LANG", "TEXT")
	tbl.AddRow("日本語", "ja")
	tbl.AddRow("English", "en")
	tbl.Render()

	output := buf.String()
	if !strings.Contains(output, "日本語") {
		t.Errorf("output missing CJK text: %q", output)
	}

	// Both rows' second column should start at the same visual column.
	// CJK characters each occupy 2 terminal columns, so we measure the visual
	// width of the prefix before the second token rather than the byte offset.
	lines := strings.Split(strings.TrimRight(output, "\n"), "\n")
	if len(lines) < 4 {
		t.Fatalf("expected at least 4 lines, got %d", len(lines))
	}

	visualPrefixWidth := func(line, needle string) int {
		idx := strings.Index(line, needle)
		if idx < 0 {
			return -1
		}
		return runewidth.StringWidth(line[:idx])
	}

	cjkVis := visualPrefixWidth(lines[2], "ja")
	engVis := visualPrefixWidth(lines[3], "en")
	if cjkVis < 0 {
		t.Fatalf("'ja' not found in line: %q", lines[2])
	}
	if engVis < 0 {
		t.Fatalf("'en' not found in line: %q", lines[3])
	}
	if cjkVis != engVis {
		t.Errorf("CJK cell misaligned: CJK row col2 starts at visual col %d, English row at %d\n%s",
			cjkVis, engVis, output)
	}
}

func TestTable_ANSI_Width(t *testing.T) {
	var buf bytes.Buffer
	tbl := NewTable(&buf, "NAME", "STATUS")
	// Color-wrapped cell: visible text is "active" (6 chars), but raw string is longer.
	coloredActive := "\033[32mactive\033[0m"
	tbl.AddRow("alpha", coloredActive)
	tbl.AddRow("beta", "inactive")
	tbl.Render()

	output := buf.String()
	lines := strings.Split(strings.TrimRight(output, "\n"), "\n")
	if len(lines) < 4 {
		t.Fatalf("expected at least 4 lines, got %d:\n%s", len(lines), output)
	}

	// The second column width should be based on visible width of "inactive" (8),
	// not the raw byte length of the ANSI sequence. Check alignment by stripping
	// ANSI from both rows before comparing positions.
	stripped2 := stripANSI(lines[2])
	stripped3 := stripANSI(lines[3])
	// After stripping, "active" and "inactive" second-column positions must match.
	idx2 := strings.Index(stripped2, "active")
	idx3 := strings.Index(stripped3, "inactive")
	if idx2 != idx3 {
		t.Errorf("ANSI-colored cell misaligned after stripping: col2 at %d vs %d\n%s\n%s",
			idx2, idx3, stripped2, stripped3)
	}
}

func TestTable_Empty(t *testing.T) {
	var buf bytes.Buffer
	tbl := NewTable(&buf)
	// Should not panic
	tbl.Render()
	// Output should be empty (no headers, no rows)
	if buf.Len() != 0 {
		t.Errorf("empty table produced output: %q", buf.String())
	}
}

func TestTable_Header(t *testing.T) {
	var buf bytes.Buffer
	tbl := NewTable(&buf, "COLUMN_A", "COLUMN_B")
	tbl.AddRow("val1", "val2")
	tbl.Render()

	output := buf.String()
	lines := strings.Split(strings.TrimRight(output, "\n"), "\n")
	if len(lines) < 3 {
		t.Fatalf("expected at least 3 lines, got %d:\n%s", len(lines), output)
	}

	headerLine := lines[0]
	if !strings.Contains(headerLine, "COLUMN_A") {
		t.Errorf("header line %q does not contain 'COLUMN_A'", headerLine)
	}
	if !strings.Contains(headerLine, "COLUMN_B") {
		t.Errorf("header line %q does not contain 'COLUMN_B'", headerLine)
	}

	// Separator line should follow immediately after the header.
	separatorLine := lines[1]
	if !strings.Contains(separatorLine, "─") {
		t.Errorf("separator line %q does not contain '─'", separatorLine)
	}
}
