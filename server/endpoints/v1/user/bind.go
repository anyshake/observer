package user

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

// @Summary User Management
// @Description This API is used to manage user accounts, including creating, removing, and editing user profiles. This API only available in restricted mode and requires a valid JWT token.
// @Router /user [post]
// @Produce application/json
// @Security ApiKeyAuth
// @Param action formData string true "Specifies the action to be performed. Use `preauth` to get a Base64 RSA public key in PEM format, `profile` to get profile of current user, `list` to get list of all users (admin only), `create` to create a new user (admin only), `remove` to remove a user (admin only), and `edit` to edit a user (admin only)."
// @Param nonce formData string false "A unique string used to prevent replay attacks, required for the `create`, `remove`, `edit` actions and left empty for other actions. The nonce is the SHA-1 hash of the RSA public key from the pre-authentication stage and becomes invalid once the request is sent. It also expires if unused within the time-to-live (TTL) period, which is set during the pre-authentication stage."
// @Param user_id formData string false "The user ID to be removed or edited, required for the `remove` and `edit` actions and left empty for other actions. The user ID is encrypted with the RSA public key."
// @Param admin formData bool false "Specifies whether the user is an administrator, required for the `create` and `edit` actions and set to false in other actions."
// @Param username formData string false "The username of the user to be created or edited, required for the `create` and `edit` actions and left empty for other actions. The username is encrypted with the RSA public key."
// @Param password formData string false "The password of the user to be created or edited, required for the `create` and `edit` actions and left empty for other actions. The password is encrypted with the RSA public key."
// @Param Authorization header string true "Bearer JWT token."
func (h *User) Bind(rg *gin.RouterGroup, jwtHandler *jwt.GinJWTMiddleware, options *services.Options) error {
	h.keyPairDataPool = haxmap.New[string, keyPairData]()

	rg.POST("/user",
		func(c *gin.Context) {
			// Return 403 Forbidden if the server is not in restricted mode
			if !options.Config.Server.Restrict {
				err := errors.New("server is not in restricted mode")
				logger.GetLogger(h.GetApiName()).Errorln(err)
				response.Message(c, options.TimeSource, err.Error(), http.StatusForbidden, nil)
				return
			}
			// Remove all expired nonce
			h.keyPairDataPool.ForEach(func(key string, nc keyPairData) bool {
				if !nc.valid() {
					h.keyPairDataPool.Del(key)
				}
				return true
			})
		},
		jwtHandler.MiddlewareFunc(),
		func(c *gin.Context) {
			var req request
			err := c.ShouldBind(&req)
			if err != nil {
				logger.GetLogger(h.GetApiName()).Errorln(err)
				response.Message(c, options.TimeSource, "request body is not valid", http.StatusBadRequest, nil)
				return
			}

			userId := int64(jwt.ExtractClaims(c)["user_id"].(float64))
			switch req.Action {
			case "preauth":
				// Set 10 seconds of expiration
				// This action is used to get an public key
				// For the client to encrypt the username and password
				code, msg, res, err := h.handlePreauth(10 * time.Second)
				if err != nil {
					logger.GetLogger(h.GetApiName()).Errorln(err)
					response.Message(c, options.TimeSource, msg, code, nil)
					return
				}

				response.Message(c, options.TimeSource, msg, code, res)
			case "profile":
				code, msg, res, err := h.handleProfile(options.Database, userId)
				if err != nil {
					logger.GetLogger(h.GetApiName()).Errorln(err)
					response.Message(c, options.TimeSource, msg, code, nil)
					return
				}

				response.Message(c, options.TimeSource, msg, code, res)
			case "list":
				code, msg, res, err := h.handleList(options.Database, userId)
				if err != nil {
					logger.GetLogger(h.GetApiName()).Errorln(err)
					response.Message(c, options.TimeSource, msg, code, nil)
					return
				}

				response.Message(c, options.TimeSource, msg, code, res)
			case "create":
				code, msg, err := h.handleCreate(options.Database, userId, &req)
				if err != nil {
					logger.GetLogger(h.GetApiName()).Errorln(err)
					response.Message(c, options.TimeSource, msg, code, nil)
					return
				}

				response.Message(c, options.TimeSource, msg, code, nil)
			case "remove":
				code, msg, err := h.handleRemove(options.Database, userId, &req)
				if err != nil {
					logger.GetLogger(h.GetApiName()).Errorln(err)
					response.Message(c, options.TimeSource, msg, code, nil)
					return
				}

				response.Message(c, options.TimeSource, msg, code, nil)
			case "edit":
				code, msg, err := h.handleEdit(options.Database, userId, &req)
				if err != nil {
					logger.GetLogger(h.GetApiName()).Errorln(err)
					response.Message(c, options.TimeSource, msg, code, nil)
					return
				}

				response.Message(c, options.TimeSource, msg, code, nil)
			}
		})
	return nil
}
