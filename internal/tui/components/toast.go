package components

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// toastTickMsg is the internal message used to advance the auto-hide timer.
type toastTickMsg struct{}

var (
	toastSuccessStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#ffffff")).
				Background(lipgloss.Color("#16a34a")).
				Padding(0, 2).
				Bold(true)

	toastErrorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ffffff")).
			Background(lipgloss.Color("#dc2626")).
			Padding(0, 2).
			Bold(true)

	toastInfoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ffffff")).
			Background(lipgloss.Color("#2563eb")).
			Padding(0, 2)
)

// Toast is a transient notification banner that auto-hides after 3 seconds.
type Toast struct {
	message string
	kind    string // "success" | "error" | "info"
	visible bool
}

// NewToast creates a new (hidden) Toast.
func NewToast() Toast {
	return Toast{}
}

// Show displays a notification with the given message and kind.
// kind should be "success", "error", or "info".
func (t *Toast) Show(msg string, kind string) {
	t.message = msg
	t.kind = kind
	t.visible = true
}

// Visible reports whether the toast is currently shown.
func (t Toast) Visible() bool {
	return t.visible
}

// Tick returns a tea.Cmd that fires the hide event after 3 seconds.
// Call this whenever Show is called.
func (t Toast) Tick() tea.Cmd {
	return tea.Tick(3*time.Second, func(_ time.Time) tea.Msg {
		return toastTickMsg{}
	})
}

// Update handles the internal tick that hides the toast.
func (t Toast) Update(msg tea.Msg) (Toast, tea.Cmd) {
	if _, ok := msg.(toastTickMsg); ok {
		t.visible = false
		return t, nil
	}
	return t, nil
}

// View renders the notification banner, or an empty string when not visible.
func (t Toast) View() string {
	if !t.visible {
		return ""
	}

	var style lipgloss.Style
	switch t.kind {
	case "success":
		style = toastSuccessStyle
	case "error":
		style = toastErrorStyle
	default:
		style = toastInfoStyle
	}

	return style.Render(t.message)
}
