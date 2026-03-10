package output

import (
	"bytes"
	"strings"
	"testing"
)

func TestProgressBar_Empty(t *testing.T) {
	result := ProgressBar(0, 10, 20, true)
	if !strings.Contains(result, "0%") {
		t.Errorf("ProgressBar(0,10) = %q, want to contain \"0%%\"", result)
	}
	// No filled blocks expected when current=0
	if strings.Contains(result, "█") {
		t.Errorf("ProgressBar(0,10) should contain no filled chars, got %q", result)
	}
}

func TestProgressBar_Full(t *testing.T) {
	result := ProgressBar(10, 10, 20, true)
	if !strings.Contains(result, "100%") {
		t.Errorf("ProgressBar(10,10) = %q, want to contain \"100%%\"", result)
	}
}

func TestProgressBar_Half(t *testing.T) {
	result := ProgressBar(5, 10, 20, true)
	if !strings.Contains(result, "50%") {
		t.Errorf("ProgressBar(5,10) = %q, want to contain \"50%%\"", result)
	}
}

func TestProgressBar_NoColor(t *testing.T) {
	result := ProgressBar(5, 10, 20, true)
	if strings.Contains(result, "\033[") {
		t.Errorf("ProgressBar with noColor=true should not contain ANSI codes, got %q", result)
	}
}

func TestProgressBar_WithColor(t *testing.T) {
	result := ProgressBar(5, 10, 20, false)
	if !strings.Contains(result, "\033[0m") {
		t.Errorf("ProgressBar with noColor=false and current>0 should contain ANSI reset code, got %q", result)
	}
}

func TestProgressBar_ZeroTotal(t *testing.T) {
	// Should not panic; total=0 is treated as 1 internally
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("ProgressBar(0,0) panicked: %v", r)
		}
	}()
	result := ProgressBar(0, 0, 20, true)
	if result == "" {
		t.Error("ProgressBar(0,0) returned empty string")
	}
}

func TestProgressBar_Width(t *testing.T) {
	const width = 10
	// Use noColor=true to avoid ANSI codes interfering with char counting.
	result := ProgressBar(10, 10, width, true)
	// The bar portion lives inside [...]; count █ + ░ characters.
	inner := result[strings.Index(result, "[")+1 : strings.Index(result, "]")]
	// Count runes because █ and ░ are multi-byte.
	runes := []rune(inner)
	if len(runes) != width {
		t.Errorf("ProgressBar width=%d: bar has %d chars, want %d (bar=%q)", width, len(runes), width, inner)
	}
}

func TestStreamWriter_Write(t *testing.T) {
	var buf bytes.Buffer
	sw := NewStreamWriter(&buf)

	if err := sw.Write("hello"); err != nil {
		t.Fatalf("Write() error = %v", err)
	}
	if got := buf.String(); got != "hello" {
		t.Errorf("Write() buf = %q, want %q", got, "hello")
	}
}

func TestStreamWriter_Flush(t *testing.T) {
	var buf bytes.Buffer
	sw := NewStreamWriter(&buf)

	if err := sw.Flush(); err != nil {
		t.Fatalf("Flush() error = %v", err)
	}
	if got := buf.String(); got != "\n" {
		t.Errorf("Flush() buf = %q, want newline", got)
	}
}

func TestStreamWriter_Write_Multiple(t *testing.T) {
	var buf bytes.Buffer
	sw := NewStreamWriter(&buf)

	chunks := []string{"foo", "bar", "baz"}
	for _, c := range chunks {
		if err := sw.Write(c); err != nil {
			t.Fatalf("Write(%q) error = %v", c, err)
		}
	}
	got := buf.String()
	for _, c := range chunks {
		if !strings.Contains(got, c) {
			t.Errorf("Write() buf = %q, want to contain %q", got, c)
		}
	}
}
