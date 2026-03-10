package components

import "github.com/charmbracelet/lipgloss"

// PriorityColor returns the lipgloss color for a given priority zone.
func PriorityColor(priority string) lipgloss.Color {
	switch priority {
	case "CRITICAL":
		return lipgloss.Color("#ef4444")
	case "OPPORTUNITY_NOW":
		return lipgloss.Color("#f59e0b")
	case "OVER_THE_HORIZON":
		return lipgloss.Color("#60a5fa")
	case "PARKING_LOT":
		return lipgloss.Color("#94a3b8")
	default:
		return lipgloss.Color("#94a3b8")
	}
}

// PriorityBadge returns a colored symbol for a given priority zone.
func PriorityBadge(priority string) string {
	switch priority {
	case "CRITICAL":
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#ef4444")).Render("●")
	case "OPPORTUNITY_NOW":
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#f59e0b")).Render("○")
	case "OVER_THE_HORIZON":
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#60a5fa")).Render("◌")
	case "PARKING_LOT":
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#94a3b8")).Render("·")
	default:
		return " "
	}
}
