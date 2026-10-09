package auth

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/alphadose/haxmap"
	"github.com/anyshake/observer/internal/dao/action"
	"github.com/anyshake/observer/internal/testsupport"
	"github.com/anyshake/observer/pkg/logger"
	lru "github.com/hashicorp/golang-lru/v2"
)

func TestLoginFlow(t *testing.T) {
	logger.Init()
	previous := verifyCaptchaString
	verifyCaptchaString = func(id, value string) bool { return value == "000000" }
	t.Cleanup(func() { verifyCaptchaString = previous })

	_, handler := testsupport.OpenDAO(t)
	if _, err := handler.SysUserCreate("alice", "GoodPass1!", true); err != nil {
		t.Fatal(err)
	}
	session := newLoginSession(t, handler)
	code, userID, err := session.login(mustCredential(t, "alice", "GoodPass1!"), "000000", "agent", "127.0.0.1")
	if err != nil || code != http.StatusOK || userID == "" {
		t.Fatalf("login = %d %s %v", code, userID, err)
	}

	code, _, err = session.login(mustCredential(t, "alice", "WrongPass1!"), "000000", "", "")
	if err == nil || code != http.StatusUnauthorized {
		t.Fatalf("bad password = %d %v", code, err)
	}
	code, _, err = session.login([]byte(`{"username":"alice"}`), "000000", "", "")
	if err == nil || code != http.StatusUnauthorized {
		t.Fatalf("partial credential = %d %v", code, err)
	}
	code, _, err = session.login([]byte("not-json"), "000000", "", "")
	if err == nil || code != http.StatusUnauthorized {
		t.Fatalf("invalid credential = %d %v", code, err)
	}

	secret, nonce, _ := session.material([]byte("x"))
	code, _, err = session.attempt(secret, nonce, "!!!", "000000")
	if err == nil || code != http.StatusBadRequest {
		t.Fatalf("bad payload = %d %v", code, err)
	}
	secret, nonce, _ = session.material([]byte("x"))
	code, _, err = session.attempt(secret, nonce, base64.StdEncoding.EncodeToString([]byte("not-ciphertext")), "000000")
	if err == nil || code != http.StatusUnauthorized {
		t.Fatalf("undecryptable payload = %d %v", code, err)
	}

	secret, _, payload := session.material(mustCredential(t, "alice", "GoodPass1!"))
	code, _, err = session.attempt(secret, base64.StdEncoding.EncodeToString([]byte("short")), payload, "000000")
	if err == nil || code != http.StatusForbidden {
		t.Fatalf("malformed nonce = %d %v", code, err)
	}

	plain, err := session.pair.rsaKeyPair.Encrypt([]byte("!!!"), true)
	if err != nil {
		t.Fatal(err)
	}
	_, nonce, payload = session.material(mustCredential(t, "alice", "GoodPass1!"))
	code, _, err = session.attempt(string(plain), nonce, payload, "000000")
	if err == nil || code != http.StatusUnauthorized {
		t.Fatalf("non-base64 secret = %d %v", code, err)
	}
	_, nonce, payload = session.material(mustCredential(t, "alice", "GoodPass1!"))
	code, _, err = session.attempt("YQ==", nonce, payload, "000000")
	if err == nil || code != http.StatusUnauthorized {
		t.Fatalf("undecryptable secret = %d %v", code, err)
	}

	_, nonce, payload = session.material(mustCredential(t, "alice", "GoodPass1!"))
	session.rotateChallenge()
	code, _, err = session.auth.login("missing-session", secret, nonce, session.challengeID, session.solution, "captcha", "000000", payload, "", "")
	if err == nil || code != http.StatusUnauthorized {
		t.Fatalf("unknown session = %d %v", code, err)
	}

	code, _, err = session.login(mustCredential(t, "alice", "GoodPass1!"), "999999", "", "")
	if err == nil || code != http.StatusUnauthorized {
		t.Fatalf("captcha = %d %v", code, err)
	}
	session.rotateChallenge()
	code, _, err = session.auth.login(session.id, "c2VjcmV0", uniqueNonce(), session.challengeID, "1:"+strings.Repeat("00", sha512.Size), "captcha", "000000", "e30=", "", "")
	if err == nil || code != http.StatusUnauthorized {
		t.Fatalf("bad proof = %d %v", code, err)
	}

	expired := newLoginSession(t, handler)
	expired.challenge.createdAt = time.Now().Add(-time.Hour)
	secret, nonce, payload = expired.material(mustCredential(t, "alice", "GoodPass1!"))
	code, _, err = expired.auth.login(expired.id, secret, nonce, expired.challengeID, expired.solution, "captcha", "000000", payload, "", "")
	if err == nil || code != http.StatusUnauthorized {
		t.Fatalf("expired challenge = %d %v", code, err)
	}

	stale := newLoginSession(t, handler)
	stale.pair.createdAt = time.Now().Add(-time.Hour)
	code, _, err = stale.login(mustCredential(t, "alice", "GoodPass1!"), "000000", "", "")
	if err == nil || code != http.StatusUnauthorized {
		t.Fatalf("expired session = %d %v", code, err)
	}
}

