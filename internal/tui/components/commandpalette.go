package components

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// CommandSelectedMsg is sent when the user selects a command.
type CommandSelectedMsg struct {
	Command string
}

// CommandPaletteDismissMsg is sent when the palette is dismissed without selection.
type CommandPaletteDismissMsg struct{}

// CommandEntry is a command in the palette.
type CommandEntry struct {
	Name        string
	Description string
}

// allCommands is the full list of available commands.
var allCommands = []CommandEntry{
	{Name: "goto now", Description: "Go to Now tab"},
	{Name: "goto inbox", Description: "Go to Inbox tab"},
	{Name: "goto tasks", Description: "Go to Tasks tab"},
	{Name: "goto habits", Description: "Go to Habits tab"},
	{Name: "goto chores", Description: "Go to Chores tab"},
	{Name: "goto calendar", Description: "Go to Calendar tab"},
	{Name: "goto timers", Description: "Go to Timers tab"},
	{Name: "goto grocery", Description: "Go to Grocery tab"},
	{Name: "goto settings", Description: "Go to Settings tab"},
	{Name: "goto compass", Description: "Go to Compass briefing"},
	{Name: "goto stats", Description: "Go to Stats & Achievements"},
	{Name: "ai chat", Description: "Open AI Chat with Kaia"},
	{Name: "pomodoro", Description: "Open Pomodoro Focus Mode"},
	{Name: "search", Description: "Open unified search"},
	{Name: "notifications", Description: "Open notifications"},
	{Name: "add task", Description: "Add a new task"},
	{Name: "add habit", Description: "Add a new habit"},
	{Name: "add chore", Description: "Add a new chore"},
	{Name: "timer", Description: "Start a timer"},
	{Name: "quit", Description: "Quit the application"},
}

// CommandPalette is the vim-style command overlay.
type CommandPalette struct {
	input     textinput.Model
	filtered  []CommandEntry
	cursor    int
	Width     int
	Height    int
}

var (
	paletteOverlayStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("#7c3aed")).
				Background(lipgloss.Color("#0f172a")).
				Padding(0, 1)

	paletteInputStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#e2e8f0")).
				Background(lipgloss.Color("#0f172a"))

	paletteItemStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#94a3b8")).
				Background(lipgloss.Color("#0f172a")).
				Padding(0, 1)

	paletteActiveItemStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#ffffff")).
				Background(lipgloss.Color("#4c1d95")).
				Bold(true).
				Padding(0, 1)

	paletteCmdStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#a78bfa")).
			Background(lipgloss.Color("#0f172a"))

	paletteDescStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#6b7280")).
				Background(lipgloss.Color("#0f172a"))
)

// NewCommandPalette creates a new CommandPalette.
func NewCommandPalette() CommandPalette {
	ti := textinput.New()
	ti.Placeholder = "type a command..."
	ti.Prompt = ": "
	ti.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#7c3aed"))
	ti.TextStyle = paletteInputStyle
	ti.Focus()

	return CommandPalette{
		input:    ti,
		filtered: allCommands,
		cursor:   0,
	}
}

// Reset clears the palette input and resets the filter.
func (c *CommandPalette) Reset() {
	c.input.Reset()
	c.filtered = allCommands
	c.cursor = 0
	c.input.Focus()
}

// SetSize updates the palette dimensions.
func (c *CommandPalette) SetSize(w, h int) {
	c.Width = w
	c.Height = h
}

// filter narrows the command list based on the input text.
func (c *CommandPalette) filter(query string) {
	if query == "" {
		c.filtered = allCommands
		c.cursor = 0
		return
	}
	q := strings.ToLower(query)
	result := make([]CommandEntry, 0, len(allCommands))
	for _, cmd := range allCommands {
		if strings.Contains(strings.ToLower(cmd.Name), q) ||
			strings.Contains(strings.ToLower(cmd.Description), q) {
			result = append(result, cmd)
		}
	}
	c.filtered = result
	if c.cursor >= len(c.filtered) {
		c.cursor = max(0, len(c.filtered)-1)
	}
}

// Update handles key events for the command palette.
func (c CommandPalette) Update(msg tea.Msg) (CommandPalette, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return c, func() tea.Msg { return CommandPaletteDismissMsg{} }
		case "enter":
			if len(c.filtered) > 0 {
				selected := c.filtered[c.cursor].Name
				return c, func() tea.Msg { return CommandSelectedMsg{Command: selected} }
			}
			return c, func() tea.Msg { return CommandPaletteDismissMsg{} }
		case "up", "ctrl+p":
			if c.cursor > 0 {
				c.cursor--
			}
			return c, nil
		case "down", "ctrl+n":
			if c.cursor < len(c.filtered)-1 {
				c.cursor++
			}
			return c, nil
		case "tab":
			// Tab completion: fill with current selection
			if len(c.filtered) > 0 {
				c.input.SetValue(c.filtered[c.cursor].Name)
				c.input.CursorEnd()
				c.filter(c.filtered[c.cursor].Name)
			}
			return c, nil
		}
	}

	prevVal := c.input.Value()
	c.input, cmd = c.input.Update(msg)
	newVal := c.input.Value()
	if newVal != prevVal {
		c.filter(newVal)
	}
	return c, cmd
}

// View renders the command palette overlay.
func (c CommandPalette) View() string {
	paletteWidth := c.Width - 10
	if paletteWidth < 40 {
		paletteWidth = 40
	}
	if paletteWidth > 80 {
		paletteWidth = 80
	}

	maxVisible := 8
	if maxVisible > len(c.filtered) {
		maxVisible = len(c.filtered)
	}

	// Input line
	inputLine := c.input.View()

	// Separator
	separator := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#374151")).
		Background(lipgloss.Color("#0f172a")).
		Render(strings.Repeat("─", paletteWidth))

	// Command list
	lines := []string{inputLine, separator}

	start := 0
	if c.cursor >= maxVisible {
		start = c.cursor - maxVisible + 1
	}
	end := start + maxVisible
	if end > len(c.filtered) {
		end = len(c.filtered)
	}

	for i := start; i < end; i++ {
		entry := c.filtered[i]
		cmdPart := paletteCmdStyle.Render(entry.Name)
		descPart := paletteDescStyle.Render("  " + entry.Description)
		row := cmdPart + descPart
		if i == c.cursor {
			row = paletteActiveItemStyle.Width(paletteWidth - 2).Render(entry.Name + "  " + entry.Description)
		} else {
			row = paletteItemStyle.Width(paletteWidth - 2).Render(row)
		}
		lines = append(lines, row)
	}

	if len(c.filtered) == 0 {
		lines = append(lines, paletteDescStyle.Render("  no matching commands"))
	}

	content := strings.Join(lines, "\n")
	palette := paletteOverlayStyle.Width(paletteWidth).Render(content)

	return palette
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
