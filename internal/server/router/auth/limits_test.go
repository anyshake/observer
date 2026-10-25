package auth

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"testing/synctest"
	"time"

	"github.com/gin-gonic/gin"
)

func TestInvalidLoginDoesNotRetainNonce(t *testing.T) {
	service := newAuthService(t, nil)
	for _, test := range []struct {
		nonce string
		code  int
	}{
		{strings.Repeat("A", 100000), http.StatusBadRequest},
		{base64.StdEncoding.EncodeToString(make([]byte, encryptedNonceBytes-1)), http.StatusBadRequest},
		{base64.StdEncoding.EncodeToString(make([]byte, encryptedNonceBytes)), http.StatusUnauthorized},
	} {
		for range 2 {
			code, _, err := service.login(t.Context(), "session", "secret", test.nonce,
				"missing", "1:aa", "captcha", "123456", "payload", "", "")
			if err == nil || code != test.code {
				t.Fatalf("invalid login = %d, %v; want %d", code, err, test.code)
			}
			if len(service.nonceCache) != 0 {
				t.Fatal("invalid login retained a nonce")
			}
		}
	}
}

func TestNonceReplayCapacityAndExpiry(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		service := newAuthService(t, nil)
		ctx, cancel := context.WithCancel(t.Context())
		defer cancel()
		go service.runCleanup(ctx)
		for i := range maxRememberedNonces {
			if _, err := service.rememberNonce([]byte(fmt.Sprint(i))); err != nil {
				t.Fatal(err)
			}
		}
		if code, _ := service.rememberNonce([]byte("extra")); code != http.StatusServiceUnavailable {
			t.Fatalf("full nonce cache = %d", code)
		}
		if code, _ := service.rememberNonce([]byte("0")); code != http.StatusForbidden {
			t.Fatalf("replayed nonce = %d", code)
		}
		time.Sleep(nonceTTL + authCleanupInterval)
		synctest.Wait()
		service.stateMu.Lock()
		remaining := len(service.nonceCache)
		service.stateMu.Unlock()
		if remaining != 0 {
			t.Fatalf("expired nonces were retained: %d", remaining)
		}
	})
}

func TestChallengeSourceIsolationAndRate(t *testing.T) {
	service := newAuthService(t, nil)
	for range maxPendingPerSource {
		if _, err := service.reserveChallenge("source-a"); err != nil {
			t.Fatal(err)
		}
	}
	if code, _ := service.reserveChallenge("source-a"); code != http.StatusTooManyRequests {
		t.Fatalf("source limit = %d", code)
	}
	if _, err := service.reserveChallenge("source-b"); err != nil {
		t.Fatalf("one source blocked another: %v", err)
	}
	service.finishChallenge("source-b", "", nil)
	for range maxPendingPerSource {
		service.finishChallenge("source-a", "", nil)
	}
	for i := maxPendingPerSource; i < maxPreAuthPerMinute; i++ {
		if _, err := service.reserveChallenge("source-a"); err != nil {
			t.Fatal(err)
		}
		service.finishChallenge("source-a", "", nil)
	}
	if code, _ := service.reserveChallenge("source-a"); code != http.StatusTooManyRequests {
		t.Fatalf("released slots bypassed the rate limit: %d", code)
	}
}

func TestConcurrentChallengeReservationsStayBounded(t *testing.T) {
	service := newAuthService(t, nil)
	var wg sync.WaitGroup
	for i := range maxPendingChallenges * 2 {
		wg.Go(func() { _, _ = service.reserveChallenge(fmt.Sprint(i)) })
	}
	wg.Wait()
	if service.reservations != maxPendingChallenges {
		t.Fatalf("concurrent reservations = %d", service.reservations)
	}
}

func TestChallengeExpiryReleasesSourceCapacity(t *testing.T) {
	service := newAuthService(t, nil)
	for i := range maxPendingPerSource {
		if _, err := service.reserveChallenge("source"); err != nil {
			t.Fatal(err)
		}
		service.finishChallenge("source", fmt.Sprint(i), &authChallenge{
			createdAt: time.Now().Add(-time.Minute), ttl: time.Second,
		})
	}
	service.cleanupExpiredAuthState()
	if service.authChallengePool.Len() != 0 || service.sources["source"].pending != 0 {
		t.Fatal("expired challenges retained source capacity")
	}
	if _, err := service.reserveChallenge("source"); err != nil {
		t.Fatalf("expired source capacity was not reusable: %v", err)
	}
}

func TestAuthRejectsOversizedBodies(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	noop := func(*gin.Context) {}
	Setup(router.Group("/"), nil, noop, noop, noop, t.Context())
	for _, contentLength := range []int64{-1, maxAuthBodyBytes + 1} {
		request := httptest.NewRequest(http.MethodPost, "/auth", strings.NewReader(
			`{"action":"preauth","nonce":"`+strings.Repeat("A", maxAuthBodyBytes)+`"}`))
		request.ContentLength = contentLength
		request.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, request)
		if recorder.Code != http.StatusRequestEntityTooLarge {
			t.Fatalf("body length %d: status = %d", contentLength, recorder.Code)
		}
	}
}

func TestPreAuthQuotaIgnoresForwardedHeaders(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	noop := func(*gin.Context) {}
	Setup(router.Group("/"), nil, noop, noop, noop, t.Context())
	for i := range maxPendingPerSource + 2 {
		request := httptest.NewRequest(http.MethodPost, "/auth", strings.NewReader(`{"action":"preauth"}`))
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("X-Forwarded-For", fmt.Sprintf("198.51.100.%d", i+1))
		request.Header.Set("X-Real-IP", fmt.Sprintf("198.51.100.%d", i+1))
		request.RemoteAddr = "203.0.113.1:12345"
		want := http.StatusOK
		if i == maxPendingPerSource {
			want = http.StatusTooManyRequests
		} else if i > maxPendingPerSource {
			request.RemoteAddr = "203.0.113.2:12345"
		}
		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, request)
		if recorder.Code != want {
			t.Fatalf("request %d: status = %d, want %d; %s", i, recorder.Code, want, recorder.Body.String())
		}
	}
}
