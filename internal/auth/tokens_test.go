package auth

import (
	"testing"
	"time"
)

// ---- clampTTL tests ----

func TestClampTTL_Normal(t *testing.T) {
	d, clamped := clampTTL(3600)
	if clamped {
		t.Error("clampTTL(3600): expected clamped=false")
	}
	if d != 3600*time.Second {
		t.Errorf("clampTTL(3600) = %v, want %v", d, 3600*time.Second)
	}
}

func TestClampTTL_TooSmall(t *testing.T) {
	d, clamped := clampTTL(0)
	if !clamped {
		t.Error("clampTTL(0): expected clamped=true")
	}
	if d != time.Duration(minTokenTTL)*time.Second {
		t.Errorf("clampTTL(0) = %v, want %v", d, time.Duration(minTokenTTL)*time.Second)
	}
}

func TestClampTTL_Negative(t *testing.T) {
	d, clamped := clampTTL(-1)
	if !clamped {
		t.Error("clampTTL(-1): expected clamped=true")
	}
	if d != time.Duration(minTokenTTL)*time.Second {
		t.Errorf("clampTTL(-1) = %v, want %v", d, time.Duration(minTokenTTL)*time.Second)
	}
}

func TestClampTTL_TooLarge(t *testing.T) {
	d, clamped := clampTTL(99999)
	if !clamped {
		t.Error("clampTTL(99999): expected clamped=true")
	}
	if d != time.Duration(maxTokenTTL)*time.Second {
		t.Errorf("clampTTL(99999) = %v, want %v", d, time.Duration(maxTokenTTL)*time.Second)
	}
}

func TestClampTTL_Boundary_Min(t *testing.T) {
	d, clamped := clampTTL(minTokenTTL)
	if clamped {
		t.Errorf("clampTTL(%d): expected clamped=false (exact boundary)", minTokenTTL)
	}
	if d != time.Duration(minTokenTTL)*time.Second {
		t.Errorf("clampTTL(%d) = %v, want %v", minTokenTTL, d, time.Duration(minTokenTTL)*time.Second)
	}
}

func TestClampTTL_Boundary_Max(t *testing.T) {
	d, clamped := clampTTL(maxTokenTTL)
	if clamped {
		t.Errorf("clampTTL(%d): expected clamped=false (exact boundary)", maxTokenTTL)
	}
	if d != time.Duration(maxTokenTTL)*time.Second {
		t.Errorf("clampTTL(%d) = %v, want %v", maxTokenTTL, d, time.Duration(maxTokenTTL)*time.Second)
	}
}

// ---- TokenCache tests ----

func TestTokenCache_SetAndGet(t *testing.T) {
	tc := NewTokenCache(nil, nil)
	tc.SetAccessToken("my-token", 10*time.Minute)

	if tc.IsExpired() {
		t.Error("IsExpired() = true immediately after SetAccessToken with 10m TTL; want false")
	}
}

func TestTokenCache_Expired(t *testing.T) {
	tc := NewTokenCache(nil, nil)
	// Set a TTL so small that the internal expiry (TTL - 60s buffer) is in the past.
	// Using 1 nanosecond means expiresAt = now + 1ns - 60s, which is already past.
	tc.SetAccessToken("short-lived", 1*time.Nanosecond)

	// No sleep needed — the expiry is already in the past because the 60 s buffer
	// is subtracted from the 1 ns TTL, yielding a time well before now.
	if !tc.IsExpired() {
		t.Error("IsExpired() = false for effectively-zero TTL; want true")
	}
}

func TestTokenCache_AuthMethod(t *testing.T) {
	tc := NewTokenCache(nil, nil)
	tc.SetAuthMethod(AuthMethodAPIKey)

	got := tc.GetAuthMethod()
	if got != AuthMethodAPIKey {
		t.Errorf("GetAuthMethod() = %q, want %q", got, AuthMethodAPIKey)
	}
}
