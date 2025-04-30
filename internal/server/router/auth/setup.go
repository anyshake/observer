package auth

import (
	"net/http"
	"time"

	"github.com/alphadose/haxmap"
	"github.com/anyshake/observer/internal/dao/action"
	"github.com/anyshake/observer/internal/server/response"
	"github.com/anyshake/observer/pkg/logger"
	"github.com/gin-gonic/gin"
)

func Setup(routerGroup *gin.RouterGroup, actionHandler *action.Handler, jwtMiddleware, loginCallback, refreshCallback gin.HandlerFunc) {
	h := auth{
		actionHandler:   actionHandler,
		keyPairDataPool: haxmap.New[string, keyPair](),
	}
	routerGroup.GET("/auth", jwtMiddleware, func(ctx *gin.Context) {
		response.Data(ctx, http.StatusOK, "user token is still valid", nil)
	})
	routerGroup.POST("/auth",
		func(c *gin.Context) {
			h.keyPairDataPool.ForEach(func(key string, nc keyPair) bool {
				if !nc.isValid() {
					h.keyPairDataPool.Del(key)
				}
				return true
			})
		},
		func(c *gin.Context) {
			var requestModel struct {
				Action     string `form:"action" json:"action" xml:"action" binding:"required,oneof=preauth login refresh"`
				Nonce      string `form:"nonce" json:"nonce" xml:"nonce"`
				Credential string `form:"credential" json:"credential" xml:"credential"` // Encrypted with RSA public key
			}
			if err := c.ShouldBind(&requestModel); err != nil {
				logger.GetLogger(LOG_PREFIX).Errorf("request body is not valid: %v", err)
				response.Error(c, http.StatusBadRequest, "request body is not valid")
				return
			}

			switch requestModel.Action {
			case "preauth":
				// Set 30 seconds expiration for the nonce
				// This is to prevent the nonce from being used for a long time
				// The client must request a new nonce if the nonce is expired
				// When 30 seconds reached, the nonce will be removed from the pool
				code, msg, res, err := h.preauth(30 * time.Second)
				if err != nil {
					logger.GetLogger(LOG_PREFIX).Errorln(err)
					response.Error(c, code, msg)
					return
				}
				response.Data(c, code, msg, res)
			case "login":
				code, userId, err := h.login(requestModel.Nonce, requestModel.Credential, c.GetHeader("User-Agent"), c.ClientIP())
				if err != nil {
					logger.GetLogger(LOG_PREFIX).Errorln(err)
					response.Error(c, code, err.Error())
					return
				}
				c.Set("user_id", userId)
				loginCallback(c)
			case "refresh":
				refreshCallback(c)
			}
		})
}
