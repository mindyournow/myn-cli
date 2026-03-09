package parse

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// DateParser handles natural language date parsing.
type DateParser struct {
	reference time.Time
}

// NewDateParser creates a new date parser with the given reference time.
func NewDateParser(reference time.Time) *DateParser {
	return &DateParser{reference: reference}
}

// Parse parses a natural language date string.
// Returns the parsed date and a boolean indicating if parsing was successful.
func (p *DateParser) Parse(input string) (time.Time, error) {
	input = strings.ToLower(strings.TrimSpace(input))

	// Try special keywords
	if t, ok := p.parseSpecial(input); ok {
		return t, nil
	}

	// Try relative dates
	if t, ok := p.parseRelative(input); ok {
		return t, nil
	}

	// Try absolute dates
	if t, ok := p.parseAbsolute(input); ok {
		return t, nil
	}

	return time.Time{}, fmt.Errorf("unable to parse date: %s", input)
}

// parseSpecial handles special keywords like "today", "tomorrow", "yesterday".
func (p *DateParser) parseSpecial(input string) (time.Time, bool) {
	switch input {
	case "today", "now":
		return p.reference, true
	case "tomorrow", "tmrw":
		return p.reference.AddDate(0, 0, 1), true
	case "yesterday":
		return p.reference.AddDate(0, 0, -1), true
	case "next week":
		return p.reference.AddDate(0, 0, 7), true
	case "last week":
		return p.reference.AddDate(0, 0, -7), true
	case "next month":
		return p.reference.AddDate(0, 1, 0), true
	case "last month":
		return p.reference.AddDate(0, -1, 0), true
	default:
		return time.Time{}, false
	}
}

// parseRelative handles relative dates like "in 3 days", "2 weeks ago".
func (p *DateParser) parseRelative(input string) (time.Time, bool) {
	// Pattern: "in X days/weeks/months"
	inPattern := regexp.MustCompile(`^in\s+(\d+)\s+(day|days|week|weeks|month|months|year|years)$`)
	if matches := inPattern.FindStringSubmatch(input); matches != nil {
		n, _ := strconv.Atoi(matches[1])
		unit := matches[2]
		switch {
		case strings.HasPrefix(unit, "day"):
			return p.reference.AddDate(0, 0, n), true
		case strings.HasPrefix(unit, "week"):
			return p.reference.AddDate(0, 0, n*7), true
		case strings.HasPrefix(unit, "month"):
			return p.reference.AddDate(0, n, 0), true
		case strings.HasPrefix(unit, "year"):
			return p.reference.AddDate(n, 0, 0), true
		}
	}

	// Pattern: "X days/weeks ago"
	agoPattern := regexp.MustCompile(`^(\d+)\s+(day|days|week|weeks|month|months|year|years)\s+ago$`)
	if matches := agoPattern.FindStringSubmatch(input); matches != nil {
		n, _ := strconv.Atoi(matches[1])
		unit := matches[2]
		switch {
		case strings.HasPrefix(unit, "day"):
			return p.reference.AddDate(0, 0, -n), true
		case strings.HasPrefix(unit, "week"):
			return p.reference.AddDate(0, 0, -n*7), true
		case strings.HasPrefix(unit, "month"):
			return p.reference.AddDate(0, -n, 0), true
		case strings.HasPrefix(unit, "year"):
			return p.reference.AddDate(-n, 0, 0), true
		}
	}

	return time.Time{}, false
}

// parseAbsolute handles absolute dates like "2024-03-15", "Mar 15", "15th".
func (p *DateParser) parseAbsolute(input string) (time.Time, bool) {
	formats := []string{
		"2006-01-02",
		"2006-01-2",
		"2006/01/02",
		"Jan 2",
		"Jan 2, 2006",
		"January 2",
		"January 2, 2006",
		"2 Jan",
		"2 Jan 2006",
		"2 January",
		"2 January 2006",
		"Mon",
		"Monday",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, input); err == nil {
			// If no year specified, use reference year
			if t.Year() == 0 {
				t = t.AddDate(p.reference.Year(), 0, 0)
			}
			return t, true
		}
	}

	return time.Time{}, false
}

// ParseDate is a convenience function to parse a date string.
func ParseDate(input string) (time.Time, error) {
	return NewDateParser(time.Now()).Parse(input)
}

// FormatRelative returns a relative time string like "2 hours ago" or "in 3 days".
func FormatRelative(t time.Time) string {
	now := time.Now()
	diff := t.Sub(now)

	if diff < 0 {
		diff = -diff
		switch {
		case diff < time.Minute:
			return "just now"
		case diff < time.Hour:
			return fmt.Sprintf("%d minutes ago", int(diff.Minutes()))
		case diff < 24*time.Hour:
			return fmt.Sprintf("%d hours ago", int(diff.Hours()))
		case diff < 7*24*time.Hour:
			return fmt.Sprintf("%d days ago", int(diff.Hours()/24))
		case diff < 30*24*time.Hour:
			return fmt.Sprintf("%d weeks ago", int(diff.Hours()/24/7))
		default:
			return fmt.Sprintf("%d months ago", int(diff.Hours()/24/30))
		}
	}

	switch {
	case diff < time.Minute:
		return "now"
	case diff < time.Hour:
		return fmt.Sprintf("in %d minutes", int(diff.Minutes()))
	case diff < 24*time.Hour:
		return fmt.Sprintf("in %d hours", int(diff.Hours()))
	case diff < 7*24*time.Hour:
		return fmt.Sprintf("in %d days", int(diff.Hours()/24))
	case diff < 30*24*time.Hour:
		return fmt.Sprintf("in %d weeks", int(diff.Hours()/24/7))
	default:
		return fmt.Sprintf("in %d months", int(diff.Hours()/24/30))
	}
}
