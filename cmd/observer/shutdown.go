package main

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"sync"
	"time"

	graph_resolver "github.com/anyshake/observer/internal/server/router/graph"
	"github.com/anyshake/observer/internal/service"
	"github.com/anyshake/observer/pkg/logger"
)

const shutdownTimeout = 30 * time.Second

// applicationShutdown owns shutdown ordering. Shared resources are released
// only after accepted management operations, modules and background jobs finish.
type applicationShutdown struct {
	resolver  *graph_resolver.Resolver
	stopHTTP  func(context.Context) error
	stopPprof func(context.Context) error
	services  map[string]service.IService
	cleanup   func(context.Context) error
	pending   sync.Map
}

func (s *applicationShutdown) Run(ctx context.Context) error {
	var stopping sync.WaitGroup
	errs := make(chan error, len(s.services)+2)
	stopModule := func(name string, stop func(context.Context) error) {
		s.pending.Store(name, struct{}{})
		stopping.Add(1)
		go func() {
			defer stopping.Done()
			if err := stop(ctx); err != nil {
				if ctx.Err() == nil {
					s.pending.Delete(name)
				}
				logger.GetLogger(name).Errorf("failed to stop %s: %v", name, err)
				errs <- fmt.Errorf("%s: %w", name, err)
			} else {
				s.pending.Delete(name)
				logger.GetLogger(name).Infof("%s has been stopped", name)
			}
		}()
	}
	if s.stopHTTP != nil {
		stopModule("http_server", s.stopHTTP)
	}
	if s.stopPprof != nil {
		stopModule("pprof_server", s.stopPprof)
	}
	s.pending.Store("service_operations", struct{}{})
	if err := s.resolver.WaitServiceOperations(ctx); err != nil {
		return err
	}
	s.pending.Delete("service_operations")
	// Even a module whose worker has stopped may still own resources.
	for name, module := range s.services {
		stopModule(name, module.Stop)
	}
	if err := s.waitPurgeJob(ctx); err != nil {
		return err
	}
	stopping.Wait()
	close(errs)
	var stopErrs []error
	for err := range errs {
		stopErrs = append(stopErrs, err)
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	s.pending.Store("resource_cleanup", struct{}{})
	stopErrs = append(stopErrs, s.cleanup(ctx))
	s.pending.Delete("resource_cleanup")
	return errors.Join(stopErrs...)
}

func (s *applicationShutdown) waitPurgeJob(ctx context.Context) error {
	if s.resolver.DataPurgeJob == nil {
		return nil
	}
	s.pending.Store("data_purge_job", struct{}{})
	warningCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	if err := s.resolver.DataPurgeJob.Wait(warningCtx); err != nil {
		logger.GetLogger(main).Warnf("cleanup job is still stopping: %v", err)
	}
	cancel()
	if err := s.resolver.DataPurgeJob.Wait(ctx); err != nil {
		return err
	}
	s.pending.Delete("data_purge_job")
	return nil
}

func (s *applicationShutdown) Pending() []string {
	var names []string
	s.pending.Range(func(name, _ any) bool {
		names = append(names, name.(string))
		return true
	})
	sort.Strings(names)
	return names
}
