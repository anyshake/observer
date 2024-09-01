package trace

import (
	"strconv"
)

func string2Float(num string) float64 {
	r, err := strconv.ParseFloat(num, 64)
	if err != nil {
		return 0.0
	}

	return r
}

func isMapKeysEmpty(m map[string]any, k []string) bool {
	for _, v := range k {
		switch m[v].(type) {
		case string:
			if len(m[v].(string)) == 0 {
				return false
			}
		default:
			continue
		}
	}

	return true
}

func isMapHasKeys(m map[string]any, k []string) bool {
	for _, v := range k {
		if _, ok := m[v]; !ok {
			return false
		}
	}

	return true
}
