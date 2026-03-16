package auth

import (
	"net/http"
	"time"

	"github.com/alphadose/haxmap"
	"github.com/anyshake/observer/internal/dao/action"
	"github.com/anyshake/observer/internal/server/middleware/auth_jwt"
	"github.com/anyshake/observer/internal/server/response"
	"github.com/anyshake/observer/pkg/logger"
	"github.com/gin-gonic/gin"
	lru "github.com/hashicorp/golang-lru/v2"
)

func Setup(routerGroup *gin.RouterGroup, actionHandler *action.Handler, jwtMiddleware, loginCallback, refreshCallback gin.HandlerFunc) {
	nonceCache, _ := lru.New[string, time.Time](100000)
	h := auth{
		actionHandler:     actionHandler,
		nonceCache:        nonceCache,
		keyPairDataPool:   haxmap.New[string, *keyPair](),
		authChallengePool: haxmap.New[string, *authChallenge](),
	}
	routerGroup.GET("/auth", jwtMiddleware, func(ctx *gin.Context) {
		response.Data(ctx, http.StatusOK, "user token is still valid", nil)
	})
	routerGroup.POST("/auth",
		func(c *gin.Context) {
			for key, pair := range h.keyPairDataPool.Iterator() {
				if !pair.isKeyPairAlive() {
					h.keyPairDataPool.Del(key)
				}
			}
			for key, attempt := range h.authChallengePool.Iterator() {
				if !attempt.isChallengeAlive() {
					h.authChallengePool.Del(key)
				}
			}
		},
		func(c *gin.Context) {
			var requestModel struct {
				Action            string `form:"action" json:"action" xml:"action" binding:"required"`
				Session           string `form:"session" json:"session" xml:"session"`
				Secret            string `form:"secret" json:"secret" xml:"secret"`                                     // AES secret encrypted with RSA public key
				Nonce             string `form:"nonce" json:"nonce" xml:"nonce"`                                        // nonce encrypted with AES secret
				ChallengeId       string `form:"challenge_id" json:"challenge_id" xml:"challenge_id"`                   // ID of the PoW challenge issued during pre-auth
				ChallengeSolution string `form:"challenge_solution" json:"challenge_solution" xml:"challenge_solution"` // solution to the PoW challenge
				CaptchaId         string `form:"captcha_id" json:"captcha_id" xml:"captcha_id"`                         // ID of the captcha issued during pre-auth
				CaptchaVal        string `form:"captcha_val" json:"captcha_val" xml:"captcha_val"`                      // solution to the captcha
				Payload           string `form:"payload" json:"payload" xml:"payload"`                                  // credential encrypted with AES secret
			}
			if err := c.ShouldBind(&requestModel); err != nil {
				logger.GetLogger(LOG_PREFIX).Errorf("request body is not valid: %v", err)
				response.Error(c, http.StatusBadRequest, "request body is not valid")
				return
			}

			switch requestModel.Action {
			case "preauth":
				code, msg, res, err := h.preAuth(30 * time.Second)
				if err != nil {
					logger.GetLogger(LOG_PREFIX).Errorln(err)
					response.Error(c, code, msg)
					return
				}
				response.Data(c, code, msg, res)
			case "login":
				code, userId, err := h.login(
					requestModel.Session,
					requestModel.Secret,
					requestModel.Nonce,
					requestModel.ChallengeId,
					requestModel.ChallengeSolution,
					requestModel.CaptchaId,
					requestModel.CaptchaVal,
					requestModel.Payload,
					c.GetHeader("User-Agent"),
					c.ClientIP(),
				)
				if err != nil {
					logger.GetLogger(LOG_PREFIX).Errorln(err)
					response.Error(c, code, err.Error())
					return
				}
				c.Set(auth_jwt.UserIdKey, userId)
				loginCallback(c)
			case "refresh":
				refreshCallback(c)
			default:
				response.Error(c, http.StatusBadRequest, "invalid action")
			}
		})
}
