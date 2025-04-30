package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/dchest/captcha"
)

func (h *auth) login(nonce, credential, userAgent, userIp string) (code int, userId string, err error) {
	nc, ok := h.keyPairDataPool.GetAndDel(nonce)
	if !ok {
		return http.StatusUnauthorized, "", errors.New("failed to login: provided nonce is invalid")
	}

	if !nc.isValid() {
		return http.StatusUnauthorized, "", errors.New("failed to login: provided nonce has expired")
	}

	credentialBytes, err := nc.rsaKeyPair.Decrypt([]byte(credential), true)
	if err != nil {
		return http.StatusUnauthorized, "", fmt.Errorf("failed to login: failed to decrypt credential: %w", err)
	}
	var credentialMap map[string]any
	err = json.Unmarshal(credentialBytes, &credentialMap)
	if err != nil {
		return http.StatusUnauthorized, "", fmt.Errorf("failed to login: failed to unserialize credential: %w", err)
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
