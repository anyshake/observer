package service

import (
	"context"
	"runtime/debug"

	"github.com/anyshake/observer/pkg/logger"
)

// RecoverCallback must be deferred directly by a subscription callback. A
// recovered panic cancels the run so its worker stops the service.
func RecoverCallback(id string, cancel context.CancelFunc) {
	if r := recover(); r != nil {
		stopAfterPanic(id, "subscription callback", r, cancel)
	}
}

// RecoverWorker must be deferred directly by a goroutine started for a run. A
// recovered panic cancels the run so its worker stops the service.
func RecoverWorker(id string, cancel context.CancelFunc) {
	if r := recover(); r != nil {
		stopAfterPanic(id, "worker goroutine", r, cancel)
	}
}

func stopAfterPanic(id, source string, value any, cancel context.CancelFunc) {
	logger.GetLogger(id).Errorf("%s unexpectedly crashed, stopping service: %v\n%s", source, value, debug.Stack())
	cancel()
}
