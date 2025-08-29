package csv2json

import (
	"strconv"
)

// convertToType converts the input string `val` to a specified type `t` such as "int", "float", or "bool".
// Returns the converted value as `any` or an error if the conversion fails.
func convertToType(t, val string) (any, error) {
	switch t {
	case "int":
		i, err := strconv.Atoi(val)
		if err != nil {
			return nil, err
		}
		return i, nil
	case "float":
		f, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return nil, err
		}
		return f, nil
	case "bool":
		b, err := strconv.ParseBool(val)
		if err != nil {
			return nil, err
		}
		return b, nil
	}
	return val, nil
}

// setValue creates and maps nested dictionaries based on a hierarchy of keys, assigning a final value.
func setValue(hierarchy []string, value any, data map[string]interface{}) map[string]interface{} {
	v := setValueInternal(hierarchy, value, data)
	data[hierarchy[0]] = v
	return data
}

// setValueInternal recursively creates and maps nested dictionaries based on a hierarchy of keys, assigning a final value.
func setValueInternal(hierarchy []string, value any, inside map[string]any) any {
	if len(hierarchy) == 1 {
		return value
	}
	v := make(map[string]any)
	if val, ok := inside[hierarchy[0]]; ok {
		if reflected, ok := val.(map[string]any); ok {
			v = reflected
		}
	}
	v[hierarchy[1]] = setValueInternal(hierarchy[1:], value, v)
	return v
}
