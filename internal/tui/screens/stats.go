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

type statsLoadedMsg struct {
	streaks      []api.Streak
	achievements []api.Achievement
}
type statsErrMsg struct{ err error }

// StatsScreen shows streaks and achievements.
type StatsScreen struct {
	app    *app.App
	width  int
	height int
	cursor int
	tab    int // 0=streaks, 1=achievements

	streaks      []api.Streak
	achievements []api.Achievement
	loading      bool
	err          error
}

// NewStatsScreen creates the Stats screen model.
func NewStatsScreen(application *app.App) StatsScreen {
	return StatsScreen{app: application, loading: true}
}

// Init implements tea.Model.
func (s StatsScreen) Init() tea.Cmd {
	return s.loadData()
}

func (s StatsScreen) loadData() tea.Cmd {
	return func() tea.Msg {
		if s.app == nil {
			return statsLoadedMsg{}
		}
		ctx := context.Background()
		streaks, err := s.app.Client.GetStreaks(ctx)
		if err != nil {
			return statsErrMsg{err}
		}
		achievements, err := s.app.Client.GetAchievements(ctx, false)
		if err != nil {
			return statsErrMsg{err}
		}
		return statsLoadedMsg{streaks, achievements}
	}
}

// Update implements tea.Model.
func (s StatsScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case statsLoadedMsg:
		s.streaks = msg.streaks
		s.achievements = msg.achievements
		s.loading = false

	case statsErrMsg:
		s.err = msg.err
		s.loading = false

	case tea.WindowSizeMsg:
		s.width = msg.Width
		s.height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "tab":
			s.tab = (s.tab + 1) % 2
			s.cursor = 0
		case "j", "down":
			if s.tab == 1 && s.cursor < len(s.achievements)-1 {
				s.cursor++
			}
		case "k", "up":
			if s.tab == 1 && s.cursor > 0 {
				s.cursor--
			}
		case "g":
			s.loading = true
			return s, s.loadData()
		}
	}
	return s, nil
}

// View implements tea.Model.
func (s StatsScreen) View() string {
	title := titleStyle.Render("STATS & ACHIEVEMENTS")

	if s.loading {
		return lipgloss.JoinVertical(lipgloss.Left, title, dimStyle.Render("  Loading..."))
	}
	if s.err != nil {
		return lipgloss.JoinVertical(lipgloss.Left, title, formatError(s.err))
	}

	// Tab header
	var tab0Style, tab1Style lipgloss.Style
	if s.tab == 0 {
		tab0Style = selectedItemStyle
	} else {
		tab0Style = dimStyle
	}
	if s.tab == 1 {
		tab1Style = selectedItemStyle
	} else {
		tab1Style = dimStyle
	}
	tabHeader := tab0Style.Render("[1] Streaks") + "  " + tab1Style.Render("[2] Achievements")

	var rows []string
	rows = append(rows, title, tabHeader, "")

	if s.tab == 0 {
		// Streaks tab
		if len(s.streaks) == 0 {
			rows = append(rows, dimStyle.Render("  No active streaks."))
		} else {
			rows = append(rows, sectionStyle.Render("ACTIVE STREAKS"))
			bars := make([]components.BarEntry, len(s.streaks))
			for i, st := range s.streaks {
				label := st.Title
				if label == "" {
					label = st.Type
				}
				bars[i] = components.BarEntry{Label: truncate(label, 20), Value: float64(st.Count)}
			}
			chart := components.NewBarChart("Streak Days", bars)
			chartWidth := s.width - 4
			if chartWidth < 20 {
				chartWidth = 20
			}
			chartHeight := s.height - 10
			if chartHeight < 5 {
				chartHeight = 5
			}
			chart.SetSize(chartWidth, chartHeight)
			rows = append(rows, chart.View())
		}
	} else {
		// Achievements tab
		if len(s.achievements) == 0 {
			rows = append(rows, dimStyle.Render("  No achievements found."))
		} else {
			rows = append(rows, sectionStyle.Render("ACHIEVEMENTS"))
			for i, a := range s.achievements {
				lockIcon := "🔒"
				if a.IsUnlocked {
					lockIcon = "✓"
				}
				row := fmt.Sprintf("%s %-30s", lockIcon, truncate(a.Title, 30))
				if a.Total > 0 {
					pct := float64(a.Progress) / float64(a.Total)
					bar := components.NewProgressBar("", pct*float64(a.Total), float64(a.Total))
					barWidth := s.width - 40
					if barWidth < 10 {
						barWidth = 10
					}
					bar.SetSize(barWidth)
					row += "  " + bar.View()
				}
				if i == s.cursor {
					rows = append(rows, selectedItemStyle.Render("► "+row))
				} else if a.IsUnlocked {
					rows = append(rows, successStyle.Render("  "+row))
				} else {
					rows = append(rows, dimStyle.Render("  "+row))
				}
			}
		}
	}

	rows = append(rows, "")
	rows = append(rows, dimStyle.Render("  Tab=switch tabs  j/k=navigate  g=refresh"))

	return lipgloss.JoinVertical(lipgloss.Left, rows...)
}
