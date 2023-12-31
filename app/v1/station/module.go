package station

import (
	"github.com/anyshake/observer/app"
	"github.com/anyshake/observer/server/response"
	"github.com/gin-gonic/gin"
)

// @Summary AnyShake Observer station status
// @Description Get Observer station status including system information, memory usage, disk usage, CPU usage, ADC information, geophone information, and location information
// @Router /station [get]
// @Produce application/json
// @Success 200 {object} response.HttpResponse{data=System} "Successfully read station status"
func (s *Station) RegisterModule(rg *gin.RouterGroup, options *app.ServerOptions) {
	rg.GET("/station", func(c *gin.Context) {
		response.Message(c, "Successfully read station information", getSystem(options.FeatureOptions))
	})
}
