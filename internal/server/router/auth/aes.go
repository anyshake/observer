package auth

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha512"
	"errors"

	"golang.org/x/crypto/hkdf"
)

type aes_256_gcm struct {
	gcm cipher.AEAD
}

func newAES256GCM(secret []byte) *aes_256_gcm {
	key := make([]byte, 32)
	_, err := hkdf.New(sha512.New, secret, nil, nil).Read(key)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil
	}

	return &aes_256_gcm{gcm: gcm}
}

func (impl *aes_256_gcm) decrypt(data, aad []byte) ([]byte, error) {
	nonceSize := impl.gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	return impl.gcm.Open(nil, nonce, ciphertext, aad)
}
