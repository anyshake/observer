package config

import "fmt"

func (s *services) GetValue(serviceName, keyField, typeName string) (any, error) {
	serviceConfig, ok := (*s)[serviceName]
	if !ok {
		return "", fmt.Errorf("service configuration \"%s\" not found", serviceName)
	}

	var err error
	switch typeName {
	case "bool":
		val, ok := serviceConfig.(map[string]any)[keyField].(bool)
		if !ok {
			return false, fmt.Errorf("key \"%s\" not found or is not a bool value", keyField)
		}
		return val, nil
	case "float64":
		val, ok := serviceConfig.(map[string]any)[keyField].(float64)
		if !ok {
			return 0.0, fmt.Errorf("key \"%s\" not found or is not a float64 value", keyField)
		}
		return val, nil
	case "int":
		val, ok := serviceConfig.(map[string]any)[keyField].(float64)
		if !ok {
			return 0, fmt.Errorf("key \"%s\" not found or is not a float64 value", keyField)
		}
		return int(val), nil
	case "string":
		val, ok := serviceConfig.(map[string]any)[keyField].(string)
		if !ok {
			return "", fmt.Errorf("key \"%s\" not found or is not a string value", keyField)
		}
		return val, nil
	default:
		err = fmt.Errorf("type \"%s\" is not supported", typeName)
	}

	return "", err
}
