package user

import (
	"net/http"
	"time"
)

func (h *User) handlePreauth(expiration time.Duration) (code int, msg string, res any, err error) {
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

	return http.StatusOK, "successfully created preauth key", map[string]any{
		"ttl":         expiration.Milliseconds(),
		"encrypt_key": pemPubKey,
	}, nil
}
