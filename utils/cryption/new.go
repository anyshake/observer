package cryption

import (
	"crypto/rand"
	"crypto/rsa"
)

func New(bits int) (RsaKeyPair, error) {
	privKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return RsaKeyPair{}, err
	}

	return RsaKeyPair{
		PrivateKey: privKey,
		PublicKey:  &privKey.PublicKey,
	}, nil
}
