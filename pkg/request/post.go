package request

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"net/http"
	"net/textproto"
	"strings"
	"time"
)

func POST(url, payload, contentType string, timeout, retryInterval time.Duration, maxRetries int, trimSpace bool, customTransport http.RoundTripper, headers ...map[string]string) ([]byte, error) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true, MinVersion: tls.VersionTLS12}
	client := http.Client{Timeout: timeout, Transport: customTransport}
	req, err := http.NewRequest("POST", url, strings.NewReader(payload))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", contentType)
	for _, header := range headers {
		for key, value := range header {
			req.Header.Set(key, value)
		}
	}
	requestHostname := req.Header.Get(textproto.CanonicalMIMEHeaderKey("Host"))
	if len(requestHostname) > 0 && len(req.URL.Host) > 0 {
		req.Host = requestHostname
	}

	for retries := 0; retries <= maxRetries; retries++ {
		resp, err := client.Do(req)
		if err != nil {
			time.Sleep(time.Second)
			continue
		}
		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("unexpected status: %s", resp.Status)
		}

		var buf bytes.Buffer
		_, _ = buf.ReadFrom(resp.Body)
		resp.Body.Close()
		b := buf.Bytes()

		if trimSpace {
			for i := 0; i < len(b); i++ {
				if b[i] == ' ' {
					b = append(b[:i], b[i+1:]...)
				}
			}
		}

		return b, nil
	}

	return nil, fmt.Errorf("POST request failed after %d retries", maxRetries)
}
