package parse

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// DurationParser handles natural language duration parsing.
type DurationParser struct{}

// NewDurationParser creates a new duration parser.
func NewDurationParser() *DurationParser {
	return &DurationParser{}
}

// Parse parses a natural language duration string.
// Examples: "30m", "2h", "1d", "1h30m", "2 weeks"
func (p *DurationParser) Parse(input string) (time.Duration, error) {
	input = strings.ToLower(strings.TrimSpace(input))

	// Try compact format (e.g., "30m", "2h30m")
	if d, ok := p.parseCompact(input); ok {
		return d, nil
	}

	// Try word format (e.g., "30 minutes", "2 hours 30 minutes")
	if d, ok := p.parseWords(input); ok {
		return d, nil
	}

	// Try short format (e.g., "2d", "1w")
	if d, ok := p.parseShort(input); ok {
		return d, nil
	}

	return 0, fmt.Errorf("unable to parse duration: %s", input)
}

// parseCompact handles Go-style duration strings like "2h30m".
func (p *DurationParser) parseCompact(input string) (time.Duration, bool) {
	// Check if it looks like a compact duration
	if matched, _ := regexp.MatchString(`^(\d+[hms])+\d*[hms]?$`, input); !matched {
		return 0, false
	}

	d, err := time.ParseDuration(input)
	if err != nil {
		return 0, false
	}
	return d, true
}

// parseWords handles word-based duration strings like "2 hours 30 minutes".
func (p *DurationParser) parseWords(input string) (time.Duration, bool) {
	total := time.Duration(0)

	// Pattern: "X hours", "X minutes", "X seconds"
	hourPattern := regexp.MustCompile(`(\d+)\s*(?:hour|hours|hr|hrs)`)
	minPattern := regexp.MustCompile(`(\d+)\s*(?:minute|minutes|min|mins|m)`)
	secPattern := regexp.MustCompile(`(\d+)\s*(?:second|seconds|sec|secs|s)`)
	dayPattern := regexp.MustCompile(`(\d+)\s*(?:day|days|d)`)
	weekPattern := regexp.MustCompile(`(\d+)\s*(?:week|weeks|wk|wks|w)`)

	if matches := hourPattern.FindStringSubmatch(input); matches != nil {
		h, _ := strconv.Atoi(matches[1])
		total += time.Duration(h) * time.Hour
	}

	if matches := minPattern.FindStringSubmatch(input); matches != nil {
		m, _ := strconv.Atoi(matches[1])
		total += time.Duration(m) * time.Minute
	}

	if matches := secPattern.FindStringSubmatch(input); matches != nil {
		s, _ := strconv.Atoi(matches[1])
		total += time.Duration(s) * time.Second
	}

	if matches := dayPattern.FindStringSubmatch(input); matches != nil {
		d, _ := strconv.Atoi(matches[1])
		total += time.Duration(d) * 24 * time.Hour
	}

	if matches := weekPattern.FindStringSubmatch(input); matches != nil {
		w, _ := strconv.Atoi(matches[1])
		total += time.Duration(w) * 7 * 24 * time.Hour
	}

	if total > 0 {
		return total, true
	}
	return 0, false
}

// parseShort handles short duration strings like "2d", "1w".
func (p *DurationParser) parseShort(input string) (time.Duration, bool) {
	// Single unit patterns
	patterns := []struct {
		pattern *regexp.Regexp
		unit    time.Duration
	}{
		{regexp.MustCompile(`^(\d+)d$`), 24 * time.Hour},
		{regexp.MustCompile(`^(\d+)w$`), 7 * 24 * time.Hour},
		{regexp.MustCompile(`^(\d+)mo$`), 30 * 24 * time.Hour},
		{regexp.MustCompile(`^(\d+)y$`), 365 * 24 * time.Hour},
	}

	for _, p := range patterns {
		if matches := p.pattern.FindStringSubmatch(input); matches != nil {
			n, _ := strconv.Atoi(matches[1])
			return time.Duration(n) * p.unit, true
		}
	}

	return 0, false
}

// ParseDuration is a convenience function to parse a duration string.
func ParseDuration(input string) (time.Duration, error) {
	return NewDurationParser().Parse(input)
}

// FormatDuration returns a human-readable duration string.
func FormatDuration(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%ds", int(d.Seconds()))
	}
	if d < time.Hour {
		return fmt.Sprintf("%dm", int(d.Minutes()))
	}
	if d < 24*time.Hour {
		return fmt.Sprintf("%dh", int(d.Hours()))
	}
	if d < 7*24*time.Hour {
		return fmt.Sprintf("%dd", int(d.Hours()/24))
	}
	return fmt.Sprintf("%dw", int(d.Hours()/24/7))
}

// FormatDurationLong returns a detailed human-readable duration string.
func FormatDurationLong(d time.Duration) string {
	days := int(d.Hours() / 24)
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60

	parts := []string{}
	if days > 0 {
		if days == 1 {
			parts = append(parts, "1 day")
		} else {
			parts = append(parts, fmt.Sprintf("%d days", days))
		}
	}
	if hours > 0 {
		if hours == 1 {
			parts = append(parts, "1 hour")
		} else {
			parts = append(parts, fmt.Sprintf("%d hours", hours))
		}
	}
	if minutes > 0 {
		if minutes == 1 {
			parts = append(parts, "1 minute")
		} else {
			parts = append(parts, fmt.Sprintf("%d minutes", minutes))
		}
	}

	if len(parts) == 0 {
		return "0 minutes"
	}

	return strings.Join(parts, ", ")
}
