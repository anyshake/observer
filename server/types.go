package server

import (
	"github.com/anyshake/observer/services"
	"github.com/sirupsen/logrus"
)

type Server interface {
	Start(*logrus.Entry, *services.Options)
}
