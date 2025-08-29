package csv2json

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v3"
)

// OptionFunc defines a function signature for configuring a Mapper instance with specific options or parameters.
type OptionFunc func(*Mapper) error

// WithArray sets the "array" field of the Mapper instance to the provided boolean value.
func WithArray(array bool) OptionFunc {
	return func(mapper *Mapper) error {
		mapper.array = array
		return nil
	}
}

// WithNamed sets the "named" field of the Mapper instance to the provided boolean value.
func WithNamed(named bool) OptionFunc {
	return func(mapper *Mapper) error {
		mapper.named = named
		return nil
	}
}

// WithSeparator returns an OptionFunc that sets the character delimiter for CSV records in the Mapper instance.
func WithSeparator(separator string) OptionFunc {
	return func(mapper *Mapper) error {
		if len(separator) != 1 {
			return errors.New(fmt.Sprintf("separator must be 1 character long (%q)", separator))
		}
		mapper.separator = rune(separator[0])
		return nil
	}
}

// WithOutputType sets the specified output type for marshaling data in a Mapper instance.
func WithOutputType(outputType string) OptionFunc {
	return func(mapper *Mapper) error {
		mapper.marshalWith = outputType
		switch outputType {
		case "json":
		case "yaml":
		case "toml":
			break
		case "":
			mapper.marshalWith = "json"
			break
		default:
			return errors.New(fmt.Sprintf("unknown marshaling type %q", outputType))
		}
		return nil
	}
}

// WithOptions sets the options for the mapper
func WithOptions(mapping []byte) OptionFunc {
	return func(mapper *Mapper) error {
		return json.Unmarshal(mapping, &mapper.configuration)
	}
}

// WithNestedPropertyName sets the property name for TOML array output.
func WithNestedPropertyName(propertyName string) OptionFunc {
	return func(mapper *Mapper) error {
		mapper.nestedPropertyName = propertyName
		return nil
	}
}

// NewMapper creates and initializes a new Mapper instance using the provided OptionFunc configurations.
func NewMapper(options ...OptionFunc) (*Mapper, error) {
	mapper := &Mapper{separator: ','}
	for _, option := range options {
		if err := option(mapper); err != nil {
			return nil, err
		}
	}
	switch mapper.marshalWith {
	case "json":
		mapper.marshaler = json.Marshal
		break
	case "yaml":
		mapper.array = true
		mapper.marshaler = yaml.Marshal
		break
	case "toml":
		mapper.array = true
		mapper.marshaler = toml.Marshal
		break
	}

	for _, configuration := range mapper.configuration.Mapping {
		if !configuration.IsValid() {
			return nil, errors.New("invalid mapping configuration")
		}
	}

	return mapper, nil
}

// Map processes input CSV data, maps it to JSON according to the configuration, and writes the result to the output destination.
func (m *Mapper) Map(in []byte) ([]byte, error) {
	out := make([]byte, 0)
	reader, writer, err := m.initialize(in, out)
	if err != nil {
		return nil, err
	}

	// Cast writer to *bytes.Buffer to retrieve the result later
	buffer, ok := writer.(*bytes.Buffer)
	if !ok {
		return nil, errors.New("writer is not a bytes.Buffer")
	}

	csvIn := csv.NewReader(reader)
	csvIn.Comma = m.separator
	csvIn.ReuseRecord = false

	var (
		arrResult []map[string]any
		header    []string
	)

	// Read header if needed
	if m.named {
		header, err = csvIn.Read()
		if err != nil {
			return nil, err
		}
	}
	// from now on we can reuse the record
	csvIn.ReuseRecord = true
	if m.array {
		arrResult = make([]map[string]any, 0)
	}
	recordNumber := 0
	// Read all records
	for {
		record, err := csvIn.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		out := make(map[string]interface{})
		out, err = m.mapCSVFields(record, header, out)
		if err != nil {
			return nil, err
		}
		// calculated fields
		out, err = m.applyCalculatedFields(record, header, recordNumber, out, "record")
		if err != nil {
			return nil, err
		}
		if m.array {
			arrResult = append(arrResult, out)
		} else {
			d, err := m.marshaler(out)
			if err != nil {
				return nil, err
			}
			if recordNumber > 0 {
				_, _ = writer.Write([]byte("\n"))
			}
			_, err = writer.Write(d)
			if err != nil {
				return nil, err
			}
		}
		recordNumber++
	}
	if m.array {
		var d []byte
		if m.marshalWith == "toml" || m.nestedPropertyName != "" {
			// Set default property name if not specified
			propertyName := "data"
			if m.nestedPropertyName != "" {
				propertyName = m.nestedPropertyName
			}

			// Create a map with the custom property name as the key
			outputData := map[string]any{
				propertyName: arrResult,
			}
			outputData, err = m.applyCalculatedFields(nil, nil, recordNumber, outputData, "document")
			if err != nil {
				return nil, err
			}
			d, err = m.marshaler(outputData)
			if err != nil {
				return nil, err
			}
		} else {
			d, err = m.marshaler(arrResult)
			if err != nil {
				return nil, err
			}
		}
		_, _ = writer.Write(d)
	}

	// Return the actual content from the buffer
	return buffer.Bytes(), nil
}

