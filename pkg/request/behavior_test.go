package request

import (
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestRoundTripperOverride(t *testing.T) {
	calls := 0
	restore := SetRoundTripperOverrideForTest(roundTripperFunc(func(req *http.Request) (*http.Response, error) {
		calls++
		return textResponse(req, http.StatusOK, "override"), nil
	}))
	defer restore()

	restoreNil := SetRoundTripperOverrideForTest(nil)
	restoreNil()

	body, err := GET("https://example.invalid/override", time.Second, 0, 0, false, nil)
	if err != nil {
		t.Fatal(err)
	}
	if string(body) != "override" || calls != 1 {
		t.Fatalf("override body = %q, calls = %d", body, calls)
	}

	body, err = GET("https://example.invalid/explicit", time.Second, 0, 0, false, roundTripperFunc(func(req *http.Request) (*http.Response, error) {
		return textResponse(req, http.StatusOK, "explicit"), nil
	}))
	if err != nil {
		t.Fatal(err)
	}
	if string(body) != "explicit" || calls != 1 {
		t.Fatalf("explicit transport body = %q, calls = %d", body, calls)
	}
}

func TestGETAndPOSTEdges(t *testing.T) {
	if _, err := GET("://bad", time.Second, 0, 0, false, nil); err == nil {
		t.Fatal("invalid GET URL accepted")
	}
	if _, err := POST("://bad", "", "text/plain", time.Second, 0, 0, false, nil); err == nil {
		t.Fatal("invalid POST URL accepted")
	}

	attempts := 0
	flaky := roundTripperFunc(func(req *http.Request) (*http.Response, error) {
		attempts++
		if attempts == 1 {
			return nil, io.ErrUnexpectedEOF
		}
		return textResponse(req, http.StatusOK, "a b"), nil
	})
	body, err := GET("https://example.invalid/retry", time.Second, 0, 1, true, flaky)
	if err != nil {
		t.Fatal(err)
	}
	if string(body) != "ab" || attempts != 2 {
		t.Fatalf("retried GET = %q after %d attempts", body, attempts)
	}

	statusTransport := roundTripperFunc(func(req *http.Request) (*http.Response, error) {
		return textResponse(req, http.StatusTeapot, "nope"), nil
	})
	if _, err := GET("https://example.invalid/status", time.Second, 0, 2, false, statusTransport); err == nil || !strings.Contains(err.Error(), "teapot") {
		t.Fatalf("status error = %v", err)
	}

	hostSeen := ""
	hostTransport := roundTripperFunc(func(req *http.Request) (*http.Response, error) {
		hostSeen = req.Host
		payload, _ := io.ReadAll(req.Body)
		return textResponse(req, http.StatusOK, string(payload)), nil
	})
	body, err = POST(
		"https://example.invalid/post",
		"hello",
		"text/plain",
		time.Second,
		0,
		0,
		false,
		hostTransport,
		map[string]string{"Host": "station.example"},
	)
	if err != nil {
		t.Fatal(err)
	}
	if string(body) != "hello" || hostSeen != "station.example" {
		t.Fatalf("POST body = %q host = %q", body, hostSeen)
	}

	if _, err := GET("https://example.invalid/give-up", time.Second, 0, 0, false, roundTripperFunc(func(*http.Request) (*http.Response, error) {
		return nil, io.ErrClosedPipe
	})); err == nil || !strings.Contains(err.Error(), "failed after 0 retries") {
		t.Fatalf("exhausted GET error = %v", err)
	}
	if _, err := POST("https://example.invalid/give-up", "x", "text/plain", time.Second, 0, 0, false, roundTripperFunc(func(*http.Request) (*http.Response, error) {
		return nil, io.ErrClosedPipe
	})); err == nil || !strings.Contains(err.Error(), "failed after 0 retries") {
		t.Fatalf("exhausted POST error = %v", err)
	}
}

func textResponse(req *http.Request, status int, body string) *http.Response {
	return &http.Response{
		StatusCode: status,
		Status:     http.StatusText(status),
		Header:     make(http.Header),
		Body:       io.NopCloser(strings.NewReader(body)),
		Request:    req,
	}
}
