package tiles

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/anyshake/observer/internal/server/response"
	"github.com/anyshake/observer/web"
	"github.com/gin-gonic/gin"
	lru "github.com/hashicorp/golang-lru/v2"
)

func Setup(routerGroup *gin.RouterGroup, version string, jwtMiddleware gin.HandlerFunc) {
	lruCache, _ := lru.New[string, *mapTileData](256)
	fs, baseDir := web.NewMapTiles()
	h := &mapTilesHandler{
		tilesFs: fs,
		baseDir: baseDir,
		version: version,
		cache:   lruCache,
	}

	routerGroup.GET("/tiles", jwtMiddleware, func(ctx *gin.Context) {
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

		data, err := h.readMapTile(z, x, y)
		if err != nil {
			response.Error(ctx, http.StatusNotFound, "tile not found")
			return
		}

		if len(data) > 2 && data[0] == 0x1f && data[1] == 0x8b {
			ctx.Header("Content-Encoding", "gzip")
		}
		ctx.Header("Cache-Control", "public, max-age=31536000, immutable")
		ctx.Header("ETag", h.version)
		response.Blob(ctx, fmt.Sprintf("maptile-z%d-x%d-y%d", z, x, y), "application/vnd.mapbox-vector-tile", data)
	})
}
