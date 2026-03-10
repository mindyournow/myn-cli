// Package util provides shared parsing utilities for the MYN CLI.
package util

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// ParseDate parses a date string per Spec Appendix C.
// Supported formats: today, tomorrow, yesterday, +3d, +1w, next week,
// monday..sunday, 2026-03-15, Mar 15, 3/15.
// Returns the parsed date at midnight in local time.
func ParseDate(s string) (time.Time, error) {
	s = strings.TrimSpace(strings.ToLower(s))
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	switch s {
	case "today":
		return today, nil
	case "tomorrow":
		return today.AddDate(0, 0, 1), nil
	case "yesterday":
		return today.AddDate(0, 0, -1), nil
	case "next week":
		// Monday of next week
		daysUntilMonday := int(time.Monday-today.Weekday()+7)%7 + 7
		if daysUntilMonday == 7 {
			daysUntilMonday = 7
		}
		return today.AddDate(0, 0, daysUntilMonday), nil
	}

	// Weekday names: "monday".."sunday"
	if d, ok := parseWeekday(s); ok {
		daysUntil := int(d-today.Weekday()+7) % 7
		if daysUntil == 0 {
			daysUntil = 7 // next occurrence (not today)
		}
		return today.AddDate(0, 0, daysUntil), nil
	}

	// +3d, +1w
	if strings.HasPrefix(s, "+") {
		return parseRelative(s[1:], today)
	}

	// ISO 8601: 2026-03-15
	if t, err := time.ParseInLocation("2006-01-02", s, now.Location()); err == nil {
		return t, nil
	}

	// "Mar 15" — case-insensitive
	for _, layout := range []string{"Jan 2", "Jan _2"} {
		orig := strings.Title(s) //nolint:staticcheck // Title is fine for short month abbrevs
		if t, err := time.ParseInLocation(layout, orig, now.Location()); err == nil {
			t = time.Date(now.Year(), t.Month(), t.Day(), 0, 0, 0, 0, now.Location())
			return t, nil
		}
	}

	// "3/15"
	parts := strings.Split(s, "/")
	if len(parts) == 2 {
		month, err1 := strconv.Atoi(parts[0])
		day, err2 := strconv.Atoi(parts[1])
		if err1 == nil && err2 == nil && month >= 1 && month <= 12 && day >= 1 && day <= 31 {
			return time.Date(now.Year(), time.Month(month), day, 0, 0, 0, 0, now.Location()), nil
		}
	}

	return time.Time{}, fmt.Errorf("cannot parse date %q — try: today, tomorrow, +3d, +1w, 2026-03-15, Mar 15, 3/15", s)
}

func parseWeekday(s string) (time.Weekday, bool) {
	days := map[string]time.Weekday{
		"sunday": time.Sunday, "monday": time.Monday, "tuesday": time.Tuesday,
		"wednesday": time.Wednesday, "thursday": time.Thursday, "friday": time.Friday,
		"saturday": time.Saturday,
	}
	if d, ok := days[s]; ok {
		return d, true
	}
	return 0, false
}

// parseRelative parses "3d", "1w" relative offsets.
func parseRelative(s string, from time.Time) (time.Time, error) {
	if len(s) < 2 {
		return time.Time{}, fmt.Errorf("invalid relative date %q", s)
	}
	unit := s[len(s)-1]
	n, err := strconv.Atoi(s[:len(s)-1])
	if err != nil || n <= 0 {
		return time.Time{}, fmt.Errorf("invalid relative date %q", s)
	}
	switch unit {
	case 'd':
		return from.AddDate(0, 0, n), nil
	case 'w':
		return from.AddDate(0, 0, n*7), nil
	case 'm':
		return from.AddDate(0, n, 0), nil
	case 'y':
		return from.AddDate(n, 0, 0), nil
	}
	return time.Time{}, fmt.Errorf("unknown date unit %q in %q", string(unit), s)
}
