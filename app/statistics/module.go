package statistics

import (
	"net/http"

	"com.geophone.observer/app"
	"com.geophone.observer/server/response"
	"github.com/gin-gonic/gin"
)

func (s *Statistics) RegisterModule(rg *gin.RouterGroup, options *app.ServerOptions) {
	rg.GET("/statistics", func(c *gin.Context) {
		c.JSON(http.StatusOK, response.MessageHandler(c, "成功取得系统状态", GetStatistics()))
	})
}
