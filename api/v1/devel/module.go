package devel

import (
	"net/http"

	v1 "github.com/anyshake/observer/api/v1"
	_ "github.com/anyshake/observer/docs"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	gs "github.com/swaggo/gin-swagger"
)

func (d *Devel) Register(rg *gin.RouterGroup, resolver *v1.Resolver) error {
	rg.GET("/devel/*any", func(ctx *gin.Context) {
		if ctx.Param("any") == "/" {
			url := ctx.Request.URL
			ctx.Redirect(http.StatusMovedPermanently, url.Path+"/index.html")
		}
	}, gs.WrapHandler(swaggerFiles.Handler))

	return nil
}
