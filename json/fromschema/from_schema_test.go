package fromschema

import (
	"encoding/json"
	"log/slog"
	"os"
	"reflect"
	"testing"
)

func TestNewSchemaParser(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	tests := []struct {
		name       string
		jsonSchema string
		wantErr    bool
	}{
		{
			name: "Valid simple schema",
			jsonSchema: `{
				"type": "object",
				"properties": {
					"name": { "type": "string" },
					"age": { "type": "integer" }
				},
				"required": ["name"]
			}`,
			wantErr: false,
		},
		{
			name: "Valid complex schema",
			jsonSchema: `{
				"type": "object",
				"properties": {
					"person": {
						"type": "object",
						"properties": {
							"name": { "type": "string" },
							"age": { "type": "integer" },
							"address": {
								"type": "object",
								"properties": {
									"street": { "type": "string" },
									"city": { "type": "string" }
								},
								"required": ["street"]
							}
						},
						"required": ["name", "address"]
					},
					"tags": {
						"type": "array",
						"items": { "type": "string" }
					}
				},
				"required": ["person"]
			}`,
			wantErr: false,
		},
		{
			name:       "Invalid schema",
			jsonSchema: `{"type": "object", "properties": {`,
			wantErr:    true,
		},
		{
			name:       "Empty schema",
			jsonSchema: `{}`,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser, err := NewSchemaParser(logger, []byte(tt.jsonSchema))
			if (err != nil) != tt.wantErr {
				t.Errorf("NewSchemaParser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && parser == nil {
				t.Errorf("NewSchemaParser() returned nil parser but no error")
			}
		})
	}
}

func TestCreateEmptyJSONDocument(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	tests := []struct {
		name       string
		jsonSchema string
		expected   map[string]interface{}
		wantErr    bool
	}{
		{
			name: "Simple schema with required fields",
			jsonSchema: `{
				"type": "object",
				"properties": {
					"name": { "type": "string" },
					"age": { "type": "integer" }
				},
				"required": ["name"]
			}`,
			expected: map[string]interface{}{
				"name": "",
			},
			wantErr: false,
		},
		{
			name: "Schema with nested objects",
			jsonSchema: `{
				"type": "object",
				"properties": {
					"person": {
						"type": "object",
						"properties": {
							"name": { "type": "string" },
							"age": { "type": "integer" }
						},
						"required": ["name"]
					}
				},
				"required": ["person"]
			}`,
			expected: map[string]interface{}{
				"person": map[string]interface{}{
					"name": "",
				},
			},
			wantErr: false,
		},
		{
			name: "Schema with arrays",
			jsonSchema: `{
				"type": "object",
				"properties": {
					"tags": {
						"type": "array",
						"items": { "type": "string" }
					}
				},
				"required": ["tags"]
			}`,
			expected: map[string]interface{}{
				"tags": []interface{}{""},
			},
			wantErr: false,
		},
		{
			name: "Schema with array of objects",
			jsonSchema: `{
				"type": "object",
				"properties": {
					"people": {
						"type": "array",
						"items": {
							"type": "object",
							"properties": {
								"name": { "type": "string" },
								"age": { "type": "integer" }
							},
							"required": ["name"]
						}
					}
				},
				"required": ["people"]
			}`,
			expected: map[string]interface{}{
				"people": []interface{}{
					map[string]interface{}{
						"name": "",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Schema with all primitive types",
			jsonSchema: `{
				"type": "object",
				"properties": {
					"string_field": { "type": "string" },
					"number_field": { "type": "number" },
					"integer_field": { "type": "integer" },
					"boolean_field": { "type": "boolean" },
					"null_field": { "type": "null" }
				},
				"required": ["string_field", "number_field", "integer_field", "boolean_field", "null_field"]
			}`,
			expected: map[string]interface{}{
				"string_field":  "",
				"number_field":  float64(0),
				"integer_field": float64(0),
				"boolean_field": false,
				"null_field":    nil,
			},
			wantErr: false,
		},
		{
			name: "Schema with allOf",
			jsonSchema: `{
				"type": "object",
				"allOf": [
					{
						"type": "object",
						"properties": {
							"name": { "type": "string" }
						},
						"required": ["name"]
					},
					{
						"type": "object",
						"properties": {
							"age": { "type": "integer" }
						},
						"required": ["age"]
					}
				]
			}`,
			expected: map[string]interface{}{
				"name": "",
				"age":  float64(0),
			},
			wantErr: false,
		},
		{
			name:       "Empty schema",
			jsonSchema: `{}`,
			expected:   map[string]interface{}{},
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser, err := NewSchemaParser(logger, []byte(tt.jsonSchema))
			if err != nil {
				if !tt.wantErr {
					t.Errorf("Failed to create schema parser: %v", err)
				}
				return
			}

			doc, err := parser.CreateEmptyJSONDocument()
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateEmptyJSONDocument() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Convert to JSON and back to normalize the structure for comparison
			docJSON, err := json.Marshal(doc)
			if err != nil {
				t.Fatalf("Failed to marshal result: %v", err)
			}

			expectedJSON, err := json.Marshal(tt.expected)
			if err != nil {
				t.Fatalf("Failed to marshal expected: %v", err)
			}

			var normalizedDoc, normalizedExpected map[string]interface{}
			if err := json.Unmarshal(docJSON, &normalizedDoc); err != nil {
				t.Fatalf("Failed to unmarshal result: %v", err)
			}
			if err := json.Unmarshal(expectedJSON, &normalizedExpected); err != nil {
				t.Fatalf("Failed to unmarshal expected: %v", err)
			}

			if !reflect.DeepEqual(normalizedDoc, normalizedExpected) {
				t.Errorf("CreateEmptyJSONDocument() = %v, want %v", normalizedDoc, normalizedExpected)
			}
		})
	}
}
