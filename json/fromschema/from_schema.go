package fromschema

import (
	"fmt"
	"log/slog"
	"slices"

	"github.com/sascha-andres/jsonedit/json"

	"github.com/santhosh-tekuri/jsonschema/v6"
)

// SchemaParser is a struct to hold the compiled schema and methods for parsing
type SchemaParser struct {
	compiler *jsonschema.Compiler
	schema   *jsonschema.Schema

	logger *slog.Logger
}

// NewSchemaParser creates a new SchemaParser instance
func NewSchemaParser(logger *slog.Logger, jsonSchema []byte) (*SchemaParser, error) {
	// Initialize the compiler and add the schema file
	compiler := jsonschema.NewCompiler()
	compiler.UseLoader(json.InMemoryLoader{Doc: jsonSchema})

	schema, err := compiler.Compile("//")
	if err != nil {
		return nil, fmt.Errorf("failed to compile schema: %w", err)
	}

	return &SchemaParser{
		compiler: compiler,
		schema:   schema,
		logger:   logger,
	}, nil
}

// CreateEmptyJSONDocument generates a JSON document with all required fields added empty
func (sp *SchemaParser) CreateEmptyJSONDocument() (interface{}, error) {
	doc := make(map[string]interface{})
	sp.processSchema(sp.schema, doc)
	return doc, nil
}

// processSchema recursively processes schema nodes to identify required fields
func (sp *SchemaParser) processSchema(schema *jsonschema.Schema, currentDoc map[string]interface{}) {
	for _, reqField := range schema.Required {
		if _, exists := currentDoc[reqField]; exists {
			continue
		}

		if propSchema, ok := schema.Properties[reqField]; ok {
			sp.logger.Debug("Adding empty field", "field", reqField)
			sp.addEmptyField(propSchema, currentDoc, reqField)
		}
	}

	for propName, propSchema := range schema.Properties {
		types := propSchema.Types.ToStrings()
		if slices.Contains(types, "object") {
			sp.logger.Debug("Recurse into 'properties' for nested required fields", "field", propName)
			if _, exists := currentDoc[propName]; !exists {
				currentDoc[propName] = make(map[string]interface{})
			}
			if nestedDoc, isNestedMap := currentDoc[propName].(map[string]interface{}); isNestedMap {
				sp.processSchema(propSchema, nestedDoc)
			}
		}
	}

	for _, subSchema := range schema.AllOf {
		sp.logger.Debug("Recurse into 'allOf' for nested required fields", "field", subSchema.ID)
		sp.processSchema(subSchema, currentDoc)
	}
}

// addEmptyField adds an empty field to the document based on the schema type
func (sp *SchemaParser) addEmptyField(fieldSchema *jsonschema.Schema, currentDoc map[string]interface{}, fieldName string) {
	types := fieldSchema.Types.ToStrings()
	if slices.Contains(types, "object") {
		nestedObject := make(map[string]interface{})
		currentDoc[fieldName] = nestedObject
		sp.processSchema(fieldSchema, nestedObject)
	} else if slices.Contains(fieldSchema.Types.ToStrings(), "array") {
		emptyArray := []interface{}{}
		if fieldSchema.Items2020 != nil {
			sp.logger.Debug("Adding empty array item (2020)", "field", fieldName)
			types = fieldSchema.Items2020.Types.ToStrings()
			sp.logger.Debug("Adding array item (2020)", "field", fieldName, "types", types)
			if slices.Contains(types, "object") {
				tempNestedObject := make(map[string]interface{})
				sp.processSchema(fieldSchema.Items2020, tempNestedObject)
				if len(tempNestedObject) > 0 {
					emptyArray = []interface{}{tempNestedObject}
				}
			} else {
				if slices.Contains(types, "string") {
					emptyArray = append(emptyArray, "")
				} else if slices.Contains(types, "number") || slices.Contains(types, "integer") {
					emptyArray = append(emptyArray, 0)
				} else if slices.Contains(types, "boolean") {
					emptyArray = append(emptyArray, false)
				} else {
					emptyArray = append(emptyArray, nil)
				}
			}
		} else if fieldSchema.Items != nil {
			sp.logger.Debug("Adding empty array item (Draft7 & older)", "field", fieldName)
			switch fieldSchema.Items.(type) {
			case *jsonschema.Schema:
				s := fieldSchema.Items.(*jsonschema.Schema)
				types = s.Types.ToStrings()
				sp.logger.Debug("Adding array item (Draft7 & older)", "field", fieldName, "types", types)
				if slices.Contains(types, "object") {
					tempNestedObject := make(map[string]interface{})
					sp.processSchema(s, tempNestedObject)
					if len(tempNestedObject) > 0 {
						emptyArray = []interface{}{tempNestedObject}
					}
				} else {
					if slices.Contains(types, "string") {
						emptyArray = append(emptyArray, "")
					} else if slices.Contains(types, "number") || slices.Contains(types, "integer") {
						emptyArray = append(emptyArray, 0)
					} else if slices.Contains(types, "boolean") {
						emptyArray = append(emptyArray, false)
					} else {
						emptyArray = append(emptyArray, nil)
					}
				}
			default:
				sp.logger.Warn("nil or []*jsonschema.Schema in Draft7 & older", "field", fieldName)
			}
		}
		currentDoc[fieldName] = emptyArray
	} else if slices.Contains(types, "string") {
		currentDoc[fieldName] = ""
	} else if slices.Contains(types, "number") || slices.Contains(types, "integer") {
		currentDoc[fieldName] = 0
	} else if slices.Contains(types, "boolean") {
		currentDoc[fieldName] = false
	} else if slices.Contains(types, "null") {
		currentDoc[fieldName] = nil
	} else {
		sp.logger.Warn("Unknown field type", "field", fieldName, "types", types)
		currentDoc[fieldName] = nil
	}
}
