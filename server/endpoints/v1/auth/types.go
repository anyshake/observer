package auth

import (
	"time"

	"github.com/alphadose/haxmap"
	"github.com/anyshake/observer/utils/cryption"
)

type Auth struct {
	keyPairDataPool *haxmap.Map[string, keyPairData] // key: SHA-1 hash of the public key, value: keyPairData
}

type keyPairData struct {
	expiration time.Duration
	createAt   time.Time
	rsaKeyPair cryption.RsaKeyPair
}

type request struct {
	Action     string `form:"action" json:"action" xml:"action" binding:"required,oneof=inspect login preauth refresh"`
	Nonce      string `form:"nonce" json:"nonce" xml:"nonce"`
	Credential string `form:"credential" json:"credential" xml:"credential"` // Encrypted with RSA public key
}
