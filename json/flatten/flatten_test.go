package flatten

import (
	"reflect"
	"testing"
)

func TestFlattenJSON(t *testing.T) {
	tests := []struct {
		name     string
		jsonDoc  string
		expected []string
		wantErr  bool
	}{
		{
			name:     "Empty JSON object",
			jsonDoc:  "{}",
			expected: []string{},
			wantErr:  false,
		},
		{
			name:     "Simple JSON object",
			jsonDoc:  `{"name": "John", "age": 30}`,
			expected: []string{"age: 30", "name: John"},
			wantErr:  false,
		},
		{
			name:     "Nested JSON object",
			jsonDoc:  `{"person": {"name": "John", "age": 30}}`,
			expected: []string{"person.age: 30", "person.name: John"},
			wantErr:  false,
		},
		{
			name:     "JSON with array",
			jsonDoc:  `{"numbers": [1, 2, 3]}`,
			expected: []string{"numbers/0: 1", "numbers/1: 2", "numbers/2: 3"},
			wantErr:  false,
		},
		{
			name:    "JSON with array requiring padding",
			jsonDoc: `{"numbers": [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]}`,
			expected: []string{
				"numbers/00: 1", "numbers/01: 2", "numbers/02: 3", "numbers/03: 4", "numbers/04: 5",
				"numbers/05: 6", "numbers/06: 7", "numbers/07: 8", "numbers/08: 9", "numbers/09: 10",
			},
			wantErr: false,
		},
		{
			name: "Complex JSON with nested objects and arrays",
			jsonDoc: `{
				"name": "John",
				"age": 30,
				"address": {
					"street": "123 Main St",
					"city": "Anytown"
				},
				"phones": [
					{"type": "home", "number": "555-1234"},
					{"type": "work", "number": "555-5678"}
				]
			}`,
			expected: []string{
				"address.city: Anytown",
				"address.street: 123 Main St",
				"age: 30",
				"name: John",
				"phones/0.number: 555-1234",
				"phones/0.type: home",
				"phones/1.number: 555-5678",
				"phones/1.type: work",
			},
			wantErr: false,
		},
		{
			name:     "JSON with null values",
			jsonDoc:  `{"name": "John", "age": null, "address": {"street": null}}`,
			expected: []string{"address.street: nil", "age: nil", "name: John"},
			wantErr:  false,
		},
		{
			name:     "Invalid JSON",
			jsonDoc:  `{"name": "John", "age": }`,
			expected: nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FlattenJSON([]byte(tt.jsonDoc))
			if (err != nil) != tt.wantErr {
				t.Errorf("FlattenJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("FlattenJSON() = %v, want %v", got, tt.expected)
			}
		})
	}
}
