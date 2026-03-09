package output

import (
	"fmt"
	"io"

	"github.com/charmbracelet/glamour"
)

// MarkdownRenderer handles markdown rendering with Glamour.
type MarkdownRenderer struct {
	renderer *glamour.TermRenderer
	width    uint
}

// NewMarkdownRenderer creates a new markdown renderer.
func NewMarkdownRenderer() (*MarkdownRenderer, error) {
	style := "auto"
	if !ColorEnabled {
		style = "notty"
	}

	r, err := glamour.NewTermRenderer(
		glamour.WithStylePath(style),
		glamour.WithWordWrap(80),
	)
	if err != nil {
		return nil, fmt.Errorf("creating markdown renderer: %w", err)
	}

	return &MarkdownRenderer{
		renderer: r,
		width:    80,
	}, nil
}

// NewMarkdownRendererWithWidth creates a renderer with a specific width.
func NewMarkdownRendererWithWidth(width uint) (*MarkdownRenderer, error) {
	style := "auto"
	if !ColorEnabled {
		style = "notty"
	}

	r, err := glamour.NewTermRenderer(
		glamour.WithStylePath(style),
		glamour.WithWordWrap(int(width)),
	)
	if err != nil {
		return nil, fmt.Errorf("creating markdown renderer: %w", err)
	}

	return &MarkdownRenderer{
		renderer: r,
		width:    width,
	}, nil
}

// Render renders markdown to styled text.
func (m *MarkdownRenderer) Render(markdown string) (string, error) {
	return m.renderer.Render(markdown)
}

// RenderTo renders markdown directly to a writer.
func (m *MarkdownRenderer) RenderTo(w io.Writer, markdown string) error {
	out, err := m.Render(markdown)
	if err != nil {
		return err
	}
	_, err = w.Write([]byte(out))
	return err
}

// RenderSimple renders markdown without error handling (returns raw on error).
func RenderSimple(markdown string) string {
	r, err := NewMarkdownRenderer()
	if err != nil {
		return markdown
	}
	out, err := r.Render(markdown)
	if err != nil {
		return markdown
	}
	return out
}

// PrintMarkdown renders and prints markdown to stdout.
func PrintMarkdown(markdown string) error {
	r, err := NewMarkdownRenderer()
	if err != nil {
		fmt.Println(markdown)
		return nil
	}
	out, err := r.Render(markdown)
	if err != nil {
		fmt.Println(markdown)
		return nil
	}
	fmt.Println(out)
	return nil
}
