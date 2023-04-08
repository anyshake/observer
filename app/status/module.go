package status

import (
	"net/http"

	"com.geophone.observer/features/collector"
	"com.geophone.observer/server/response"
	"github.com/gin-gonic/gin"
)

func (s *Status) RegisterModule(rg *gin.RouterGroup, message *collector.Message, status *collector.Status) {
	rg.GET("/status", func(c *gin.Context) {
		c.JSON(http.StatusOK, response.MessageHandler(c, "成功取得软件状态", status))
	})
}
