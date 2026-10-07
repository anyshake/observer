package timesource

import (
	"testing"
	"time"
)

func TestMonotonicClockAdvances(t *testing.T) {
	start := Monotonic()
	if Monotonic() < start {
		t.Fatal("monotonic clock moved backwards")
	}
	now := MonotonicNow()
	if now.Before(time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)) {
		t.Fatalf("monotonic time = %v", now)
	}
}
