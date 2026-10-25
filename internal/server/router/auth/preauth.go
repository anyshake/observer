package auth

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image/color"
	"image/png"
	"net/http"
	"time"

	"github.com/dchest/captcha"
)

var captchaStore = captcha.NewMemoryStore(captcha.CollectNum, captcha.Expiration)

func init() {
	// Share the store with captcha.New and captcha.VerifyString.
	captcha.SetCustomStore(captchaStore)
}

func (h *auth) cleanupExpiredAuthState() {
	for key, pair := range h.keyPairDataPool.Iterator() {
		if !pair.isKeyPairAlive() {
			h.keyPairDataPool.Del(key)
		}
	}
	h.stateMu.Lock()
	defer h.stateMu.Unlock()
	h.cleanupStateLocked(time.Now())
}

func (h *auth) getSharedKeyPair(ttl time.Duration) (*keyPair, error) {
	h.keyPairMu.Lock()
	defer h.keyPairMu.Unlock()

	for key, pair := range h.keyPairDataPool.Iterator() {
		if pair.isKeyPairAlive() {
			return pair, nil
		}
		h.keyPairDataPool.Del(key)
	}

	kp, err := newKeyPair(ttl)
	if err != nil {
		return nil, err
	}
	h.keyPairDataPool.Set(kp.getKeyPairId(), kp)
	return kp, nil
}

func (h *auth) preAuth(source string, ttl time.Duration) (code int, msg string, res any, err error) {
	if code, err := h.reserveChallenge(source); err != nil {
		return code, err.Error(), nil, err
	}
	var challenge *authChallenge
	var challengeId, captchaId string
	defer func() {
		if err != nil {
			challenge = nil
			captchaStore.Get(captchaId, true)
		}
		h.finishChallenge(source, challengeId, challenge)
	}()

	kp, err := h.getSharedKeyPair(sharedKeyPairTTL)
	if err != nil {
		errText := "failed to generate new RSA key pair"
		return http.StatusInternalServerError, errText, nil, fmt.Errorf("%s: %w", errText, err)
	}

	_, pemPubKey, err := kp.rsaKeyPair.GetPEM(true)
	if err != nil {
		errText := "failed to create RSA public key for pre-auth"
		return http.StatusInternalServerError, errText, nil, fmt.Errorf("%s: %w", errText, err)
	}

	var buf bytes.Buffer
	captchaId = captcha.New()
	captchaDigits := captchaStore.Get(captchaId, false)
	if captchaDigits == nil {
		errText := "failed to create captcha"
		return http.StatusInternalServerError, errText, nil, fmt.Errorf("%s: %w", errText, captcha.ErrNotFound)
	}
	captchaImage := captcha.NewImage(captchaId, captchaDigits, captcha.StdWidth, captcha.StdHeight)
	// Palette index 0 is the background; make it opaque white before encoding.
	captchaImage.Palette[0] = color.White
	if err = png.Encode(&buf, captchaImage.Paletted); err != nil {
		errText := "failed to encode captcha image"
		return http.StatusInternalServerError, errText, nil, fmt.Errorf("%s: %w", errText, err)
	}

	challenge, err = newAuthChallenge(ttl)
	if err != nil {
		errText := "failed to create PoW challenge for pre-auth"
		return http.StatusInternalServerError, errText, nil, fmt.Errorf("%s: %w", errText, err)
	}
	challengeId = challenge.getChallengeId()
	challenge.captchaID = captchaId

	return http.StatusOK, "successfully created pre-auth key", map[string]any{
		"ttl":            ttl.Milliseconds(),
		"public_key":     pemPubKey,
		"challenge_id":   challengeId,
		"challenge_seed": base64.StdEncoding.EncodeToString(challenge.seed),
		"captcha_id":     captchaId,
		"captcha_img":    base64.StdEncoding.EncodeToString(buf.Bytes()),
	}, nil
}
