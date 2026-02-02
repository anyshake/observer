package auth

import (
	"time"

	"github.com/alphadose/haxmap"
	"github.com/anyshake/observer/internal/dao/action"
	lru "github.com/hashicorp/golang-lru/v2"
)

const LOG_PREFIX = "restful_api_auth"

type auth struct {
	actionHandler   *action.Handler               // action handler for accessing the database
	nonceCache      *lru.Cache[string, time.Time] // key: nonce, value: time.Time
	keyPairDataPool *haxmap.Map[string, *keyPair] // key: SHA-1 hash of the public key, value: keyPair
}
