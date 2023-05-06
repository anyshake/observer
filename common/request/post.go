package request

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
	"time"
)

func POSTRequest(url, payload, contentType string, timeout, retryInterval time.Duration, maxRetries int, trim bool) ([]byte, error) {
	client := http.Client{
		Timeout: timeout,
	}
	err := fmt.Errorf("GET request failed")

	for retries := 0; retries < maxRetries; retries++ {
		resp, err := client.Post(
			url, contentType, strings.NewReader(payload),
		)
		if err != nil {
			time.Sleep(time.Second)
			continue
		}

		var buf bytes.Buffer
		_, _ = buf.ReadFrom(resp.Body)
		resp.Body.Close()
		b := buf.Bytes()

		if trim {
			for i := 0; i < len(b); i++ {
				if b[i] == ' ' {
					b = append(b[:i], b[i+1:]...)
				}
			}
		}

		return b, nil
	}

	return nil, err
}
