package util

import "testing"

func TestParseDuration(t *testing.T) {
	tests := []struct {
		input   string
		want    int
		wantErr bool
	}{
		{"30s", 30, false},
		{"5m", 300, false},
		{"25m", 1500, false},
		{"1h", 3600, false},
		{"1h30m", 5400, false},
		{"2h", 7200, false},
		{"1h30m10s", 5410, false},
		{"", 0, true},
		{"abc", 0, true},
		{"5x", 0, true},
		{"1", 0, true},
	}
	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			got, err := ParseDuration(tc.input)
			if tc.wantErr {
				if err == nil {
					t.Errorf("ParseDuration(%q) expected error", tc.input)
				}
				return
			}
			if err != nil {
				t.Fatalf("ParseDuration(%q) error = %v", tc.input, err)
			}
			if got != tc.want {
				t.Errorf("ParseDuration(%q) = %d, want %d", tc.input, got, tc.want)
			}
		})
	}
}

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		seconds int
		want    string
	}{
		{0, "0s"},
		{30, "30s"},
		{60, "1m"},
		{90, "1m30s"},
		{3600, "1h"},
		{5400, "1h30m"},
		{3661, "1h1m1s"},
	}
	for _, tc := range tests {
		got := FormatDuration(tc.seconds)
		if got != tc.want {
			t.Errorf("FormatDuration(%d) = %q, want %q", tc.seconds, got, tc.want)
		}
	}
}
