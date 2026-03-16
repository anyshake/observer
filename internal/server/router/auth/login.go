package auth

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/dchest/captcha"
)

func (h *auth) login(sessionId, secret, nonce, challengeId, challengeSolution, captchaId, captchaVal, payload, userAgent, userIp string) (code int, userId string, err error) {
	if sessionId == "" || secret == "" || nonce == "" || challengeId == "" || challengeSolution == "" || captchaId == "" || captchaVal == "" || payload == "" {
		return http.StatusBadRequest, "", errors.New("failed to login: missing required fields in request")
	}

	// 1. Check if the nonce is valid and not reused
	nc, err := base64.StdEncoding.DecodeString(nonce)
	if err != nil {
		return http.StatusBadRequest, "", errors.New("failed to login: provided nonce is invalid")
	}
	ncStr := string(nc)
	if t, ok := h.nonceCache.Get(ncStr); ok && time.Since(t) <= time.Hour {
		return http.StatusForbidden, "", errors.New("failed to login: replay attack detected")
	}
	h.nonceCache.Add(ncStr, time.Now())

	// 2. Verify PoW challenge solution
	challenge, ok := h.authChallengePool.GetAndDel(challengeId)
	if !ok {
		return http.StatusUnauthorized, "", errors.New("failed to login: invalid pre-auth challenge ID")
	}
	if !challenge.isChallengeAlive() {
		return http.StatusUnauthorized, "", errors.New("failed to login: pre-auth challenge has expired, please try again")
	}
	if !challenge.verifyChallenge(challengeId, challengeSolution) {
		return http.StatusUnauthorized, "", errors.New("failed to login: invalid challenge solution")
	}

	// 3. Verify captcha solution
	if !captcha.VerifyString(captchaId, captchaVal) {
		return http.StatusUnauthorized, "", fmt.Errorf("failed to login: provided captcha solution %s is invalid", captchaVal)
	}

	// 4. Extract RSA key pair from session ID
	kp, ok := h.keyPairDataPool.GetAndDel(sessionId)
	if !ok {
		return http.StatusUnauthorized, "", errors.New("failed to login: provided session ID is invalid")
	}
	if !kp.isKeyPairAlive() {
		return http.StatusUnauthorized, "", errors.New("failed to login: provided session ID has expired")
	}

	// 5. Attempt to decrypt RSA encrypted AES secret from payload
	aesSecretBytesB64, err := kp.rsaKeyPair.Decrypt([]byte(secret), true)
	if err != nil {
		return http.StatusUnauthorized, "", fmt.Errorf("failed to login: failed to decrypt secret: %w", err)
	}
	aesSecretBytes, err := base64.StdEncoding.DecodeString(string(aesSecretBytesB64))
	if err != nil {
		return http.StatusUnauthorized, "", fmt.Errorf("failed to login: failed to decode secret: %w", err)
	}
	encryptor := newAES256GCM(aesSecretBytes)

	// 6. Additional check to ensure data integrity (defense in depth)
	if _, err := encryptor.decrypt(nc, []byte(sessionId)); err != nil {
		return http.StatusForbidden, "", fmt.Errorf("failed to login: malformed nonce received: %w", err)
	}

	// 7. Attempt to decrypt AES encrypted payload containing credential
	pl, err := base64.StdEncoding.DecodeString(payload)
	if err != nil {
		return http.StatusBadRequest, "", errors.New("failed to login: provided payload is invalid")
	}
	credential, err := encryptor.decrypt(pl, []byte(sessionId))
	if err != nil {
		return http.StatusUnauthorized, "", fmt.Errorf("failed to login: failed to decrypt payload: %w", err)
	}

	// 8. Unmarshal credential and perform login logic
	var credentialMap map[string]any
	if err = json.Unmarshal(credential, &credentialMap); err != nil {
		return http.StatusUnauthorized, "", fmt.Errorf("failed to login: failed to unmarshal payload: %w", err)
	}
	username, usernameOk := credentialMap["username"].(string)
	password, passwordOk := credentialMap["password"].(string)
	if !usernameOk || !passwordOk {
		return http.StatusUnauthorized, "", errors.New("failed to login: provided credential format is invalid")
	}

	userId, err = h.actionHandler.SysUserLogin(username, password, userAgent, userIp)
	if err != nil {
		return http.StatusUnauthorized, "", fmt.Errorf("user login failed: %w", err)
	}

	return http.StatusOK, userId, nil
}
