package auth_jwt

import (
	"fmt"
	"net/http"
	"time"

	"github.com/anyshake/observer/internal/dao/action"
	"github.com/anyshake/observer/internal/dao/model"
	"github.com/anyshake/observer/internal/server/response"
	"github.com/anyshake/observer/pkg/timesource"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/dgraph-io/ristretto"
	"github.com/gin-gonic/gin"
)

func New(timeSource *timesource.Source, actionHandler *action.Handler, expiration time.Duration) (*jwt.GinJWTMiddleware, error) {
	userCache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e6,
		MaxCost:     1 << 26,
		BufferItems: 64,
	})
	if err != nil {
		return nil, err
	}

	secret, err := createJwtSecret()
	if err != nil {
		return nil, fmt.Errorf("failed to create JWT secret: %w", err)
	}

	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Authorizator: func(data any, c *gin.Context) bool {
			userId, ok := data.(string)
			if !ok || userId == "" {
				return false
			}

			if v, ok := userCache.Get(userId); ok {
				if ua, ok := v.(*cache); ok {
					c.Set(IsAdminKey, ua.IsAdmin)
					return true
				}
				userCache.Del(userId)
			}

			userModel, err := actionHandler.SysUserGetByUserId(userId)
			if err != nil {
				return false
			}

			ua := &cache{IsAdmin: userModel.IsAdmin == model.ADMIN}

			cacheTTL := time.Minute
			if userModel.IsAdmin != model.ADMIN {
				cacheTTL = 10 * time.Minute
			}
			userCache.SetWithTTL(userId, ua, 1, cacheTTL)

			c.Set(IsAdminKey, ua.IsAdmin)
			return true
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			baseMessage := "access denied due to invalid authorization token"
			response.Error(c, http.StatusUnauthorized, fmt.Sprintf("%s: %s", baseMessage, message))
		},
		Authenticator: func(c *gin.Context) (any, error) {
			userId, ok := c.MustGet(UserIdKey).(string)
			if !ok {
				return nil, jwt.ErrInvalidAuthHeader
			}
			return map[string]any{UserIdKey: userId}, nil
		},
		PayloadFunc: func(data any) jwt.MapClaims {
			val, ok := data.(map[string]any)
			if !ok {
				return nil
			}
			return jwt.MapClaims{UserIdKey: val[UserIdKey]}
		},
		LoginResponse: func(c *gin.Context, code int, tokenStr string, t time.Time) {
			response.Data(c, code, "login succeed and token has been created", token{
				Token:    tokenStr,
				LifeTime: expiration.Milliseconds(),
			})
		},
		RefreshResponse: func(c *gin.Context, code int, tokenStr string, t time.Time) {
			response.Data(c, http.StatusOK, "token has been refreshed", token{
				Token:    tokenStr,
				LifeTime: expiration.Milliseconds(),
			})
		},
		Key:         secret,
		IdentityKey: UserIdKey,
		Realm:       "anyshake-observer",
		Timeout:     expiration,
		MaxRefresh:  expiration,
		TimeFunc:    timeSource.Now,
		TokenLookup: "header: Authorization, query: token",
	})
	if err != nil {
		return nil, err
	}

	err = authMiddleware.MiddlewareInit()
	if err != nil {
		return nil, err
	}

	return authMiddleware, nil
}
