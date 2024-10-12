package helicorder

import (
	"errors"
	"net/http"

	"github.com/anyshake/observer/server/response"
	"github.com/anyshake/observer/services"
	"github.com/anyshake/observer/services/helicorder"
	"github.com/anyshake/observer/utils/logger"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// @Summary Helicorder Image
// @Description This API returns a list of helicorders or exports a specific helicorder image. This API requires a valid JWT token if the server is in restricted mode.
// @Router /helicorder [post]
// @Produce application/json
// @Produce application/octet-stream
// @Security ApiKeyAuth
// @Param action formData string true "Action to be performed, Use `list` to get list of helicorders, `export` to export a specific helicorder image."
// @Param name formData string false "A valid filename of the helicorder image to be exported, only required when action is `export`."
// @Param Authorization header string false "Bearer JWT token, only required when the server is in restricted mode."
func (h *HeliCorder) Bind(rg *gin.RouterGroup, jwtHandler *jwt.GinJWTMiddleware, options *services.Options) error {
	// Get helicorder service configuration
	var helicorderService helicorder.HelicorderService
	serviceConfig, ok := options.Config.Services[helicorderService.GetServiceName()]
	if !ok {
		return errors.New("failed to get configuration for helicorder service")
	}
	enable := serviceConfig.(map[string]any)["enable"].(bool)
	basePath := serviceConfig.(map[string]any)["path"].(string)
	lifeCycle := int(serviceConfig.(map[string]any)["lifecycle"].(float64))

	var handlerFunc []gin.HandlerFunc
	if options.Config.Server.Restrict {
		handlerFunc = append(handlerFunc, jwtHandler.MiddlewareFunc())
	}
	handlerFunc = append(handlerFunc, func(c *gin.Context) {
		var req request
		err := c.ShouldBind(&req)
		if err != nil {
			logger.GetLogger(h.GetApiName()).Errorln(err)
			response.Message(c, options.TimeSource, "request body is not valid", http.StatusBadRequest, nil)
			return
		}

		if !enable {
			response.Message(c, options.TimeSource, "helicorder service is disabled by admin", http.StatusOK, nil)
			return
		}

		switch req.Action {
		case "list":
			fileList, err := h.handleList(
				basePath,
				options.Config.Stream.Station,
				options.Config.Stream.Network,
				lifeCycle,
			)
			if err != nil {
				logger.GetLogger(h.GetApiName()).Errorln(err)
				response.Message(c, options.TimeSource, "failed to list helicorder files", http.StatusInternalServerError, nil)
				return
			}
			response.Message(c, options.TimeSource, "successfully get list of helicorder files", http.StatusOK, fileList)
		case "export":
			if len(req.Name) == 0 {
				err := errors.New("name of helicorder file cannot be empty")
				logger.GetLogger(h.GetApiName()).Errorln(err)
				response.Message(c, options.TimeSource, err.Error(), http.StatusBadRequest, nil)
				return
			}

			fileBytes, err := h.handleExport(basePath, req.Name)
			if err != nil {
				logger.GetLogger(h.GetApiName()).Errorln(err)
				response.Message(c, options.TimeSource, "failed to export helicorder image", http.StatusInternalServerError, nil)
				return
			}

			response.File(c, req.Name, fileBytes)
		}
	})

	rg.POST("/helicorder", handlerFunc...)
	return nil
}
