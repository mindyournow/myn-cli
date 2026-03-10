package tui

import "github.com/charmbracelet/bubbles/key"

// KeyMap holds all global key bindings for the TUI.
type KeyMap struct {
	// Tab switching
	Tab1     key.Binding
	Tab2     key.Binding
	Tab3     key.Binding
	Tab4     key.Binding
	Tab5     key.Binding
	Tab6     key.Binding
	Tab7     key.Binding
	Tab8     key.Binding
	Tab9     key.Binding
	NextTab  key.Binding
	PrevTab  key.Binding

	// Navigation
	Up     key.Binding
	Down   key.Binding
	Enter  key.Binding
	Back   key.Binding

	// Actions
	New     key.Binding
	Done    key.Binding
	Refresh key.Binding

	// Overlays
	Help           key.Binding
	Search         key.Binding
	CommandPalette key.Binding

	// Quit
	Quit key.Binding
}

// DefaultKeyMap returns the default key bindings.
func DefaultKeyMap() KeyMap {
	return KeyMap{
		Tab1: key.NewBinding(
			key.WithKeys("1"),
			key.WithHelp("1", "Now"),
		),
		Tab2: key.NewBinding(
			key.WithKeys("2"),
			key.WithHelp("2", "Inbox"),
		),
		Tab3: key.NewBinding(
			key.WithKeys("3"),
			key.WithHelp("3", "Tasks"),
		),
		Tab4: key.NewBinding(
			key.WithKeys("4"),
			key.WithHelp("4", "Habits"),
		),
		Tab5: key.NewBinding(
			key.WithKeys("5"),
			key.WithHelp("5", "Chores"),
		),
		Tab6: key.NewBinding(
			key.WithKeys("6"),
			key.WithHelp("6", "Cal"),
		),
		Tab7: key.NewBinding(
			key.WithKeys("7"),
			key.WithHelp("7", "Timers"),
		),
		Tab8: key.NewBinding(
			key.WithKeys("8"),
			key.WithHelp("8", "Grocery"),
		),
		Tab9: key.NewBinding(
			key.WithKeys("9"),
			key.WithHelp("9", "Settings"),
		),
		NextTab: key.NewBinding(
			key.WithKeys("tab", "]"),
			key.WithHelp("tab/]", "next tab"),
		),
		PrevTab: key.NewBinding(
			key.WithKeys("shift+tab", "["),
			key.WithHelp("shift+tab/[", "prev tab"),
		),
		Up: key.NewBinding(
			key.WithKeys("k", "up"),
			key.WithHelp("k/↑", "up"),
		),
		Down: key.NewBinding(
			key.WithKeys("j", "down"),
			key.WithHelp("j/↓", "down"),
		),
		Enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "select"),
		),
		Back: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "back"),
		),
		New: key.NewBinding(
			key.WithKeys("n"),
			key.WithHelp("n", "new item"),
		),
		Done: key.NewBinding(
			key.WithKeys("d"),
			key.WithHelp("d", "done/delete"),
		),
		Refresh: key.NewBinding(
			key.WithKeys("g"),
			key.WithHelp("g", "refresh"),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "help"),
		),
		Search: key.NewBinding(
			key.WithKeys("/"),
			key.WithHelp("/", "search"),
		),
		CommandPalette: key.NewBinding(
			key.WithKeys(":"),
			key.WithHelp(":", "command palette"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
	}
}

// ShortHelp returns a compact key help list.
func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit, k.CommandPalette, k.Search}
}

// FullHelp returns the full key help list.
func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Tab1, k.Tab2, k.Tab3, k.Tab4, k.Tab5},
		{k.Tab6, k.Tab7, k.Tab8, k.Tab9},
		{k.NextTab, k.PrevTab},
		{k.Up, k.Down, k.Enter, k.Back},
		{k.New, k.Done, k.Refresh},
		{k.Help, k.Search, k.CommandPalette, k.Quit},
	}
}
