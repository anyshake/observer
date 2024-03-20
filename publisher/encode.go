package publisher

import (
	"strconv"
)

func EncodeInt32Array(data []int32) string {
	var rawText string
	for ii, vv := range data {
		rawText += strconv.Itoa(int(vv))
		if ii != len(data)-1 {
			rawText += "|"
		}
	}

	return rawText
}
