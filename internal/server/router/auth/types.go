package auth

import (
	"sync"
	"time"

	"github.com/alphadose/haxmap"
	"github.com/anyshake/observer/internal/dao/action"
)

const LOG_PREFIX = "restful_api_auth"

// maxPendingChallenges caps in-flight PoW challenges to avoid unbounded memory growth.
const maxPendingChallenges = 256

const (
	maxAuthBodyBytes    = 8 << 10
	encryptedNonceBytes = 12 + 16 + 16 // GCM IV, random plaintext, authentication tag
	maxRememberedNonces = 4096
	nonceTTL            = time.Hour
	maxPendingPerSource = 8
	maxPreAuthPerMinute = 20
	maxAuthSources      = 4096
	authCleanupInterval = time.Second
)

// sharedKeyPairTTL controls how long a reused RSA key pair remains valid.
// This avoids generating a 2048-bit key on every preauth request.
const sharedKeyPairTTL = 5 * time.Minute

type auth struct {
	actionHandler     *action.Handler                     // action handler for accessing the database
	nonceCache        map[[32]byte]time.Time              // SHA-256 digest -> expiry
	keyPairDataPool   *haxmap.Map[string, *keyPair]       // key: SHA-512 hash of the public key, value: keyPair
	authChallengePool *haxmap.Map[string, *authChallenge] // key: Random string as ID, value: authChallenge
	keyPairMu         sync.Mutex
	stateMu           sync.Mutex
	sources           map[string]*authSource
	reservations      int // Includes challenges still generating their CAPTCHA.
}

type authSource struct {
	resetAt  time.Time
	requests int
	pending  int
}
