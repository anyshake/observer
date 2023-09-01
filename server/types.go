package server

import (
	"github.com/bclswl0827/observer/app"
	"github.com/gin-gonic/gin"
)

type ApiServices interface {
	RegisterModule(rg *gin.RouterGroup, options *app.ServerOptions)
}
