package server

import (
	"context"
	"net"

	"github.com/anyshake/observer/internal/notification"
	graph_resolver "github.com/anyshake/observer/internal/server/router/graph"
	"github.com/anyshake/observer/pkg/logger"
	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.ReleaseMode)
}

func New(parentCtx context.Context, debug, cors bool, resolver *graph_resolver.Resolver, notificationHub *notification.Hub, logger *logger.Adapter) *HttpServer {
	shutdownCtx, shutdownCancel := context.WithCancel(parentCtx)
	obj := &HttpServer{
		debug:           debug,
		cors:            cors,
		log:             logger,
		resolver:        resolver,
		notificationHub: notificationHub,
		shutdownCtx:     shutdownCtx,
		shutdownCancel:  shutdownCancel,
	}
	obj.server.BaseContext = func(net.Listener) context.Context { return shutdownCtx }
	return obj
}
