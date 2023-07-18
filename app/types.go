package app

import "com.geophone.observer/feature"

type ServerOptions struct {
	Gzip           int
	CORS           bool
	Version        string
	WebPrefix      string
	APIPrefix      string
	FeatureOptions *feature.FeatureOptions
}
