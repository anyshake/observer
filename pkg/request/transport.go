package request

import (
	"crypto/tls"
	"net/http"
	"sync/atomic"
)

type roundTripperBox struct {
	rt http.RoundTripper
}

// roundTripperOverride is unset in production. Tests use it to stub callers
// that do not pass their own transport.
var roundTripperOverride atomic.Pointer[roundTripperBox]

func selectTransport(custom http.RoundTripper) http.RoundTripper {
	if custom != nil {
		return custom
	}
	if box := roundTripperOverride.Load(); box != nil && box.rt != nil {
		return box.rt
	}
	return defaultTransport
}

// SetRoundTripperOverrideForTest installs a process-wide transport used only
// when a caller passes a nil transport. The returned function restores the
// previous value. Do not call it from tests that run in parallel with each other.
func SetRoundTripperOverrideForTest(rt http.RoundTripper) func() {
	previous := roundTripperOverride.Load()
	if rt == nil {
		roundTripperOverride.Store(nil)
	} else {
		roundTripperOverride.Store(&roundTripperBox{rt: rt})
	}
	return func() {
		roundTripperOverride.Store(previous)
	}
}

// Keep the existing TLS policy local to this package and reuse its connection pool.
var defaultTransport = func() *http.Transport {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.TLSClientConfig = &tls.Config{
		InsecureSkipVerify: true,
		MinVersion:         tls.VersionTLS12,
	}
	return transport
}()
