package endpoints

import (
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	v1 "github.com/anyshake/observer/server/endpoints/v1"
	"github.com/anyshake/observer/server/endpoints/v1/auth"
	"github.com/anyshake/observer/server/endpoints/v1/devel"
	"github.com/anyshake/observer/server/endpoints/v1/helicorder"
	"github.com/anyshake/observer/server/endpoints/v1/history"
	"github.com/anyshake/observer/server/endpoints/v1/inventory"
	"github.com/anyshake/observer/server/endpoints/v1/miniseed"
	"github.com/anyshake/observer/server/endpoints/v1/socket"
	"github.com/anyshake/observer/server/endpoints/v1/station"
	"github.com/anyshake/observer/server/endpoints/v1/trace"
	"github.com/anyshake/observer/server/endpoints/v1/user"
	v2 "github.com/anyshake/observer/server/endpoints/v2"
	"github.com/anyshake/observer/server/middlewares/limit"
	"github.com/anyshake/observer/services"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

func SetupApiV1(routerGroup *gin.RouterGroup, jwtHandler *jwt.GinJWTMiddleware, options *services.Options) error {
	rateLimitFactor := int64(options.Config.Server.Rate)
	if rateLimitFactor > 0 {
		routerGroup.Use(limit.New(options.TimeSource, time.Second, rateLimitFactor, rateLimitFactor))
	}

	services := []v1.Endpoint{
		&station.Station{},
		&history.History{},
		&socket.Socket{},
		&trace.Trace{},
		&miniseed.MiniSEED{},
		&inventory.Inventory{},
		&user.User{},
		&auth.Auth{},
		&helicorder.HeliCorder{},
	}
	for _, s := range services {
		err := s.Bind(routerGroup, jwtHandler, options)
		if err != nil {
			return err
		}
	}

	if options.Config.Server.Debug {
		err := v1.Endpoint(&devel.Devel{}).Bind(routerGroup, jwtHandler, options)
		if err != nil {
			return err
		}
	}

	return nil
}

func SetupApiV2(routerGroup *gin.RouterGroup, options *services.Options) {
	rateLimitFactor := int64(options.Config.Server.Rate)
	if rateLimitFactor > 0 {
		routerGroup.Use(limit.New(options.TimeSource, time.Second, rateLimitFactor, rateLimitFactor))
	}

	apiEndpoint := handler.NewDefaultServer(v2.NewExecutableSchema(v2.Config{
		Resolvers: &v2.Resolver{Options: options},
	}))
	routerGroup.POST("/graph", func(ctx *gin.Context) {
		apiEndpoint.ServeHTTP(ctx.Writer, ctx.Request)
	})
	if options.Config.Server.Debug {
		routerGroup.GET("/graph", func(ctx *gin.Context) {
			currentPath := ctx.Request.URL.Path
			playground.Handler("AnyShake Observer API v2", currentPath).ServeHTTP(ctx.Writer, ctx.Request)
		})
	}
}
