package auth

import (
	"context"
	"crypto/sha256"
	"errors"
	"net/http"
	"time"
)

// All challenge admission, consumption and expiry share stateMu. Reservations
// count toward both limits before expensive key/CAPTCHA generation starts.
func (h *auth) reserveChallenge(source string) (int, error) {
	h.stateMu.Lock()
	defer h.stateMu.Unlock()
	now := time.Now()
	h.cleanupStateLocked(now)

	budget := h.sources[source]
	if budget == nil {
		if len(h.sources) >= maxAuthSources {
			return http.StatusServiceUnavailable, errors.New("too many login sources")
		}
		budget = &authSource{resetAt: now.Add(time.Minute)}
		h.sources[source] = budget
	}
	if !now.Before(budget.resetAt) {
		budget.requests = 0
		budget.resetAt = now.Add(time.Minute)
	}
	if budget.requests >= maxPreAuthPerMinute || budget.pending >= maxPendingPerSource {
		return http.StatusTooManyRequests, errors.New("too many login challenges from this source")
	}
	if int(h.authChallengePool.Len())+h.reservations >= maxPendingChallenges {
		return http.StatusServiceUnavailable, errors.New("too many pending login challenges")
	}
	budget.requests++
	budget.pending++
	h.reservations++
	return http.StatusOK, nil
}

func (h *auth) finishChallenge(source, id string, challenge *authChallenge) {
	h.stateMu.Lock()
	defer h.stateMu.Unlock()
	h.reservations--
	if challenge == nil {
		h.sources[source].pending--
		return
	}
	challenge.source = source
	h.authChallengePool.Set(id, challenge)
}

func (h *auth) takeChallenge(id string) (*authChallenge, bool) {
	h.stateMu.Lock()
	defer h.stateMu.Unlock()
	challenge, ok := h.authChallengePool.GetAndDel(id)
	if ok {
		h.releaseChallengeLocked(challenge)
	}
	return challenge, ok
}

func (h *auth) releaseChallengeLocked(challenge *authChallenge) {
	if budget := h.sources[challenge.source]; budget != nil {
		budget.pending--
	}
}

func (h *auth) rememberNonce(nonce []byte) (int, error) {
	h.stateMu.Lock()
	defer h.stateMu.Unlock()
	now := time.Now()
	h.cleanupStateLocked(now)
	key := sha256.Sum256(nonce)
	if _, exists := h.nonceCache[key]; exists {
		return http.StatusForbidden, errAuthenticationFailed
	}
	// Do not evict live replay records to admit new ones.
	if len(h.nonceCache) >= maxRememberedNonces {
		return http.StatusServiceUnavailable, errors.New("too many active login nonces")
	}
	h.nonceCache[key] = now.Add(nonceTTL)
	return http.StatusOK, nil
}

func (h *auth) cleanupStateLocked(now time.Time) {
	for key, expiry := range h.nonceCache {
		if !now.Before(expiry) {
			delete(h.nonceCache, key)
		}
	}
	for id, challenge := range h.authChallengePool.Iterator() {
		if !challenge.isChallengeAlive() {
			h.authChallengePool.Del(id)
			h.releaseChallengeLocked(challenge)
			captchaStore.Get(challenge.captchaID, true)
		}
	}
	for source, budget := range h.sources {
		if budget.pending == 0 && !now.Before(budget.resetAt) {
			delete(h.sources, source)
		}
	}
}

func (h *auth) runCleanup(ctx context.Context) {
	ticker := time.NewTicker(authCleanupInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			h.cleanupExpiredAuthState()
		}
	}
}
