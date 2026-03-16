package auth

import (
	"crypto/rand"
	"math/big"
	"time"
)

type authChallenge struct {
	ttl       time.Duration
	createdAt time.Time
	seed      []byte
}

func newAuthChallenge(ttl time.Duration) (*authChallenge, error) {
	seed := make([]byte, 1+96) // 1 byte for difficulty, 96 bytes for random data
	seed[0] = 0x04
	if _, err := rand.Read(seed[1:]); err != nil {
		return nil, err
	}

	ch := &authChallenge{
		createdAt: time.Now(),
		ttl:       ttl,
		seed:      seed,
	}
	return ch, nil
}

func (la *authChallenge) getChallengeId() string {
	const (
		chars  = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
		length = 20
	)
	b := make([]byte, length)
	for i := 0; i < length; i++ {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		b[i] = chars[n.Int64()]
	}

	return string(b)
}

func (la *authChallenge) isChallengeAlive() bool {
	return time.Now().Before(la.createdAt.Add(la.ttl))
}

func (la *authChallenge) verifyChallenge(id, solution string) bool {
	for i := byte(0); i < la.seed[0]; i++ {
		if solution[i] != '0' {
			return false
		}
	}

	return true
}
