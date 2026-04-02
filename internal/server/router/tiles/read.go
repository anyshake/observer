package tiles

import (
	"fmt"
	"sort"
)

func (h *mapTilesHandler) readMapTile(z, x, y int) ([]byte, error) {
	tf, err := h.loadMapTile(z, x)
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
