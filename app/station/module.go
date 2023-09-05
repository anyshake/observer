package station

import (
	"github.com/bclswl0827/observer/app"
	"github.com/bclswl0827/observer/server/response"
	"github.com/gin-gonic/gin"
)

func (s *Station) RegisterModule(rg *gin.RouterGroup, options *app.ServerOptions) {
	rg.GET("/station", func(c *gin.Context) {
		response.Message(c, "Successfully read station status", getStation(options.FeatureOptions))
	})
}
