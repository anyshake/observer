package publisher

import (
	"database/sql/driver"
)

func (i *Int32Array) Scan(val any) error {
	intArr, err := DecodeInt32Array(val)
	if err != nil {
		return err
	}

	*i = intArr
	return nil
}

func (i Int32Array) Value() (driver.Value, error) {
	return EncodeInt32Array(i), nil
}
