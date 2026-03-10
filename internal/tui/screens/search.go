package screens

import (
	"context"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/mindyournow/myn-cli/internal/api"
	"github.com/mindyournow/myn-cli/internal/app"
	"github.com/mindyournow/myn-cli/internal/tui/components"
)

type searchResultsMsg struct {
	results []api.SearchResult
	query   string
}
type searchErrMsg struct{ err error }

// SearchScreen is a unified search overlay.
type SearchScreen struct {
	app       *app.App
	width     int
	height    int
	input     components.Input
	results   []api.SearchResult
	cursor    int
	loading   bool
	err       error
	lastQuery string
}

// NewSearchScreen creates the Search screen model.
func NewSearchScreen(application *app.App) SearchScreen {
	inp := components.NewInput("search tasks, habits, chores...")
	inp.Focus()
	return SearchScreen{
		app:   application,
		input: inp,
	}
}

// Init implements tea.Model.
func (s SearchScreen) Init() tea.Cmd {
	return nil
}

func (s SearchScreen) doSearch(query string) tea.Cmd {
	return func() tea.Msg {
		if query == "" || s.app == nil {
			return searchResultsMsg{nil, ""}
		}
		ctx := context.Background()
		results, err := s.app.Client.Search(ctx, api.SearchParams{
			Query: query,
			Limit: 20,
		})
		if err != nil {
			return searchErrMsg{err}
		}
		return searchResultsMsg{results, query}
	}
}

// Update implements tea.Model.
func (s SearchScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case searchResultsMsg:
		s.results = msg.results
		s.loading = false
		s.cursor = 0
		s.lastQuery = msg.query
		return s, nil

	case searchErrMsg:
		s.err = msg.err
		s.loading = false
		return s, nil

	case tea.WindowSizeMsg:
		s.width = msg.Width
		s.height = msg.Height
		return s, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			s.results = nil
			s.cursor = 0
			s.input.SetValue("")
			return s, nil
		case "enter":
			// Item selection not yet implemented in TUI.
			return s, nil
		case "j", "down":
			if s.cursor < len(s.results)-1 {
				s.cursor++
			}
			return s, nil
		case "k", "up":
			if s.cursor > 0 {
				s.cursor--
			}
			return s, nil
		default:
			// Forward key to input, then check for query change.
			var inputCmd tea.Cmd
			s.input, inputCmd = s.input.Update(msg)
			val := s.input.Value()
			if val != s.lastQuery && len(val) >= 2 {
				s.loading = true
				return s, tea.Batch(inputCmd, s.doSearch(val))
			}
			return s, inputCmd
		}
	}

	// For all other messages, forward to input.
	var inputCmd tea.Cmd
	s.input, inputCmd = s.input.Update(msg)
	return s, inputCmd
}

// View implements tea.Model.
func (s SearchScreen) View() string {
	title := titleStyle.Render("SEARCH")

	var rows []string
	rows = append(rows, title)
	rows = append(rows, s.input.View())
	rows = append(rows, "")

	if s.loading {
		rows = append(rows, dimStyle.Render("  Searching..."))
	} else if s.err != nil {
		rows = append(rows, formatError(s.err))
	} else if len(s.results) == 0 && s.lastQuery != "" {
		rows = append(rows, dimStyle.Render("  No results found."))
	} else {
		for i, r := range s.results {
			typeLabel := r.Type
			row := fmt.Sprintf("%-8s  %s", typeLabel, r.Title)
			if i == s.cursor {
				rows = append(rows, selectedItemStyle.Render("► "+row))
			} else {
				rows = append(rows, itemStyle.Render("  "+row))
			}
		}
	}

	rows = append(rows, "")
	rows = append(rows, dimStyle.Render("  Type to search  j/k navigate  Enter=open  Esc=close"))

	return lipgloss.JoinVertical(lipgloss.Left, rows...)
}
