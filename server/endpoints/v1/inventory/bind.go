package inventory

import (
	"net/http"

	"github.com/anyshake/observer/drivers/explorer"
	"github.com/anyshake/observer/server/response"
	"github.com/anyshake/observer/services"
	"github.com/anyshake/observer/utils/logger"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// @Summary Station Inventory
// @Description Get SeisComP XML inventory, which contains meta data of the station. This API requires a valid JWT token if the server is in restricted mode.
// @Router /inventory [get]
// @Produce application/json, application/xml
// @Security ApiKeyAuth
// @Param format query string false "Format of the inventory, available options are `json` or `xml`", default is `xml` if not specified."
// @Param Authorization header string false "Bearer JWT token, only required when the server is in restricted mode."
func (i *Inventory) Bind(rg *gin.RouterGroup, jwtHandler *jwt.GinJWTMiddleware, options *services.Options) error {
	var explorerDeps *explorer.ExplorerDependency
	err := options.Dependency.Invoke(func(deps *explorer.ExplorerDependency) error {
		explorerDeps = deps
		return nil
	})
	if err != nil {
		return err
	}

	var handlerFunc []gin.HandlerFunc
	if options.Config.Server.Restrict {
		handlerFunc = append(handlerFunc, jwtHandler.MiddlewareFunc())
	}
	handlerFunc = append(handlerFunc, func(c *gin.Context) {
		var req request
		err := c.ShouldBind(&req)
		if err != nil {
			logger.GetLogger(i.GetApiName()).Errorln(err)
			response.Message(c, options.TimeSource, "request body is not valid", http.StatusBadRequest, nil)
			return
		}

		inventory := i.handleInventory(options.Config, explorerDeps)
		if req.Format == "json" {
			response.Message(c, options.TimeSource, "successfully get station inventory in JSON format", http.StatusOK, inventory)
			return
		}
		c.Data(http.StatusOK, "application/xml", []byte(inventory))
	})

	rg.GET("/inventory", handlerFunc...)
	return nil
}
