package components

import (
	"fmt"
	"time"

	"github.com/charmbracelet/lipgloss"
)

var (
	timerLabelStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#94a3b8"))

	timerGreenStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#22c55e")).
			Bold(true)

	timerAmberStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#f59e0b")).
			Bold(true)

	timerRedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ef4444")).
			Bold(true)
)

// TimerDisplay renders a countdown / pomodoro timer with a label and MM:SS readout.
type TimerDisplay struct {
	remaining time.Duration
	label     string
}

// NewTimerDisplay creates a zeroed TimerDisplay.
func NewTimerDisplay() TimerDisplay {
	return TimerDisplay{}
}

// SetRemaining updates the remaining duration.
func (td *TimerDisplay) SetRemaining(d time.Duration) {
	td.remaining = d
}

// SetLabel updates the descriptive label shown above the timer.
func (td *TimerDisplay) SetLabel(label string) {
	td.label = label
}

// View renders the label and time display.
// Color: green when > 5 min, amber when 1–5 min, red when < 1 min.
func (td TimerDisplay) View() string {
	r := td.remaining
	if r < 0 {
		r = 0
	}
	minutes := int(r.Minutes())
	seconds := int(r.Seconds()) % 60
	clock := fmt.Sprintf("%02d:%02d", minutes, seconds)

	var timeStyle lipgloss.Style
	switch {
	case r >= 5*time.Minute:
		timeStyle = timerGreenStyle
	case r >= time.Minute:
		timeStyle = timerAmberStyle
	default:
		timeStyle = timerRedStyle
	}

	label := timerLabelStyle.Render(td.label)
	display := timeStyle.Render(clock)
	return label + "\n" + display
}
