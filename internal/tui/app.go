package tui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/mindyournow/myn-cli/internal/app"
	"github.com/mindyournow/myn-cli/internal/tui/components"
	"github.com/mindyournow/myn-cli/internal/tui/screens"
)

// Screen identifies which tab/screen is active.
type Screen int

const (
	ScreenNow Screen = iota
	ScreenInbox
	ScreenTasks
	ScreenHabits
	ScreenChores
	ScreenCalendar
	ScreenTimers
	ScreenGrocery
	ScreenSettings
	ScreenCompass
	ScreenStats
	ScreenAIChat
	ScreenPomodoro
	screenCount = 13
)

// Model is the root Bubble Tea model for the TUI.
type Model struct {
	width, height      int
	activeScreen       Screen
	tabBar             components.TabBar
	statusBar          components.StatusBar
	commandPalette     components.CommandPalette
	showCommandPalette bool
	showHelp           bool
	showSearch         bool
	showNotifications  bool
	searchScreen       screens.SearchScreen
	notifScreen        screens.NotificationsScreen
	application        *app.App
	screens            [screenCount]tea.Model
}

var contentStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#e2e8f0"))

// New creates a new root TUI Model.
func New(application *app.App) *Model {
	m := &Model{
		activeScreen:   ScreenNow,
		tabBar:         components.NewTabBar(),
		statusBar:      components.NewStatusBar(),
		commandPalette: components.NewCommandPalette(),
		application:    application,
	}

	// Populate screen slots
	m.screens[ScreenNow] = screens.NewNowScreen(application)
	m.screens[ScreenInbox] = screens.NewInboxScreen(application)
	m.screens[ScreenTasks] = screens.NewTasksScreen(application)
	m.screens[ScreenHabits] = screens.NewHabitsScreen(application)
	m.screens[ScreenChores] = screens.NewChoresScreen(application)
	m.screens[ScreenCalendar] = screens.NewCalendarScreen(application)
	m.screens[ScreenTimers] = screens.NewTimersScreen(application)
	m.screens[ScreenGrocery] = screens.NewGroceryScreen(application)
	m.screens[ScreenSettings] = screens.NewSettingsScreen(application)
	m.screens[ScreenCompass] = screens.NewCompassScreen(application)
	m.screens[ScreenStats] = screens.NewStatsScreen(application)
	m.screens[ScreenAIChat] = screens.NewAIChatScreen(application)
	m.screens[ScreenPomodoro] = screens.NewPomodoroScreen(application)

	// Overlay screens (not in main tab slots)
	m.searchScreen = screens.NewSearchScreen(application)
	m.notifScreen = screens.NewNotificationsScreen(application)

	return m
}

// Init implements tea.Model.
func (m Model) Init() tea.Cmd {
	var cmds []tea.Cmd
	for _, s := range m.screens {
		cmds = append(cmds, s.Init())
	}
	cmds = append(cmds, m.searchScreen.Init(), m.notifScreen.Init())
	return tea.Batch(cmds...)
}

// Update implements tea.Model.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.tabBar.SetWidth(msg.Width)
		m.statusBar.SetWidth(msg.Width)
		m.commandPalette.SetSize(msg.Width, msg.Height)
		// Forward window size to all screens so they can size themselves
		for i, s := range m.screens {
			updated, cmd := s.Update(msg)
			m.screens[i] = updated
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
		}

	case components.CommandSelectedMsg:
		m.showCommandPalette = false
		m.commandPalette.Reset()
		cmd := m.handleCommand(msg.Command)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}

	case components.CommandPaletteDismissMsg:
		m.showCommandPalette = false
		m.commandPalette.Reset()

	case tea.KeyMsg:
		// Command palette captures all keys when open
		if m.showCommandPalette {
			var paletteCmd tea.Cmd
			m.commandPalette, paletteCmd = m.commandPalette.Update(msg)
			cmds = append(cmds, paletteCmd)
			return m, tea.Batch(cmds...)
		}

		// Help overlay: any key dismisses it
		if m.showHelp {
			m.showHelp = false
			return m, nil
		}

		// Search overlay captures all keys when open
		if m.showSearch {
			if msg.String() == "esc" {
				m.showSearch = false
				return m, nil
			}
			updated, cmd := m.searchScreen.Update(msg)
			if s, ok := updated.(screens.SearchScreen); ok {
				m.searchScreen = s
			}
			cmds = append(cmds, cmd)
			return m, tea.Batch(cmds...)
		}

		// Notifications overlay captures all keys when open
		if m.showNotifications {
			if msg.String() == "esc" {
				m.showNotifications = false
				return m, nil
			}
			updated, cmd := m.notifScreen.Update(msg)
			if s, ok := updated.(screens.NotificationsScreen); ok {
				m.notifScreen = s
			}
			cmds = append(cmds, cmd)
			return m, tea.Batch(cmds...)
		}

		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit

		case "q":
			return m, tea.Quit

		case ":":
			m.showCommandPalette = true
			m.commandPalette.Reset()
			return m, nil

		case "?":
			m.showHelp = true
			return m, nil

		case "/":
			m.showSearch = true
			return m, nil

		case "N":
			m.showNotifications = true
			return m, nil

		case "1":
			m.setScreen(ScreenNow)
		case "2":
			m.setScreen(ScreenInbox)
		case "3":
			m.setScreen(ScreenTasks)
		case "4":
			m.setScreen(ScreenHabits)
		case "5":
			m.setScreen(ScreenChores)
		case "6":
			m.setScreen(ScreenCalendar)
		case "7":
			m.setScreen(ScreenTimers)
		case "8":
			m.setScreen(ScreenGrocery)
		case "9":
			m.setScreen(ScreenSettings)

		case "tab":
			m.setScreen(Screen((int(m.activeScreen) + 1) % screenCount))

		case "shift+tab":
			m.setScreen(Screen((int(m.activeScreen) + screenCount - 1) % screenCount))

		default:
			// Delegate remaining key events to the active screen
			updated, cmd := m.screens[m.activeScreen].Update(msg)
			m.screens[m.activeScreen] = updated
			cmds = append(cmds, cmd)
		}

	default:
		// Forward all other messages to the active screen
		updated, cmd := m.screens[m.activeScreen].Update(msg)
		m.screens[m.activeScreen] = updated
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

