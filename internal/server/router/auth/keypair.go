package auth

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/anyshake/observer/pkg/cryption"
)

type keyPair struct {
	ttl        time.Duration
	createAt   time.Time
	rsaKeyPair cryption.RsaKeyPair
}

func newKeyPair(ttl time.Duration) (*keyPair, error) {
	rsaKeyPair, err := cryption.New(2048)
	if err != nil {
		return nil, fmt.Errorf("failed to create RSA key pair: %w", err)
	}

	return &keyPair{
		ttl:        ttl,
		createAt:   time.Now(),
		rsaKeyPair: rsaKeyPair,
	}, nil
}

func (n *keyPair) getKeyPairId() string {
	_, pemPubKey, _ := n.rsaKeyPair.GetPEM(false)
	h := sha512.New()
	h.Write([]byte(pemPubKey))

	return hex.EncodeToString(h.Sum(nil))
}

func (n *keyPair) isKeyPairAlive() bool {
	return time.Now().Before(n.createAt.Add(n.ttl))
}
