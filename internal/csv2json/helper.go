package csv2json

import (
	"strconv"
	"strings"
	"time"
)

// convertToType converts the input string `val` to a specified type `t` such as "int", "float", or "bool".
// Returns the converted value as `any` or an error if the conversion fails.
func convertToType(t TypeInformation, val string) (any, error) {
	switch t.Type {
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
	if strings.HasPrefix(t.Type, "date") || strings.HasPrefix(t.Type, "time") {
		if t.ParseWith == "" && t.FormatWith == "" {
			format := "2006-01-02"
			if strings.HasPrefix(t.Type, "time") {
				format = "15:04:05"
			}
			d, err := time.Parse(format, val)
			if err != nil {
				return nil, err
			}
			return d, nil
		}
		if t.ParseWith != "" && t.FormatWith == "" {
			d, err := time.Parse(t.ParseWith, val)
			if err != nil {
				return nil, err
			}
			return d, nil
		}
		d, err := time.Parse(t.ParseWith, val)
		if err != nil {
			return nil, err
		}
		return d.Format(t.FormatWith), nil
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
