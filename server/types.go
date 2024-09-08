package server

import (
	"github.com/anyshake/observer/services"
)

type Options struct {
	CORS            bool
	DebugMode       bool
	GzipLevel       int
	RateFactor      int
	WebPrefix       string
	ApiPrefix       string
	ServicesOptions *services.Options
}
