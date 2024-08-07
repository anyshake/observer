package array

import (
	"fmt"
	"strconv"
	"strings"
)

func (data Int32Array) Decode(val any) ([]int32, error) {
	var strArr []string
	switch v := val.(type) {
	case string:
		strArr = strings.Split(v, "|")
	case []byte:
		strArr = strings.Split(string(v), "|")
	default:
		return nil, fmt.Errorf("unsupported type: %T", v)
	}

	var intArr []int32
	for _, v := range strArr {
		intData, err := strconv.Atoi(v)
		if err != nil {
			return nil, err
		}
		intArr = append(intArr, int32(intData))
	}

	return intArr, nil
}
