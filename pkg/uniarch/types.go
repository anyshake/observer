package uniarch

import (
	"embed"
)

//go:embed arch_map.json
var archMapAsset embed.FS

type ArchMap struct {
	Name   string            `json:"name"`
	GOOS   string            `json:"goos"`
	GOARM  string            `json:"goarm"`
	GOMIPS string            `json:"gomips"`
	GOARCH string            `json:"goarch"`
	Flags  map[string]string `json:"flags"`
}