func TestPreAuthLimits(t *testing.T) {
	logger.Init()
	_, handler := testsupport.OpenDAO(t)
	service := newAuthService(t, handler)
	deadPair, err := newKeyPair(time.Second)
	if err != nil {
		t.Fatal(err)
	}
	deadPair.createdAt = time.Now().Add(-time.Hour)
	service.keyPairDataPool.Set(deadPair.getKeyPairId(), deadPair)
	deadChallenge, err := newAuthChallenge(time.Second)
	if err != nil {
		t.Fatal(err)
	}
	deadChallenge.createdAt = time.Now().Add(-time.Hour)
	service.authChallengePool.Set("dead", deadChallenge)
	service.cleanupExpiredAuthState()
	if service.keyPairDataPool.Len() != 0 || service.authChallengePool.Len() != 0 {
		t.Fatal("expired auth state was kept")
	}

	first, err := service.getSharedKeyPair(time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	second, err := service.getSharedKeyPair(time.Minute)
	if err != nil || second.getKeyPairId() != first.getKeyPairId() {
		t.Fatal(err)
	}
	code, _, _, err := service.preAuth(time.Minute)
	if err != nil || code != http.StatusOK {
		t.Fatalf("preauth = %d %v", code, err)
	}

	for range maxPendingChallenges {
		challenge, err := newAuthChallenge(time.Minute)
		if err != nil {
			t.Fatal(err)
		}
		service.authChallengePool.Set(challenge.getChallengeId(), challenge)
	}
	code, _, _, err = service.preAuth(time.Minute)
	if err == nil || code != http.StatusServiceUnavailable {
		t.Fatalf("full challenge pool = %d %v", code, err)
	}
}

type loginSession struct {
	auth        *auth
	pair        *keyPair
	id          string
	challenge   *authChallenge
	challengeID string
	solution    string
}

func newAuthService(t *testing.T, handler *action.Handler) *auth {
	t.Helper()
	cache, err := lru.New[string, time.Time](100)
	if err != nil {
		t.Fatal(err)
	}
	return &auth{
		actionHandler:     handler,
		nonceCache:        cache,
		keyPairDataPool:   haxmap.New[string, *keyPair](),
		authChallengePool: haxmap.New[string, *authChallenge](),
	}
}

func newLoginSession(t *testing.T, handler *action.Handler) *loginSession {
	t.Helper()
	service := newAuthService(t, handler)
	pair, err := newKeyPair(time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	challenge, err := newAuthChallenge(time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	session := &loginSession{
		auth:        service,
		pair:        pair,
		id:          pair.getKeyPairId(),
		challenge:   challenge,
		challengeID: challenge.getChallengeId(),
		solution:    solveChallenge(challenge),
	}
	service.keyPairDataPool.Set(session.id, pair)
	service.authChallengePool.Set(session.challengeID, challenge)
	return session
}

func (s *loginSession) rotateChallenge() {
	challenge, err := newAuthChallenge(time.Minute)
	if err != nil {
		panic(err)
	}
	s.challenge = challenge
	s.challengeID = challenge.getChallengeId()
	s.solution = solveChallenge(challenge)
	s.auth.authChallengePool.Set(s.challengeID, challenge)
}

func (s *loginSession) login(credential []byte, captcha, agent, ip string) (int, string, error) {
	s.rotateChallenge()
	secret, nonce, payload := s.material(credential)
	return s.auth.login(s.id, secret, nonce, s.challengeID, s.solution, "captcha", captcha, payload, agent, ip)
}

func (s *loginSession) attempt(secret, nonce, payload, captcha string) (int, string, error) {
	s.rotateChallenge()
	return s.auth.login(s.id, secret, nonce, s.challengeID, s.solution, "captcha", captcha, payload, "", "")
}

func (s *loginSession) material(credential []byte) (string, string, string) {
	raw := make([]byte, 32)
	if _, err := rand.Read(raw); err != nil {
		panic(err)
	}
	encrypted, err := s.pair.rsaKeyPair.Encrypt([]byte(base64.StdEncoding.EncodeToString(raw)), true)
	if err != nil {
		panic(err)
	}
	return string(encrypted),
		base64.StdEncoding.EncodeToString(sealLogin(raw, []byte("ok"), []byte(s.id))),
		base64.StdEncoding.EncodeToString(sealLogin(raw, credential, []byte(s.id)))
}

func sealLogin(secret, plaintext, aad []byte) []byte {
	encryptor := newAES256GCM(secret)
	nonce := make([]byte, encryptor.gcm.NonceSize())
	_, _ = rand.Read(nonce)
	return encryptor.gcm.Seal(nonce, nonce, plaintext, aad)
}

func solveChallenge(challenge *authChallenge) string {
	prefix := strings.Repeat("0", int(challenge.seed[0]))
	body := challenge.seed[1:]
	for i := uint64(0); ; i++ {
		nonce := strconv.FormatUint(i, 10)
		var text strings.Builder
		text.Grow(len(body) + len(nonce))
		for _, value := range body {
			text.WriteRune(rune(value))
		}
		text.WriteString(nonce)
		sum := sha512.Sum512([]byte(text.String()))
		encoded := hex.EncodeToString(sum[:])
		if strings.HasPrefix(encoded, prefix) {
			return nonce + ":" + encoded
		}
	}
}

func mustCredential(t *testing.T, username, password string) []byte {
	t.Helper()
	payload, err := json.Marshal(map[string]string{"username": username, "password": password})
	if err != nil {
		t.Fatal(err)
	}
	return payload
}

func uniqueNonce() string {
	buf := make([]byte, 18)
	_, _ = rand.Read(buf)
	return base64.StdEncoding.EncodeToString(buf)
}
