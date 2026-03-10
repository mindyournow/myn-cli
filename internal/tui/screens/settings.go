package screens

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/mindyournow/myn-cli/internal/app"
)

// derefBool safely dereferences a *bool, returning "false" if nil.
func derefBool(b *bool) string {
	if b == nil {
		return "false"
	}
	return fmt.Sprintf("%v", *b)
}

// settingItem represents a single config key/value row.
type settingItem struct {
	key   string
	value string
}

// SettingsScreen displays application configuration.
type SettingsScreen struct {
	app    *app.App
	width  int
	height int
	cursor int
	items  []settingItem
}

// NewSettingsScreen creates the Settings screen model.
func NewSettingsScreen(application *app.App) SettingsScreen {
	s := SettingsScreen{app: application}
	s.buildItems()
	return s
}

func (s *SettingsScreen) buildItems() {
	s.items = nil
	if s.app == nil || s.app.Config == nil {
		s.items = []settingItem{{key: "config", value: "(not loaded)"}}
		return
	}
	cfg := s.app.Config
	s.items = []settingItem{
		{key: "api.url", value: cfg.API.URL},
		{key: "api.timeout", value: cfg.API.Timeout},
		{key: "api.retries", value: fmt.Sprintf("%d", cfg.API.Retries)},
		{key: "auth.method", value: cfg.Auth.Method},
		{key: "auth.keyring", value: cfg.Auth.Keyring},
		{key: "display.color", value: cfg.Display.Color},
		{key: "display.date_format", value: cfg.Display.DateFormat},
		{key: "display.time_format", value: cfg.Display.TimeFormat},
		{key: "display.default_output", value: cfg.Display.DefaultOutput},
		{key: "tui.theme", value: cfg.TUI.Theme},
		{key: "tui.refresh_interval", value: cfg.TUI.RefreshInterval},
		{key: "tui.vim_keys", value: derefBool(cfg.TUI.VimKeys)},
		{key: "tui.mouse", value: derefBool(cfg.TUI.Mouse)},
		{key: "tui.animations", value: derefBool(cfg.TUI.Animations)},
		{key: "defaults.priority", value: cfg.Defaults.Priority},
		{key: "defaults.task_type", value: cfg.Defaults.TaskType},
		{key: "defaults.calendar_days", value: fmt.Sprintf("%d", cfg.Defaults.CalendarDays)},
		{key: "defaults.habit_schedule_days", value: fmt.Sprintf("%d", cfg.Defaults.HabitScheduleDays)},
	}
	if cfg.ConfigFile != "" {
		s.items = append([]settingItem{{key: "config_file", value: cfg.ConfigFile}}, s.items...)
	}
}

// Init implements tea.Model.
func (s SettingsScreen) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model.
func (s SettingsScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.width = msg.Width
		s.height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "j", "down":
			if s.cursor < len(s.items)-1 {
				s.cursor++
			}
		case "k", "up":
			if s.cursor > 0 {
				s.cursor--
			}
		}
	}
	return s, nil
}

// View implements tea.Model.
func (s SettingsScreen) View() string {
	title := titleStyle.Render("SETTINGS")

	var rows []string
	rows = append(rows, title)

	if len(s.items) == 0 {
		rows = append(rows, dimStyle.Render("  No configuration loaded."))
		return lipgloss.JoinVertical(lipgloss.Left, rows...)
	}

	for i, item := range s.items {
		row := fmt.Sprintf("%-36s  %s", item.key, item.value)
		if i == s.cursor {
			rows = append(rows, selectedItemStyle.Render("► "+row))
		} else {
			rows = append(rows, itemStyle.Render("  "+row))
		}
	}

	rows = append(rows, "")
	rows = append(rows, dimStyle.Render("  j/k navigate  (use 'mynow config set <key> <value>' to edit)"))

	return lipgloss.JoinVertical(lipgloss.Left, rows...)
}
