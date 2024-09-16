package auth

import (
	"errors"
	"net/http"
	"time"

	"github.com/alphadose/haxmap"
	"github.com/anyshake/observer/server/response"
	"github.com/anyshake/observer/services"
	"github.com/anyshake/observer/utils/logger"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// @Summary User Authentication
// @Description In restricted mode, the client must log in to access other APIs. This API is used to checks the server's authentication status, issues an RSA public key for credential encryption, generates a captcha, authenticates the client, and signs or refreshes the JWT token. This API requires a valid JWT token if action is `refresh`.
// @Router /auth [post]
// @Produce application/json
// @Security ApiKeyAuth
// @Param action formData string true "Specifies the action to be performed. Use `inspect` to check the server's restriction status, `preauth` to get a Base64 RSA public key in PEM format and generate a Base64 captcha PNG image, `login` to authenticate the client using encrypted credentials, and `refresh` to refresh the JWT token."
// @Param nonce formData string false "A unique string used to prevent replay attacks, required for the `login` action and left empty for other actions. The nonce is the SHA-1 hash of the RSA public key from the pre-authentication stage and becomes invalid once the request is sent. It also expires if unused within the time-to-live (TTL) period, which is set during the pre-authentication stage."
// @Param credential formData string false "Base64 encrypted credential using the RSA public key, required for the `login` action and left empty for other actions. The decrypted credential is a JSON object that includes the username, password, captcha ID, and captcha solution. Example: `{ username: admin, password: admin, captcha_id: 123, captcha_solution: abc }`."
// @Param Authorization header string false "Bearer JWT token, only required for the `refresh` action."
func (h *Auth) Bind(rg *gin.RouterGroup, jwtHandler *jwt.GinJWTMiddleware, options *services.Options) error {
	h.keyPairDataPool = haxmap.New[string, keyPairData]()

	rg.POST("/auth",
		func(c *gin.Context) {
			// Remove all expired nonce
			h.keyPairDataPool.ForEach(func(key string, nc keyPairData) bool {
				if !nc.valid() {
					h.keyPairDataPool.Del(key)
				}
				return true
			})
		},
		func(c *gin.Context) {
			var req request
			err := c.ShouldBind(&req)
			if err != nil {
				logger.GetLogger(h.GetApiName()).Errorln(err)
				response.Message(c, options.TimeSource, "request body is not valid", http.StatusBadRequest, nil)
				return
			}

			switch req.Action {
			case "inspect":
				msg, res := h.handleInspect(options.Config.Server.Restrict)
				response.Message(c, options.TimeSource, msg, http.StatusOK, res)
			case "preauth":
				if !options.Config.Server.Restrict {
					err := errors.New("server is not in restricted mode")
					logger.GetLogger(h.GetApiName()).Errorln(err)
					response.Message(c, options.TimeSource, err.Error(), http.StatusForbidden, nil)
					return
				}

				// Set 30 seconds expiration for the nonce
				// This is to prevent the nonce from being used for a long time
				// The client must request a new nonce if the nonce is expired
				// When 30 seconds reached, the nonce will be removed from the pool
				code, msg, res, err := h.handlePreauth(30 * time.Second)
				if err != nil {
					logger.GetLogger(h.GetApiName()).Errorln(err)
					response.Message(c, options.TimeSource, msg, code, nil)
					return
				}

				response.Message(c, options.TimeSource, msg, code, res)
			case "login":
				if !options.Config.Server.Restrict {
					err := errors.New("server is not in restricted mode")
					logger.GetLogger(h.GetApiName()).Errorln(err)
					response.Message(c, options.TimeSource, err.Error(), http.StatusForbidden, nil)
					return
				}

				code, userId, err := h.handleLogin(options.Database, &req, c.ClientIP(), c.GetHeader("User-Agent"))
				if err != nil {
					logger.GetLogger(h.GetApiName()).Errorln(err)
					response.Message(c, options.TimeSource, err.Error(), code, nil)
					return
				}

				c.Set("user_id", userId) // Set user_id to context so that jwt can use it
				jwtHandler.LoginHandler(c)
			case "refresh":
				if !options.Config.Server.Restrict {
					err := errors.New("server is not in restricted mode")
					logger.GetLogger(h.GetApiName()).Errorln(err)
					response.Message(c, options.TimeSource, err.Error(), http.StatusForbidden, nil)
					return
				}

				jwtHandler.RefreshHandler(c)
			}
		})

	return nil
}
