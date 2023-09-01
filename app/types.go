package app

import "github.com/bclswl0827/observer/feature"

type ServerOptions struct {
	Gzip           int
	CORS           bool
	Version        string
	WebPrefix      string
	APIPrefix      string
	FeatureOptions *feature.FeatureOptions
}
