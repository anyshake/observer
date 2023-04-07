package server

import (
	"fmt"
	"io/fs"
	"log"
	"net/http"

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

func ServerDaemon(options *ServerOptions) error {
	r := gin.New()
	r.NoRoute(func(c *gin.Context) {
		response.ErrorHandler(c, http.StatusNotFound)
	})

	RegisterRouter(r.Group(
		fmt.Sprintf("/api/%s", options.Version),
	), options)

	r.Use(gzip.Gzip(options.Gzip),
		gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
			return fmt.Sprintf("%s [%s] %s %d %s %s\n",
				param.TimeStamp.Format("2006/01/02 15:04:05"),
				param.ClientIP, param.Method, param.StatusCode,
				param.Path, param.ErrorMessage,
			)
		}),
		static.ServeEmbed(&static.LocalFileSystem{
			Root: "/", Prefix: "/",
			FileSystem: http.FS(func(path string, f fs.FS) fs.FS {
				p, _ := fs.Sub(f, path)
				return p
			}("dist", frontend.Dist)),
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

	err := r.Run(fmt.Sprintf("%s:%d", options.Host, options.Port))
	if err != nil {
		log.Println("failed to start HTTP server", err)
	}

	return err
}
