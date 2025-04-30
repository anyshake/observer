package uniarch

import (
	"encoding/json"
	"reflect"
)

func GetArchMap() []ArchMap {
	dataBytes, err := archMapAsset.ReadFile("arch_map.json")
	if err != nil {
		return nil
	}

	var rawMap map[string]map[string]any
	if err = json.Unmarshal(dataBytes, &rawMap); err != nil {
		return nil
	}

	var result []ArchMap
	for _, v := range rawMap {
		arch := ArchMap{
			Flags: make(map[string]string),
		}

		archVal := reflect.ValueOf(&arch).Elem()
		archType := archVal.Type()
		for i := 0; i < archType.NumField(); i++ {
			field := archType.Field(i)
			tag := field.Tag.Get("json")
			if tag == "flags" {
				continue
			}
			if val, ok := v[tag]; ok {
				if strVal, ok := val.(string); ok {
					archVal.FieldByName(field.Name).SetString(strVal)
				}
			}
		}

		for key, value := range v {
			if _, ok := arch.Flags[key]; !ok && !archVal.FieldByNameFunc(func(n string) bool {
				field, _ := archType.FieldByName(n)
				return field.Tag.Get("json") == key
			}).IsValid() {
				if strVal, ok := value.(string); ok {
					arch.Flags[key] = strVal
				}
			}
		}

		result = append(result, arch)
	}

	return result
}
