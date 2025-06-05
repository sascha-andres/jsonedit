package jsonedit

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
)

// FlattenJSON accepts a JSON document as a string and returns a sorted slice of strings
// where each string represents a flattened property:value pair.
// Nested objects are prefixed with parent property names.
// Array items are suffixed with /index where index is padded with leading zeros.
func FlattenJSON(jsonDoc string) ([]string, error) {
	var data interface{}
	if err := json.Unmarshal([]byte(jsonDoc), &data); err != nil {
		return nil, fmt.Errorf("error parsing JSON: %w", err)
	}

	lines := make([]string, 0)
	flattenValue("", data, &lines)

	// Sort the lines in ascending order
	sort.Strings(lines)

	return lines, nil
}

// flattenValue recursively processes JSON values and adds flattened property:value pairs to lines
func flattenValue(prefix string, value interface{}, lines *[]string) {
	switch v := value.(type) {
	case map[string]interface{}:
		// Handle objects
		for k, val := range v {
			newPrefix := k
			if prefix != "" {
				newPrefix = prefix + "." + k
			}
			flattenValue(newPrefix, val, lines)
		}
	case []interface{}:
		// Handle arrays
		padding := len(strconv.Itoa(len(v)))
		for i, val := range v {
			// Format index with leading zeros
			indexStr := fmt.Sprintf("%0*d", padding, i)
			newPrefix := prefix + "/" + indexStr
			flattenValue(newPrefix, val, lines)
		}
	case nil:
		// Handle null values
		*lines = append(*lines, fmt.Sprintf("%s: nil", prefix))
	default:
		// Handle primitive values
		*lines = append(*lines, fmt.Sprintf("%s: %v", prefix, v))
	}
}
