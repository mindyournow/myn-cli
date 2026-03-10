package output

import "testing"

func TestPrioritySymbol(t *testing.T) {
	tests := []struct {
		priority string
		want     string
	}{
		{"CRITICAL", SymbolCritical},
		{"OPPORTUNITY_NOW", SymbolOpportunity},
		{"OVER_THE_HORIZON", SymbolHorizon},
		{"PARKING_LOT", SymbolParking},
		{"", SymbolInbox},
		{"UNKNOWN", SymbolInbox},
	}
	for _, tc := range tests {
		got := PrioritySymbol(tc.priority)
		if got != tc.want {
			t.Errorf("PrioritySymbol(%q) = %q, want %q", tc.priority, got, tc.want)
		}
	}
}

func TestColorString_NoColor(t *testing.T) {
	s := ColorString("hello", colorRed, true)
	if s != "hello" {
		t.Errorf("noColor: got %q, want %q", s, "hello")
	}
}

func TestColorString_WithColor(t *testing.T) {
	s := ColorString("hello", colorRed, false)
	if s == "hello" {
		t.Error("should be wrapped with ANSI codes")
	}
	if s != colorRed+"hello"+colorReset {
		t.Errorf("unexpected color wrap: %q", s)
	}
}

func TestStripANSI(t *testing.T) {
	s := "\033[31mhello\033[0m world"
	got := stripANSI(s)
	if got != "hello world" {
		t.Errorf("stripANSI = %q, want %q", got, "hello world")
	}
}
