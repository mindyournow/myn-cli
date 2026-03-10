package output

import (
	"strings"
	"testing"
)

func TestRenderMarkdown_Basic(t *testing.T) {
	out := RenderMarkdown("# Hello", false)
	if out == "" {
		t.Fatal("RenderMarkdown returned empty string")
	}
	if !strings.Contains(out, "Hello") {
		t.Errorf("RenderMarkdown(%q) = %q, want output containing 'Hello'", "# Hello", out)
	}
}

func TestRenderMarkdown_FallbackOnEmpty(t *testing.T) {
	out := RenderMarkdown("", false)
	// Empty input should not panic and should return an empty or whitespace-only string.
	stripped := strings.TrimSpace(out)
	if stripped != "" {
		// glamour may emit a single newline for empty input; that is acceptable.
		// Only fail if substantive content is returned.
		if len(stripped) > 5 {
			t.Errorf("RenderMarkdown(%q) = %q, expected empty or near-empty", "", out)
		}
	}
}

func TestRenderMarkdown_NoColor(t *testing.T) {
	out := RenderMarkdown("# Hello", true)
	// Should not panic and should still contain the heading text.
	if !strings.Contains(out, "Hello") {
		t.Errorf("RenderMarkdown noColor: output %q does not contain 'Hello'", out)
	}
}

func TestRenderMarkdown_Bold(t *testing.T) {
	out := RenderMarkdown("**bold**", false)
	if !strings.Contains(out, "bold") {
		t.Errorf("RenderMarkdown(%q) = %q, want output containing 'bold'", "**bold**", out)
	}
}

func TestRenderMarkdown_FallbackPreservesInput(t *testing.T) {
	input := "Just plain text, no markdown syntax."
	out := RenderMarkdown(input, false)
	if !strings.Contains(out, "Just plain text") {
		t.Errorf("RenderMarkdown plain text: output %q does not contain original text", out)
	}
}
