package request

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"testing/synctest"
	"time"
)

func TestRetryWaitStopsWhenContextCanceled(t *testing.T) {
	for _, method := range []string{http.MethodGet, http.MethodPost} {
		t.Run(method, func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				ctx, cancel := context.WithCancel(t.Context())
				defer cancel()
				attempts := 0
				transport := roundTripperFunc(func(*http.Request) (*http.Response, error) {
					attempts++
					return nil, io.ErrClosedPipe
				})
				done := make(chan error, 1)
				go func() {
					var err error
					if method == http.MethodGet {
						_, err = GET(ctx, "https://example.invalid/retry", time.Second, time.Hour, 3, false, transport)
					} else {
						_, err = POST(ctx, "https://example.invalid/retry", "x", "text/plain", time.Second, time.Hour, 3, false, transport)
					}
					done <- err
				}()
				synctest.Wait()
				cancel()
				if err := <-done; !errors.Is(err, context.Canceled) {
					t.Fatalf("error = %v, want context.Canceled", err)
				}
				if attempts != 1 {
					t.Fatalf("attempts = %d, want 1", attempts)
				}
			})
		})
	}
}

func TestCanceledContextSkipsRequest(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		calls.Add(1)
	}))
	defer server.Close()
	transport := server.Client().Transport
	if _, err := GET(ctx, server.URL, time.Second, 0, 3, false, transport); !errors.Is(err, context.Canceled) {
		t.Fatalf("GET error = %v, want context.Canceled", err)
	}
	if _, err := POST(ctx, server.URL, "x", "text/plain", time.Second, 0, 3, false, transport); !errors.Is(err, context.Canceled) {
		t.Fatalf("POST error = %v, want context.Canceled", err)
	}
	if calls.Load() != 0 {
		t.Fatalf("server calls = %d, want 0", calls.Load())
	}
}

func TestBodyReadErrorIsReturned(t *testing.T) {
	transport := roundTripperFunc(func(req *http.Request) (*http.Response, error) {
		resp := textResponse(req, http.StatusOK, "")
		resp.Body = io.NopCloser(io.MultiReader(strings.NewReader("partial"), errorReader{io.ErrUnexpectedEOF}))
		return resp, nil
	})
	if body, err := GET(t.Context(), "https://example.invalid/truncated", time.Second, 0, 0, false, transport); !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("GET body = %q, error = %v", body, err)
	}
	if body, err := POST(t.Context(), "https://example.invalid/truncated", "x", "text/plain", time.Second, 0, 0, false, transport); !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("POST body = %q, error = %v", body, err)
	}
}

type errorReader struct {
	err error
}

func (r errorReader) Read([]byte) (int, error) {
	return 0, r.err
}
