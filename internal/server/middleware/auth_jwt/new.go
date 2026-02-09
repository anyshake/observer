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
	"github.com/gin-gonic/gin"
)

func New(timeSource *timesource.Source, actionHandler *action.Handler, expiration time.Duration) (*jwt.GinJWTMiddleware, error) {
	secret, err := createJwtSecret()
	if err != nil {
		return nil, fmt.Errorf("failed to create JWT secret: %w", err)
	}

	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Authorizator: func(data any, c *gin.Context) bool {
			userId, ok := data.(string)
			if !ok {
				return false
			}
			userModel, err := actionHandler.SysUserGetByUserId(userId)
			if userModel.UserId != "" && err == nil {
				c.Set(IsAdminKey, userModel.IsAdmin == model.ADMIN)
				return true
			}
			return false
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
