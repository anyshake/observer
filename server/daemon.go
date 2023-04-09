package server

import (
	"fmt"
	"io/fs"
	"log"
	"net/http"

	"com.geophone.observer/app"
	"com.geophone.observer/frontend"
	"com.geophone.observer/server/middleware/cors"
	"com.geophone.observer/server/middleware/static"
	"com.geophone.observer/server/response"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.ReleaseMode)
}

func ServerDaemon(host string, port int, options *app.ServerOptions) error {
	r := gin.New()
	r.Use(gzip.Gzip(options.Gzip),
		gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
			return fmt.Sprintf("%s [%s] %s %d %s %s\n",
				param.TimeStamp.Format("2006/01/02 15:04:05"),
				param.ClientIP, param.Method, param.StatusCode,
				param.Path, param.ErrorMessage,
			)
		}))
	if options.Cors {
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
		response.ErrorHandler(c, http.StatusNotFound)
	})

	RegisterRouter(r.Group(
		fmt.Sprintf("/%s/%s", options.ApiPrefix, options.Version),
	), options)

	r.Use(static.ServeEmbed(&static.LocalFileSystem{
		Root: options.WebPrefix, Prefix: options.WebPrefix,
		FileSystem: http.FS(func(path string, f fs.FS) fs.FS {
			p, _ := fs.Sub(f, path)
			return p
		}("dist", frontend.Dist)),
	}))

	err := r.Run(fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		log.Println("failed to start HTTP server", err)
	}

	return err
}
