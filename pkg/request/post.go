package request

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/textproto"
	"strings"
	"time"
)

func POST(ctx context.Context, url, payload, contentType string, timeout, retryInterval time.Duration, maxRetries int, trimSpace bool, customTransport http.RoundTripper, headers ...map[string]string) ([]byte, error) {
	customTransport = selectTransport(customTransport)
	client := http.Client{Timeout: timeout, Transport: customTransport}

	for retries := 0; retries <= maxRetries; retries++ {
		req, err := http.NewRequestWithContext(ctx, "POST", url, strings.NewReader(payload))
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

		resp, err := client.Do(req)
		if err != nil {
			if ctxErr := ctx.Err(); ctxErr != nil {
				return nil, ctxErr
			}
			if retries < maxRetries {
				if err := waitRetry(ctx, retryInterval); err != nil {
					return nil, err
				}
			}
			continue
		}
		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			return nil, fmt.Errorf("unexpected status: %s", resp.Status)
		}

		var buf bytes.Buffer
		_, err = buf.ReadFrom(resp.Body)
		resp.Body.Close()
		if err != nil {
			if ctxErr := ctx.Err(); ctxErr != nil {
				return nil, ctxErr
			}
			return nil, fmt.Errorf("failed to read response body: %w", err)
		}
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
