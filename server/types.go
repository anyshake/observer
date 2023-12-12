package server

import (
	"github.com/anyshake/observer/app"
	"github.com/gin-gonic/gin"
)

const CAPACITY = 30 // Capacity for rate limiter to prevent from being attacked

type ApiServices interface {
	RegisterModule(rg *gin.RouterGroup, options *app.ServerOptions)
}
