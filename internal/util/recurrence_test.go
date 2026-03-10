package util

import "testing"

func TestParseRecurrence(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"daily", "FREQ=DAILY"},
		{"DAILY", "FREQ=DAILY"},
		{"weekly", "FREQ=WEEKLY"},
		{"monthly", "FREQ=MONTHLY"},
		{"weekdays", "FREQ=WEEKLY;BYDAY=MO,TU,WE,TH,FR"},
		{"MWF", "FREQ=WEEKLY;BYDAY=MO,WE,FR"},
		{"mwf", "FREQ=WEEKLY;BYDAY=MO,WE,FR"},
		{"TTh", "FREQ=WEEKLY;BYDAY=TU,TH"},
		{"tth", "FREQ=WEEKLY;BYDAY=TU,TH"},
		{"FREQ=DAILY;INTERVAL=2", "FREQ=DAILY;INTERVAL=2"}, // passthrough
	}
	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			got := ParseRecurrence(tc.input)
			if got != tc.want {
				t.Errorf("ParseRecurrence(%q) = %q, want %q", tc.input, got, tc.want)
			}
		})
	}
}
