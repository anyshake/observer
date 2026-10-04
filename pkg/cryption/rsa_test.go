package cryption_test

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"testing"

	"github.com/anyshake/observer/pkg/cryption"
)

func TestRSAEncryptionAndPEM(t *testing.T) {
	t.Parallel()
	keys, err := cryption.New(2048)
	if err != nil {
		t.Fatal(err)
	}
	plaintext := []byte("observer login secret")
	for _, encoded := range []bool{false, true} {
		ciphertext, err := keys.Encrypt(plaintext, encoded)
		if err != nil {
			t.Fatal(err)
		}
		got, err := keys.Decrypt(ciphertext, encoded)
		if err != nil || !bytes.Equal(got, plaintext) {
			t.Fatalf("round trip (base64=%v) = %q, %v", encoded, got, err)
		}
		priv, pub, err := keys.GetPEM(encoded)
		if err != nil {
			t.Fatal(err)
		}
		if encoded {
			privateBytes, err := base64.StdEncoding.DecodeString(priv)
			if err != nil {
				t.Fatal(err)
			}
			publicBytes, err := base64.StdEncoding.DecodeString(pub)
			if err != nil {
				t.Fatal(err)
			}
			priv, pub = string(privateBytes), string(publicBytes)
		}
		privateBlock, _ := pem.Decode([]byte(priv))
		publicBlock, _ := pem.Decode([]byte(pub))
		if privateBlock == nil || publicBlock == nil {
			t.Fatal("GetPEM() returned invalid PEM")
		}
		privateKey, err := x509.ParsePKCS1PrivateKey(privateBlock.Bytes)
		if err != nil {
			t.Fatal(err)
		}
		publicKey, err := x509.ParsePKIXPublicKey(publicBlock.Bytes)
		if err != nil {
			t.Fatal(err)
		}
		if !keys.PrivateKey.Equal(privateKey) || !keys.PublicKey.Equal(publicKey) {
			t.Fatal("PEM keys differ from generated key pair")
		}
	}
	// Verify compatibility with the browser's RSA-OAEP/SHA-1 convention.
	standardCiphertext, err := rsa.EncryptOAEP(sha1.New(), rand.Reader, keys.PublicKey, plaintext, nil)
	if err != nil {
		t.Fatal(err)
	}
	if got, err := keys.Decrypt(standardCiphertext, false); err != nil || !bytes.Equal(got, plaintext) {
		t.Fatalf("standard OAEP decryption = %q, %v", got, err)
	}
	if _, err := keys.Encrypt(make([]byte, keys.PublicKey.Size()), false); err == nil {
		t.Fatal("oversized plaintext was accepted")
	}
	standardCiphertext[0] ^= 1
	if got, err := keys.Decrypt(standardCiphertext, false); err == nil || len(got) != 0 {
		t.Fatal("tampered ciphertext was accepted")
	}
	if _, err := keys.Decrypt([]byte("not-base64!"), true); err == nil {
		t.Fatal("invalid base64 was accepted")
	}
	if _, err := keys.Decrypt([]byte{1, 2, 3}, false); err == nil {
		t.Fatal("truncated ciphertext was accepted")
	}
}
