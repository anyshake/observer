package web

import (
	"embed"
)

//go:embed dist
var dist embed.FS

func NewEmbedFs() (embed.FS, string) {
	return dist, "dist"
}
