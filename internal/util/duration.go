package util

import (
	"fmt"
	"strconv"
	"strings"
)

// ParseDuration parses a human-friendly duration per Spec Appendix D.
// Supported: 30s, 5m, 25m, 1h, 1h30m, 2h → seconds.
func ParseDuration(s string) (int, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, fmt.Errorf("empty duration")
	}

	total := 0
	rest := s
	for rest != "" {
		// parse digits
		i := 0
		for i < len(rest) && rest[i] >= '0' && rest[i] <= '9' {
			i++
		}
		if i == 0 {
			return 0, fmt.Errorf("invalid duration %q: expected number", s)
		}
		n, err := strconv.Atoi(rest[:i])
		if err != nil {
			return 0, fmt.Errorf("invalid duration %q", s)
		}
		rest = rest[i:]

		if rest == "" {
			return 0, fmt.Errorf("invalid duration %q: missing unit (s/m/h)", s)
		}

		unit := rest[0]
		rest = rest[1:]
		switch unit {
		case 's':
			total += n
		case 'm':
			total += n * 60
		case 'h':
			total += n * 3600
		default:
			return 0, fmt.Errorf("invalid duration unit %q in %q", string(unit), s)
		}
	}
	return total, nil
}

// FormatDuration formats seconds as a human-friendly duration string.
func FormatDuration(seconds int) string {
	if seconds <= 0 {
		return "0s"
	}
	h := seconds / 3600
	m := (seconds % 3600) / 60
	s := seconds % 60
	switch {
	case h > 0 && m > 0 && s > 0:
		return fmt.Sprintf("%dh%dm%ds", h, m, s)
	case h > 0 && m > 0:
		return fmt.Sprintf("%dh%dm", h, m)
	case h > 0:
		return fmt.Sprintf("%dh", h)
	case m > 0 && s > 0:
		return fmt.Sprintf("%dm%ds", m, s)
	case m > 0:
		return fmt.Sprintf("%dm", m)
	default:
		return fmt.Sprintf("%ds", s)
	}
}
