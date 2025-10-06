package config

import (
	"encoding/json"
	"errors"
	"strconv"
)

func GetConfigValInt64(val any) (int64, error) {
	switch v := val.(type) {
	case json.Number:
		return v.Int64()
	case string:
		return strconv.ParseInt(v, 10, 64)
	case int64:
		return v, nil
	case float64:
		return int64(v), nil
	}
	return 0, errors.New("integer expected")
}

func GetConfigValInt64Array(val any) ([]int64, error) {
	if val == nil {
		return []int64{}, nil
	}
	arr, ok := val.([]any)
	if !ok {
		return nil, errors.New("integer array expected")
	}
	var res []int64
	for _, v := range arr {
		i, err := GetConfigValInt64(v)
		if err != nil {
			return nil, err
		}
		res = append(res, i)
	}
	return res, nil
}

func GetConfigValFloat64(val any) (float64, error) {
	switch v := val.(type) {
	case json.Number:
		return v.Float64()
	case string:
		return strconv.ParseFloat(v, 64)
	case float64:
		return v, nil
	}
	return 0, errors.New("float expected")
}

func GetConfigValFloat64Array(val any) ([]float64, error) {
	if val == nil {
		return []float64{}, nil
	}
	arr, ok := val.([]any)
	if !ok {
		return nil, errors.New("float array expected")
	}
	var res []float64
	for _, v := range arr {
		f, err := GetConfigValFloat64(v)
		if err != nil {
			return nil, err
		}
		res = append(res, f)
	}
	return res, nil
}

func GetConfigValString(val any) (string, error) {
	strVal, ok := val.(string)
	if !ok {
		return "", errors.New("string expected")
	}
	return strVal, nil
}

func GetConfigValStringArray(val any) ([]string, error) {
	if val == nil {
		return []string{}, nil
	}
	arr, ok := val.([]any)
	if !ok {
		return nil, errors.New("string array expected")
	}
	var res []string
	for _, v := range arr {
		s, err := GetConfigValString(v)
		if err != nil {
			return nil, err
		}
		res = append(res, s)
	}
	return res, nil
}

func GetConfigValBool(val any) (bool, error) {
	boolVal, ok := val.(bool)
	if !ok {
		return false, errors.New("boolean expected")
	}
	return boolVal, nil
}
