package main

import (
	"context"
	"errors"
	"os"
	"syscall"
	"testing"
	"testing/synctest"
)

func TestOpenHardwareCancelsOnSignal(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		signalChan := make(chan os.Signal, 1)
		openStarted := make(chan struct{})
		go func() {
			<-openStarted
			signalChan <- syscall.SIGINT
		}()
		_, _, interrupted, err := openHardware(signalChan, func(ctx context.Context) (context.Context, context.CancelFunc, error) {
			close(openStarted)
			<-ctx.Done()
			return nil, nil, ctx.Err()
		})
		if !interrupted {
			t.Fatal("expected startup to be interrupted")
		}
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("expected context.Canceled, got %v", err)
		}
	})
}

func TestOpenHardwareKeepsContextAndSignalAfterSuccess(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		signalChan := make(chan os.Signal, 1)
		hardwareCtx, abort, interrupted, err := openHardware(signalChan, func(ctx context.Context) (context.Context, context.CancelFunc, error) {
			subCtx, cancel := context.WithCancel(ctx)
			return subCtx, cancel, nil
		})
		if err != nil || interrupted {
			t.Fatalf("unexpected result: interrupted=%v err=%v", interrupted, err)
		}
		defer abort()
		signalChan <- syscall.SIGINT
		synctest.Wait()
		if hardwareCtx.Err() != nil {
			t.Fatalf("hardware context canceled after startup: %v", hardwareCtx.Err())
		}
		if len(signalChan) != 1 {
			t.Fatal("signal consumed after startup")
		}
	})
}
