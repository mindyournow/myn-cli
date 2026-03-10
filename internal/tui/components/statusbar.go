package components

import (
	"fmt"
	"time"

	"github.com/charmbracelet/lipgloss"
)

// StatusBar renders the bottom status bar.
type StatusBar struct {
	Width        int
	UserName     string
	InboxCount   int
	HabitsCount  int
	NotifCount   int
	TimerRunning bool
	TimerElapsed time.Duration
	Message      string // transient message shown in the middle
}

var (
	statusBarStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#1e3a5f")).
			Foreground(lipgloss.Color("#e2e8f0")).
			Padding(0, 1)

	statusBarAccentStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("#1e3a5f")).
				Foreground(lipgloss.Color("#60a5fa")).
				Bold(true)

	statusBarDimStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("#1e3a5f")).
				Foreground(lipgloss.Color("#94a3b8"))

	statusBarHelpStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("#1e3a5f")).
				Foreground(lipgloss.Color("#64748b"))
)

// NewStatusBar creates a new StatusBar.
func NewStatusBar() StatusBar {
	return StatusBar{}
}

// SetWidth updates the status bar width.
func (s *StatusBar) SetWidth(w int) {
	s.Width = w
}

// View renders the status bar.
func (s StatusBar) View() string {
	// Left section: username + inbox count
	left := ""
	if s.UserName != "" {
		left += statusBarAccentStyle.Render(s.UserName)
		left += statusBarStyle.Render("  ")
	}
	left += statusBarStyle.Render(fmt.Sprintf("● %d inbox", s.InboxCount))
	left += statusBarDimStyle.Render("  ")
	left += statusBarStyle.Render(fmt.Sprintf("%d habits", s.HabitsCount))

	// Middle section: timer or message
	middle := ""
	if s.TimerRunning {
		h := int(s.TimerElapsed.Hours())
		m := int(s.TimerElapsed.Minutes()) % 60
		middle = statusBarAccentStyle.Render(fmt.Sprintf("⏱ %02d:%02d", h, m))
	} else if s.Message != "" {
		middle = statusBarDimStyle.Render(s.Message)
	}

	// Right section: notifications + help
	right := ""
	if s.NotifCount > 0 {
		right += statusBarAccentStyle.Render(fmt.Sprintf("🔔 %d", s.NotifCount))
		right += statusBarStyle.Render("  ")
	}
	right += statusBarHelpStyle.Render("?=help  /:search  ::cmd")

	// Calculate widths for layout
	leftWidth := lipgloss.Width(left)
	middleWidth := lipgloss.Width(middle)
	rightWidth := lipgloss.Width(right)
	totalWidth := s.Width - 2 // account for padding

	if totalWidth <= 0 {
		return statusBarStyle.Width(s.Width).Render(left + "  " + middle + "  " + right)
	}

	// Pad between sections
	leftPad := (totalWidth - leftWidth - middleWidth - rightWidth) / 2
	if leftPad < 1 {
		leftPad = 1
	}
	rightPad := totalWidth - leftWidth - middleWidth - rightWidth - leftPad
	if rightPad < 1 {
		rightPad = 1
	}

	lp := statusBarStyle.Render(fmt.Sprintf("%*s", leftPad, ""))
	rp := statusBarStyle.Render(fmt.Sprintf("%*s", rightPad, ""))

	content := left + lp + middle + rp + right
	return statusBarStyle.Width(s.Width).Render(content)
}
