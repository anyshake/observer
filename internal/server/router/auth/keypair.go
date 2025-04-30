package auth

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/anyshake/observer/pkg/cryption"
)

type keyPair struct {
	expiration time.Duration
	createAt   time.Time
	rsaKeyPair cryption.RsaKeyPair
}

func (n *keyPair) init(expiration time.Duration) error {
	rsaKeyPair, err := cryption.New(2048)
	if err != nil {
		return fmt.Errorf("failed to create RSA key pair: %w", err)
	}

	n.createAt = time.Now()
	n.expiration = expiration
	n.rsaKeyPair = rsaKeyPair

	return nil
}

func (n *keyPair) getNonce() string {
	_, pemPubKey, _ := n.rsaKeyPair.GetPEM(false)
	h := sha1.New()
	h.Write([]byte(pemPubKey))

	return hex.EncodeToString(h.Sum(nil))
}

func (n *keyPair) isValid() bool {
	return time.Now().Before(n.createAt.Add(n.expiration))
}
