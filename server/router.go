package server

import (
	"com.geophone.observer/app/statistics"
	"com.geophone.observer/app/status"
	"com.geophone.observer/server/socket"
	"github.com/gin-gonic/gin"
)

func RegisterRouter(rg *gin.RouterGroup, options *ServerOptions) {
	rg.GET("/socket", func(c *gin.Context) {
		socket.WebsocketHandler(c, options.Message)
	})

	services := []ApiServices{
		&status.Status{},
		&statistics.Statistics{},
	}
	for _, s := range services {
		s.RegisterModule(rg, options.Message, options.Status)
	}
}
