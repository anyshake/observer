package auth_jwt

const (
	UserIdKey  = "user_id"
	IsAdminKey = "is_admin"
)

type token struct {
	Token    string `json:"token"`
	LifeTime int64  `json:"life_time"`
}
