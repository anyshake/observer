package user

import (
	"crypto/sha1"
	"encoding/hex"
	"time"

	"github.com/anyshake/observer/utils/cryption"
)

func (n *keyPairData) init(expiration time.Duration) error {
	rsaKeyPair, err := cryption.New(2048)
	if err != nil {
		return err
	}

	n.expiration = expiration
	n.createAt = time.Now()
	n.rsaKeyPair = rsaKeyPair

	return nil
}

func (n *keyPairData) nonce() string {
	_, pemPubKey, _ := n.rsaKeyPair.GetPEM(false)
	h := sha1.New()
	h.Write([]byte(pemPubKey))

	return hex.EncodeToString(h.Sum(nil))
}

func (n *keyPairData) valid() bool {
	return time.Now().Before(n.createAt.Add(n.expiration))
}
