package auth

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/anyshake/observer/pkg/logger"
	"github.com/dchest/captcha"
)

// verifyCaptchaString is the captcha check used by login. Tests replace it.
var verifyCaptchaString = captcha.VerifyString

var (
	errInvalidLoginRequest  = errors.New("invalid login request")
	errAuthenticationFailed = errors.New("authentication failed")
)

func (h *auth) login(ctx context.Context, sessionId, secret, nonce, challengeId, challengeSolution, captchaId, captchaVal, payload, userAgent, userIp string) (code int, userId string, err error) {
	fail := func(status int, public error, internal string, args ...any) (int, string, error) {
		if internal != "" {
			logger.GetLogger(LOG_PREFIX).Errorf(internal, args...)
		}
		return status, "", public
	}

	if sessionId == "" || secret == "" || nonce == "" || challengeId == "" || challengeSolution == "" || captchaId == "" || captchaVal == "" || payload == "" {
		return fail(http.StatusBadRequest, errInvalidLoginRequest, "login rejected: missing required fields")
	}
	if len(sessionId) > 128 || len(secret) > 344 || len(nonce) != base64.StdEncoding.EncodedLen(encryptedNonceBytes) ||
		len(challengeId) > 20 || len(challengeSolution) > 149 || len(captchaId) > 64 || len(captchaVal) > 6 || len(payload) > 4096 {
		return fail(http.StatusBadRequest, errInvalidLoginRequest, "login rejected: invalid field length")
	}

	// 1. Bound the nonce before decoding; unauthenticated data is never cached.
	nc, err := base64.StdEncoding.DecodeString(nonce)
	if err != nil || len(nc) != encryptedNonceBytes {
		return fail(http.StatusBadRequest, errInvalidLoginRequest, "login rejected: invalid nonce encoding: %v", err)
	}

	// 2. Verify PoW challenge solution
	challenge, ok := h.takeChallenge(challengeId)
	if !ok {
		return fail(http.StatusUnauthorized, errAuthenticationFailed, "login rejected: unknown challenge id")
	}
	defer captchaStore.Get(challenge.captchaID, true)
	if !challenge.isChallengeAlive() {
		return fail(http.StatusUnauthorized, errAuthenticationFailed, "login rejected: challenge expired")
	}
	if !challenge.verifyChallenge(challengeSolution) {
		return fail(http.StatusUnauthorized, errAuthenticationFailed, "login rejected: invalid PoW solution")
	}

	// 3. Verify captcha solution
	if !verifyCaptchaString(captchaId, captchaVal) {
		return fail(http.StatusUnauthorized, errAuthenticationFailed, "login rejected: invalid captcha")
	}

	// 4. Extract RSA key pair from session ID
	// Key pairs are shared across preauth requests within TTL, so do not delete on use.
	kp, ok := h.keyPairDataPool.Get(sessionId)
	if !ok {
		return fail(http.StatusUnauthorized, errAuthenticationFailed, "login rejected: unknown session id")
	}
	if !kp.isKeyPairAlive() {
		h.keyPairDataPool.Del(sessionId)
		return fail(http.StatusUnauthorized, errAuthenticationFailed, "login rejected: session expired")
	}

	// 5. Attempt to decrypt RSA encrypted AES secret from payload
	aesSecretBytesB64, err := kp.rsaKeyPair.Decrypt([]byte(secret), true)
	if err != nil {
		return fail(http.StatusUnauthorized, errAuthenticationFailed, "login rejected: secret decrypt failed: %v", err)
	}
	aesSecretBytes, err := base64.StdEncoding.DecodeString(string(aesSecretBytesB64))
	if err != nil || len(aesSecretBytes) != 32 {
		return fail(http.StatusUnauthorized, errAuthenticationFailed, "login rejected: secret decode failed: %v", err)
	}
	encryptor := newAES256GCM(aesSecretBytes)
	if encryptor == nil {
		return fail(http.StatusUnauthorized, errAuthenticationFailed, "login rejected: failed to init AES-GCM")
	}

	// 6. Additional check to ensure data integrity (defense in depth)
	if _, err := encryptor.decrypt(nc, []byte(sessionId)); err != nil {
		return fail(http.StatusForbidden, errAuthenticationFailed, "login rejected: malformed nonce: %v", err)
	}

	// 7. Attempt to decrypt AES encrypted payload containing credential
	pl, err := base64.StdEncoding.DecodeString(payload)
	if err != nil {
		return fail(http.StatusBadRequest, errInvalidLoginRequest, "login rejected: invalid payload encoding: %v", err)
	}
	credential, err := encryptor.decrypt(pl, []byte(sessionId))
	if err != nil {
		return fail(http.StatusUnauthorized, errAuthenticationFailed, "login rejected: payload decrypt failed: %v", err)
	}

	// 8. Unmarshal credential and perform login logic
	var credentialMap map[string]any
	if err = json.Unmarshal(credential, &credentialMap); err != nil {
		return fail(http.StatusUnauthorized, errAuthenticationFailed, "login rejected: credential unmarshal failed: %v", err)
	}
	username, usernameOk := credentialMap["username"].(string)
	password, passwordOk := credentialMap["password"].(string)
	if !usernameOk || !passwordOk {
		return fail(http.StatusUnauthorized, errAuthenticationFailed, "login rejected: credential format invalid")
	}

	userId, err = h.actionHandler.SysUserLogin(ctx, username, password, userAgent, userIp)
	if err != nil {
		return fail(http.StatusUnauthorized, errAuthenticationFailed, "login rejected: user auth failed for %s: %v", username, err)
	}
	if code, err := h.rememberNonce(nc); err != nil {
		return fail(code, err, "login rejected: nonce replay or cache capacity exceeded")
	}

	return http.StatusOK, userId, nil
}
