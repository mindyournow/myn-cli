package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	streakFilledStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#22c55e"))

	streakEmptyStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#475569"))

	streakCountStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#94a3b8"))
)

// StreakBar visualises a 7-day habit history grid and current streak count.
type StreakBar struct {
	// History holds completion status for the last 7 days.
	// Index 0 is the oldest day, index 6 is the most recent.
	History       []bool
	CurrentStreak int
}

// NewStreakBar creates an empty StreakBar with a 7-day history slice.
func NewStreakBar() StreakBar {
	return StreakBar{
		History: make([]bool, 7),
	}
}

// View renders the 7-day grid followed by the streak count.
// Example output: "■ ■ □ ■ ■ ■ ■  streak: 4"
func (s StreakBar) View() string {
	cells := make([]string, 7)
	history := s.History
	if len(history) > 7 {
		history = history[len(history)-7:]
	}
	for i := 0; i < 7; i++ {
		var done bool
		if i < len(history) {
			done = history[i]
		}
		if done {
			cells[i] = streakFilledStyle.Render("■")
		} else {
			cells[i] = streakEmptyStyle.Render("□")
		}
	}
	grid := strings.Join(cells, " ")
	count := streakCountStyle.Render(fmt.Sprintf("  streak: %d", s.CurrentStreak))
	return grid + count
}
