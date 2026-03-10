package components

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	inputFocusedStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("#7c3aed")).
				Padding(0, 1)

	inputBlurredStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("#374151")).
				Padding(0, 1)
)

// Input wraps charmbracelet/bubbles textinput with styled borders.
type Input struct {
	model textinput.Model
}

// NewInput creates a new Input with the given placeholder text.
func NewInput(placeholder string) Input {
	ti := textinput.New()
	ti.Placeholder = placeholder
	ti.TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#e2e8f0"))
	ti.PlaceholderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#64748b"))
	ti.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#7c3aed"))
	return Input{model: ti}
}

// Focus gives the input keyboard focus.
func (i *Input) Focus() {
	i.model.Focus()
}

// Blur removes keyboard focus from the input.
func (i *Input) Blur() {
	i.model.Blur()
}

// Focused reports whether the input currently has focus.
func (i Input) Focused() bool {
	return i.model.Focused()
}

// Value returns the current text value.
func (i Input) Value() string {
	return i.model.Value()
}

// SetValue replaces the current text value.
func (i *Input) SetValue(s string) {
	i.model.SetValue(s)
}

// Update handles incoming messages and returns the updated Input and any command.
func (i Input) Update(msg tea.Msg) (Input, tea.Cmd) {
	var cmd tea.Cmd
	i.model, cmd = i.model.Update(msg)
	return i, cmd
}

// View renders the input field with an appropriate border style.
func (i Input) View() string {
	inner := i.model.View()
	if i.model.Focused() {
		return inputFocusedStyle.Render(inner)
	}
	return inputBlurredStyle.Render(inner)
}
