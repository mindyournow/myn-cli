package util

import (
	"testing"
	"time"
)

func TestParseDate(t *testing.T) {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	tests := []struct {
		input   string
		wantErr bool
		check   func(t *testing.T, got time.Time)
	}{
		{
			input: "today",
			check: func(t *testing.T, got time.Time) {
				if !got.Equal(today) {
					t.Errorf("today: got %v, want %v", got, today)
				}
			},
		},
		{
			input: "tomorrow",
			check: func(t *testing.T, got time.Time) {
				want := today.AddDate(0, 0, 1)
				if !got.Equal(want) {
					t.Errorf("tomorrow: got %v, want %v", got, want)
				}
			},
		},
		{
			input: "yesterday",
			check: func(t *testing.T, got time.Time) {
				want := today.AddDate(0, 0, -1)
				if !got.Equal(want) {
					t.Errorf("yesterday: got %v, want %v", got, want)
				}
			},
		},
		{
			input: "+3d",
			check: func(t *testing.T, got time.Time) {
				want := today.AddDate(0, 0, 3)
				if !got.Equal(want) {
					t.Errorf("+3d: got %v, want %v", got, want)
				}
			},
		},
		{
			input: "+1w",
			check: func(t *testing.T, got time.Time) {
				want := today.AddDate(0, 0, 7)
				if !got.Equal(want) {
					t.Errorf("+1w: got %v, want %v", got, want)
				}
			},
		},
		{
			input: "2026-03-15",
			check: func(t *testing.T, got time.Time) {
				if got.Year() != 2026 || got.Month() != 3 || got.Day() != 15 {
					t.Errorf("ISO date: got %v", got)
				}
			},
		},
		{
			input: "3/15",
			check: func(t *testing.T, got time.Time) {
				if got.Month() != 3 || got.Day() != 15 {
					t.Errorf("3/15: got %v", got)
				}
			},
		},
		{
			input:   "not-a-date",
			wantErr: true,
		},
		{
			input:   "",
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			got, err := ParseDate(tc.input)
			if tc.wantErr {
				if err == nil {
					t.Errorf("ParseDate(%q) expected error, got nil", tc.input)
				}
				return
			}
			if err != nil {
				t.Fatalf("ParseDate(%q) error = %v", tc.input, err)
			}
			tc.check(t, got)
		})
	}
}
