package miniseed

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/anyshake/observer/internal/testsupport"
	"github.com/anyshake/observer/pkg/logger"
	"github.com/anyshake/observer/pkg/timesource"
)

func TestSubscriptionFailureCanStop(t *testing.T) {
	logger.Init()
	hw := testsupport.NewHardware()
	hw.SubscribeErr = errors.New("subscription failed")
	svc := New(t.Context(), hw, nil, timesource.New(nil))
	if err := svc.Start(t.Context()); err == nil {
		t.Fatal("Start succeeded despite a subscription failure")
	}
	if svc.GetStatus().GetIsRunning() {
		t.Fatal("failed subscription left the service running")
	}
	ctx, cancel := context.WithTimeout(t.Context(), time.Second)
	defer cancel()
	if err := svc.Stop(ctx); err != nil {
		t.Fatalf("Stop after subscription failure: %v", err)
	}
}
