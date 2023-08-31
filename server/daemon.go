package server

import (
	"fmt"
	"net/http"

	"com.geophone.observer/app"
	"com.geophone.observer/frontend"
	"com.geophone.observer/server/middleware/cors"
	"com.geophone.observer/server/middleware/static"
	"com.geophone.observer/server/response"
	"com.geophone.observer/utils/logger"
	"github.com/fatih/color"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.ReleaseMode)
}

func ServerDaemon(host string, port int, options *app.ServerOptions) error {
	r := gin.New()
	r.Use(
		gzip.Gzip(options.Gzip),
		gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
			w := color.New(color.FgCyan).SprintFunc()
			text := w(fmt.Sprintf("%s [server] %s %d %s %s %s\n",
				param.TimeStamp.Format("2006/01/02 15:04:05"),
				param.Method, param.StatusCode, param.ClientIP,
				param.Path, param.ErrorMessage,
			))

			return text
		}),
	)
	if options.CORS {
		r.Use(cors.AllowCros([]cors.HttpHeader{
			{
				Header: "Access-Control-Allow-Origin",
				Value:  "*",
			}, {
				Header: "Access-Control-Allow-Methods",
				Value:  "POST, OPTIONS, GET",
			}, {
				Header: "Access-Control-Allow-Headers",
				Value:  "Authorization, Content-Type, Version",
			},
		}))
	}

	r.NoRoute(func(c *gin.Context) {
		response.Error(c, http.StatusNotFound)
	})

	RegisterRouter(r.Group(
		fmt.Sprintf("/%s/%s", options.APIPrefix, options.Version),
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
