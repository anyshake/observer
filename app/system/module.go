package system

import (
	"com.geophone.observer/features/collector"
	"com.geophone.observer/server/response"
	"github.com/gin-gonic/gin"
)

func (s *System) RegisterModule(rg *gin.RouterGroup, message *collector.Message, status *collector.Status) {
	rg.GET("/system", func(c *gin.Context) {
		c.JSON(200, response.MessageHandler(
			c, "成功取得系统资讯", struct {
				Version   string `json:"version"`
				Processor string `json:"processor"`
				Memory    string `json:"memory"`
				Distro    string `json:"distro"`
			}{
				Version:   "1.0.0",
				Processor: "Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz",
				Memory:    "16GB",
				Distro:    "Debian GNU/Linux 10 (buster)",
			},
		))
	})
}
