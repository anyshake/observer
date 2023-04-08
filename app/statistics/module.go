package statistics

import (
	"net/http"

	"com.geophone.observer/features/collector"
	"com.geophone.observer/server/response"
	"github.com/gin-gonic/gin"
)

func (s *Statistics) RegisterModule(rg *gin.RouterGroup, message *collector.Message, status *collector.Status) {
	rg.GET("/statistics", func(c *gin.Context) {
		c.JSON(http.StatusOK, response.MessageHandler(c, "成功取得系统资讯", GetStatistics()))
	})
}
