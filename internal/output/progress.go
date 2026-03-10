package output

import (
	"fmt"
	"io"
	"strings"
)

// ProgressBar renders a simple horizontal ASCII progress bar.
func ProgressBar(current, total int, width int, noColor bool) string {
	if total <= 0 {
		total = 1
	}
	if width <= 0 {
		width = 20
	}
	filled := int(float64(current) / float64(total) * float64(width))
	if filled > width {
		filled = width
	}
	bar := strings.Repeat("█", filled) + strings.Repeat("░", width-filled)
	pct := int(float64(current) / float64(total) * 100)
	label := fmt.Sprintf("[%s] %d%%", bar, pct)
	if !noColor && filled > 0 {
		bar = colorGreen + strings.Repeat("█", filled) + colorReset + strings.Repeat("░", width-filled)
		label = fmt.Sprintf("[%s] %d%%", bar, pct)
	}
	return label
}

// StreamWriter writes text chunks to w as they arrive (for SSE streaming output).
type StreamWriter struct {
	w io.Writer
}

// NewStreamWriter creates a StreamWriter wrapping w.
func NewStreamWriter(w io.Writer) *StreamWriter {
	return &StreamWriter{w: w}
}

// Write outputs a chunk. Returns an error if the write fails.
func (s *StreamWriter) Write(chunk string) error {
	_, err := fmt.Fprint(s.w, chunk)
	return err
}

// Flush writes a newline to terminate the stream.
func (s *StreamWriter) Flush() error {
	_, err := fmt.Fprintln(s.w)
	return err
}
