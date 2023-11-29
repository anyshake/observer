package publisher

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"strings"
)

func (i *int32Array) Scan(val any) error {
	var strArr []string
	switch v := val.(type) {
	case string:
		strArr = strings.Split(v, "|")
	case []byte:
		strArr = strings.Split(string(v), "|")
	default:
		err := fmt.Errorf("unsupported type: %T", v)
		return err
	}

	for _, v := range strArr {
		intData, err := strconv.Atoi(v)
		if err != nil {
			return err
		}
		*i = append(*i, int32(intData))
	}

	return nil
}

func (i int32Array) Value() (driver.Value, error) {
	var rawText string
	for ii, vv := range i {
		rawText += strconv.Itoa(int(vv))
		if ii != len(i)-1 {
			rawText += "|"
		}
	}

	return rawText, nil
}
