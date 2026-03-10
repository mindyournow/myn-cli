package errors

import (
	"errors"
	"fmt"
	"testing"
)

func TestExitCode(t *testing.T) {
	tests := []struct {
		err  error
		want int
	}{
		{nil, ExitSuccess},
		{fmt.Errorf("generic error"), ExitGeneral},
		{New(ExitAuth, "auth error"), ExitAuth},
		{New(ExitNetwork, "net error"), ExitNetwork},
		{New(ExitRateLimit, "rate limit"), ExitRateLimit},
		{Wrap(ExitAPI, "api err", fmt.Errorf("cause")), ExitAPI},
	}
	for _, tc := range tests {
		got := ExitCode(tc.err)
		if got != tc.want {
			t.Errorf("ExitCode(%v) = %d, want %d", tc.err, got, tc.want)
		}
	}
}

func TestWrap_Unwrap(t *testing.T) {
	cause := fmt.Errorf("original")
	wrapped := Wrap(ExitNetwork, "network failed", cause)
	if !errors.Is(wrapped, cause) {
		t.Error("Wrap should allow errors.Is to find cause")
	}
	if wrapped.Code != ExitNetwork {
		t.Errorf("code = %d, want %d", wrapped.Code, ExitNetwork)
	}
}

func TestJSONError(t *testing.T) {
	err := Auth("not authenticated", nil).WithHint("Run 'mynow login'")
	b, jsonErr := JSONError(err)
	if jsonErr != nil {
		t.Fatalf("JSONError() error = %v", jsonErr)
	}
	s := string(b)
	if !containsAll(s, `"error"`, `"code"`, `"hint"`) {
		t.Errorf("JSONError output missing fields: %s", s)
	}
}

func TestHelpers(t *testing.T) {
	if !IsAuth(Auth("x", nil)) {
		t.Error("IsAuth should be true for Auth error")
	}
	if !IsNetwork(Network(nil)) {
		t.Error("IsNetwork should be true for Network error")
	}
	if !IsRateLimit(RateLimit()) {
		t.Error("IsRateLimit should be true for RateLimit error")
	}
}

func containsAll(s string, subs ...string) bool {
	for _, sub := range subs {
		found := false
		for i := 0; i <= len(s)-len(sub); i++ {
			if s[i:i+len(sub)] == sub {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}
