package testsupport

import (
	"context"
	"testing"
	"time"
)

func StopService(t *testing.T, svc interface{ Stop(context.Context) error }) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := svc.Stop(ctx); err != nil {
		t.Errorf("failed to stop service: %v", err)
	}
}
