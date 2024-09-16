package user

import (
	"time"

	"github.com/alphadose/haxmap"
	"github.com/anyshake/observer/utils/cryption"
)

type User struct {
	keyPairDataPool *haxmap.Map[string, keyPairData] // key: SHA-1 hash of the public key, value: keyPairData
}

type keyPairData struct {
	expiration time.Duration
	createAt   time.Time
	rsaKeyPair cryption.RsaKeyPair
}

type request struct {
	Action   string `form:"action" json:"action" xml:"action" binding:"required,oneof=profile list create remove edit preauth"`
	Password string `form:"password" json:"password" xml:"password"` // Encrypted with RSA public key
	Username string `form:"username" json:"username" xml:"username"` // Encrypted with RSA public key
	UserId   string `form:"user_id" json:"user_id" xml:"user_id"`    // Encrypted with RSA public key, required only for edit action
	Admin    bool   `form:"admin" json:"admin" xml:"admin"`
	Nonce    string `form:"nonce" json:"nonce" xml:"nonce"`
}

type user struct {
	Admin     bool   `json:"admin"`
	UserId    string `json:"user_id"` // To prevent integer overflow in JavaScript
	Username  string `json:"username"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
	LastLogin int64  `json:"last_login"`
	UserIp    string `json:"user_ip"`
	UserAgent string `json:"user_agent"`
}
