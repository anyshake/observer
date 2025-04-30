package server

import (
	"net/http"
	"time"

	graph_resolver "github.com/anyshake/observer/internal/server/router/graph"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const AUTH_TIMEOUT = 24 * time.Hour

type httpServer struct {
	debug bool
	cors  bool

	resolver *graph_resolver.Resolver
	log      *logrus.Entry
	engine   *gin.Engine
	server   http.Server
}
