package v1

import (
	"github.com/anyshake/observer/services"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

type Endpoint interface {
	Bind(*gin.RouterGroup, *jwt.GinJWTMiddleware, *services.Options) error
	GetApiName() string
}
