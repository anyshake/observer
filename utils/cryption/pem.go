package cryption

import (
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
)

func (r *RsaKeyPair) GetPEM(base64Output bool) (privKey string, pubKey string, err error) {
	privKey = string(pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(r.PrivateKey),
	}))

	pubKeyRaw, err := x509.MarshalPKIXPublicKey(r.PublicKey)
	if err != nil {
		return "", "", err
	}
	pubKey = string(pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pubKeyRaw,
	}))

	if base64Output {
		privKey = base64.StdEncoding.EncodeToString([]byte(privKey))
		pubKey = base64.StdEncoding.EncodeToString([]byte(pubKey))
	}
	return privKey, pubKey, nil
}
