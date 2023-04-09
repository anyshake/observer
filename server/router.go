package server

import (
	"com.geophone.observer/app"
	"com.geophone.observer/app/history"
	"com.geophone.observer/app/socket"
	"com.geophone.observer/app/station"
	"com.geophone.observer/app/statistics"
	"github.com/gin-gonic/gin"
)

func RegisterRouter(rg *gin.RouterGroup, options *app.ServerOptions) {
	services := []ApiServices{
		&station.Station{},
		&statistics.Statistics{},
		&history.History{},
		&socket.Socket{},
	}
	for _, s := range services {
		s.RegisterModule(rg, options)
	}
}
