package cryption

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"encoding/base64"
)

func (r *RsaKeyPair) Decrypt(cipherText []byte, base64Input bool) ([]byte, error) {
	if base64Input {
		data, err := base64.StdEncoding.DecodeString(string(cipherText))
		if err != nil {
			return nil, err
		}

		cipherText = data
	}

	return rsa.DecryptOAEP(sha1.New(), rand.Reader, r.PrivateKey, cipherText, nil)
}