// mapCSVFields maps CSV records to a nested output structure using a header and mapping configuration. Returns the updated map or an error.
func (m *Mapper) mapCSVFields(record []string, header []string, out map[string]any) (map[string]any, error) {
	for i := range record {
		key := fmt.Sprintf("%d", i)
		if m.named {
			key = header[i]
		}
		var (
			v  ColumnConfiguration
			ok bool
		)
		if v, ok = m.configuration.Mapping[key]; !ok {
			return out, nil
		}
		val, err := convertToType(v.Type, record[i])
		if err != nil {
			return nil, err
		}
		out = setValue(strings.Split(v.Property, "."), val, out)
	}
	return out, nil
}

// applyCalculatedFields applies calculated fields to the output based on the configuration and specified record number.
func (m *Mapper) applyCalculatedFields(record, header []string, recordNumber int, out map[string]any, loc string) (map[string]any, error) {
	var err error

	for _, field := range m.configuration.Calculated {
		if field.Location != loc {
			continue
		}
		var val any
		switch field.Kind {
		case "application":
			val, err = m.getApplicationValue(field, recordNumber)
			if err != nil {
				return nil, err
			}
			break
		case "datetime":
			val, err = m.getDateTimeValue(field)
			if err != nil {
				return nil, err
			}
			break
		case "environment":
			e := os.Getenv(field.Format)
			val, err = convertToType(field.Type, e)
			if err != nil {
				return nil, err
			}
			break
		case "extra":
			e, ok := m.configuration.ExtraVariables[field.Format]
			if !ok {
				return nil, errors.New("extra variable " + field.Format + " not found")
			}
			val, err = convertToType(field.Type, e.Value)
			if err != nil {
				return nil, err
			}
		case "mapping":
			if record == nil {
				continue
			}
			splitFormat := strings.Split(field.Format, ":")
			if len(splitFormat) != 2 {
				return nil, errors.New(fmt.Sprintf("expected format field:mapping list, %q", field.Format))
			}
			var currentValue string
			if m.named {
				if !slices.Contains(header, splitFormat[0]) {
					return nil, errors.New("mapping field " + splitFormat[0] + " not found in header")
				}
				currentValue = record[slices.Index(header, splitFormat[0])]
			} else {
				var i int
				if i, err = strconv.Atoi(splitFormat[0]); err != nil {
					return nil, errors.New("mapping field " + splitFormat[0] + " not found as it is an invalid index")
				} else if i >= len(record) {
					return nil, errors.New("mapping field " + splitFormat[0] + " not found as it does not exist in the record")
				}
				currentValue = record[i]
			}
			splitMappings := strings.Split(splitFormat[1], ",")
			var (
				defaultMapping *string
				isSet          bool
			)
			for _, splitMapping := range splitMappings {
				splitMapping := strings.Split(splitMapping, "=")
				if len(splitMapping) != 2 {
					return nil, errors.New(fmt.Sprintf("expected format from=to list, %q", splitMapping))
				}
				if splitMapping[0] == currentValue {
					val, err = convertToType(field.Type, splitMapping[1])
					isSet = true
					break
				}
				if splitMapping[0] == "default" {
					defaultMapping = &splitMapping[1]
					break
				}
			}
			if !isSet && defaultMapping != nil {
				val, err = convertToType(field.Type, *defaultMapping)
			}
		default:
			return nil, errors.New("unknown kind " + field.Kind)
		}
		out = setValue(strings.Split(field.Property, "."), val, out)
	}
	return out, nil
}

// getDateTimeValue generates a date and time value formatted based on the Format field of the CalculatedField structure.
func (m *Mapper) getDateTimeValue(field CalculatedField) (any, error) {
	return time.Now().Format(field.Format), nil
}

// getApplicationValue computes and returns a value based on the specified CalculatedField and index.
// Returns the computed value or an error if the field format is unknown.
func (m *Mapper) getApplicationValue(field CalculatedField, i int) (any, error) {
	switch field.Format {
	case "record":
		return convertToType("int", strconv.Itoa(i))
	case "records":
		return convertToType("int", strconv.Itoa(i))
	}
	return nil, errors.New("unknown format " + field.Format)
}

// initialize initializes the Mapper instance by reading the mapping file and opening the input and output files.
func (m *Mapper) initialize(in, out []byte) (io.Reader, io.Writer, error) {
	return bytes.NewReader(in), bytes.NewBuffer(out), nil
}
