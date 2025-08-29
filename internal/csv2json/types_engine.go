package csv2json

import "log/slog"

type (

	// Mapper defines a structure for mapping input data to output data, applying configuration and marshaling as needed.
	Mapper struct {

		// in specifies the input file path or '-' for standard input in the mapping process.
		in string

		// out specifies the output file path or '-' for standard output in the mapping process.
		out string

		// array indicates whether the JSON output should be wrapped in an array format during the mapping process.
		array bool

		// named indicates whether the CSV header should be used for mapping column names to JSON fields.
		named bool

		// mappingFile specifies the path to a JSON file containing the mapping configuration for the data transformation process.
		mappingFile string

		// marshalWith specifies a custom Marshaler to be used for serializing data during the mapping process. json, yaml or toml
		marshalWith string

		// nestedPropertyName specifies the property name to use for TOML array output (defaults to "data")
		nestedPropertyName string

		// marshaler defines a custom function for serializing a value of any type into a byte slice with error handling.
		marshaler func(v any) ([]byte, error)

		// configuration holds the mapping configuration used during the data transformation process.
		configuration Configuration

		// separator defines the byte value used as a delimiter or boundary in certain operations within the Mapper.
		separator rune

		// logger holds a reference to an slog.Logger for logging messages and errors within the Mapper's operations.
		logger *slog.Logger
	}

	// PropertyConfiguration defines the mapping configuration for a single property, including its name and data type.
	PropertyConfiguration struct {
		// Property specifies the name of the column property in the mapping configuration.
		Property string `json:"property"`

		// Type specifies the data type of the column in the mapping configuration.
		Type string `json:"type"`
	}

	// ColumnConfiguration defines the structure for configuring a column's property and type in a mapping.
	ColumnConfiguration struct {
		// Properties represents a list of property configurations defining column mappings in the mapping structure.
		Properties []PropertyConfiguration `json:"properties"`

		// Property specifies the name of the column property in the mapping configuration.
		Property string `json:"property"`

		// Type specifies the data type of the column in the mapping configuration.
		Type string `json:"type"`
	}

	// CalculatedField defines a structure for representing dynamically computed fields within a configuration.
	CalculatedField struct {

		// Property specifies the name of the property in the generated document.
		Property string `json:"property"`

		// Kind denotes what kind of calculate has to be treated
		Kind string `json:"kind"`

		// Format denotes a formatting value or the value to acquire
		Format string `json:"format"`

		// Type informa about the type represented
		Type string `json:"type"`

		// Location specifies the location or context where the calculated field applies in the mapping.
		//Either document or record
		Location string `json:"location"`
	}

	// ExtraVariable is used to store a static extra variable in the mapping
	ExtraVariable struct {

		// Value contains the string representation of the variable
		Value string `json:"value"`
	}

	// Configuration represents a mapping configuration where keys map to ColumnConfiguration structures.
	Configuration struct {

		// ExtraVariables is the list of all defines variables
		ExtraVariables map[string]ExtraVariable `json:"extra_variables"`

		// Calculated represents a slice of calculated fields to be included in the configuration, each defined by CalculatedField.
		Calculated []CalculatedField `json:"calculated"`

		// Mapping represents a map of keys to their corresponding column configurations in the mapping structure.
		Mapping map[string]ColumnConfiguration `json:"mapping"`
	}
)
