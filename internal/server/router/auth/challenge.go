package auth

import (
	"crypto/rand"
	"crypto/sha512"
	"crypto/subtle"
	"encoding/hex"
	"math/big"
	"strconv"
	"strings"
	"time"
)

type authChallenge struct {
	ttl       time.Duration
	createdAt time.Time
	seed      []byte
	source    string
	captchaID string
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
	for i := range length {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		if err != nil {
			var rb [1]byte
			_, _ = rand.Read(rb[:])
			b[i] = chars[int(rb[0])%len(chars)]
			continue
		}
		b[i] = chars[n.Int64()]
	}

	return string(b)
}

func (la *authChallenge) isChallengeAlive() bool {
	return time.Now().Before(la.createdAt.Add(la.ttl))
}

func (la *authChallenge) verifyChallenge(solution string) bool {
	if len(la.seed) < 2 {
		return false
	}

	difficulty := int(la.seed[0])
	if difficulty < 0 || difficulty > hex.EncodedLen(sha512.Size) {
		return false
	}

	nonceStr, hashHex, ok := strings.Cut(solution, ":")
	if !ok || nonceStr == "" || hashHex == "" {
		return false
	}
	if _, err := strconv.ParseUint(nonceStr, 10, 64); err != nil {
		return false
	}

	expected, err := hex.DecodeString(hashHex)
	if err != nil || len(expected) != sha512.Size {
		return false
	}

	challenge := la.seed[1:]
	var textBuilder strings.Builder
	textBuilder.Grow(len(challenge) + len(nonceStr))
	for _, b := range challenge {
		textBuilder.WriteRune(rune(b))
	}
	textBuilder.WriteString(nonceStr)

	sum := sha512.Sum512([]byte(textBuilder.String()))
	if subtle.ConstantTimeCompare(expected, sum[:]) != 1 {
		return false
	}

	actualHex := hex.EncodeToString(sum[:])
	return strings.HasPrefix(actualHex, strings.Repeat("0", difficulty))
}
