package compare

import (
	"testing"
)

// TestGetJSONComparison tests the getJSONComparison method of the App struct
func TestGetJSONComparison(t *testing.T) {
	// Test case 1: Identical objects
	t.Run("IdenticalObjects", func(t *testing.T) {
		obj1 := map[string]interface{}{
			"name": "John",
			"age":  30,
			"address": map[string]interface{}{
				"city":  "New York",
				"state": "NY",
			},
		}
		obj2 := map[string]interface{}{
			"name": "John",
			"age":  30,
			"address": map[string]interface{}{
				"city":  "New York",
				"state": "NY",
			},
		}

		diff, err := GetJSONComparison(obj1, obj2, "  ")
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if diff != "" {
			t.Errorf("Expected empty diff for identical objects, got: %s", diff)
		}
	})

	// Test case 2: Different objects
	t.Run("DifferentObjects", func(t *testing.T) {
		obj1 := map[string]interface{}{
			"name": "John",
			"age":  30,
		}
		obj2 := map[string]interface{}{
			"name": "Jane",
			"age":  25,
		}

		diff, err := GetJSONComparison(obj1, obj2, "  ")
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if diff == "" {
			t.Errorf("Expected non-empty diff for different objects, got empty string")
		}
	})

	// Test case 3: Nested different objects
	t.Run("NestedDifferentObjects", func(t *testing.T) {
		obj1 := map[string]interface{}{
			"name": "John",
			"address": map[string]interface{}{
				"city":  "New York",
				"state": "NY",
			},
		}
		obj2 := map[string]interface{}{
			"name": "John",
			"address": map[string]interface{}{
				"city":  "Boston",
				"state": "MA",
			},
		}

		diff, err := GetJSONComparison(obj1, obj2, "  ")
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if diff == "" {
			t.Errorf("Expected non-empty diff for nested different objects, got empty string")
		}
	})

	// Test case 4: Array objects
	t.Run("ArrayObjects", func(t *testing.T) {
		obj1 := []interface{}{1, 2, 3}
		obj2 := []interface{}{1, 2, 4}

		diff, err := GetJSONComparison(obj1, obj2, "  ")
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if diff == "" {
			t.Errorf("Expected non-empty diff for different arrays, got empty string")
		}
	})

	// Test case 5: Unmarshalable object
	t.Run("UnmarshalableObject", func(t *testing.T) {
		// Create a circular reference that can't be marshaled to JSON
		obj1 := make(map[string]interface{})
		obj1["self"] = obj1

		obj2 := map[string]interface{}{"name": "John"}

		_, err := GetJSONComparison(obj1, obj2, "  ")
		if err == nil {
			t.Errorf("Expected error for unmarshalable object, got nil")
		}
	})

	// Test case 6: Different indentation
	t.Run("DifferentIndentation", func(t *testing.T) {
		// Create a new App with different indentation
		obj1 := map[string]interface{}{"name": "John"}
		obj2 := map[string]interface{}{"name": "John"}

		diff, err := GetJSONComparison(obj1, obj2, "    ")
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if diff != "" {
			t.Errorf("Expected empty diff for identical objects with different indentation, got: %s", diff)
		}
	})
}
