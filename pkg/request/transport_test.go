package request

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func TestRequestsDoNotModifyDefaultTransport(t *testing.T) {
	globalTransport := http.DefaultTransport.(*http.Transport)
	globalTLSConfig := globalTransport.TLSClientConfig
	if globalTLSConfig != nil && globalTLSConfig.InsecureSkipVerify {
		t.Fatal("global transport already disables certificate verification")
	}

	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(r.Method))
	}))
	defer server.Close()
	defer defaultTransport.CloseIdleConnections()

	customTransport := server.Client().Transport.(*http.Transport)
	defer customTransport.CloseIdleConnections()
	customTLSConfig := customTransport.TLSClientConfig
	customUsed := false
	trackedTransport := roundTripperFunc(func(req *http.Request) (*http.Response, error) {
		customUsed = true
		return customTransport.RoundTrip(req)
	})

	for _, method := range []string{http.MethodGet, http.MethodPost} {
		for _, tc := range []struct {
			name      string
			transport http.RoundTripper
		}{
			{name: "default"},
			{name: "custom", transport: trackedTransport},
		} {
			t.Run(method+"/"+tc.name, func(t *testing.T) {
				customUsed = false
				var body []byte
				var err error
				if method == http.MethodGet {
					body, err = GET(server.URL, 5*time.Second, 0, 0, false, tc.transport)
				} else {
					body, err = POST(server.URL, "", "application/json", 5*time.Second, 0, 0, false, tc.transport)
				}
				if err != nil {
					t.Fatalf("%s request failed: %v", method, err)
				}
				if string(body) != method {
					t.Fatalf("response = %q, want %q", body, method)
				}
				if customUsed != (tc.transport != nil) {
					t.Fatal("request did not use the expected transport")
				}
				if http.DefaultTransport != globalTransport || globalTransport.TLSClientConfig != globalTLSConfig {
					t.Fatal("request changed the global transport or its TLS config")
				}
				if globalTLSConfig != nil && globalTLSConfig.InsecureSkipVerify {
					t.Fatal("request disabled global certificate verification")
				}
				if customTransport.TLSClientConfig != customTLSConfig || customTLSConfig.InsecureSkipVerify {
					t.Fatal("request changed the custom transport's TLS policy")
				}
			})
		}
	}
}
