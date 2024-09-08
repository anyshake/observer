package server

import (
	"fmt"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	v1 "github.com/anyshake/observer/api/v1"
	"github.com/anyshake/observer/api/v1/devel"
	"github.com/anyshake/observer/api/v1/history"
	"github.com/anyshake/observer/api/v1/inventory"
	"github.com/anyshake/observer/api/v1/mseed"
	"github.com/anyshake/observer/api/v1/socket"
	"github.com/anyshake/observer/api/v1/station"
	"github.com/anyshake/observer/api/v1/trace"
	v2 "github.com/anyshake/observer/api/v2"
	"github.com/anyshake/observer/server/middleware/limit"
	"github.com/gin-gonic/gin"
)

func registerEndpointsV1(routerGroup *gin.RouterGroup, options *Options) error {
	if options.RateFactor > 0 {
		rateFactor := int64(options.RateFactor)
		routerGroup.Use(limit.RateLimit(time.Second, rateFactor, rateFactor))
	}
	resolver := &v1.Resolver{Options: options.ServicesOptions}
	services := []v1.Endpoint{
		&station.Station{},
		&history.History{},
		&socket.Socket{},
		&trace.Trace{},
		&mseed.MSeed{},
		&inventory.Inventory{},
	}
	for _, s := range services {
		err := s.Register(routerGroup, resolver)
		if err != nil {
			return err
		}
	}
	if options.DebugMode {
		err := v1.Endpoint(&devel.Devel{}).Register(routerGroup, resolver)
		if err != nil {
			return err
		}
	}

	return nil
}

func registerEndpointsV2(routerGroup *gin.RouterGroup, options *Options) {
	apiEndpoint := handler.NewDefaultServer(v2.NewExecutableSchema(v2.Config{
		Resolvers: &v2.Resolver{Options: options.ServicesOptions},
	}))
	routerGroup.POST("/graph", func(ctx *gin.Context) {
		apiEndpoint.ServeHTTP(ctx.Writer, ctx.Request)
	})
	if options.DebugMode {
		routerGroup.GET("/graph", func(ctx *gin.Context) {
			playground.Handler("AnyShake Observer APIv2", fmt.Sprintf("%s/v2/graph", options.ApiPrefix)).ServeHTTP(ctx.Writer, ctx.Request)
		})
	}
}
