package station

import (
	"com.geophone.observer/app"
	"com.geophone.observer/server/response"
	"github.com/gin-gonic/gin"
)

func (s *Station) RegisterModule(rg *gin.RouterGroup, options *app.ServerOptions) {
	rg.GET("/station", func(c *gin.Context) {
		response.Message(c, "成功取得测站状态", getStation(options.FeatureOptions))
	})
}
