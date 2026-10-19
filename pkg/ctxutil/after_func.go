package ctxutil

import (
	"context"
	"sync"
)

func AfterFunc(ctx context.Context, f func()) func() {
	done := make(chan struct{})
	stop := context.AfterFunc(ctx, func() {
		defer close(done)
		f()
	})
	return sync.OnceFunc(func() {
		if !stop() {
			<-done
		}
	})
}
