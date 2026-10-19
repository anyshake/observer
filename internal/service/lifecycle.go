package service

import (
	"context"
	"sync"
)

// Lifecycle serializes initialization, startup and cleanup. Its zero value is
// ready to use. A Lifecycle must not be copied after first use.
type Lifecycle struct {
	once sync.Once
	gate chan struct{}
}

func (l *Lifecycle) Lock(ctx context.Context) error {
	l.once.Do(func() { l.gate = make(chan struct{}, 1) })
	select {
	case l.gate <- struct{}{}:
		if err := ctx.Err(); err != nil {
			l.Unlock()
			return err
		}
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (l *Lifecycle) Unlock() {
	<-l.gate
}

// Stop keeps the lifecycle locked until cleanup actually finishes, even when
// the caller's deadline expires. Later operations must wait for that cleanup.
func (l *Lifecycle) Stop(ctx context.Context, cleanup func() error) error {
	if err := l.Lock(ctx); err != nil {
		return err
	}
	done := make(chan error, 1)
	go func() {
		defer l.Unlock()
		done <- cleanup()
	}()
	select {
	case err := <-done:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}
