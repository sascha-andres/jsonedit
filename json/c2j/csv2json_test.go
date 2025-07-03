package c2j

import (
	"strings"
	"testing"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v3"
)

func TestMapCSV2JSON(t *testing.T) {
	// Test case 1: Basic CSV to JSON conversion
	t.Run("BasicCSVToJSON", func(t *testing.T) {
		csvData := "name,age\nJohn,30\nJane,25"
		options := C2JOptions{
			Array:              false,
			Named:              true,
			OutputType:         "json",
			NestedPropertyName: "data",
			Separator:          ",",
		}

		result, contentType, err := MapCSV2JSON(options, []byte(csvData), []byte("{}"))
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if contentType != "application/json" {
			t.Errorf("Expected content type 'application/json', got '%s'", contentType)
		}

		// Verify the result is a valid JSON string
		if string(result) != "{}\n{}" {
			t.Errorf("Unexpected JSON result: %s", string(result))
		}
	})

	// Test case 2: Array output
	t.Run("ArrayOutput", func(t *testing.T) {
		csvData := "name,age\nJohn,30\nJane,25"
		options := C2JOptions{
			Array:              true,
			Named:              true,
			OutputType:         "json",
			NestedPropertyName: "data",
			Separator:          ",",
		}

		result, contentType, err := MapCSV2JSON(options, []byte(csvData), []byte("{}"))
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if contentType != "application/json" {
			t.Errorf("Expected content type 'application/json', got '%s'", contentType)
		}

		// Verify the result is a valid JSON string with expected format
		expectedResult := `{"data":[{},{}]}`
		if string(result) != expectedResult {
			t.Errorf("Expected result '%s', got '%s'", expectedResult, string(result))
		}
	})

	// Test case 3: YAML output
	t.Run("YAMLOutput", func(t *testing.T) {
		csvData := "name,age\nJohn,30"
		options := C2JOptions{
			Array:              false,
			Named:              true,
			OutputType:         "yaml",
			NestedPropertyName: "data",
			Separator:          ",",
		}

		result, contentType, err := MapCSV2JSON(options, []byte(csvData), []byte("{}"))
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if contentType != "application/yaml" {
			t.Errorf("Expected content type 'application/yaml', got '%s'", contentType)
		}

		// Verify the result is a valid YAML string with expected format
		// The exact format might vary slightly due to whitespace, so we'll check for key parts
		resultStr := string(result)
		if !strings.Contains(resultStr, "data:") {
			t.Errorf("Expected YAML result to contain 'data:', got: %s", resultStr)
		}

		// Verify the result can be unmarshaled as YAML
		var yamlResult map[string]interface{}
		err = yaml.Unmarshal(result, &yamlResult)
		if err != nil {
			t.Errorf("Failed to unmarshal result as YAML: %v", err)
		}

		// Verify the data field exists and is an array
		data, ok := yamlResult["data"]
		if !ok {
			t.Errorf("Expected 'data' field in YAML result, got: %v", yamlResult)
			return
		}

		// Check that data is a slice
		dataSlice, ok := data.([]interface{})
		if !ok {
			t.Errorf("Expected 'data' to be a slice, got: %T", data)
			return
		}

		// Check that we have at least one item in the data slice
		if len(dataSlice) == 0 {
			t.Errorf("Expected non-empty 'data' slice")
		}
	})

	// Test case 4: TOML output
	t.Run("TOMLOutput", func(t *testing.T) {
		csvData := "name,age\nJohn,30"
		options := C2JOptions{
			Array:              false,
			Named:              true,
			OutputType:         "toml",
			NestedPropertyName: "data",
			Separator:          ",",
		}

		result, contentType, err := MapCSV2JSON(options, []byte(csvData), []byte("{}"))
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if contentType != "application/toml" {
			t.Errorf("Expected content type 'application/toml', got '%s'", contentType)
		}

		// Verify the result is a valid TOML string with expected format
		resultStr := string(result)
		if !strings.Contains(resultStr, "[[data]]") {
			t.Errorf("Expected TOML result to contain '[[data]]', got: %s", resultStr)
		}

		// Verify the result can be unmarshaled as TOML
		var tomlResult map[string]interface{}
		err = toml.Unmarshal(result, &tomlResult)
		if err != nil {
			t.Errorf("Failed to unmarshal result as TOML: %v", err)
		}
	})

	// Test case 5: Custom separator
	t.Run("CustomSeparator", func(t *testing.T) {
		csvData := "name;age\nJohn;30"
		options := C2JOptions{
			Array:              false,
			Named:              true,
			OutputType:         "json",
			NestedPropertyName: "data",
			Separator:          ";",
		}

		result, contentType, err := MapCSV2JSON(options, []byte(csvData), []byte("{}"))
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if contentType != "application/json" {
			t.Errorf("Expected content type 'application/json', got '%s'", contentType)
		}

		// Verify the result is a valid JSON string
		if string(result) != "{}" {
			t.Errorf("Unexpected JSON result: %s", string(result))
		}
	})

	// Test case 6: Error - Invalid CSV
	t.Run("InvalidCSV", func(t *testing.T) {
		csvData := "name,age\nJohn,30,extra"
		options := C2JOptions{
			Array:              false,
			Named:              true,
			OutputType:         "json",
			NestedPropertyName: "data",
			Separator:          ",",
		}

		_, _, err := MapCSV2JSON(options, []byte(csvData), []byte("{}"))
		if err == nil {
			t.Errorf("Expected error for invalid CSV, got nil")
		}
	})

	// Test case 7: Error - Invalid mapping
	t.Run("InvalidMapping", func(t *testing.T) {
		csvData := "name,age\nJohn,30"
		options := C2JOptions{
			Array:              false,
			Named:              true,
			OutputType:         "json",
			NestedPropertyName: "data",
			Separator:          ",",
		}

		_, _, err := MapCSV2JSON(options, []byte(csvData), []byte("{invalid json"))
		if err == nil {
			t.Errorf("Expected error for invalid mapping, got nil")
		}
	})

	// Test case 8: Unnamed columns
	t.Run("UnnamedColumns", func(t *testing.T) {
		csvData := "John,30\nJane,25"
		options := C2JOptions{
			Array:              false,
			Named:              false,
			OutputType:         "json",
			NestedPropertyName: "data",
			Separator:          ",",
		}

		result, contentType, err := MapCSV2JSON(options, []byte(csvData), []byte("{}"))
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if contentType != "application/json" {
			t.Errorf("Expected content type 'application/json', got '%s'", contentType)
		}

		// Verify the result is a valid JSON string
		if string(result) != "{}\n{}" {
			t.Errorf("Unexpected JSON result: %s", string(result))
		}
	})

	// Test case 9: Unknown output type
	t.Run("UnknownOutputType", func(t *testing.T) {
		csvData := "name,age\nJohn,30"
		options := C2JOptions{
			Array:              false,
			Named:              true,
			OutputType:         "unknown",
			NestedPropertyName: "data",
			Separator:          ",",
		}

		_, _, err := MapCSV2JSON(options, []byte(csvData), []byte("{}"))
		if err == nil {
			t.Errorf("Expected error for unknown output type, got nil")
		}

		if !strings.Contains(err.Error(), "unknown marshaling type") {
			t.Errorf("Expected error message to contain 'unknown marshaling type', got: %v", err)
		}
	})
}
