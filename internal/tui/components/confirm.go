package components

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ConfirmYesMsg is sent when the user confirms the dialog.
type ConfirmYesMsg struct{}

// ConfirmNoMsg is sent when the user cancels the dialog.
type ConfirmNoMsg struct{}

var (
	confirmOverlayStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("#7c3aed")).
				Background(lipgloss.Color("#0f172a")).
				Padding(1, 2)

	confirmMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#e2e8f0")).
				Background(lipgloss.Color("#0f172a")).
				Bold(true)

	confirmHintStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#94a3b8")).
				Background(lipgloss.Color("#0f172a"))

	confirmAccentStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#a78bfa")).
				Background(lipgloss.Color("#0f172a")).
				Bold(true)
)

// Confirm is a modal yes/no confirmation dialog.
type Confirm struct {
	message string
	width   int
	height  int
}

// NewConfirm creates a new Confirm dialog with the given message.
func NewConfirm(message string) Confirm {
	return Confirm{message: message}
}

// SetSize updates the dialog's available area (used for centering by the caller).
func (c *Confirm) SetSize(w, h int) {
	c.width = w
	c.height = h
}

// Update handles key events for the confirmation dialog.
func (c Confirm) Update(msg tea.Msg) (Confirm, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "y", "Y", "enter":
			return c, func() tea.Msg { return ConfirmYesMsg{} }
		case "n", "N", "esc":
			return c, func() tea.Msg { return ConfirmNoMsg{} }
		}
	}
	return c, nil
}

// View renders the confirmation dialog box.
func (c Confirm) View() string {
	dialogWidth := c.width - 10
	if dialogWidth < 40 {
		dialogWidth = 40
	}
	if dialogWidth > 70 {
		dialogWidth = 70
	}

	msgLine := confirmMessageStyle.Render(c.message)

	separator := confirmHintStyle.Render(strings.Repeat("─", dialogWidth-4))

	yesKey := confirmAccentStyle.Render("[y]")
	noKey := confirmAccentStyle.Render("[n]")
	escKey := confirmHintStyle.Render("[Esc]")

	hint := "  " + yesKey + confirmHintStyle.Render(" Yes  ") +
		noKey + confirmHintStyle.Render(" No  ") +
		escKey + confirmHintStyle.Render(" Cancel  ")

	content := strings.Join([]string{msgLine, separator, hint}, "\n")
	return confirmOverlayStyle.Width(dialogWidth).Render(content)
}
