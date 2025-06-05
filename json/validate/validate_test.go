package validate

import (
	"io"
	"log/slog"
	"os"
	"strings"
	"testing"
)

func TestNewJSONValidator(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	tests := []struct {
		name    string
		opts    []JsonValidatorOption
		wantErr bool
	}{
		{
			name:    "No options",
			opts:    []JsonValidatorOption{},
			wantErr: false,
		},
		{
			name: "With logger",
			opts: []JsonValidatorOption{
				WithLogger(logger),
			},
			wantErr: false,
		},
		{
			name: "With valid JSON document",
			opts: []JsonValidatorOption{
				WithJSONDocument([]byte(`{"name": "test"}`)),
			},
			wantErr: false,
		},
		{
			name: "With valid JSON schema",
			opts: []JsonValidatorOption{
				WithJSONSchema([]byte(`{
					"type": "object",
					"properties": {
						"name": { "type": "string" }
					}
				}`)),
			},
			wantErr: false,
		},
		{
			name: "With invalid JSON schema",
			opts: []JsonValidatorOption{
				WithJSONSchema([]byte(`{
					"type": "object",
					"properties": {
				}`)),
			},
			wantErr: true,
		},
		{
			name: "With all valid options",
			opts: []JsonValidatorOption{
				WithLogger(logger),
				WithJSONDocument([]byte(`{"name": "test"}`)),
				WithJSONSchema([]byte(`{
					"type": "object",
					"properties": {
						"name": { "type": "string" }
					}
				}`)),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator, err := NewJSONValidator(tt.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewJSONValidator() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && validator == nil {
				t.Errorf("NewJSONValidator() returned nil validator but no error")
			}
		})
	}
}

func TestJSONValidator_Validate(t *testing.T) {
	tests := []struct {
		name       string
		jsonSchema string
		jsonDoc    string
		wantErr    bool
	}{
		{
			name: "Valid JSON against schema",
			jsonSchema: `{
				"type": "object",
				"properties": {
					"name": { "type": "string" },
					"age": { "type": "integer" }
				},
				"required": ["name"]
			}`,
			jsonDoc:    `{"name": "John", "age": 30}`,
			wantErr:    false,
		},
		{
			name: "Missing required field",
			jsonSchema: `{
				"type": "object",
				"properties": {
					"name": { "type": "string" },
					"age": { "type": "integer" }
				},
				"required": ["name"]
			}`,
			jsonDoc:    `{"age": 30}`,
			wantErr:    true,
		},
		{
			name: "Wrong type for field",
			jsonSchema: `{
				"type": "object",
				"properties": {
					"name": { "type": "string" },
					"age": { "type": "integer" }
				}
			}`,
			jsonDoc:    `{"name": "John", "age": "thirty"}`,
			wantErr:    true,
		},
		{
			name: "Complex schema validation",
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
					}
				},
				"required": ["person"]
			}`,
			jsonDoc: `{
				"person": {
					"name": "John",
					"age": 30,
					"address": {
						"street": "123 Main St",
						"city": "Anytown"
					}
				}
			}`,
			wantErr: false,
		},
		{
			name: "Complex schema validation with error",
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
					}
				},
				"required": ["person"]
			}`,
			jsonDoc: `{
				"person": {
					"name": "John",
					"age": 30,
					"address": {
						"city": "Anytown"
					}
				}
			}`,
			wantErr: true,
		},
		{
			name: "Array validation",
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
			jsonDoc:    `{"tags": ["tag1", "tag2", "tag3"]}`,
			wantErr:    false,
		},
		{
			name: "Array validation with error",
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
			jsonDoc:    `{"tags": ["tag1", 2, "tag3"]}`,
			wantErr:    true,
		},
		{
			name: "Empty object validation",
			jsonSchema: `{
				"type": "object"
			}`,
			jsonDoc:    `{}`,
			wantErr:    false,
		},
		{
			name: "Invalid JSON document",
			jsonSchema: `{
				"type": "object",
				"properties": {
					"name": { "type": "string" }
				}
			}`,
			jsonDoc:    `{"name": "test"`,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// No need to skip the invalid JSON document test - we should test this case

			// Create a validator with the schema
			validator, err := NewJSONValidator(
				WithJSONSchema([]byte(tt.jsonSchema)),
			)
			if err != nil {
				t.Fatalf("Failed to create validator: %v", err)
			}

			// Set the document on the validator regardless of whether it's valid JSON
			// The validator should handle invalid JSON appropriately

			// Set the document on the validator
			validator.document = io.NopCloser(strings.NewReader(tt.jsonDoc))

			// Validate
			err = validator.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWithJSONDocument(t *testing.T) {
	tests := []struct {
		name    string
		doc     []byte
		wantErr bool
	}{
		{
			name:    "Valid JSON document",
			doc:     []byte(`{"name": "test"}`),
			wantErr: false,
		},
		{
			name:    "Empty JSON document",
			doc:     []byte(`{}`),
			wantErr: false,
		},
		{
			name:    "Nil document",
			doc:     nil,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			option := WithJSONDocument(tt.doc)
			validator := &JSONValidator{}
			err := option(validator)
			if (err != nil) != tt.wantErr {
				t.Errorf("WithJSONDocument() error = %v, wantErr %v", err, tt.wantErr)
			}
			if validator.document == nil {
				t.Errorf("WithJSONDocument() did not set document")
			}
		})
	}
}

func TestWithLogger(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	tests := []struct {
		name    string
		logger  *slog.Logger
		wantErr bool
	}{
		{
			name:    "Valid logger",
			logger:  logger,
			wantErr: false,
		},
		{
			name:    "Nil logger",
			logger:  nil,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			option := WithLogger(tt.logger)
			validator := &JSONValidator{}
			err := option(validator)
			if (err != nil) != tt.wantErr {
				t.Errorf("WithLogger() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.logger != nil && validator.logger != tt.logger {
				t.Errorf("WithLogger() did not set logger correctly")
			}
		})
	}
}

func TestWithJSONSchema(t *testing.T) {
	tests := []struct {
		name       string
		jsonSchema []byte
		wantErr    bool
	}{
		{
			name: "Valid JSON schema",
			jsonSchema: []byte(`{
				"type": "object",
				"properties": {
					"name": { "type": "string" }
				}
			}`),
			wantErr: false,
		},
		{
			name: "Invalid JSON schema",
			jsonSchema: []byte(`{
				"type": "object",
				"properties": {
			}`),
			wantErr: true,
		},
		{
			name:       "Empty JSON schema",
			jsonSchema: []byte(`{}`),
			wantErr:    false,
		},
		{
			name:       "Nil schema",
			jsonSchema: nil,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			option := WithJSONSchema(tt.jsonSchema)
			validator := &JSONValidator{}
			err := option(validator)
			if (err != nil) != tt.wantErr {
				t.Errorf("WithJSONSchema() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				if validator.schema == nil {
					t.Errorf("WithJSONSchema() did not set schema")
				}
				if validator.compiler == nil {
					t.Errorf("WithJSONSchema() did not set compiler")
				}
			}
		})
	}
}
