package output

import (
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/glamour/styles"
)

// RenderMarkdown renders a markdown string for terminal display.
// Falls back to raw text if rendering fails.
func RenderMarkdown(md string, noColor bool) string {
	opts := []glamour.TermRendererOption{
		glamour.WithWordWrap(80),
	}
	if noColor {
		opts = append(opts, glamour.WithStyles(styles.NoTTYStyleConfig))
	} else {
		opts = append(opts, glamour.WithAutoStyle())
	}
	r, err := glamour.NewTermRenderer(opts...)
	if err != nil {
		return md
	}
	out, err := r.Render(md)
	if err != nil {
		return md
	}
	return out
}
