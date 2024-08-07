package array

import (
	"strconv"
)

func (data Int32Array) Encode() string {
	var rawText string
	for ii, vv := range data {
		rawText += strconv.Itoa(int(vv))
		if ii != len(data)-1 {
			rawText += "|"
		}
	}

	return rawText
}
