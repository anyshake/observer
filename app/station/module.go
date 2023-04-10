package station

import (
	"net/http"

	"com.geophone.observer/app"
	"com.geophone.observer/server/response"
	"github.com/gin-gonic/gin"
)

func (s *Station) RegisterModule(rg *gin.RouterGroup, options *app.ServerOptions) {
	rg.GET("/station", func(c *gin.Context) {
		c.JSON(http.StatusOK, response.MessageHandler(c, "成功取得测站状态", GetStation(
			options.Message.UUID, options.Message.Station, *options.Status,
		)))
	})
}
