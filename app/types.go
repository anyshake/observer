package app

import "github.com/anyshake/observer/feature"

type ServerOptions struct {
	Gzip           int
	CORS           bool
	Version        string
	WebPrefix      string
	APIPrefix      string
	FeatureOptions *feature.FeatureOptions
}
