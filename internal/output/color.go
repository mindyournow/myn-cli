package output

import (
	"fmt"
	"os"
	"strings"
)

// ANSI color codes
const (
	reset   = "\033[0m"
	red     = "\033[31m"
	green   = "\033[32m"
	yellow  = "\033[33m"
	blue    = "\033[34m"
	magenta = "\033[35m"
	cyan    = "\033[36m"
	gray    = "\033[90m"
	bold    = "\033[1m"
	dim     = "\033[2m"
)

// ColorEnabled determines if colors should be used.
var ColorEnabled = !isNoColor()

func isNoColor() bool {
	// Check NO_COLOR environment variable (standard)
	if os.Getenv("NO_COLOR") != "" {
		return true
	}
	// Check if not a terminal
	if fi, _ := os.Stdout.Stat(); fi != nil {
		return (fi.Mode() & os.ModeCharDevice) == 0
	}
	return false
}

// DisableColors disables all color output.
func DisableColors() {
	ColorEnabled = false
}

// EnableColors enables color output.
func EnableColors() {
	ColorEnabled = true
}

// colorize wraps text with color codes if enabled.
func colorize(text, color string) string {
	if !ColorEnabled {
		return text
	}
	return color + text + reset
}

// Red returns red colored text.
func Red(text string) string {
	return colorize(text, red)
}

// Green returns green colored text.
func Green(text string) string {
	return colorize(text, green)
}

// Yellow returns yellow colored text.
func Yellow(text string) string {
	return colorize(text, yellow)
}

// Blue returns blue colored text.
func Blue(text string) string {
	return colorize(text, blue)
}

// Magenta returns magenta colored text.
func Magenta(text string) string {
	return colorize(text, magenta)
}

// Cyan returns cyan colored text.
func Cyan(text string) string {
	return colorize(text, cyan)
}

// Gray returns gray colored text.
func Gray(text string) string {
	return colorize(text, gray)
}

// Bold returns bold text.
func Bold(text string) string {
	return colorize(text, bold)
}

// Dim returns dimmed text.
func Dim(text string) string {
	return colorize(text, dim)
}

// Strip removes all ANSI color codes from text.
func Strip(text string) string {
	// Simple regex-free implementation
	result := text
	for _, code := range []string{
		reset, red, green, yellow, blue, magenta, cyan, gray,
		bold, dim,
	} {
		result = strings.ReplaceAll(result, code, "")
	}
	return result
}

// SprintColored returns a formatted string with color support.
func SprintColored(format string, a ...any) string {
	return fmt.Sprintf(format, a...)
}

// PriorityColor returns the appropriate color for a priority level.
func PriorityColor(priority string) string {
	switch strings.ToUpper(priority) {
	case "CRITICAL_NOW", "CRITICAL":
		return Red(priority)
	case "OPPORTUNITY_NOW", "OPPORTUNITY":
		return Yellow(priority)
	case "OVER_THE_HORIZON", "HORIZON":
		return Blue(priority)
	case "PARKING_LOT", "PARKING":
		return Gray(priority)
	default:
		return priority
	}
}

// Bool returns a colored yes/no string.
func Bool(val bool) string {
	if val {
		return Green("yes")
	}
	return Red("no")
}
