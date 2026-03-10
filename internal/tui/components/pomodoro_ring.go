package components

import (
	"fmt"
	"math"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	pomWorkStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ef4444")).
			Background(lipgloss.Color("#0f172a"))

	pomShortBreakStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#22c55e")).
				Background(lipgloss.Color("#0f172a"))

	pomLongBreakStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#60a5fa")).
				Background(lipgloss.Color("#0f172a"))

	pomDimRingStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#475569")).
			Background(lipgloss.Color("#0f172a"))

	pomLabelStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#e2e8f0")).
			Background(lipgloss.Color("#0f172a")).
			Bold(true)

	pomBgStyle2 = lipgloss.NewStyle().
			Background(lipgloss.Color("#0f172a"))
)

// PomodoroRing renders a text-art progress ring for a Pomodoro timer.
type PomodoroRing struct {
	progress float64 // 0.0 - 1.0
	phase    string  // "work", "short_break", "long_break"
	width    int
}

// NewPomodoroRing creates a new PomodoroRing.
func NewPomodoroRing() PomodoroRing {
	return PomodoroRing{phase: "work"}
}

// SetProgress sets the current progress fraction (0.0–1.0).
func (p *PomodoroRing) SetProgress(fraction float64) {
	if fraction < 0 {
		fraction = 0
	}
	if fraction > 1 {
		fraction = 1
	}
	p.progress = fraction
}

// SetPhase sets the current phase label.
func (p *PomodoroRing) SetPhase(phase string) {
	p.phase = phase
}

// SetSize sets the available width.
func (p *PomodoroRing) SetSize(w int) {
	p.width = w
}

// phaseStyle returns the lipgloss style for the current phase.
func (p PomodoroRing) phaseStyle() lipgloss.Style {
	switch p.phase {
	case "short_break":
		return pomShortBreakStyle
	case "long_break":
		return pomLongBreakStyle
	default:
		return pomWorkStyle
	}
}

// phaseLabel returns a short display label for the phase.
func (p PomodoroRing) phaseLabel() string {
	switch p.phase {
	case "short_break":
		return "Short Break"
	case "long_break":
		return "Long Break"
	default:
		return "Focus"
	}
}

// View renders the pomodoro ring (or a progress bar fallback for narrow widths).
func (p PomodoroRing) View() string {
	w := p.width
	if w < 1 {
		w = 40
	}

	// Narrow fallback: render a simple progress bar
	if w < 20 {
		pb := NewProgressBar(p.phaseLabel(), p.progress, 1.0)
		pb.SetSize(w)
		return pb.View()
	}

	style := p.phaseStyle()

	// Ring parameters
	// We use a grid of ~21 cols x 11 rows. Each cell represents an angle.
	// We'll use a circle formula and mark positions as filled/empty/center.
	const (
		rows    = 11
		cols    = 21
		centerR = 5   // rows/2
		centerC = 10  // cols/2
		radiusR = 4.5 // vertical radius (rows)
		radiusC = 9.5 // horizontal radius (cols, wider because terminal chars are taller)
	)

	// Characters used
	const (
		filled  = "█"
		partial = "▓"
		empty   = "░"
		space   = " "
	)

	// Total arc points: we'll sample 360 degrees and mark each cell
	// Build a 2D grid
	grid := make([][]string, rows)
	for i := range grid {
		grid[i] = make([]string, cols)
		for j := range grid[i] {
			grid[i][j] = space
		}
	}

	// Draw full ring (empty chars) then overlay filled arc
	steps := 360
	for step := 0; step < steps; step++ {
		// Start from top (-90 degrees), go clockwise
		angleDeg := float64(step) - 90.0
		angleRad := angleDeg * math.Pi / 180.0
		r := centerR + radiusR*math.Sin(angleRad)
		c := centerC + radiusC*math.Cos(angleRad)
		ri := int(math.Round(r))
		ci := int(math.Round(c))
		if ri >= 0 && ri < rows && ci >= 0 && ci < cols {
			fraction := float64(step) / float64(steps)
			if fraction <= p.progress {
				grid[ri][ci] = filled
			} else {
				if grid[ri][ci] != filled {
					grid[ri][ci] = empty
				}
			}
		}
	}

	// Build center labels
	pctStr := fmt.Sprintf("%d%%", int(math.Round(p.progress*100)))
	phaseStr := p.phaseLabel()

	// Build center label rows
	centerRow := rows / 2
	labelRow := centerRow - 1
	pctRow := centerRow

	// Render rows with center text overlay
	finalRows := make([]string, rows)
	for ri := 0; ri < rows; ri++ {
		var rowStr strings.Builder
		for ci := 0; ci < cols; ci++ {
			ch := grid[ri][ci]
			// For center rows, replace chars 6-14 with centered text
			if ri == labelRow && ci == 6 {
				label := fmt.Sprintf("%-9s", phaseStr)
				if len(label) > 9 {
					label = label[:9]
				}
				rowStr.WriteString(pomLabelStyle.Render(label))
				continue
			}
			if ri == labelRow && ci > 6 && ci <= 14 {
				continue // already written
			}
			if ri == pctRow && ci == 6 {
				pctLabel := fmt.Sprintf("%-9s", pctStr)
				if len(pctLabel) > 9 {
					pctLabel = pctLabel[:9]
				}
				rowStr.WriteString(style.Render(pctLabel))
				continue
			}
			if ri == pctRow && ci > 6 && ci <= 14 {
				continue
			}

			switch ch {
			case filled:
				rowStr.WriteString(style.Render(ch))
			case empty:
				rowStr.WriteString(pomDimRingStyle.Render(ch))
			default:
				rowStr.WriteString(pomBgStyle2.Render(ch))
			}
		}
		finalRows[ri] = rowStr.String()
	}

	// Center the ring horizontally
	ringWidth := cols // visual char width (lipgloss width of a single rendered char row)
	var output []string
	leftPad := (w - ringWidth) / 2
	if leftPad < 0 {
		leftPad = 0
	}
	padding := pomBgStyle2.Render(strings.Repeat(" ", leftPad))
	for _, row := range finalRows {
		output = append(output, padding+row)
	}

	return pomBgStyle2.Render(strings.Join(output, "\n"))
}
