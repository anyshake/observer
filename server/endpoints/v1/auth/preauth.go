package auth

import (
	"bytes"
	"encoding/base64"
	"net/http"
	"time"

	"github.com/dchest/captcha"
)

func (h *Auth) handlePreauth(expiration time.Duration) (code int, msg string, res any, err error) {
	// Create a nonce for the client
	var nc keyPairData
	err = nc.init(expiration)
	if err != nil {
		return http.StatusInternalServerError, "failed to create nonce", nil, err
	}
	h.keyPairDataPool.Set(nc.nonce(), nc)

	// Return base64 encoded public key to the client
	_, pemPubKey, err := nc.rsaKeyPair.GetPEM(true)
	if err != nil {
		return http.StatusInternalServerError, "failed to create preauth key", nil, err
	}

	// Generate a captcha for the client
	var buf bytes.Buffer
	captchaId := captcha.New()
	err = captcha.WriteImage(&buf, captchaId, captcha.StdWidth, captcha.StdHeight)
	if err != nil {
		return http.StatusInternalServerError, "failed to create captcha", nil, err
	}

	return http.StatusOK, "successfully created preauth key", map[string]any{
		"ttl":         expiration.Milliseconds(),
		"encrypt_key": pemPubKey,
		"captcha_id":  captchaId,
		"captcha_img": base64.StdEncoding.EncodeToString(buf.Bytes()),
	}, nil
}
