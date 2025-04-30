package auth_jwt

type token struct {
	Token    string `json:"token"`
	LifeTime int64  `json:"life_time"`
}
