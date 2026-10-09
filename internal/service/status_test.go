package service

import (
	"testing"
	"time"
)

func TestStatusRoundTrip(t *testing.T) {
	t.Parallel()

	var status Status
	started := time.Unix(100, 0).UTC()
	stopped := time.Unix(200, 0).UTC()
	updated := time.Unix(300, 0).UTC()

	status.SetIsRunning(true)
	status.SetRestarts(2)
	status.IncrementRestarts()
	status.SetStartedAt(started)
	status.SetStoppedAt(stopped)
	status.SetUpdatedAt(updated)

	if !status.GetIsRunning() || status.GetRestarts() != 3 {
		t.Fatalf("running = %v, restarts = %d", status.GetIsRunning(), status.GetRestarts())
	}
	if !status.GetStartedAt().Equal(started) || !status.GetStoppedAt().Equal(stopped) || !status.GetUpdatedAt().Equal(updated) {
		t.Fatalf("times = %v %v %v", status.GetStartedAt(), status.GetStoppedAt(), status.GetUpdatedAt())
	}

	status.SetIsRunning(false)
	if status.GetIsRunning() {
		t.Fatal("status stayed running")
	}
}
