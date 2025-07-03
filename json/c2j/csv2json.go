package c2j

import (
	"github.com/sascha-andres/jsonedit/internal/csv2json"
)

type C2JOptions struct {
	// Array indicates whether the JSON output should be wrapped in an array format during the mapping process.
	Array bool

	// Named indicates whether the CSV header should be used for mapping column names to JSON fields.
	Named bool

	// OutputType specifies a custom Marshaler to be used for serializing data during the mapping process. json, yaml or toml
	OutputType string

	// NestedPropertyName specifies the property name to use for TOML array output (defaults to "data")
	NestedPropertyName string

	// separator defines the byte value used as a delimiter or boundary in certain operations within the Mapper.
	Separator string
}

// MapCSV2JSON converts CSV data to JSON, applying mapping rules and configurations provided via options and mappings.
func MapCSV2JSON(options C2JOptions, in, mapping []byte) ([]byte, string, error) {
	opts := make([]csv2json.OptionFunc, 0)
	opts = append(opts, csv2json.WithArray(options.Array))
	opts = append(opts, csv2json.WithNestedPropertyName(options.NestedPropertyName))
	opts = append(opts, csv2json.WithNamed(options.Named))
	opts = append(opts, csv2json.WithOutputType(options.OutputType))
	opts = append(opts, csv2json.WithSeparator(options.Separator))
	opts = append(opts, csv2json.WithOptions(mapping))
	app, err := csv2json.NewMapper(opts...)
	if err != nil {
		return nil, "", err
	}
	out, err := app.Map(in)
	if err != nil {
		return nil, "", err
	}
	return out, getContentType(options), nil
}

// getContentType determines the MIME type for the output format based on the OutputType property in the C2JOptions struct.
func getContentType(options C2JOptions) string {
	switch options.OutputType {
	case "json":
		return "application/json"
	case "yaml":
		return "application/yaml"
	case "toml":
		return "application/toml"
	default:
		return "application/text"
	}
}
