package devel

import (
	"net/http"

	"github.com/anyshake/observer/app"
	_ "github.com/anyshake/observer/docs"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	gs "github.com/swaggo/gin-swagger"
)

func (d *Devel) RegisterModule(rg *gin.RouterGroup, options *app.ServerOptions) {
	if !options.FeatureOptions.Config.Server.Debug {
		return
	}

	rg.GET("/devel/*any", func(ctx *gin.Context) {
		if ctx.Param("any") == "/" {
			url := ctx.Request.URL
			ctx.Redirect(http.StatusMovedPermanently, url.Path+"/index.html")
		}
	}, gs.WrapHandler(swaggerFiles.Handler))
}
