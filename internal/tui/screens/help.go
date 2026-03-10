package screens

import (
	"strings"
)

// HelpOverlay renders the keybindings help panel as a string.
// It is displayed by the root TUI model as a centered overlay.
func HelpOverlay() string {
	inner := 43
	top := "╔" + strings.Repeat("═", inner) + "╗"
	headerText := padRight("KEYBINDINGS", inner-4) + "[?]"
	header := "║" + headerText + "║"
	sep := "╠" + strings.Repeat("═", inner) + "╣"
	bottom := "╚" + strings.Repeat("═", inner) + "╝"

	bindings := []string{
		"1-9        Switch tabs",
		"Tab/]/[    Cycle tabs",
		"j/k ↑/↓   Navigate",
		"Enter      Select/Open",
		":          Command palette",
		"/          Search",
		"n          New item",
		"d          Done/Delete",
		"g          Refresh",
		"?          Toggle help",
		"q          Quit",
	}

	lines := []string{top, header, sep}
	for _, b := range bindings {
		lines = append(lines, "║ "+padRight(b, inner-2)+" ║")
	}
	lines = append(lines, bottom)

	return strings.Join(lines, "\n")
}
