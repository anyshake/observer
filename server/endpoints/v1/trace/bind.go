package trace

import (
	"fmt"
	"net/http"
	"time"

	"github.com/anyshake/observer/drivers/explorer"
	"github.com/anyshake/observer/server/response"
	"github.com/anyshake/observer/services"
	"github.com/anyshake/observer/utils/logger"
	"github.com/anyshake/observer/utils/seisevent"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// @Summary Seismic Trace
// @Description This API retrieves seismic events from the specified data source, including essential information such as event time, location, magnitude, depth and estimated distance and arrival time from the station. This API requires a valid JWT token if the server is in restricted mode.
// @Router /trace [post]
// @Produce application/json
// @Security ApiKeyAuth
// @Param source formData string true "Use `list` to get available sources first, then choose one and request again to get events"
// @Param Authorization header string false "Bearer JWT token, only required when the server is in restricted mode."
func (t *Trace) Bind(rg *gin.RouterGroup, jwtHandler *jwt.GinJWTMiddleware, options *services.Options) error {
	var explorerDeps *explorer.ExplorerDependency
	err := options.Dependency.Invoke(func(deps *explorer.ExplorerDependency) error {
		explorerDeps = deps
		return nil
	})
	if err != nil {
		return err
	}

	seisSources := seisevent.New(30 * time.Second)
	var handlerFunc []gin.HandlerFunc
	if options.Config.Server.Restrict {
		handlerFunc = append(handlerFunc, jwtHandler.MiddlewareFunc())
	}
	handlerFunc = append(handlerFunc, func(c *gin.Context) {
		var req request
		err := c.ShouldBind(&req)
		if err != nil {
			logger.GetLogger(t.GetApiName()).Errorln(err)
			response.Message(c, options.TimeSource, "request body is not valid", http.StatusBadRequest, nil)
			return
		}

		switch req.Source {
		case "list":
			var properties []seisevent.DataSourceProperty
			for _, source := range seisSources {
				properties = append(properties, source.GetProperty())
			}
			response.Message(c, options.TimeSource, "successfully get properties of available data source", http.StatusOK, properties)
		default:
			source, ok := seisSources[req.Source]
			if !ok {
				err := fmt.Errorf("provided data source %s is not available", req.Source)
				logger.GetLogger(t.GetApiName()).Errorln(err)
				response.Message(c, options.TimeSource, err.Error(), http.StatusBadRequest, nil)
				return
			}
			var (
				latitude  = explorerDeps.Config.GetLatitude()
				longitude = explorerDeps.Config.GetLongitude()
			)
			events, err := source.GetEvents(latitude, longitude)
			if err != nil {
				logger.GetLogger(t.GetApiName()).Errorln(err)
				response.Message(c, options.TimeSource, err.Error(), http.StatusInternalServerError, nil)
				return
			}
			response.Message(c, options.TimeSource, "successfully get the list of earthquake events", http.StatusOK, events)
		}
	})

	rg.POST("/trace", handlerFunc...)
	return nil
}
