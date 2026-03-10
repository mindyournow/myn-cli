package screens

import (
	"context"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/mindyournow/myn-cli/internal/api"
	"github.com/mindyournow/myn-cli/internal/app"
)

type groceryLoadedMsg struct{ items []api.GroceryItem }
type groceryErrMsg struct{ err error }

// GroceryScreen shows the household grocery list.
type GroceryScreen struct {
	app    *app.App
	width  int
	height int
	cursor int

	items   []api.GroceryItem
	loading bool
	err     error
}

// NewGroceryScreen creates the Grocery screen model.
func NewGroceryScreen(application *app.App) GroceryScreen {
	return GroceryScreen{app: application, loading: true}
}

// Init implements tea.Model.
func (s GroceryScreen) Init() tea.Cmd {
	return s.loadData()
}

func (s GroceryScreen) loadData() tea.Cmd {
	return func() tea.Msg {
		if s.app == nil {
			return groceryLoadedMsg{}
		}
		// The grocery API requires a householdID; attempt to fetch
		// by using an empty string (the API may return the user's default household list).
		// If unavailable, the error is shown.
		ctx := context.Background()
		householdID := ""
		if s.app.Config != nil {
			// Household ID is not currently in config; use empty string to attempt default.
		}
		items, err := s.app.Client.ListGroceryItems(ctx, householdID)
		if err != nil {
			return groceryErrMsg{err}
		}
		return groceryLoadedMsg{items}
	}
}

// Update implements tea.Model.
func (s GroceryScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case groceryLoadedMsg:
		s.items = msg.items
		s.loading = false

	case groceryErrMsg:
		s.err = msg.err
		s.loading = false

	case tea.WindowSizeMsg:
		s.width = msg.Width
		s.height = msg.Height

	case tea.KeyMsg:
		unchecked := s.uncheckedItems()
		switch msg.String() {
		case "j", "down":
			if s.cursor < len(unchecked)-1 {
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

func (s GroceryScreen) uncheckedItems() []api.GroceryItem {
	var items []api.GroceryItem
	for _, item := range s.items {
		if !item.IsChecked {
			items = append(items, item)
		}
	}
	return items
}

func (s GroceryScreen) checkedItems() []api.GroceryItem {
	var items []api.GroceryItem
	for _, item := range s.items {
		if item.IsChecked {
			items = append(items, item)
		}
	}
	return items
}

// View implements tea.Model.
func (s GroceryScreen) View() string {
	title := titleStyle.Render("GROCERY LIST")

	if s.loading {
		return lipgloss.JoinVertical(lipgloss.Left, title, dimStyle.Render("  Loading..."))
	}
	if s.err != nil {
		return lipgloss.JoinVertical(lipgloss.Left, title, formatError(s.err))
	}

	unchecked := s.uncheckedItems()
	checked := s.checkedItems()

	var rows []string
	rows = append(rows, title)

	if len(unchecked) > 0 {
		rows = append(rows, sectionStyle.Render("UNCHECKED"))
		for i, item := range unchecked {
			rows = append(rows, renderGroceryRow(item, i == s.cursor))
		}
	}

	if len(checked) > 0 {
		rows = append(rows, sectionStyle.Render("CHECKED"))
		for _, item := range checked {
			rows = append(rows, renderGroceryRow(item, false))
		}
	}

	if len(s.items) == 0 {
		rows = append(rows, dimStyle.Render("  Grocery list is empty."))
	}

	rows = append(rows, "")
	rows = append(rows, dimStyle.Render("  j/k navigate  x=check/uncheck  d=delete  n=add  c=clear checked  g=refresh"))

	return lipgloss.JoinVertical(lipgloss.Left, rows...)
}

func renderGroceryRow(item api.GroceryItem, selected bool) string {
	qty := ""
	if item.Quantity > 0 {
		if item.Unit != "" {
			qty = fmt.Sprintf("%.4g%s", item.Quantity, item.Unit)
		} else {
			qty = fmt.Sprintf("%.4g", item.Quantity)
		}
	}

	cat := item.Category
	if cat == "" {
		cat = "—"
	}

	prefix := "  "
	if item.IsChecked {
		prefix = "✓ "
	}

	row := fmt.Sprintf("%s%-28s  %-10s  %s",
		prefix,
		truncate(item.Name, 28),
		truncate(qty, 10),
		truncate(cat, 16),
	)

	if item.IsChecked {
		return dimStyle.Render("    " + row)
	}
	if selected {
		return selectedItemStyle.Render("► " + row)
	}
	return itemStyle.Render("  " + row)
}
