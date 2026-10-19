package graph_resolver

import (
	"context"
	"testing"
	"testing/synctest"
)

func TestShutdownRejectsNewServiceOperationsAndWaitsForAcceptedOperation(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		appCtx, cancel := context.WithCancel(t.Context())
		defer cancel()
		resolver := &Resolver{AppContext: appCtx}
		operationCtx, release, err := resolver.beginServiceOperation(t.Context())
		if err != nil {
			t.Fatal(err)
		}
		resolver.BeginShutdown()
		cancel()
		<-operationCtx.Done()
		if _, release, err := resolver.beginServiceOperation(t.Context()); err == nil {
			release()
			t.Fatal("accepted a service operation during shutdown")
		}
		finished := make(chan error, 1)
		go func() { finished <- resolver.WaitServiceOperations(t.Context()) }()
		synctest.Wait()
		select {
		case <-finished:
			t.Fatal("shutdown passed an active service operation")
		default:
		}
		release()
		if err := <-finished; err != nil {
			t.Fatal(err)
		}
	})
}
