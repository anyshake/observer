package action

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"strings"
	"unicode"

	"github.com/anyshake/observer/internal/dao/model"
)

type SettingType string

const (
	String      SettingType = "string"
	StringArray SettingType = "string[]"
	Bool        SettingType = "bool"
	Int         SettingType = "int"
	IntArray    SettingType = "int[]"
	Float       SettingType = "float"
	FloatArray  SettingType = "float[]"
)

func (h *Handler) SettingsGet(namespace, key string) (any, SettingType, int, error) {
	if h.daoObj == nil {
		return nil, "", 0, errors.New("database is not opened")
	}

	var settings model.UserSettings
	err := h.daoObj.Database.
		Model(settings).
		Where("namespace = ? AND config_key = ?", namespace, key).
		First(&settings).
		Error
	if err != nil {
		return nil, "", 0, fmt.Errorf("failed to get settings for namespace %s, key %s: %w", namespace, key, err)
	}

	switch SettingType(settings.ConfigType) {
	case String:
		return h.removeInvisible(string(settings.ConfigValue)), String, settings.Version, nil
	case Bool:
		return settings.ConfigValue[0] == 1, Bool, settings.Version, nil
	case Int:
		var result int64
		if err := gob.NewDecoder(bytes.NewReader(settings.ConfigValue)).Decode(&result); err != nil {
			return nil, "", 0, fmt.Errorf("failed to decode int: %w", err)
		}
		return result, Int, settings.Version, nil
	case Float:
		var result float64
		if err := gob.NewDecoder(bytes.NewReader(settings.ConfigValue)).Decode(&result); err != nil {
			return nil, "", 0, fmt.Errorf("failed to decode float: %w", err)
		}
		return result, Float, settings.Version, nil
	case StringArray:
		var result []string
		if err := gob.NewDecoder(bytes.NewReader(settings.ConfigValue)).Decode(&result); err != nil {
			return nil, "", 0, fmt.Errorf("failed to decode string array: %w", err)
		}
		return result, StringArray, settings.Version, nil
	case IntArray:
		var result []int64
		if err := gob.NewDecoder(bytes.NewReader(settings.ConfigValue)).Decode(&result); err != nil {
			return nil, "", 0, fmt.Errorf("failed to decode int array: %w", err)
		}
		return result, IntArray, settings.Version, nil
	case FloatArray:
		var result []float64
		if err := gob.NewDecoder(bytes.NewReader(settings.ConfigValue)).Decode(&result); err != nil {
			return nil, "", 0, fmt.Errorf("failed to decode float array: %w", err)
		}
		return result, FloatArray, settings.Version, nil
	}

	return nil, "", 0, fmt.Errorf("unknown data type for namespace %s, key %s", namespace, key)
}

func (h *Handler) SettingsSet(namespace, key string, valueType SettingType, version int, value any) error {
	if h.daoObj == nil {
		return errors.New("database is not opened")
	}

	var dataValBytes []byte
	buf := new(bytes.Buffer)
	encoder := gob.NewEncoder(buf)

	switch valueType {
	case String:
		str, ok := value.(string)
		if !ok {
			return fmt.Errorf("invalid value type for %s: expected string", valueType)
		}
		dataValBytes = []byte(str)
		if len(dataValBytes) == 0 {
			dataValBytes = []byte{0}
		}
	case Bool:
		boolVal, ok := value.(bool)
		if !ok {
			return fmt.Errorf("invalid value type for %s: expected bool", valueType)
		}
		if boolVal {
			dataValBytes = []byte{1}
		} else {
			dataValBytes = []byte{0}
		}
	case Int, Float, StringArray, IntArray, FloatArray:
		if err := encoder.Encode(value); err != nil {
			return fmt.Errorf("failed to encode %s: %w", valueType, err)
		}
		dataValBytes = buf.Bytes()
	default:
		return fmt.Errorf("unsupported SettingType: %s", valueType)
	}

	settings := model.UserSettings{
		Namespace:   namespace,
		ConfigKey:   key,
		ConfigValue: dataValBytes,
		Version:     version,
		ConfigType:  string(valueType),
	}
	err := h.daoObj.Database.
		Model(settings).
		Where("namespace = ? AND config_key = ?", namespace, key).
		Assign(settings).
		FirstOrCreate(&settings).
		Error
	if err != nil {
		return fmt.Errorf("failed to set settings for namespace %s, key %s: %w", namespace, key, err)
	}

	return nil
}

func (h *Handler) SettingsInit(namespace, key string, valueType SettingType, version int, value any) (bool, error) {
	if h.daoObj == nil {
		return false, errors.New("database is not opened")
	}

	settingsVal, readValueType, versionDB, err := h.SettingsGet(namespace, key)
	if settingsVal != nil && valueType == readValueType && err == nil {
		return false, nil
	}

	if versionDB != version {
		return true, nil
	}

	err = h.SettingsSet(namespace, key, valueType, version, value)
	if err != nil {
		return false, fmt.Errorf("failed to initialize settings for namespace %s, key %s: %w", namespace, key, err)
	}

	return false, nil
}

func (h *Handler) removeInvisible(s string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsControl(r) {
			return -1
		}
		switch r {
		case '\u200B', '\u200C', '\u200D', '\uFEFF':
			return -1
		}
		return r
	}, s)
}
