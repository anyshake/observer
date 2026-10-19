package server

import (
	"context"
	"net"
	"net/http"
	"testing"
	"testing/synctest"
	"time"

	"github.com/anyshake/observer/pkg/logger"
)

func newTestServer(t *testing.T, addr string, handler http.Handler) *HttpServer {
	t.Helper()
	srv := New(t.Context(), false, false, nil, nil, logger.GetLogger("http_server_test"))
	srv.server.Addr = addr
	srv.server.Handler = handler
	return srv
}

func TestListenReportsAddressInUse(t *testing.T) {
	occupied, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer occupied.Close()

	srv := newTestServer(t, occupied.Addr().String(), http.NotFoundHandler())
	if err := srv.Listen(); err == nil {
		t.Fatal("Listen succeeded on an occupied address")
	}
	if err := srv.Start(); err == nil {
		t.Fatal("Start succeeded without a listener")
	}
}

func TestStopCancelsInFlightRequestContext(t *testing.T) {
	entered := make(chan struct{})
	srv := newTestServer(t, "127.0.0.1:0", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		close(entered)
		<-r.Context().Done()
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	if err := srv.Listen(); err != nil {
		t.Fatal(err)
	}
	served := make(chan error, 1)
	go func() { served <- srv.Start() }()

	requestDone := make(chan struct{})
	go func() {
		defer close(requestDone)
		resp, err := http.Get("http://" + srv.listener.Addr().String())
		if err == nil {
			_ = resp.Body.Close()
		}
	}()
	select {
	case <-entered:
	case <-time.After(2 * time.Second):
		t.Fatal("request did not reach the handler")
	}

	ctx, cancel := context.WithTimeout(t.Context(), 2*time.Second)
	defer cancel()
	if err := srv.Stop(ctx); err != nil {
		t.Fatalf("Stop waited for a request that ignores shutdown: %v", err)
	}
	if err := <-served; err != nil {
		t.Fatalf("Start() error = %v", err)
	}
	<-requestDone
}

func TestStopClosesListenerWithoutServe(t *testing.T) {
	srv := newTestServer(t, "127.0.0.1:0", http.NotFoundHandler())
	if err := srv.Listen(); err != nil {
		t.Fatal(err)
	}
	addr := srv.listener.Addr().String()

	ctx, cancel := context.WithTimeout(t.Context(), time.Second)
	defer cancel()
	if err := srv.Stop(ctx); err != nil {
		t.Fatalf("Stop() error = %v", err)
	}
	if conn, err := net.DialTimeout("tcp", addr, time.Second); err == nil {
		_ = conn.Close()
		t.Fatal("listener stayed open after Stop")
	}
}

func TestDebugWarningStopsOnShutdown(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		srv := New(t.Context(), true, false, nil, nil, logger.GetLogger("http_server_test"))
		finished := make(chan struct{})
		go func() {
			defer close(finished)
			srv.warnDebugMode()
		}()
		time.Sleep(3 * time.Minute)
		synctest.Wait()
		select {
		case <-finished:
			t.Fatal("debug warning stopped before shutdown")
		default:
		}
		srv.shutdownCancel()
		synctest.Wait()
		select {
		case <-finished:
		default:
			t.Fatal("debug warning kept running after shutdown")
		}
	})
}
