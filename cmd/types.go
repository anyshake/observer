package main

const (
	GZIP_LEVEL = 9
	WEB_PREFIX = "/"
	API_PREFIX = "/api"
)

type arguments struct {
	Path    string // Path to config file
	Version bool   // Show version information
}
