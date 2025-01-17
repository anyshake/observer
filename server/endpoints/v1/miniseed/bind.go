package miniseed

import (
	"errors"
	"net/http"

	"github.com/anyshake/observer/server/response"
	"github.com/anyshake/observer/services"
	"github.com/anyshake/observer/services/miniseed"
	"github.com/anyshake/observer/utils/logger"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// @Summary MiniSEED Data
// @Description This API returns a list of MiniSEED files or exports a specific MiniSEED file. This API requires a valid JWT token if the server is in restricted mode.
// @Router /miniseed [post]
// @Produce application/json
// @Produce application/octet-stream
// @Security ApiKeyAuth
// @Param action formData string true "Action to be performed, Use `list` to get list of MiniSEED files, `export` to export a specific MiniSEED file."
// @Param name formData string false "A valid filename of the MiniSEED file to be exported, only required when action is `export`."
// @Param Authorization header string false "Bearer JWT token, only required when the server is in restricted mode."
func (h *MiniSEED) Bind(rg *gin.RouterGroup, jwtHandler *jwt.GinJWTMiddleware, options *services.Options) error {
	// Get MiniSEED service configuration
	var miniseedService miniseed.MiniSeedService
	serviceConfig, ok := options.Config.Services[miniseedService.GetServiceName()]
	if !ok {
		// Set default configuration
		serviceConfig = map[string]any{
			"enable":    false,
			"path":      "",
			"lifecycle": 0.0,
		}
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
			response.Message(c, options.TimeSource, "miniSEED service is disabled by admin", http.StatusOK, nil)
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
				response.Message(c, options.TimeSource, "failed to list miniSEED files", http.StatusInternalServerError, nil)
				return
			}
			response.Message(c, options.TimeSource, "successfully get list of miniSEED files", http.StatusOK, fileList)
		case "export":
			if len(req.Name) == 0 {
				err := errors.New("name of miniSEED file cannot be empty")
				logger.GetLogger(h.GetApiName()).Errorln(err)
				response.Message(c, options.TimeSource, err.Error(), http.StatusBadRequest, nil)
				return
			}

			fileBytes, err := h.handleExport(basePath, req.Name)
			if err != nil {
				logger.GetLogger(h.GetApiName()).Errorln(err)
				response.Message(c, options.TimeSource, "failed to export miniSEED file", http.StatusInternalServerError, nil)
				return
			}

			response.File(c, req.Name, fileBytes)
		}
	})

	rg.POST("/miniseed", handlerFunc...)
	return nil
}
