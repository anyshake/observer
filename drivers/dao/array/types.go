package array

import "database/sql/driver"

type Int32Array []int32

func (i *Int32Array) Scan(val any) error {
	intArr, err := i.Decode(val)
	if err != nil {
		return err
	}

	*i = intArr
	return nil
}

func (i Int32Array) Value() (driver.Value, error) {
	return i.Encode(), nil
}
