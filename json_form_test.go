package jsonedit

import (
	"log/slog"
	"os"
	"strings"
	"testing"
)

// TestGenerateJSONForm tests the GenerateJSONForm method of the App struct
func TestGenerateJSONForm(t *testing.T) {
	// Create a new App instance with default settings and a test logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))

	// Test case 1: Empty object
	t.Run("EmptyObject", func(t *testing.T) {
		data := map[string]interface{}{}
		result := GenerateJSONForm(logger, false, data, "root", 0)

		// Check for expected elements in the output
		if !strings.Contains(result, "<em>Empty object</em>") {
			t.Errorf("Expected empty object message, not found in result")
		}
		if !strings.Contains(result, "class=\"add-property-btn\"") {
			t.Errorf("Expected add property button, not found in result")
		}
		if !strings.Contains(result, "data-path=\"root\"") {
			t.Errorf("Expected correct data path, not found in result")
		}
	})

	// Test case 2: Object with simple properties
	t.Run("ObjectWithSimpleProperties", func(t *testing.T) {
		data := map[string]interface{}{
			"name": "John",
			"age":  30,
		}
		result := GenerateJSONForm(logger, false, data, "person", 0)

		// Check for expected elements in the output
		if !strings.Contains(result, "<label for=\"person.name\">name:</label>") {
			t.Errorf("Expected label for name property, not found in result")
		}
		if !strings.Contains(result, "<input type=\"text\" name=\"person.name\" id=\"person.name\" value=\"John\">") {
			t.Errorf("Expected input field for name property, not found in result")
		}
		if !strings.Contains(result, "<label for=\"person.age\">age:</label>") {
			t.Errorf("Expected label for age property, not found in result")
		}
		if !strings.Contains(result, "<input type=\"text\" name=\"person.age\" id=\"person.age\" value=\"30\">") {
			t.Errorf("Expected input field for age property, not found in result")
		}
		if !strings.Contains(result, "class=\"delete-property-btn\"") {
			t.Errorf("Expected delete property button, not found in result")
		}
		if !strings.Contains(result, "class=\"add-property-btn\"") {
			t.Errorf("Expected add property button, not found in result")
		}
	})

	// Test case 3: Nested object
	t.Run("NestedObject", func(t *testing.T) {
		data := map[string]interface{}{
			"person": map[string]interface{}{
				"name": "John",
				"age":  30,
			},
		}
		result := GenerateJSONForm(logger, false, data, "root", 0)

		// Check for expected elements in the output
		if !strings.Contains(result, "<label for=\"root.person\">person:</label>") {
			t.Errorf("Expected label for person property, not found in result")
		}
		if !strings.Contains(result, "<label for=\"root.person.name\">name:</label>") {
			t.Errorf("Expected label for nested name property, not found in result")
		}
		if !strings.Contains(result, "<input type=\"text\" name=\"root.person.name\" id=\"root.person.name\" value=\"John\">") {
			t.Errorf("Expected input field for nested name property, not found in result")
		}
		// Check for proper indentation
		if !strings.Contains(result, "margin-left: 20px") {
			t.Errorf("Expected indentation for nested object, not found in result")
		}
	})

	// Test case 4: Empty array
	t.Run("EmptyArray", func(t *testing.T) {
		data := []interface{}{}
		result := GenerateJSONForm(logger, false, data, "items", 0)

		// Check for expected elements in the output
		if !strings.Contains(result, "<em>Empty array</em>") {
			t.Errorf("Expected empty array message, not found in result")
		}
		if !strings.Contains(result, "class=\"add-array-item-btn\"") {
			t.Errorf("Expected add array item button, not found in result")
		}
		if !strings.Contains(result, "data-path=\"items\"") {
			t.Errorf("Expected correct data path, not found in result")
		}
	})

	// Test case 5: Array with simple values
	t.Run("ArrayWithSimpleValues", func(t *testing.T) {
		data := []interface{}{1, "two", true}
		result := GenerateJSONForm(logger, false, data, "items", 0)

		// Check for expected elements in the output
		if !strings.Contains(result, "<label for=\"items[0]\">[0]:</label>") {
			t.Errorf("Expected label for first array item, not found in result")
		}
		if !strings.Contains(result, "<input type=\"text\" name=\"items[0]\" id=\"items[0]\" value=\"1\">") {
			t.Errorf("Expected input field for first array item, not found in result")
		}
		if !strings.Contains(result, "<label for=\"items[1]\">[1]:</label>") {
			t.Errorf("Expected label for second array item, not found in result")
		}
		if !strings.Contains(result, "<input type=\"text\" name=\"items[1]\" id=\"items[1]\" value=\"two\">") {
			t.Errorf("Expected input field for second array item, not found in result")
		}
		if !strings.Contains(result, "<label for=\"items[2]\">[2]:</label>") {
			t.Errorf("Expected label for third array item, not found in result")
		}
		if !strings.Contains(result, "<input type=\"text\" name=\"items[2]\" id=\"items[2]\" value=\"true\">") {
			t.Errorf("Expected input field for third array item, not found in result")
		}
		if !strings.Contains(result, "class=\"delete-array-item-btn\"") {
			t.Errorf("Expected delete array item button, not found in result")
		}
		if !strings.Contains(result, "class=\"add-array-item-btn\"") {
			t.Errorf("Expected add array item button, not found in result")
		}
	})

	// Test case 6: Array with nested objects
	t.Run("ArrayWithNestedObjects", func(t *testing.T) {
		data := []interface{}{
			map[string]interface{}{
				"name": "John",
				"age":  30,
			},
			map[string]interface{}{
				"name": "Jane",
				"age":  25,
			},
		}
		result := GenerateJSONForm(logger, false, data, "people", 0)

		// Check for expected elements in the output
		if !strings.Contains(result, "<label for=\"people[0]\">[0]:</label>") {
			t.Errorf("Expected label for first array item, not found in result")
		}
		if !strings.Contains(result, "<label for=\"people[0].name\">name:</label>") {
			t.Errorf("Expected label for name property in first object, not found in result")
		}
		if !strings.Contains(result, "<input type=\"text\" name=\"people[0].name\" id=\"people[0].name\" value=\"John\">") {
			t.Errorf("Expected input field for name property in first object, not found in result")
		}
		if !strings.Contains(result, "<label for=\"people[1]\">[1]:</label>") {
			t.Errorf("Expected label for second array item, not found in result")
		}
		if !strings.Contains(result, "<label for=\"people[1].name\">name:</label>") {
			t.Errorf("Expected label for name property in second object, not found in result")
		}
		if !strings.Contains(result, "<input type=\"text\" name=\"people[1].name\" id=\"people[1].name\" value=\"Jane\">") {
			t.Errorf("Expected input field for name property in second object, not found in result")
		}
		// Check for proper indentation
		if !strings.Contains(result, "margin-left: 20px") {
			t.Errorf("Expected indentation for array items, not found in result")
		}
	})

	// Test case 7: Primitive value
	t.Run("PrimitiveValue", func(t *testing.T) {
		data := "test string"
		result := GenerateJSONForm(logger, false, data, "value", 0)

		// Check for expected elements in the output
		if !strings.Contains(result, "<label for=\"value\">Value:</label>") {
			t.Errorf("Expected label for primitive value, not found in result")
		}
		if !strings.Contains(result, "<input type=\"text\" name=\"value\" id=\"value\" value=\"test string\">") {
			t.Errorf("Expected input field for primitive value, not found in result")
		}
	})

	// Test case 8: Read-only mode
	t.Run("ReadOnlyMode", func(t *testing.T) {
		data := map[string]interface{}{
			"name": "John",
			"age":  30,
		}
		result := GenerateJSONForm(logger, true, data, "person", 0)

		// Check for expected elements in the output
		if !strings.Contains(result, "<label>name:</label>") {
			t.Errorf("Expected label for name property, not found in result")
		}
		if !strings.Contains(result, "<span style=\"font-family: 'CustomMonoFont', monospace;\">John</span>") {
			t.Errorf("Expected span for name value, not found in result")
		}
		if !strings.Contains(result, "<label>age:</label>") {
			t.Errorf("Expected label for age property, not found in result")
		}
		if !strings.Contains(result, "<span style=\"font-family: 'CustomMonoFont', monospace;\">30</span>") {
			t.Errorf("Expected span for age value, not found in result")
		}

		// Check that edit controls are not present
		if strings.Contains(result, "class=\"add-property-btn\"") {
			t.Errorf("Found add property button in read-only mode, should not be present")
		}
		if strings.Contains(result, "class=\"delete-property-btn\"") {
			t.Errorf("Found delete property button in read-only mode, should not be present")
		}
		if strings.Contains(result, "<input") {
			t.Errorf("Found input field in read-only mode, should not be present")
		}
	})

	// Test case 9: HTML escaping
	t.Run("HTMLEscaping", func(t *testing.T) {
		data := map[string]interface{}{
			"html": "<script>alert('XSS')</script>",
		}
		result := GenerateJSONForm(logger, false, data, "root", 0)

		// Check that HTML is properly escaped
		if strings.Contains(result, "<script>alert('XSS')</script>") {
			t.Errorf("Found unescaped HTML in result, should be escaped")
		}
		if !strings.Contains(result, "&lt;script&gt;alert(&#39;XSS&#39;)&lt;/script&gt;") {
			t.Errorf("Expected escaped HTML, not found in result")
		}
	})

	// Test case 10: Null values
	t.Run("NullValues", func(t *testing.T) {
		data := map[string]interface{}{
			"nullValue": nil,
		}
		result := GenerateJSONForm(logger, false, data, "root", 0)

		// Check for expected elements in the output
		if !strings.Contains(result, "<label for=\"root.nullValue\">nullValue:</label>") {
			t.Errorf("Expected label for null value property, not found in result")
		}
		if !strings.Contains(result, "<input type=\"text\" name=\"root.nullValue\" id=\"root.nullValue\" value=\"\">") {
			t.Errorf("Expected empty input field for null value, not found in result")
		}
	})
}
