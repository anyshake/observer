package tiles

import (
	"fmt"
	"strconv"
	"strings"
)

func parseGlyphRange(glyphRange string) (int, int, error) {
	const glyphRangeSize = 256

	startStr, endStr, ok := strings.Cut(glyphRange, "-")
	if !ok {
		return 0, 0, fmt.Errorf("malformed range")
	}
	start, err := strconv.Atoi(startStr)
	if err != nil || start < 0 || start%glyphRangeSize != 0 {
		return 0, 0, fmt.Errorf("invalid range start")
	}
	end, err := strconv.Atoi(endStr)
	if err != nil || end != start+glyphRangeSize-1 {
		return 0, 0, fmt.Errorf("invalid range end")
	}
	return start, end, nil
}

func (h *mapTilesHandler) readGlyphRange(fontStack string, start, end int) ([]byte, error) {
	return h.tilesFs.ReadFile(fmt.Sprintf("%s/glyphs/%s/%d-%d.pbf", h.baseDir, fontStack, start, end))
}
