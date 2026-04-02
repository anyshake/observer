package web

import (
	"embed"
	"encoding/binary"
	"fmt"
	"net/http"
	"sort"
	"strconv"

	"github.com/anyshake/observer/internal/server/response"
	"github.com/gin-gonic/gin"
	lru "github.com/hashicorp/golang-lru/v2"
)

//go:embed tiles
var tiles embed.FS

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
	cache *lru.Cache[string, *mapTileData]
}

func NewMapTilesHandler(cacheSize int) (gin.HandlerFunc, error) {
	lruCache, err := lru.New[string, *mapTileData](cacheSize)
	if err != nil {
		return nil, err
	}

	return (&mapTilesHandler{cache: lruCache}).handle, nil
}

func (h *mapTilesHandler) handle(ctx *gin.Context) {
	z, err := strconv.Atoi(ctx.Query("z"))
	if err != nil || z < 0 {
		response.Error(ctx, http.StatusBadRequest, "invalid Z index")
		return
	}
	x, err := strconv.Atoi(ctx.Query("x"))
	if err != nil || x < 0 {
		response.Error(ctx, http.StatusBadRequest, "invalid X index")
		return
	}
	y, err := strconv.Atoi(ctx.Query("y"))
	if err != nil || y < 0 {
		response.Error(ctx, http.StatusBadRequest, "invalid Y index")
		return
	}

	data, err := h.readTile(z, x, y)
	if err != nil {
		response.Error(ctx, http.StatusNotFound, "tile not found")
		return
	}

	if len(data) > 2 && data[0] == 0x1f && data[1] == 0x8b {
		ctx.Header("Content-Encoding", "gzip")
	}
	response.Blob(ctx, fmt.Sprintf("maptile-z%d-x%d-y%d", z, x, y), "application/x-protobuf", data)
}

func (h *mapTilesHandler) readTile(z, x, y int) ([]byte, error) {
	tf, err := h.loadTile(z, x)
	if err != nil {
		return nil, err
	}

	i := sort.Search(len(tf.objects), func(i int) bool {
		return tf.objects[i].y >= uint32(y)
	})

	if i >= len(tf.objects) || tf.objects[i].y != uint32(y) {
		return nil, fmt.Errorf("tile not found")
	}

	e := tf.objects[i]
	offset := tf.base + int64(e.offset)
	end := offset + int64(e.size)

	if end > int64(len(tf.data)) {
		return nil, fmt.Errorf("data out of range")
	}

	return tf.data[offset:end], nil
}

func (h *mapTilesHandler) loadTile(z, x int) (*mapTileData, error) {
	cacheKey := fmt.Sprintf("z%d/x%d", z, x)
	if v, ok := h.cache.Get(cacheKey); ok {
		return v, nil
	}

	data, err := tiles.ReadFile(fmt.Sprintf("tiles/z%d/x%d.bin", z, x))
	if err != nil {
		return nil, err
	}

	if len(data) < 4 {
		return nil, fmt.Errorf("corrupt file")
	}
	count := binary.LittleEndian.Uint32(data[:4])
	if len(data) < int(4+count*12) {
		return nil, fmt.Errorf("corrupt header")
	}
	header := data[4 : 4+count*12]

	objects := make([]mapTileDataObject, count)
	for i := uint32(0); i < count; i++ {
		base := i * 12
		objects[i] = mapTileDataObject{
			y:      binary.LittleEndian.Uint32(header[base:]),
			offset: binary.LittleEndian.Uint32(header[base+4:]),
			size:   binary.LittleEndian.Uint32(header[base+8:]),
		}
	}

	for i := 1; i < len(objects); i++ {
		if objects[i].y <= objects[i-1].y {
			return nil, fmt.Errorf("bad index order: %d/%d", z, x)
		}
	}

	tf := &mapTileData{
		data:    data,
		objects: objects,
		base:    int64(4 + len(objects)*12),
	}

	h.cache.Add(cacheKey, tf)
	return tf, nil
}
