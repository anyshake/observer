package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/anyshake/observer/frontend"
	"github.com/anyshake/observer/server/endpoints"
	"github.com/anyshake/observer/server/middlewares/cors"
	"github.com/anyshake/observer/server/middlewares/jwt"
	loggerware "github.com/anyshake/observer/server/middlewares/logger"
	"github.com/anyshake/observer/server/middlewares/static"
	"github.com/anyshake/observer/server/response"
	"github.com/anyshake/observer/services"
	loggerutil "github.com/anyshake/observer/utils/logger"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type ServerImpl struct {
	GzipLevel int
	ApiPrefix string
	WebPrefix string
}

func init() {
	gin.SetMode(gin.ReleaseMode)
}

func (s *ServerImpl) Start(logger *logrus.Entry, options *services.Options) {
	if options.Config.Server.Debug {
		logger.Warnln("PLEASE NOTE THAT DEBUG MODE IS ENABLED")
	}
	httpSrv := gin.New()

	// Setup Gzip & logger
	httpSrv.Use(
		gzip.Gzip(s.GzipLevel, gzip.WithExcludedPaths([]string{s.ApiPrefix})),
		loggerware.New(logger),
	)

	// Setup Cross-Origin Resource Sharing (CORS)
	if options.Config.Server.CORS {
		httpSrv.Use(cors.New([]cors.HttpHeader{
			{Header: "Access-Control-Allow-Origin", Value: "*"},
			{Header: "Access-Control-Allow-Methods", Value: "POST, OPTIONS, GET"},
			{Header: "Access-Control-Allow-Headers", Value: "Content-Type, Authorization"},
			{Header: "Access-Control-Expose-Headers", Value: "Content-Length"},
		}))
	}

	// Setup 404 error handler
	httpSrv.NoRoute(func(ctx *gin.Context) {
		response.Message(ctx, options.TimeSource, "requested resource is not found", http.StatusNotFound, nil)
	})

	// Setup API routers
	jwtHandler, err := jwt.New(options.TimeSource, loggerutil.GetLogger("server_jwt"), options.Database, 24*time.Hour)
	if err != nil {
		logger.Fatalln(err)
	}
	err = endpoints.SetupApiV1(httpSrv.Group(fmt.Sprintf("/%s/v1", s.ApiPrefix)), jwtHandler, options)
	if err != nil {
		logger.Fatalln(err)
	}
	endpoints.SetupApiV2(httpSrv.Group(fmt.Sprintf("/%s/v2", s.ApiPrefix)), options)

	// Setup static file serve
	httpSrv.Use(static.NewEmbed(options.TimeSource, &static.LocalFileSystem{
		Root: s.WebPrefix, Prefix: s.WebPrefix,
		FileSystem: static.NewFilesystem(frontend.Dist, "dist"),
	}))

	// Start web server
	var (
		host = options.Config.Server.Host
		port = options.Config.Server.Port
	)
	logger.Infof("web server is listening on %s:%d", host, port)
	httpSrv.Run(fmt.Sprintf("%s:%d", host, port))
}
