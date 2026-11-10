package tiles

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"github.com/anyshake/observer/internal/server/response"
	"github.com/gin-gonic/gin"
)

func (h *mapTilesHandler) setCacheHeaders(ctx *gin.Context) {
	ctx.Header("Cache-Control", "public, max-age=31536000, immutable")
	ctx.Header("ETag", h.version)
}

func (h *mapTilesHandler) serveMapTile(ctx *gin.Context) {
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
	h.setCacheHeaders(ctx)
	response.Blob(ctx, fmt.Sprintf("maptile-z%d-x%d-y%d", z, x, y), "application/vnd.mapbox-vector-tile", data)
}

var glyphFontStackPattern = regexp.MustCompile(`^[A-Za-z0-9_]+$`)

func (h *mapTilesHandler) serveGlyphRange(ctx *gin.Context) {
	fontStack := ctx.Query("fontstack")
	if !glyphFontStackPattern.MatchString(fontStack) {
		response.Error(ctx, http.StatusBadRequest, "invalid font stack")
		return
	}
	start, end, err := parseGlyphRange(ctx.Query("range"))
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid glyph range")
		return
	}

	data, err := h.readGlyphRange(fontStack, start, end)
	if err != nil {
		response.Error(ctx, http.StatusNotFound, "glyph range not found")
		return
	}

	if len(data) > 2 && data[0] == 0x1f && data[1] == 0x8b {
		ctx.Header("Content-Encoding", "gzip")
	}
	h.setCacheHeaders(ctx)
	response.Blob(ctx, fmt.Sprintf("glyph-%s-%d-%d", fontStack, start, end), "application/x-protobuf", data)
}
