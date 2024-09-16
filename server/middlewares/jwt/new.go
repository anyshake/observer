package jwt

import (
	"net/http"
	"os"
	"time"

	"github.com/anyshake/observer/drivers/dao/tables"
	"github.com/anyshake/observer/server/response"
	"github.com/anyshake/observer/utils/timesource"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func New(timeSource *timesource.Source, logger logrus.FieldLogger, db *gorm.DB, expiration time.Duration) (*jwt.GinJWTMiddleware, error) {
	// Hostname as Realm
	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	// Unique secret key for signing
	secret, err := createSecret()
	if err != nil {
		return nil, err
	}

	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Authenticator: func(c *gin.Context) (any, error) {
			userId, ok := c.MustGet("user_id").(int64)
			if !ok {
				return nil, jwt.ErrMissingLoginValues
			}
			return map[string]any{"user_id": userId}, nil
		},
		PayloadFunc: func(data any) jwt.MapClaims {
			val, ok := data.(map[string]any)
			if !ok {
				return jwt.MapClaims{}
			}
			return jwt.MapClaims{"user_id": val["user_id"]}
		},
		Authorizator: func(data any, c *gin.Context) bool {
			userId := int64(data.(float64))
			// To check if user ID is in database
			var sysUserModel tables.SysUser
			err := db.
				Table(sysUserModel.GetName()).
				Where("user_id = ?", userId).
				First(&sysUserModel)
			return sysUserModel.UserId != 0 && err.Error == nil
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			logger.Warnln(message)
			response.Message(c, timeSource, "access denied due to invalid authorization token", http.StatusUnauthorized, nil)
		},
		LoginResponse: func(c *gin.Context, code int, token string, t time.Time) {
			response.Message(c, timeSource, "login succeed and token has been created", http.StatusOK, Token{
				ExpiresAt: t.UnixMilli(),
				Token:     token,
			})
		},
		RefreshResponse: func(c *gin.Context, code int, token string, t time.Time) {
			response.Message(c, timeSource, "token has been refreshed", http.StatusOK, Token{
				ExpiresAt: t.UnixMilli(),
				Token:     token,
			})
		},
		Key:         secret,
		Realm:       hostname,
		IdentityKey: "user_id",
		Timeout:     expiration,
		MaxRefresh:  expiration,
		TimeFunc:    timeSource.Get,
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
