package auth_jwt

const (
	UserIdKey  = "user_id"
	IsAdminKey = "is_admin"
)

type cache struct {
	IsAdmin bool
}

type token struct {
	Token    string `json:"token"`
	LifeTime int64  `json:"life_time"`
}
