package components

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	sseContentStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#e2e8f0")).
			Background(lipgloss.Color("#0f172a"))

	sseCursorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7c3aed")).
			Background(lipgloss.Color("#0f172a"))

	sseErrStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ef4444")).
			Background(lipgloss.Color("#0f172a"))

	sseBgStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#0f172a"))
)

// SSEChunkMsg carries a single streamed text chunk.
type SSEChunkMsg struct {
	Chunk string
}

// SSEDoneMsg signals the stream has ended.
type SSEDoneMsg struct{}

// SSEErrMsg carries a stream error.
type SSEErrMsg struct {
	Err error
}

// SSEReader accumulates SSE stream chunks and renders them.
type SSEReader struct {
	content strings.Builder
	loading bool
	err     error
}

// NewSSEReader creates a new SSEReader.
func NewSSEReader() SSEReader {
	return SSEReader{loading: false}
}

// AppendChunk adds a new text chunk to the accumulated content.
func (s *SSEReader) AppendChunk(chunk string) {
	s.content.WriteString(chunk)
	s.loading = true
}

// SetDone marks the stream as complete.
func (s *SSEReader) SetDone() {
	s.loading = false
}

// SetError records a stream error and stops loading.
func (s *SSEReader) SetError(err error) {
	s.err = err
	s.loading = false
}

// Reset clears all accumulated content and state.
func (s *SSEReader) Reset() {
	s.content.Reset()
	s.loading = false
	s.err = nil
}

// View renders the accumulated content with a blinking cursor when loading.
func (s SSEReader) View() string {
	if s.err != nil {
		return sseErrStyle.Render("Error: " + s.err.Error())
	}

	text := s.content.String()
	if s.loading {
		return sseBgStyle.Render(
			sseContentStyle.Render(text) + sseCursorStyle.Render("▋"),
		)
	}
	return sseContentStyle.Render(text)
}
