package csv2json

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
