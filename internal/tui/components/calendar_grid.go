package components

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

var (
	calHeaderStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#e2e8f0")).
			Background(lipgloss.Color("#0f172a")).
			Bold(true).
			Align(lipgloss.Center)

	calTodayStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#0f172a")).
			Background(lipgloss.Color("#7c3aed")).
			Bold(true).
			Align(lipgloss.Center)

	calDayStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#e2e8f0")).
			Background(lipgloss.Color("#0f172a")).
			Align(lipgloss.Center)

	calEventStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#60a5fa")).
			Background(lipgloss.Color("#0f172a"))

	calDimStyle2 = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#475569")).
			Background(lipgloss.Color("#0f172a"))

	calBorderStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("#475569")).
			Background(lipgloss.Color("#0f172a"))
)

// CalendarEvent represents a single event on the calendar.
type CalendarEvent struct {
	Title  string
	Date   string // "2006-01-02"
	AllDay bool
}

// CalendarGrid renders a 7-column week view calendar grid.
type CalendarGrid struct {
	weekStart time.Time
	events    []CalendarEvent
	width     int
	height    int
}

// NewCalendarGrid creates a new CalendarGrid.
func NewCalendarGrid() CalendarGrid {
	return CalendarGrid{}
}

// SetWeek sets the week starting date (should be Monday) and the events to display.
func (c *CalendarGrid) SetWeek(start time.Time, events []CalendarEvent) {
	c.weekStart = start
	c.events = events
}

// SetSize sets the available terminal dimensions.
func (c *CalendarGrid) SetSize(w, h int) {
	c.width = w
	c.height = h
}

// View renders the calendar grid.
func (c CalendarGrid) View() string {
	if c.width < 7 {
		return ""
	}

	colWidth := (c.width - 8) / 7
	if colWidth < 6 {
		colWidth = 6
	}

	today := time.Now().Format("2006-01-02")

	weekdays := []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"}

	// Build event map keyed by date string
	eventsByDate := make(map[string][]CalendarEvent)
	for _, ev := range c.events {
		eventsByDate[ev.Date] = append(eventsByDate[ev.Date], ev)
	}

	// Header row: day names
	headerCells := make([]string, 7)
	dateCells := make([]string, 7)
	days := make([]time.Time, 7)

	for i := 0; i < 7; i++ {
		day := c.weekStart.AddDate(0, 0, i)
		days[i] = day
		dateStr := day.Format("2006-01-02")
		dayNum := fmt.Sprintf("%d", day.Day())
		header := weekdays[i]

		if dateStr == today {
			headerCells[i] = calTodayStyle.Width(colWidth).Render(header)
			dateCells[i] = calTodayStyle.Width(colWidth).Render(dayNum)
		} else {
			headerCells[i] = calHeaderStyle.Width(colWidth).Render(header)
			dateCells[i] = calDayStyle.Width(colWidth).Render(dayNum)
		}
	}

	headerRow := strings.Join(headerCells, calDimStyle2.Render("│"))
	dateRow := strings.Join(dateCells, calDimStyle2.Render("│"))

	divider := calDimStyle2.Render(strings.Repeat("─", c.width))

	// Determine max events per day to know how many event rows to render
	maxEvents := 0
	for _, day := range days {
		dateStr := day.Format("2006-01-02")
		if n := len(eventsByDate[dateStr]); n > maxEvents {
			maxEvents = n
		}
	}
	if maxEvents > 5 {
		maxEvents = 5
	}

	// Build event rows
	var eventRows []string
	for row := 0; row < maxEvents; row++ {
		cells := make([]string, 7)
		for col, day := range days {
			dateStr := day.Format("2006-01-02")
			evs := eventsByDate[dateStr]
			if row < len(evs) {
				title := evs[row].Title
				if len(title) > colWidth-1 {
					title = title[:colWidth-1]
				}
				cells[col] = calEventStyle.Width(colWidth).Render(title)
			} else {
				cells[col] = calDayStyle.Width(colWidth).Render("")
			}
		}
		eventRows = append(eventRows, strings.Join(cells, calDimStyle2.Render("│")))
	}

	parts := []string{
		headerRow,
		divider,
		dateRow,
		divider,
	}
	parts = append(parts, eventRows...)

	return lipgloss.NewStyle().
		Background(lipgloss.Color("#0f172a")).
		Width(c.width).
		Render(strings.Join(parts, "\n"))
}
