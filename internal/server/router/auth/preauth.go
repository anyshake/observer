package auth

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	"github.com/dchest/captcha"
)

func (h *auth) preauth(expiration time.Duration) (code int, msg string, res any, err error) {
	// Create a nonce for the client
	var nc keyPair
	err = nc.init(expiration)
	if err != nil {
		return http.StatusInternalServerError, "failed to create nonce", nil, fmt.Errorf("failed to create nonce: %w", err)
	}
	h.keyPairDataPool.Set(nc.getNonce(), nc)

	// Return base64 encoded public key to the client
	_, pemPubKey, err := nc.rsaKeyPair.GetPEM(true)
	if err != nil {
		return http.StatusInternalServerError, "failed to create preauth key", nil, fmt.Errorf("failed to create preauth key: %w", err)
	}

	// Generate a captcha for the client
	var buf bytes.Buffer
	captchaId := captcha.New()
	err = captcha.WriteImage(&buf, captchaId, captcha.StdWidth, captcha.StdHeight)
	if err != nil {
		return http.StatusInternalServerError, "failed to create captcha", nil, fmt.Errorf("failed to create captcha: %w", err)
	}

	return http.StatusOK, "successfully created preauth key", map[string]any{
		"ttl":         expiration.Milliseconds(),
		"encrypt_key": pemPubKey,
		"captcha_id":  captchaId,
		"captcha_img": base64.StdEncoding.EncodeToString(buf.Bytes()),
	}, nil
}
