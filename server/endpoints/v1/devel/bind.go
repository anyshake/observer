package devel

import (
	"net/http"

	_ "github.com/anyshake/observer/docs"
	"github.com/anyshake/observer/services"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	gs "github.com/swaggo/gin-swagger"
)

func (d *Devel) Bind(rg *gin.RouterGroup, jwtHandler *jwt.GinJWTMiddleware, options *services.Options) error {
	rg.GET("/devel/*any", func(ctx *gin.Context) {
		if ctx.Param("any") == "/" {
			url := ctx.Request.URL
			ctx.Redirect(http.StatusMovedPermanently, url.Path+"/index.html")
		}
	}, gs.WrapHandler(swaggerFiles.Handler))

	return nil
}
