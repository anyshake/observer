package tiles

import (
	"embed"

	lru "github.com/hashicorp/golang-lru/v2"
)

type mapTileDataObject struct {
	y      uint32
	offset uint32
	size   uint32
}

type mapTileData struct {
	objects []mapTileDataObject
	data    []byte
	base    int64
}

type mapTilesHandler struct {
	tilesFs embed.FS
	baseDir string
	version string // for http cache invalidation
	cache   *lru.Cache[string, *mapTileData]
}
