package v1

import (
	"github.com/anyshake/observer/services"
	"github.com/gin-gonic/gin"
)

type Resolver struct {
	*services.Options
}

type Endpoint interface {
	Register(*gin.RouterGroup, *Resolver) error
	GetApiName() string
}
