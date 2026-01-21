package server

import (
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/anyshake/observer/internal/server/middleware/auth_jwt"
	"github.com/anyshake/observer/internal/server/middleware/date_header"
	"github.com/anyshake/observer/internal/server/middleware/httplog"
	"github.com/anyshake/observer/internal/server/middleware/recovery"
	"github.com/anyshake/observer/internal/server/response"
	"github.com/anyshake/observer/internal/server/router/auth"
	"github.com/anyshake/observer/internal/server/router/export"
	"github.com/anyshake/observer/internal/server/router/files"
	graph_resolver "github.com/anyshake/observer/internal/server/router/graph"
	"github.com/anyshake/observer/internal/server/router/socket"
	"github.com/anyshake/observer/web"
	"github.com/gin-contrib/cors"
	gzipHandler "github.com/gin-contrib/gzip"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/contrib/secure"
	"github.com/gin-gonic/gin"
)

func (s *httpServer) Setup(listen string) error {
	s.engine = gin.New()

	s.engine.Use(date_header.New(s.resolver.TimeSource))
	s.engine.Use(recovery.New(s.log.Logger))
	s.engine.Use(httplog.New(s.log))
	s.engine.Use(gzipHandler.Gzip(gzip.BestCompression))
	s.engine.Use(secure.Secure(secure.Options{
		FrameDeny:             true,
		BrowserXssFilter:      true,
		ContentTypeNosniff:    true,
		ContentSecurityPolicy: "default-src 'self'; connect-src 'self' https://anyshake.org; style-src 'self' https://cdn.jsdelivr.net 'unsafe-inline'; script-src 'self' https://cdn.jsdelivr.net 'unsafe-inline' 'wasm-unsafe-eval'; font-src 'self' data:; img-src 'self' data: blob:;",
	}))
	if s.cors {
		s.engine.Use(cors.New(cors.Config{
			MaxAge:        12 * time.Hour,
			AllowOrigins:  []string{"*"},
			ExposeHeaders: []string{"Content-Length", "Content-Disposition"},
			AllowMethods:  []string{"GET", "POST", "PATCH"},
			AllowHeaders:  []string{"Content-Type", "Authorization"},
		}))
	}

	s.engine.NoRoute(func(ctx *gin.Context) {
		uri := ctx.Request.RequestURI
		msg := fmt.Sprintf("resource not found on %s", uri)
		response.Error(ctx, http.StatusNotFound, msg)
	})

	jwtHandler, err := auth_jwt.New(s.resolver.TimeSource, s.resolver.ActionHandler, AUTH_TIMEOUT)
	if err != nil {
		return fmt.Errorf("failed to create jwt handler: %w", err)
	}
	jwtMiddlewareFn := jwtHandler.MiddlewareFunc()

	api := s.engine.Group("/api")
	auth.Setup(api, s.resolver.ActionHandler, jwtMiddlewareFn, jwtHandler.LoginHandler, jwtHandler.RefreshHandler)
	export.Setup(api, s.resolver.ActionHandler, s.resolver.HardwareDev, jwtMiddlewareFn)
	socket.Setup(api, s.resolver.TimeSource, s.resolver.HardwareDev, jwtMiddlewareFn)
	files.Setup(api, s.resolver.ServiceMap, jwtMiddlewareFn)

	graphql := handler.NewDefaultServer(graph_resolver.NewExecutableSchema(graph_resolver.Config{Resolvers: s.resolver}))
	graphql.SetRecoverFunc(func(ctx context.Context, err any) (userMessage error) {
		s.log.Errorf("recovered from panic in GraphQL: %v\n%s", err, debug.Stack())
		return errors.New("fatal error occured")
	})

	api.POST("/graphql", jwtMiddlewareFn, func(ctx *gin.Context) {
		ctxWithUserStatus := context.WithValue(ctx.Request.Context(), graph_resolver.ContextKey("user_status"), map[string]any{
			"is_admin": ctx.GetBool("is_admin"),
			"user_id":  ctx.GetString("user_id"),
		})
		ctx.Request = ctx.Request.WithContext(ctxWithUserStatus)
		graphql.ServeHTTP(ctx.Writer, ctx.Request)
	})
	if s.debug {
		api.GET("/graphql", func(ctx *gin.Context) {
			writer, request := ctx.Writer, ctx.Request
			playground.Handler("GraphQL API Debug", request.URL.Path).ServeHTTP(writer, request)
		})
	}

	webFs, webPath := web.NewEmbedFs()
	s.engine.Use(static.Serve("/", static.EmbedFolder(webFs, webPath)))

	s.server.Addr = listen
	s.server.Handler = s.engine.Handler()

	return nil
}
