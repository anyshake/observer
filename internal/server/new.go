package server

import (
	graph_resolver "github.com/anyshake/observer/internal/server/router/graph"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func init() {
	gin.SetMode(gin.ReleaseMode)
}

func New(debug, cors bool, resolver *graph_resolver.Resolver, logger *logrus.Entry) *HttpServer {
	return &HttpServer{
		debug:    debug,
		cors:     cors,
		log:      logger,
		resolver: resolver,
	}
}
