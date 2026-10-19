package graph_resolver

import (
	"context"
	"errors"
	"sync"
)

type serviceControl struct {
	mu      sync.Mutex
	closing bool
	active  sync.WaitGroup
	done    chan struct{}
}

func (r *Resolver) BeginShutdown() {
	c := &r.serviceControl
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.closing {
		return
	}
	c.closing = true
	c.done = make(chan struct{})
	go func() {
		c.active.Wait()
		close(c.done)
	}()
}

// WaitServiceOperations ensures an accepted management request cannot restart a
// module after the application has begun stopping it.
func (r *Resolver) WaitServiceOperations(ctx context.Context) error {
	r.BeginShutdown()
	select {
	case <-r.serviceControl.done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (r *Resolver) beginServiceOperation(ctx context.Context) (context.Context, func(), error) {
	c := &r.serviceControl
	c.mu.Lock()
	if c.closing {
		c.mu.Unlock()
		return nil, nil, errors.New("application is shutting down")
	}
	c.active.Add(1)
	c.mu.Unlock()
	operationCtx, cancel := context.WithCancel(ctx)
	stopCancel := func() bool { return false }
	if r.AppContext != nil {
		stopCancel = context.AfterFunc(r.AppContext, cancel)
		if r.AppContext.Err() != nil {
			cancel()
		}
	}
	release := func() {
		stopCancel()
		cancel()
		c.active.Done()
	}
	if err := operationCtx.Err(); err != nil {
		release()
		return nil, nil, err
	}
	return operationCtx, release, nil
}
