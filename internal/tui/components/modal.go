package components

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	modalBorderStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("#7c3aed")).
				Background(lipgloss.Color("#0f172a")).
				Foreground(lipgloss.Color("#e2e8f0"))

	modalTitleStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#0f172a")).
			Foreground(lipgloss.Color("#7c3aed")).
			Bold(true).
			Padding(0, 1)

	modalBodyStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#0f172a")).
			Foreground(lipgloss.Color("#e2e8f0")).
			Padding(1, 2)
)

// Modal is a centered overlay box with a title and body content.
type Modal struct {
	title   string
	content string
	width   int
	height  int
}

// NewModal creates a new Modal with the given title and content.
func NewModal(title, content string) Modal {
	return Modal{
		title:   title,
		content: content,
	}
}

// SetSize informs the modal of the terminal dimensions so it can center itself.
func (m *Modal) SetSize(w, h int) {
	m.width = w
	m.height = h
}

// View renders the modal as a centered overlay string.
func (m Modal) View() string {
	boxWidth := m.width - 10
	if boxWidth > 80 {
		boxWidth = 80
	}
	if boxWidth < 20 {
		boxWidth = 20
	}

	innerWidth := boxWidth - 4 // border (2) + padding (2)
	if innerWidth < 1 {
		innerWidth = 1
	}

	title := modalTitleStyle.Width(innerWidth).Render(m.title)
	divider := lipgloss.NewStyle().
		Background(lipgloss.Color("#0f172a")).
		Foreground(lipgloss.Color("#7c3aed")).
		Render(lipgloss.NewStyle().Width(innerWidth).Render(""))

	body := modalBodyStyle.Width(innerWidth).Render(m.content)

	inner := lipgloss.JoinVertical(lipgloss.Left, title, divider, body)

	box := modalBorderStyle.Width(boxWidth).Render(inner)

	if m.width > 0 && m.height > 0 {
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, box)
	}
	return box
}
