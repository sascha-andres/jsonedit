package csv2json

import (
	"log/slog"
)

type (

	// RecordWithInformation represents a structured data record with a header and associated metadata.
	RecordWithInformation struct {

		// Record represents a collection of string values, typically used to store multiple related pieces of data.
		Record []string

		// Header represents a collection of strings typically used for storing key names or metadata information.
		Header []string

		// HeaderIndex maps header names to their corresponding column indices for quick lookup in a data processing context.
		HeaderIndex map[string]int
	}

	// Mapper defines a structure for mapping input data to output data, applying configuration and marshaling as needed.
	Mapper struct {

		//// in specifies the input file path or '-' for standard input in the mapping process.
		//in string
		//
		//// out specifies the output file path or '-' for standard output in the mapping process.
		//out string

		// array indicates whether the JSON output should be wrapped in an array format during the mapping process.
		array bool

		// named indicates whether the CSV header should be used for mapping column names to JSON fields.
		named bool

		// mappingFile specifies the path to a JSON file containing the mapping configuration for the data transformation process.
		mappingFile string

		// marshalWith specifies a custom Marshaler to be used for serializing data during the mapping process. json, yaml, or toml
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

		// newRecordFunc is called when a new record is being processed
		newRecordFunc NewRecordFunc

		// askForValueFunc defines a function type used to request a specific value dynamically during runtime operations.
		askForValueFunc AskForValueFunc

		// filteredNotification gets called if a record is filtered
		filteredNotification FilteredNotification

		// preProcess is called before the record is about to be processed and may be used for complex mappings
		preProcess PreProcess
	}

	// PropertyConfiguration defines the mapping configuration for a single property, including its name and data type.
	PropertyConfiguration struct {

		// Property specifies the name of the column property in the mapping configuration.
		Property string `json:"property"`

		// Type specifies the data type of the column in the mapping configuration.
		Type string `json:"type"`

		// Condition specifies the logical condition to be evaluated for a property, including value, column, and comparison operator.
		Condition *Condition `json:"condition,omitempty"`
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

	// FieldLocation represents the location or identifier of a specific field within a system or context as a string.
	FieldLocation string

	// CalculatedField defines a structure for representing dynamically computed fields within a configuration.
	CalculatedField struct {
		// Properties represent a list of property configurations defining column mappings in the mapping structure.
		Properties []PropertyConfiguration `json:"properties"`

		// Property specifies the name of the property in the generated document.
		Property string `json:"property"`

		// Type informa about the type represented
		Type string `json:"type"`

		// Kind denotes what kind of calculating has to be treated
		Kind string `json:"kind"`

		// Format denotes a formatting value or the value to acquire
		Format string `json:"format"`

		// Location specifies the location or context where the calculated field applies in the mapping.
		//Either document or record
		Location FieldLocation `json:"location"`
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

		// Filter represents a mapping of key-value pairs where each key maps to a slice of conditions for filtering data.
		Filter map[string]Conditions `json:"filter"`
	}

	// NewRecordFunc represents a function called when a new record is being processed
	NewRecordFunc func([]string, []string)

	// AskForValueFunc defines a function type that calculates a string value based on a record, header, and a calculated field.
	AskForValueFunc func(record, header []string, field CalculatedField) (string, error)

	// FilteredNotification is called if a record got filtered
	FilteredNotification func(record, header []string)

	// PreProcess is called before the record is about to be processed and may be used for complex mappings
	PreProcess func(record, header []string) ([]string, error)
)

const RecordLocation FieldLocation = "record"
const DocumentLocation FieldLocation = "document"
