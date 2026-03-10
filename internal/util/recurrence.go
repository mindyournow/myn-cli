package util

import (
	"strings"
)

// ParseRecurrence converts a recurrence shortcut to an RRULE string per Spec Appendix E.
// If s is already an RRULE (starts with "FREQ="), it is returned unchanged.
func ParseRecurrence(s string) string {
	s = strings.TrimSpace(s)
	if strings.HasPrefix(s, "FREQ=") {
		return s
	}
	switch strings.ToLower(s) {
	case "daily":
		return "FREQ=DAILY"
	case "weekly":
		return "FREQ=WEEKLY"
	case "monthly":
		return "FREQ=MONTHLY"
	case "weekdays":
		return "FREQ=WEEKLY;BYDAY=MO,TU,WE,TH,FR"
	case "mwf":
		return "FREQ=WEEKLY;BYDAY=MO,WE,FR"
	case "tth", "tuth":
		return "FREQ=WEEKLY;BYDAY=TU,TH"
	}
	return s // return as-is; validation is the caller's responsibility
}

// IsValidRRULE returns true if the string looks like a valid RRULE.
func IsValidRRULE(s string) bool {
	return strings.HasPrefix(s, "FREQ=")
}
