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

func (h *auth) login(sessionId, secret, nonce, payload, userAgent, userIp string) (code int, userId string, err error) {
	nc, err := base64.StdEncoding.DecodeString(nonce)
	if err != nil {
		return http.StatusBadRequest, "", errors.New("failed to login: provided nonce is invalid")
	}
	ncStr := string(nc)

	if t, ok := h.nonceCache.Get(ncStr); ok && time.Since(t) <= time.Hour {
		return http.StatusForbidden, "", errors.New("failed to login: replay attack detected")
	}
	h.nonceCache.Add(ncStr, time.Now())

	kp, ok := h.keyPairDataPool.GetAndDel(sessionId)
	if !ok {
		return http.StatusUnauthorized, "", errors.New("failed to login: provided session ID is invalid")
	}
	if !kp.isKeyPairAlive() {
		return http.StatusUnauthorized, "", errors.New("failed to login: provided session ID has expired")
	}

	// Attempt to decrypt RSA encrypted AES secret from payload
	aesSecretBytesB64, err := kp.rsaKeyPair.Decrypt([]byte(secret), true)
	if err != nil {
		return http.StatusUnauthorized, "", fmt.Errorf("failed to login: failed to decrypt secret: %w", err)
	}
	aesSecretBytes, err := base64.StdEncoding.DecodeString(string(aesSecretBytesB64))
	if err != nil {
		return http.StatusUnauthorized, "", fmt.Errorf("failed to login: failed to decode secret: %w", err)
	}
	encryptor := newAES256GCM(aesSecretBytes)

	// Attempt to decrypt AES encrypted nonce
	if _, err := encryptor.decrypt(nc, []byte(sessionId)); err != nil {
		return http.StatusForbidden, "", fmt.Errorf("failed to login: malformed nonce received: %w", err)
	}

	pl, err := base64.StdEncoding.DecodeString(payload)
	if err != nil {
		return http.StatusBadRequest, "", errors.New("failed to login: provided payload is invalid")
	}
	credential, err := encryptor.decrypt(pl, []byte(sessionId))
	if err != nil {
		return http.StatusUnauthorized, "", fmt.Errorf("failed to login: failed to decrypt payload: %w", err)
	}

	var credentialMap map[string]any
	if err = json.Unmarshal(credential, &credentialMap); err != nil {
		return http.StatusUnauthorized, "", fmt.Errorf("failed to login: failed to unmarshal payload: %w", err)
	}
	username, usernameOk := credentialMap["username"].(string)
	password, passwordOk := credentialMap["password"].(string)
	captchaId, captchaIdOk := credentialMap["captcha_id"].(string)
	captchaVal, captchaValOk := credentialMap["captcha_solution"].(string)
	if !usernameOk || !passwordOk || !captchaIdOk || !captchaValOk {
		return http.StatusUnauthorized, "", errors.New("failed to login: provided credential format is invalid")
	}

	if !captcha.VerifyString(captchaId, captchaVal) {
		return http.StatusUnauthorized, "", fmt.Errorf("failed to login: provided captcha solution %s is invalid", captchaVal)
	}

	userId, err = h.actionHandler.SysUserLogin(username, password, userAgent, userIp)
	if err != nil {
		return http.StatusUnauthorized, "", fmt.Errorf("user login failed: %w", err)
	}

	return http.StatusOK, userId, nil
}
