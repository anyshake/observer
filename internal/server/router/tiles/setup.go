package tiles

import (
	"net/http"

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
		switch ctx.Query("type") {
		case "maptile":
			h.serveMapTile(ctx)
		case "glyph":
			h.serveGlyphRange(ctx)
		default:
			response.Error(ctx, http.StatusBadRequest, "invalid resource type")
		}
	})
}
