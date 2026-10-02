package server

import (
	"context"

	"github.com/anyshake/observer/internal/notification"
	graph_resolver "github.com/anyshake/observer/internal/server/router/graph"
	"github.com/anyshake/observer/pkg/logger"
	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.ReleaseMode)
}

func New(debug, cors bool, resolver *graph_resolver.Resolver, notificationHub *notification.Hub, logger *logger.Adapter) *HttpServer {
	shutdownCtx, shutdownCancel := context.WithCancel(context.Background())
	return &HttpServer{
		debug:           debug,
		cors:            cors,
		log:             logger,
		resolver:        resolver,
		notificationHub: notificationHub,
		shutdownCtx:     shutdownCtx,
		shutdownCancel:  shutdownCancel,
	}
}
