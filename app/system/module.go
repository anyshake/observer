package system

import (
	"com.geophone.observer/features/collector"
	"com.geophone.observer/server/response"
	"github.com/gin-gonic/gin"
)

func (s *System) RegisterModule(rg *gin.RouterGroup, status *collector.Status) {
	rg.GET("/system", func(c *gin.Context) {
		c.JSON(200, response.MessageHandler(
			c, "成功取得系统资讯", status,
		))
	})
}
