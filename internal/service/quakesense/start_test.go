package quakesense

import (
	"testing"
	"time"
)

func TestQuakeSenseNotificationMessage(t *testing.T) {
	eventTime := time.Date(2026, time.October, 3, 0, 43, 42, 0, time.FixedZone("UTC+8", 8*60*60))

	message := quakeSenseNotificationMessage(eventTime, 1)
	want := "[2026-10-02 16:43:42 UTC] QuakeSense detected 1 seismic event(s)."
	if message != want {
		t.Fatalf("unexpected notification message: got %q, want %q", message, want)
	}
}
