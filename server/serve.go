package server

import (
	"fmt"
	"net/http"

	"github.com/anyshake/observer/frontend"
	"github.com/anyshake/observer/server/middleware/cors"
	loggerware "github.com/anyshake/observer/server/middleware/logger"
	"github.com/anyshake/observer/server/middleware/static"
	"github.com/anyshake/observer/server/response"
	"github.com/anyshake/observer/utils/logger"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.ReleaseMode)
}

func Serve(host string, port int, options *Options) error {
	if options.DebugMode {
		logger.GetLogger(Serve).Warnln("PLEASE NOTE THAT DEBUG MODE IS ENABLED")
	}
	r := gin.New()

	// Setup Gzip & logger
	r.Use(
		gzip.Gzip(options.GzipLevel, gzip.WithExcludedPaths([]string{options.ApiPrefix})),
		loggerware.WriteLog(logger.GetLogger(Serve)),
	)

	// Setup Cross-Origin Resource Sharing (CORS)
	if options.CORS {
		r.Use(cors.AllowCORS([]cors.HttpHeader{
			{
				Header: "Access-Control-Allow-Origin",
				Value:  "*",
			}, {
				Header: "Access-Control-Allow-Methods",
				Value:  "POST, OPTIONS, GET",
			}, {
				Header: "Access-Control-Allow-Headers",
				Value:  "Content-Type",
			}, {
				Header: "Access-Control-Expose-Headers",
				Value:  "Content-Length",
			},
		}))
	}

	// Setup 404 error handler
	r.NoRoute(func(c *gin.Context) {
		response.Error(c, http.StatusNotFound)
	})

	// Register API routers
	err := registerEndpointsV1(r.Group(fmt.Sprintf("/%s/v1", options.ApiPrefix)), options)
	if err != nil {
		return err
	}
	registerEndpointsV2(r.Group(fmt.Sprintf("/%s/v2", options.ApiPrefix)), options)

	// Setup static file serve
	r.Use(static.ServeEmbed(&static.LocalFileSystem{
		Root: options.WebPrefix, Prefix: options.WebPrefix,
		FileSystem: static.CreateFilesystem(frontend.Dist, "dist"),
	}))

	// Start server
	err = r.Run(fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		logger.GetLogger(Serve).Fatalf("server: %v\n", err)
	}

	return err
}
