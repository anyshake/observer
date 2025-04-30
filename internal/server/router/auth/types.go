package auth

import (
	"github.com/alphadose/haxmap"
	"github.com/anyshake/observer/internal/dao/action"
)

const LOG_PREFIX = "restful_api_auth"

type auth struct {
	actionHandler   *action.Handler              // action handler for accessing the database
	keyPairDataPool *haxmap.Map[string, keyPair] // key: SHA-1 hash of the public key, value: keyPair
}
