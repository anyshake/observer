package server

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/anyshake/observer/app"
	"github.com/anyshake/observer/frontend"
	"github.com/anyshake/observer/server/middleware/cors"
	"github.com/anyshake/observer/server/middleware/static"
	"github.com/anyshake/observer/server/response"
	"github.com/anyshake/observer/utils/logger"
	"github.com/fatih/color"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.ReleaseMode)
}

func StartDaemon(host string, port int, options *app.ServerOptions) error {
	r := gin.New()
	r.Use(
		gzip.Gzip(options.Gzip, gzip.WithExcludedPaths([]string{options.APIPrefix})),
		gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
			w := color.New(color.FgCyan).SprintFunc()
			trimmedErr := strings.TrimRight(param.ErrorMessage, "\n")
			loggerText := w(fmt.Sprintf("%s [server] %s %d %s %s %s\n",
				param.TimeStamp.Format("2006/01/02 15:04:05"),
				param.Method, param.StatusCode,
				param.ClientIP, param.Path, trimmedErr,
			))

			return loggerText
		}),
	)
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

	r.NoRoute(func(c *gin.Context) {
		response.Error(c, http.StatusNotFound)
	})

	// Register API v1 routers
	RegisterRouterV1(r.Group(
		fmt.Sprintf("/%s/v1", options.APIPrefix),
	), options)

	r.Use(static.ServeEmbed(&static.LocalFileSystem{
		Root: options.WebPrefix, Prefix: options.WebPrefix,
		FileSystem: static.CreateFilesystem(frontend.Dist, "dist"),
	}))

	err := r.Run(fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		logger.Fatal("server", err, color.FgRed)
	}

	return err
}
