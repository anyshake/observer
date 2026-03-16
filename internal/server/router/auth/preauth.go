package auth

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	"github.com/dchest/captcha"
)

func (h *auth) preAuth(ttl time.Duration) (code int, msg string, res any, err error) {
	kp, err := newKeyPair(ttl)
	if err != nil {
		errText := "failed to generate new RSA key pair"
		return http.StatusInternalServerError, errText, nil, fmt.Errorf("%s: %w", errText, err)
	}
	h.keyPairDataPool.Set(kp.getKeyPairId(), kp)

	_, pemPubKey, err := kp.rsaKeyPair.GetPEM(true)
	if err != nil {
		errText := "failed to create RSA public key for pre-auth"
		return http.StatusInternalServerError, errText, nil, fmt.Errorf("%s: %w", errText, err)
	}

	var buf bytes.Buffer
	captchaId := captcha.New()
	if err = captcha.WriteImage(&buf, captchaId, captcha.StdWidth, captcha.StdHeight); err != nil {
		errText := "failed to create captcha"
		return http.StatusInternalServerError, errText, nil, fmt.Errorf("%s: %w", errText, err)
	}

	challenge, err := newAuthChallenge(ttl)
	if err != nil {
		errText := "failed to create PoW challenge for pre-auth"
		return http.StatusInternalServerError, errText, nil, fmt.Errorf("%s: %w", errText, err)
	}
	challengeId := challenge.getChallengeId()
	h.authChallengePool.Set(challengeId, challenge)

	return http.StatusOK, "successfully created pre-auth key", map[string]any{
		"ttl":            ttl.Milliseconds(),
		"public_key":     pemPubKey,
		"challenge_id":   challengeId,
		"challenge_seed": base64.StdEncoding.EncodeToString(challenge.seed),
		"captcha_id":     captchaId,
		"captcha_img":    base64.StdEncoding.EncodeToString(buf.Bytes()),
	}, nil
}
