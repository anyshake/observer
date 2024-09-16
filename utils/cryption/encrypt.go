package cryption

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"encoding/base64"
)

func (r *RsaKeyPair) Encrypt(message []byte, base64Output bool) ([]byte, error) {
	result, err := rsa.EncryptOAEP(sha1.New(), rand.Reader, r.PublicKey, message, nil)
	if err != nil {
		return nil, err
	}

	if base64Output {
		return []byte(base64.StdEncoding.EncodeToString(result)), nil
	}

	return result, nil
}
