package tiles

import (
	"encoding/binary"
	"fmt"
)

func (h *mapTilesHandler) loadMapTile(z, x int) (*mapTileData, error) {
	cacheKey := fmt.Sprintf("z%d/x%d", z, x)
	if v, ok := h.cache.Get(cacheKey); ok {
		return v, nil
	}

	data, err := h.tilesFs.ReadFile(fmt.Sprintf("%s/z%d/x%d.bin", h.baseDir, z, x))
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
