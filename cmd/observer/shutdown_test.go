package main

import (
	"context"
	"errors"
	"testing"
	"testing/synctest"
	"time"

	graph_resolver "github.com/anyshake/observer/internal/server/router/graph"
	"github.com/anyshake/observer/internal/service"
	"github.com/anyshake/observer/pkg/jobtracker"
	"github.com/anyshake/observer/pkg/logger"
)

type shutdownTestService struct {
	service.IService
	stop func(context.Context) error
}

func (s *shutdownTestService) Stop(ctx context.Context) error { return s.stop(ctx) }

func TestShutdownWaitsForModulesAndPurgeBeforeCleanup(t *testing.T) {
	logger.Init()
	synctest.Test(t, func(t *testing.T) {
		appCtx, cancelApp := context.WithCancel(t.Context())
		defer cancelApp()
		purgeStarted := make(chan struct{})
		releasePurge := make(chan struct{})
		tracker := jobtracker.New(appCtx, "purge")
		tracker.Start(time.Now(), "purge", func(context.Context) error {
			close(purgeStarted)
			<-releasePurge
			return nil
		})
		<-purgeStarted
		moduleStarted := make(chan struct{})
		releaseModule := make(chan struct{})
		cleaned := false
		shutdown := &applicationShutdown{
			resolver: &graph_resolver.Resolver{AppContext: appCtx, DataPurgeJob: tracker},
			services: map[string]service.IService{
				"module": &shutdownTestService{stop: func(ctx context.Context) error {
					close(moduleStarted)
					select {
					case <-releaseModule:
						return nil
					case <-ctx.Done():
						return ctx.Err()
					}
				}},
			},
			cleanup: func(context.Context) error { cleaned = true; return nil },
		}
		cancelApp()
		ctx, cancel := context.WithTimeout(t.Context(), shutdownTimeout)
		defer cancel()
		done := make(chan error, 1)
		go func() { done <- shutdown.Run(ctx) }()
		<-moduleStarted
		time.Sleep(6 * time.Second)
		synctest.Wait()
		if cleaned {
			t.Fatal("shared resources closed while modules and purge were running")
		}
		close(releaseModule)
		synctest.Wait()
		if cleaned {
			t.Fatal("the five-second warning ended the wait for purge")
		}
		close(releasePurge)
		if err := <-done; err != nil {
			t.Fatal(err)
		}
		if !cleaned || len(shutdown.Pending()) != 0 {
			t.Fatalf("shutdown incomplete: cleaned=%v, pending=%v", cleaned, shutdown.Pending())
		}
	})
}

func TestShutdownDeadlineDoesNotCloseResourcesUsedByPurge(t *testing.T) {
	logger.Init()
	synctest.Test(t, func(t *testing.T) {
		tracker := jobtracker.New(t.Context(), "purge")
		started := make(chan struct{})
		release := make(chan struct{})
		tracker.Start(time.Now(), "purge", func(context.Context) error {
			close(started)
			<-release
			return nil
		})
		<-started
		cleaned := false
		shutdown := &applicationShutdown{
			resolver: &graph_resolver.Resolver{DataPurgeJob: tracker},
			cleanup:  func(context.Context) error { cleaned = true; return nil },
		}
		ctx, cancel := context.WithTimeout(t.Context(), shutdownTimeout)
		defer cancel()
		if err := shutdown.Run(ctx); !errors.Is(err, context.DeadlineExceeded) {
			t.Fatalf("shutdown while purge remained active = %v", err)
		}
		if cleaned {
			t.Fatal("shared resources closed after shutdown timed out")
		}
		close(release)
		if err := tracker.Wait(t.Context()); err != nil {
			t.Fatal(err)
		}
	})
}
