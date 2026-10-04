package request

import (
	"crypto/tls"
	"net/http"
)

// Keep the existing TLS policy local to this package and reuse its connection pool.
var defaultTransport = func() *http.Transport {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.TLSClientConfig = &tls.Config{
		InsecureSkipVerify: true,
		MinVersion:         tls.VersionTLS12,
	}
	return transport
}()
