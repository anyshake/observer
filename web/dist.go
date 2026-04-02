package web

import (
	"embed"
)

//go:embed dist
var dist embed.FS

func NewWebDist() (embed.FS, string) {
	return dist, "dist"
}
