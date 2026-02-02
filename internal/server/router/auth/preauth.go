package auth

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	"github.com/dchest/captcha"
)

func (h *auth) preauth(ttl time.Duration) (code int, msg string, res any, err error) {
	kp, err := newKeyPair(ttl)
	if err != nil {
		errText := "failed to generate new RSA key pair"
		return http.StatusInternalServerError, errText, nil, fmt.Errorf("%s: %w", errText, err)
	}
	h.keyPairDataPool.Set(kp.getKeyPairId(), kp)

	_, pemPubKey, err := kp.rsaKeyPair.GetPEM(true)
	if err != nil {
		errText := "failed to create RSA public key for preauth"
		return http.StatusInternalServerError, errText, nil, fmt.Errorf("%s: %w", errText, err)
	}

	var buf bytes.Buffer
	captchaId := captcha.New()
	err = captcha.WriteImage(&buf, captchaId, captcha.StdWidth, captcha.StdHeight)
	if err != nil {
		return http.StatusInternalServerError, "failed to create captcha", nil, fmt.Errorf("failed to create captcha: %w", err)
	}

	return http.StatusOK, "successfully created preauth key", map[string]any{
		"ttl":         ttl.Milliseconds(),
		"public_key":  pemPubKey,
		"captcha_id":  captchaId,
		"captcha_img": base64.StdEncoding.EncodeToString(buf.Bytes()),
	}, nil
}
