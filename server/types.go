package server

import (
	"com.geophone.observer/features/collector"
	"github.com/gin-gonic/gin"
)

type ServerOptions struct {
	Version string
	Host    string
	Port    int
	Gzip    int
	Cors    bool
	Message *collector.Message
	Status  *collector.Status
}

type ApiServices interface {
	RegisterModule(rg *gin.RouterGroup, m *collector.Message, s *collector.Status)
}
