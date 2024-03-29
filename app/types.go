package app

import "github.com/anyshake/observer/feature"

type ServerOptions struct {
	Gzip           int
	CORS           bool
	WebPrefix      string
	APIPrefix      string
	RateFactor     int
	FeatureOptions *feature.FeatureOptions
}
