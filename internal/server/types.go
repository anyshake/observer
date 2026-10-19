package server

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/anyshake/observer/internal/notification"
	graph_resolver "github.com/anyshake/observer/internal/server/router/graph"
	"github.com/anyshake/observer/pkg/logger"
	"github.com/gin-gonic/gin"
)

const AUTH_TIMEOUT = 24 * time.Hour

type HttpServer struct {
	debug bool
	cors  bool

	resolver        *graph_resolver.Resolver
	notificationHub *notification.Hub
	log             *logger.Adapter
	engine          *gin.Engine
	server          http.Server
	listener        net.Listener
	closeSocket     func(context.Context) error
	shutdownCtx     context.Context
	shutdownCancel  context.CancelFunc
}
