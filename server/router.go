package server

import (
	"time"

	"github.com/anyshake/observer/app"
	"github.com/anyshake/observer/app/devel"
	"github.com/anyshake/observer/app/history"
	"github.com/anyshake/observer/app/mseed"
	"github.com/anyshake/observer/app/socket"
	"github.com/anyshake/observer/app/station"
	"github.com/anyshake/observer/app/trace"
	"github.com/anyshake/observer/server/middleware/limit"
	"github.com/gin-gonic/gin"
)

func RegisterRouter(rg *gin.RouterGroup, options *app.ServerOptions) {
	rg.Use(limit.RateLimit(time.Second, CAPACITY, CAPACITY))
	services := []ApiServices{
		&station.Station{},
		&history.History{},
		&socket.Socket{},
		&trace.Trace{},
		&mseed.MSeed{},
		&devel.Devel{},
	}
	for _, s := range services {
		s.RegisterModule(rg, options)
	}
}
