package auth

import (
	"crypto/sha512"
	"encoding/hex"
	"strings"
	"testing"
	"testing/synctest"
	"time"
)

func TestVerifyChallenge(t *testing.T) {
	t.Parallel()

	// The browser encodes seed bytes as Unicode characters before hashing UTF-8.
	seed := []byte{0, 0x80, 0xff, 'A'}
	sum := sha512.Sum512([]byte("\xc2\x80\xc3\xbfA42"))
	hash := hex.EncodeToString(sum[:])
	solution := "42:" + hash
	tests := []struct {
		name     string
		seed     []byte
		solution string
		want     bool
	}{
		{"valid browser encoding", seed, solution, true},
		{"uppercase hash", seed, "42:" + strings.ToUpper(hash), true},
		{"empty seed", nil, solution, false},
		{"missing challenge bytes", []byte{0}, solution, false},
		{"difficulty exceeds hash length", []byte{129, 0x80, 0xff, 'A'}, solution, false},
		{"insufficient proof of work", []byte{128, 0x80, 0xff, 'A'}, solution, false},
		{"different seed", []byte{0, 0x80, 0xff, 'B'}, solution, false},
		{"changed nonce", seed, "43:" + hash, false},
		{"forged hash", seed, "42:" + strings.Repeat("0", 128), false},
		{"empty solution", seed, "", false},
		{"missing separator", seed, "42" + hash, false},
		{"empty nonce", seed, ":" + hash, false},
		{"negative nonce", seed, "-1:" + hash, false},
		{"fractional nonce", seed, "4.2:" + hash, false},
		{"overflowing nonce", seed, "18446744073709551616:" + hash, false},
		{"empty hash", seed, "42:", false},
		{"non-hex hash", seed, "42:" + strings.Repeat("z", 128), false},
		{"short hash", seed, "42:" + hash[:126], false},
		{"odd-length hash", seed, "42:" + hash[:127], false},
		{"trailing data", seed, solution + ":extra", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			challenge := &authChallenge{seed: tt.seed}
			if got := challenge.verifyChallenge(tt.solution); got != tt.want {
				t.Fatalf("verifyChallenge() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChallengeExpiration(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		challenge := &authChallenge{createdAt: time.Now(), ttl: time.Minute}
		if !challenge.isChallengeAlive() {
			t.Fatal("new challenge is expired")
		}
		time.Sleep(time.Minute - time.Nanosecond)
		if !challenge.isChallengeAlive() {
			t.Fatal("challenge expired before its TTL")
		}
		time.Sleep(time.Nanosecond)
		if challenge.isChallengeAlive() {
			t.Fatal("challenge remains valid at its expiration time")
		}
	})
}
