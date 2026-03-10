package components

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

// Tab represents a single tab entry.
type Tab struct {
	Label string
	Key   string
}

// TabBar renders the tab navigation bar at the top of the TUI.
type TabBar struct {
	Tabs      []Tab
	ActiveIdx int
	Width     int
}

var (
	tabBarStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#1a1a2e")).
			Width(0) // width set dynamically

	activeTabStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#ffffff")).
			Background(lipgloss.Color("#7c3aed")).
			Padding(0, 2)

	inactiveTabStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#9ca3af")).
				Background(lipgloss.Color("#1a1a2e")).
				Padding(0, 2)

	tabSeparatorStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#374151")).
				Background(lipgloss.Color("#1a1a2e"))
)

// DefaultTabs returns the standard tab layout.
func DefaultTabs() []Tab {
	return []Tab{
		{Label: "Now", Key: "1"},
		{Label: "Inbox", Key: "2"},
		{Label: "Tasks", Key: "3"},
		{Label: "Habits", Key: "4"},
		{Label: "Chores", Key: "5"},
		{Label: "Cal", Key: "6"},
		{Label: "Timers", Key: "7"},
		{Label: "Grocery", Key: "8"},
		{Label: "Settings", Key: "9"},
		{Label: "Compass", Key: ""},
		{Label: "Stats", Key: ""},
		{Label: "AI Chat", Key: ""},
		{Label: "Pomodoro", Key: ""},
	}
}

// NewTabBar creates a new TabBar with the default tabs.
func NewTabBar() TabBar {
	return TabBar{
		Tabs:      DefaultTabs(),
		ActiveIdx: 0,
	}
}

// SetWidth updates the tab bar width.
func (t *TabBar) SetWidth(w int) {
	t.Width = w
}

// View renders the tab bar.
func (t TabBar) View() string {
	rendered := make([]string, 0, len(t.Tabs))
	for i, tab := range t.Tabs {
		label := fmt.Sprintf("%s %s", tab.Key, tab.Label)
		if i == t.ActiveIdx {
			rendered = append(rendered, activeTabStyle.Render(label))
		} else {
			rendered = append(rendered, inactiveTabStyle.Render(label))
		}
		if i < len(t.Tabs)-1 {
			rendered = append(rendered, tabSeparatorStyle.Render("│"))
		}
	}

	bar := lipgloss.JoinHorizontal(lipgloss.Top, rendered...)
	// Pad to full width
	return tabBarStyle.Width(t.Width).Render(bar)
}
