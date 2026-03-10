package screens

import (
	"context"
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/mindyournow/myn-cli/internal/api"
	"github.com/mindyournow/myn-cli/internal/app"
)

type choresLoadedMsg struct{ chores []api.ChoreInstance }
type choresErrMsg struct{ err error }

// ChoresScreen shows today's chore instances.
type ChoresScreen struct {
	app    *app.App
	width  int
	height int
	cursor int

	chores  []api.ChoreInstance
	loading bool
	err     error
}

// NewChoresScreen creates the Chores screen model.
func NewChoresScreen(application *app.App) ChoresScreen {
	return ChoresScreen{app: application, loading: true}
}

// Init implements tea.Model.
func (s ChoresScreen) Init() tea.Cmd {
	return s.loadData()
}

func (s ChoresScreen) loadData() tea.Cmd {
	return func() tea.Msg {
		if s.app == nil {
			return choresLoadedMsg{}
		}
		ctx := context.Background()
		today := time.Now().Format("2006-01-02")
		chores, err := s.app.Client.ListTodayChores(ctx, today, "", "")
		if err != nil {
			return choresErrMsg{err}
		}
		return choresLoadedMsg{chores}
	}
}

// Update implements tea.Model.
func (s ChoresScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case choresLoadedMsg:
		s.chores = msg.chores
		s.loading = false

	case choresErrMsg:
		s.err = msg.err
		s.loading = false

	case tea.WindowSizeMsg:
		s.width = msg.Width
		s.height = msg.Height

	case tea.KeyMsg:
		pending := s.pendingChores()
		switch msg.String() {
		case "j", "down":
			if s.cursor < len(pending)-1 {
				s.cursor++
			}
		case "k", "up":
			if s.cursor > 0 {
				s.cursor--
			}
		case "g":
			s.loading = true
			return s, s.loadData()
		}
	}
	return s, nil
}

func (s ChoresScreen) pendingChores() []api.ChoreInstance {
	var pending []api.ChoreInstance
	for _, c := range s.chores {
		if !c.IsCompleted {
			pending = append(pending, c)
		}
	}
	return pending
}

func (s ChoresScreen) completedChores() []api.ChoreInstance {
	var done []api.ChoreInstance
	for _, c := range s.chores {
		if c.IsCompleted {
			done = append(done, c)
		}
	}
	return done
}

// View implements tea.Model.
func (s ChoresScreen) View() string {
	title := titleStyle.Render("CHORES")

	if s.loading {
		return lipgloss.JoinVertical(lipgloss.Left, title, dimStyle.Render("  Loading..."))
	}
	if s.err != nil {
		return lipgloss.JoinVertical(lipgloss.Left, title, formatError(s.err))
	}

	pending := s.pendingChores()
	completed := s.completedChores()

	var rows []string
	rows = append(rows, title)

	if len(s.chores) == 0 {
		rows = append(rows, dimStyle.Render("  No chores scheduled for today."))
	} else {
		rows = append(rows, sectionStyle.Render("TODAY'S CHORES"))
		for i, c := range pending {
			rows = append(rows, renderChoreRow(c, i == s.cursor))
		}
		for _, c := range completed {
			rows = append(rows, renderChoreRow(c, false))
		}
	}

	rows = append(rows, "")
	rows = append(rows, dimStyle.Render("  j/k navigate  Enter=done  g=refresh"))

	return lipgloss.JoinVertical(lipgloss.Left, rows...)
}

func renderChoreRow(c api.ChoreInstance, selected bool) string {
	check := "[ ]"
	if c.IsCompleted {
		check = "[✓]"
	}

	assigned := c.AssignedTo
	if assigned == "" {
		assigned = "—"
	}

	due := "today"
	if c.ScheduledDate != "" {
		due = c.ScheduledDate
	}

	row := fmt.Sprintf("%-28s  Assigned: %-12s  Due: %-10s  %s",
		truncate(c.Title, 28),
		truncate(assigned, 12),
		truncate(due, 10),
		check,
	)

	if c.IsCompleted {
		return dimStyle.Render("    " + row)
	}
	if selected {
		return selectedItemStyle.Render("► " + row)
	}
	return itemStyle.Render("  " + row)
}
