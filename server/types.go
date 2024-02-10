package server

import (
	"github.com/anyshake/observer/app"
	"github.com/gin-gonic/gin"
)

type ApiServices interface {
	RegisterModule(rg *gin.RouterGroup, options *app.ServerOptions)
}
