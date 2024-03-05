package server

import (
	"time"

	"github.com/anyshake/observer/app"
	"github.com/anyshake/observer/app/v1/devel"
	"github.com/anyshake/observer/app/v1/history"
	"github.com/anyshake/observer/app/v1/inventory"
	"github.com/anyshake/observer/app/v1/mseed"
	"github.com/anyshake/observer/app/v1/socket"
	"github.com/anyshake/observer/app/v1/station"
	"github.com/anyshake/observer/app/v1/trace"
	"github.com/anyshake/observer/server/middleware/limit"
	"github.com/gin-gonic/gin"
)

func registerRouterV1(rg *gin.RouterGroup, options *app.ServerOptions) {
	if options.RateFactor > 0 {
		rateFactor := int64(options.RateFactor)
		rg.Use(limit.RateLimit(time.Second, rateFactor, rateFactor))
	}
	services := []ApiServices{
		&station.Station{},
		&history.History{},
		&socket.Socket{},
		&trace.Trace{},
		&mseed.MSeed{},
		&devel.Devel{},
		&inventory.Inventory{},
	}
	for _, s := range services {
		s.RegisterModule(rg, options)
	}
}

// TODO: Support FDSNWS
// func registerRouterFDSNWS(rg *gin.RouterGroup, options *app.ServerOptions) {
// }
