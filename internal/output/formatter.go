package output

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

// Formatter handles CLI output in various formats.
type Formatter struct {
	JSON    bool
	Quiet   bool
	NoColor bool
	writer  io.Writer
}

// NewFormatter creates a new formatter with default settings.
func NewFormatter() *Formatter {
	return &Formatter{
		writer: os.Stdout,
	}
}

// New creates a new formatter with specific settings.
func New(json, quiet, noColor bool) *Formatter {
	return &Formatter{
		JSON:    json,
		Quiet:   quiet,
		NoColor: noColor,
		writer:  os.Stdout,
	}
}

// SetWriter sets the output writer.
func (f *Formatter) SetWriter(w io.Writer) {
	f.writer = w
}

// Print outputs data in the configured format.
func (f *Formatter) Print(data any) {
	if f.Quiet {
		return
	}

	if f.JSON {
		f.printJSON(data)
		return
	}

	switch v := data.(type) {
	case string:
		f.Println(v)
	case fmt.Stringer:
		f.Println(v.String())
	default:
		f.printJSON(data)
	}
}

// Println prints a line of text.
func (f *Formatter) Println(a ...any) {
	if f.Quiet {
		return
	}
	fmt.Fprintln(f.writer, a...)
}

// Printf prints formatted text.
func (f *Formatter) Printf(format string, a ...any) {
	if f.Quiet {
		return
	}
	fmt.Fprintf(f.writer, format, a...)
}

// Success prints a success message.
func (f *Formatter) Success(msg string) {
	if f.Quiet {
		return
	}
	if f.JSON {
		f.printJSON(map[string]string{"status": "success", "message": msg})
		return
	}
	if !f.NoColor {
		f.Println(Green("✓ " + msg))
	} else {
		f.Println("✓ " + msg)
	}
}

// Error prints an error message.
func (f *Formatter) Error(err error) {
	if f.Quiet {
		return
	}
	if f.JSON {
		f.printJSON(map[string]string{"status": "error", "error": err.Error()})
		return
	}
	if !f.NoColor {
		fmt.Fprintln(f.writer, Red("✗ "+err.Error()))
	} else {
		fmt.Fprintln(f.writer, "✗ "+err.Error())
	}
}

// Warning prints a warning message.
func (f *Formatter) Warning(msg string) {
	if f.Quiet {
		return
	}
	if f.JSON {
		f.printJSON(map[string]string{"status": "warning", "message": msg})
		return
	}
	if !f.NoColor {
		f.Println(Yellow("⚠ " + msg))
	} else {
		f.Println("⚠ " + msg)
	}
}

// Info prints an informational message.
func (f *Formatter) Info(msg string) {
	if f.Quiet {
		return
	}
	if !f.NoColor {
		f.Println(Blue("ℹ " + msg))
	} else {
		f.Println("ℹ " + msg)
	}
}

// JSONResult outputs a structured JSON result.
func (f *Formatter) JSONResult(data any) {
	f.printJSON(data)
}

// printJSON outputs data as formatted JSON.
func (f *Formatter) printJSON(data any) {
	enc := json.NewEncoder(f.writer)
	enc.SetIndent("", "  ")
	if err := enc.Encode(data); err != nil {
		fmt.Fprintf(f.writer, "Error encoding JSON: %v\n", err)
	}
}

// Table outputs data as a formatted table.
func (f *Formatter) Table(headers []string, rows [][]string) {
	if f.Quiet {
		return
	}

	if f.JSON {
		// Convert to JSON array of objects
		data := make([]map[string]string, len(rows))
		for i, row := range rows {
			obj := make(map[string]string)
			for j, header := range headers {
				if j < len(row) {
					obj[header] = row[j]
				}
			}
			data[i] = obj
		}
		f.printJSON(data)
		return
	}

	t := NewTable(headers...)
	for _, row := range rows {
		t.AddRow(row...)
	}
	f.Println(t.String())
}

// List outputs a list of items.
func (f *Formatter) List(items []string) {
	if f.Quiet {
		return
	}

	if f.JSON {
		f.printJSON(items)
		return
	}

	for _, item := range items {
		f.Println("  • " + item)
	}
}

// KeyValue outputs key-value pairs.
func (f *Formatter) KeyValue(pairs map[string]string) {
	if f.Quiet {
		return
	}

	if f.JSON {
		f.printJSON(pairs)
		return
	}

	// Find max key length for alignment
	maxLen := 0
	for k := range pairs {
		if len(k) > maxLen {
			maxLen = len(k)
		}
	}

	for k, v := range pairs {
		padding := strings.Repeat(" ", maxLen-len(k))
		f.Printf("%s%s: %s\n", k, padding, v)
	}
}

// Empty prints a message when there's no data.
func (f *Formatter) Empty(msg string) {
	if f.Quiet {
		return
	}
	if f.JSON {
		f.printJSON([]any{})
		return
	}
	f.Println(msg)
}