// setScreen switches to the given screen and updates the tab bar.
func (m *Model) setScreen(s Screen) {
	m.activeScreen = s
	m.tabBar.ActiveIdx = int(s)
}

// handleCommand executes a command palette selection.
func (m *Model) handleCommand(cmd string) tea.Cmd {
	switch cmd {
	case "goto now":
		m.setScreen(ScreenNow)
	case "goto inbox":
		m.setScreen(ScreenInbox)
	case "goto tasks":
		m.setScreen(ScreenTasks)
	case "goto habits":
		m.setScreen(ScreenHabits)
	case "goto chores":
		m.setScreen(ScreenChores)
	case "goto calendar":
		m.setScreen(ScreenCalendar)
	case "goto timers":
		m.setScreen(ScreenTimers)
	case "goto grocery":
		m.setScreen(ScreenGrocery)
	case "goto settings":
		m.setScreen(ScreenSettings)
	case "goto compass":
		m.setScreen(ScreenCompass)
	case "goto stats":
		m.setScreen(ScreenStats)
	case "goto ai", "ai chat":
		m.setScreen(ScreenAIChat)
	case "pomodoro":
		m.setScreen(ScreenPomodoro)
	case "search":
		m.showSearch = true
	case "notifications":
		m.showNotifications = true
	case "quit":
		return tea.Quit
	}
	return nil
}

// View implements tea.Model.
func (m Model) View() string {
	if m.width == 0 || m.height == 0 {
		return "Initializing..."
	}

	tabBarView := m.tabBar.View()
	statusBarView := m.statusBar.View()

	tabBarHeight := lipgloss.Height(tabBarView)
	statusBarHeight := lipgloss.Height(statusBarView)
	contentHeight := m.height - tabBarHeight - statusBarHeight
	if contentHeight < 0 {
		contentHeight = 0
	}

	screenView := m.screens[m.activeScreen].View()
	content := contentStyle.
		Width(m.width).
		Height(contentHeight).
		Render(screenView)

	view := strings.Join([]string{tabBarView, content, statusBarView}, "\n")

	// Overlay help if active
	if m.showHelp {
		helpView := screens.HelpOverlay()
		view = overlayCenter(view, helpView, m.width, m.height)
	}

	// Overlay command palette if active
	if m.showCommandPalette {
		paletteView := m.commandPalette.View()
		view = overlayCenter(view, paletteView, m.width, m.height)
	}

	// Overlay search if active
	if m.showSearch {
		searchView := m.searchScreen.View()
		view = overlayCenter(view, searchView, m.width, m.height)
	}

	// Overlay notifications if active
	if m.showNotifications {
		notifView := m.notifScreen.View()
		view = overlayCenter(view, notifView, m.width, m.height)
	}

	return view
}

// overlayCenter places an overlay string centered over the base string.
func overlayCenter(base, overlay string, w, h int) string {
	overlayLines := strings.Split(overlay, "\n")
	overlayH := len(overlayLines)
	overlayW := 0
	for _, l := range overlayLines {
		if ww := lipgloss.Width(l); ww > overlayW {
			overlayW = ww
		}
	}

	baseLines := strings.Split(base, "\n")
	startRow := (h - overlayH) / 2
	startCol := (w - overlayW) / 2
	if startRow < 0 {
		startRow = 0
	}
	if startCol < 0 {
		startCol = 0
	}

	for i, ol := range overlayLines {
		row := startRow + i
		if row >= len(baseLines) {
			baseLines = append(baseLines, strings.Repeat(" ", w))
		}
		bl := baseLines[row]
		// Pad base line if shorter than needed
		blRunes := []rune(bl)
		for len(blRunes) < startCol+overlayW {
			blRunes = append(blRunes, ' ')
		}
		// Overlay the palette content
		olRunes := []rune(ol)
		for j, r := range olRunes {
			pos := startCol + j
			if pos < len(blRunes) {
				blRunes[pos] = r
			}
		}
		baseLines[row] = string(blRunes)
	}

	return strings.Join(baseLines, "\n")
}

// Run starts the TUI program with alt-screen and mouse support.
func Run(application *app.App) error {
	m := New(application)
	p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion())
	_, err := p.Run()
	return err
}
