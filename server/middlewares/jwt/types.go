package jwt

type Token struct {
	ExpiresAt int64  `json:"expires_at"`
	Token     string `json:"token"`
}
