package output

// ANSI color codes (reset-safe, respects NoColor).

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorGray   = "\033[90m"
	colorGreen  = "\033[32m"
	colorCyan   = "\033[36m"
	colorBold   = "\033[1m"
)

// Priority zone symbols and colors per Spec Appendix B.
const (
	SymbolCritical    = "●"
	SymbolOpportunity = "○"
	SymbolHorizon     = "◌"
	SymbolParking     = "·"
	SymbolDone        = "✓"
	SymbolHabit       = "◆"
	SymbolInbox       = "?"
)

// Priority maps API values to display info.
type priorityInfo struct {
	Symbol string
	Color  string
	Label  string
}

var priorityMap = map[string]priorityInfo{
	"CRITICAL":         {SymbolCritical, colorRed, "Critical"},
	"OPPORTUNITY_NOW":  {SymbolOpportunity, colorYellow, "Opportunity"},
	"OVER_THE_HORIZON": {SymbolHorizon, colorBlue, "Horizon"},
	"PARKING_LOT":      {SymbolParking, colorGray, "Parking"},
	"":                 {SymbolInbox, "", "Inbox"},
}

// PrioritySymbol returns the symbol for a priority value.
func PrioritySymbol(priority string) string {
	if info, ok := priorityMap[priority]; ok {
		return info.Symbol
	}
	return SymbolInbox
}

// ColorString wraps s in the given ANSI color, or returns s unchanged if noColor.
func ColorString(s, color string, noColor bool) string {
	if noColor || color == "" {
		return s
	}
	return color + s + colorReset
}

// PriorityColored returns the priority symbol colored appropriately.
func PriorityColored(priority string, noColor bool) string {
	info, ok := priorityMap[priority]
	if !ok {
		info = priorityMap[""]
	}
	return ColorString(info.Symbol, info.Color, noColor)
}

// Bold wraps s in bold ANSI if !noColor.
func Bold(s string, noColor bool) string {
	return ColorString(s, colorBold, noColor)
}

// Green wraps s in green ANSI if !noColor.
func Green(s string, noColor bool) string {
	return ColorString(s, colorGreen, noColor)
}

// Red wraps s in red ANSI if !noColor.
func Red(s string, noColor bool) string {
	return ColorString(s, colorRed, noColor)
}

// Yellow wraps s in yellow ANSI if !noColor.
func Yellow(s string, noColor bool) string {
	return ColorString(s, colorYellow, noColor)
}

// Cyan wraps s in cyan ANSI if !noColor.
func Cyan(s string, noColor bool) string {
	return ColorString(s, colorCyan, noColor)
}

// Gray wraps s in gray ANSI if !noColor.
func Gray(s string, noColor bool) string {
	return ColorString(s, colorGray, noColor)
}
