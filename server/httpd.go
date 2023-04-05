package server

import (
	"fmt"
	"log"

	"com.geophone.observer/server/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.ReleaseMode)
}

func New() *gin.Engine {
	r := gin.New()
	r.Use(gin.LoggerWithFormatter(
		func(param gin.LogFormatterParams) string {
			return fmt.Sprintf("%s [%s] %s %d %s %s\n",
				param.TimeStamp.Format("2006/01/02 15:04:05"),
				param.ClientIP, param.Method, param.StatusCode,
				param.Path, param.ErrorMessage,
			)
		},
	))

	return r
}

func StartServer(r *gin.Engine, options ServerOptions) error {
	if options.Cors {
		r.Use(cors.AllowCros([]cors.HttpHeader{
			{
				Header: "Access-Control-Allow-Origin",
				Value:  "*",
			}, {
				Header: "Access-Control-Allow-Credentials",
				Value:  "true",
			}, {
				Header: "Access-Control-Allow-Methods",
				Value:  "POST, OPTIONS, GET",
			}, {
				Header: "Access-Control-Allow-Headers",
				Value:  "Authorization, Content-Type, Version",
			},
		}))
	}
	r.Use(gzip.Gzip(options.Gzip))

	err := r.Run(fmt.Sprintf("%s:%s", options.Listen, options.Port))
	if err != nil {
		log.Println("HTTP 服务启动失败")
	}

	return err
}
