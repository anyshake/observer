package service

import (
	"context"
	"testing"
)

func TestRecoverCallbackCancelsRun(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()

	func() {
		defer RecoverCallback("test", cancel)
		panic("callback failed")
	}()

	if ctx.Err() == nil {
		t.Fatal("recovered callback panic did not cancel the run")
	}
}

func TestRecoverCallbackKeepsRunWithoutPanic(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()

	func() {
		defer RecoverCallback("test", cancel)
	}()

	if ctx.Err() != nil {
		t.Fatal("callback without panic canceled the run")
	}
}

func TestRecoverWorkerCancelsRun(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()

	func() {
		defer RecoverWorker("test", cancel)
		panic("worker failed")
	}()

	if ctx.Err() == nil {
		t.Fatal("recovered worker panic did not cancel the run")
	}
}

func TestRecoverWorkerKeepsRunWithoutPanic(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()

	func() {
		defer RecoverWorker("test", cancel)
	}()

	if ctx.Err() != nil {
		t.Fatal("worker without panic canceled the run")
	}
}
