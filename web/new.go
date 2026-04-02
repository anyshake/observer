package web

import (
	"embed"
)

//go:embed dist
var dist embed.FS

func NewWebDist() (embed.FS, string) {
	return dist, "dist"
}

//go:embed tiles
var tiles embed.FS

func NewMapTiles() (embed.FS, string) {
	return tiles, "tiles"
}
