package validate

import (
	"fmt"
	"io"
	"log/slog"
	"strings"

	"github.com/santhosh-tekuri/jsonschema/v6"

	"github.com/sascha-andres/jsonedit/json"
)

type JSONValidator struct {
	compiler *jsonschema.Compiler
	schema   *jsonschema.Schema
	document io.ReadCloser

	logger *slog.Logger
}

// JsonValidatorOption defines a functional option for configuring a JSONValidator instance.
type JsonValidatorOption func(*JSONValidator) error

// WithJSONDocument sets the JSON document for the JSONValidator using the provided byte slice.
func WithJSONDocument(doc []byte) JsonValidatorOption {
	return func(v *JSONValidator) error {
		v.document = io.NopCloser(strings.NewReader(string(doc)))
		return nil
	}
}

// WithLogger provides the logger to use
func WithLogger(logger *slog.Logger) JsonValidatorOption {
	return func(v *JSONValidator) error {
		v.logger = logger
		return nil
	}
}

// WithJSONSchema sets a JSON schema for the JSONValidator using the provided schema byte slice.
// It compiles the schema and attaches it to the validator instance.
// Returns an error if schema compilation fails.
func WithJSONSchema(jsonSchema []byte) JsonValidatorOption {
	return func(v *JSONValidator) error {
		// Initialize the compiler and add the schema file
		compiler := jsonschema.NewCompiler()
		compiler.UseLoader(json.InMemoryLoader{Doc: jsonSchema})

		schema, err := compiler.Compile("//")
		if err != nil {
			return fmt.Errorf("failed to compile schema: %w", err)
		}
		v.schema = schema
		v.compiler = compiler
		return nil
	}
}

// NewJSONValidator creates a new instance to validat4e a JSON file against a document
func NewJSONValidator(opts ...JsonValidatorOption) (*JSONValidator, error) {
	v := &JSONValidator{logger: slog.Default()}
	for i := range opts {
		err := opts[i](v)
		if err != nil {
			return nil, err
		}
	}
	return v, nil
}

// Validate checks the JSON document against the schema and returns an error if validation fails.
func (v *JSONValidator) Validate() error {
	return v.schema.Validate(v.document)
}
