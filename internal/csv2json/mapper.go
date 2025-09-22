package csv2json

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v3"
)

// OptionFunc defines a function signature for configuring a Mapper instance with specific options or parameters.
type OptionFunc func(*Mapper) error

// WithLogger sets a custom logger for the Mapper instance. Returns an OptionFunc to configure the logger.
func WithLogger(logger *slog.Logger) OptionFunc {
	return func(mapper *Mapper) error {
		mapper.logger = logger
		return nil
	}
}

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

	for _, field := range mapper.configuration.Calculated {
		if !field.Location.IsValid() {
			return nil, errors.New("invalid calculated field location")
		}
	}

	for key, configuration := range mapper.configuration.Mapping {
		if !configuration.IsValid(mapper.logger, key) {
			return nil, errors.New("invalid mapping configuration")
		}
	}

	return mapper, nil
}

// SetNewRecordFunc sets the function to be called when processing a new record.
func (m *Mapper) SetNewRecordFunc(n NewRecordFunc) {
	m.newRecordFunc = n
}

// SetAskForValueFunc sets the function responsible for dynamically providing string values based on record, header, and field.
func (m *Mapper) SetAskForValueFunc(f AskForValueFunc) {
	m.askForValueFunc = f
}

// SetFilteredNotification can be used to provide a callback to be notified if a record is filtered
func (m *Mapper) SetFilteredNotification(f FilteredNotification) {
	m.filteredNotification = f
}

// MapIo processes CSV input data, applies mapping logic, and writes the mapped output to the provided writer.
func (m *Mapper) MapIo(in io.Reader, writer io.Writer) error {
	csvIn := csv.NewReader(in)
	csvIn.Comma = m.separator
	csvIn.ReuseRecord = false

	var (
		arrResult []map[string]any
		header    []string
		err       error
	)

	var headerIndex map[string]int
	// Read header if needed
	if m.named {
		header, err = csvIn.Read()
		if err != nil {
			return err
		}
		// build header index cache
		headerIndex = make(map[string]int, len(header))
		for i, h := range header {
			headerIndex[h] = i
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
			return err
		}
		recordInfo := &RecordWithInformation{
			Record: record,
			Header: header,
		}
		if headerIndex != nil {
			recordInfo.HeaderIndex = headerIndex
		}
		filtered := false
		for group, conditions := range m.configuration.Filter {
			if m.logger != nil {
				m.logger.Debug("checking filter", slog.Any("group", group), slog.Any("conditions", conditions))
			}
			if conditions.Apply(m.logger, group, recordInfo) {
				filtered = true
				if m.filteredNotification != nil {
					m.filteredNotification(record, header)
				}
				break
			}
		}
		if filtered {
			continue
		}
		if m.newRecordFunc != nil {
			m.newRecordFunc(record, header)
		}
		out := make(map[string]interface{})
		out, err = m.mapCSVFields(out, recordInfo)
		if err != nil {
			return err
		}
		// calculated fields
		out, err = m.applyCalculatedFields(recordNumber, out, "record", recordInfo)
		if err != nil {
			return err
		}
		if m.array {
			arrResult = append(arrResult, out)
		} else {
			d, err := m.marshaler(out)
			if err != nil {
				return err
			}
			if recordNumber > 0 {
				_, _ = writer.Write([]byte("\n"))
			}
			_, err = writer.Write(d)
			if err != nil {
				return err
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
			outputData, err = m.applyCalculatedFields(recordNumber, outputData, "document", nil)
			if err != nil {
				return err
			}
			d, err = m.marshaler(outputData)
			if err != nil {
				return err
			}
		} else {
			d, err = m.marshaler(arrResult)
			if err != nil {
				return err
			}
		}
		_, _ = writer.Write(d)
	}
	return nil
}

// Map processes input CSV data, maps it to JSON according to the configuration, and writes the result to the output destination.
func (m *Mapper) Map(in []byte) ([]byte, error) {
	out := make([]byte, 0)
	reader, writer, err := m.initialize(in, out)
	if err != nil {
		return nil, err
	}

	err = m.MapIo(reader, writer)
	if err != nil {
		return nil, err
	}

	// Cast writer to *bytes.Buffer to retrieve the result later
	buffer, ok := writer.(*bytes.Buffer)
	if !ok {
		return nil, errors.New("writer is not a bytes.Buffer")
	}

	// Return the actual content from the buffer
	return buffer.Bytes(), nil
}

// mapCSVFields maps CSV records to a nested output structure using a header and mapping configuration. Returns the updated map or an error.
func (m *Mapper) mapCSVFields(out map[string]any, recordInfo *RecordWithInformation) (map[string]any, error) {
	for i := range recordInfo.Record {
		key := fmt.Sprintf("%d", i)
		if recordInfo.HeaderIndex != nil {
			key = recordInfo.Header[i]
		}
		var (
			v  ColumnConfiguration
			ok bool
		)
		if v, ok = m.configuration.Mapping[key]; !ok {
			continue
		}
		if len(v.Properties) > 0 {
			for _, property := range v.Properties {
				if property.Condition != nil {
					if !property.Applies(m.logger, recordInfo) {
						continue
					}
				}
				val, err := convertToType(property.Type, recordInfo.Record[i])
				if err != nil {
					return nil, err
				}
				out = setValue(strings.Split(property.Property, "."), val, out)
			}
		} else {
			val, err := convertToType(v.Type, recordInfo.Record[i])
			if err != nil {
				return nil, err
			}
			out = setValue(strings.Split(v.Property, "."), val, out)
		}
	}
	return out, nil
}

// applyCalculatedFields applies calculated fields to the output based on the configuration and specified record number.
func (m *Mapper) applyCalculatedFields(recordNumber int, out map[string]any, loc FieldLocation, recordInfo *RecordWithInformation) (map[string]any, error) {
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
			if recordInfo == nil {
				continue
			}
			splitFormat := strings.Split(field.Format, ":")
			if len(splitFormat) != 2 {
				return nil, errors.New(fmt.Sprintf("expected format field:mapping list, %q", field.Format))
			}
			var currentValue string
			if recordInfo.HeaderIndex != nil {
				idx, ok := recordInfo.HeaderIndex[splitFormat[0]]
				if !ok {
					return nil, errors.New("mapping field " + splitFormat[0] + " not found in header")
				}
				if idx >= len(recordInfo.Record) {
					return nil, errors.New("mapping field " + splitFormat[0] + " not found as it does not exist in the record")
				}
				currentValue = recordInfo.Record[idx]
			} else {
				var i int
				if i, err = strconv.Atoi(splitFormat[0]); err != nil {
					return nil, errors.New("mapping field " + splitFormat[0] + " not found as it is an invalid index")
				} else if i >= len(recordInfo.Record) {
					return nil, errors.New("mapping field " + splitFormat[0] + " not found as it does not exist in the record")
				}
				currentValue = recordInfo.Record[i]
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
		case "ask":
			if m.askForValueFunc == nil {
				return nil, errors.New("ask value function not set")
			}
			var answer string
			if recordInfo == nil {
				answer, err = m.askForValueFunc(nil, nil, field)
			} else {
				answer, err = m.askForValueFunc(recordInfo.Record, recordInfo.Header, field)
			}
			if err != nil {
				return nil, err
			}
			val, err = convertToType(field.Type, answer)
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
