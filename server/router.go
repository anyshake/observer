package server

import (
	"github.com/bclswl0827/observer/app"
	"github.com/bclswl0827/observer/app/devel"
	"github.com/bclswl0827/observer/app/history"
	"github.com/bclswl0827/observer/app/socket"
	"github.com/bclswl0827/observer/app/station"
	"github.com/bclswl0827/observer/app/trace"
	"github.com/gin-gonic/gin"
)

func RegisterRouter(rg *gin.RouterGroup, options *app.ServerOptions) {
	services := []ApiServices{
		&station.Station{},
		&history.History{},
		&socket.Socket{},
		&trace.Trace{},
		&devel.Devel{},
	}
	for _, s := range services {
		s.RegisterModule(rg, options)
	}
}
