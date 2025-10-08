package splitter

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
)

// OptionFunc defines a function type used to configure a Splitter by modifying its internal properties.
type OptionFunc func(*Splitter) error

// WithConfigurationFile sets the configuration for the Splitter and returns an OptionFunc.
func WithConfigurationFile(cfg string) OptionFunc {
	return func(s *Splitter) error {
		s.Configuration = cfg
		var conf Configuration
		d, err := os.ReadFile(cfg)
		if err != nil {
			return err
		}
		err = json.Unmarshal(d, &conf)
		if err != nil {
			return err
		}
		s.config = conf
		return nil
	}
}

// WithConfiguration sets the configuration for the Splitter and returns an OptionFunc.
func WithConfiguration(conf Configuration) OptionFunc {
	return func(s *Splitter) error {
		s.config = conf
		return nil
	}
}

// WithLogger sets a custom slog.Logger for the Splitter to use for logging operations and returns an OptionFunc.
func WithLogger(logger *slog.Logger) OptionFunc {
	return func(s *Splitter) error {
		s.logger = logger
		return nil
	}
}

// WithOutputEmptyGroups configures the Splitter to include or exclude empty groups in the output based on the specified value.
func WithOutputEmptyGroups(outputEmptyGroups bool) OptionFunc {
	return func(s *Splitter) error {
		s.OutputEmptyGroups = outputEmptyGroups
		return nil
	}
}

// WithArrayPath sets the array path for splitting and duplicates non-array properties if the path is not empty.
func WithArrayPath(arrayPath string) OptionFunc {
	return func(s *Splitter) error {
		s.arrayPath = arrayPath
		return nil
	}
}

// NewSplitter creates a new instance of the worker
func NewSplitter(opts ...OptionFunc) (*Splitter, error) {
	app := &Splitter{}
	for _, opt := range opts {
		err := opt(app)
		if err != nil {
			return nil, err
		}
	}
	return app, nil
}

// isJSONArray checks if the provided interface{} is a JSON array
func (s *Splitter) isJSONArray(data any) bool {
	_, ok := data.([]any)
	return ok
}

// isValueType checks if the provided item is a primitive value type (string, number, boolean)
func (s *Splitter) isValueType(item any) bool {
	switch item.(type) {
	case string, float64, bool, nil:
		return true
	default:
		return false
	}
}

// allArrayElementsAreValues checks if all elements in the array are primitive values
func (s *Splitter) allArrayElementsAreValues(arr []any) bool {
	for _, item := range arr {
		if !s.isValueType(item) {
			return false
		}
	}
	return true
}

// Split reads input data as bytes, processes it, and splits it into groups defined by the configuration.
func (s *Splitter) Split(in []byte) (map[string][]any, error) {
	var input any
	err := json.Unmarshal(in, &input)
	if err != nil {
		return nil, err
	}

	dataToSplit, err := s.extractDataToSplit(input)
	if err != nil {
		return nil, err
	}

	// iterate over original arrays and add to groups
	mapOfArrays := make(map[string][]any)
	result := make(map[string][]byte)
	for k, _ := range s.config.Groups {
		mapOfArrays[k] = make([]any, 0)
		result[k] = make([]byte, 0)
	}

	if !s.allArrayElementsAreValues(dataToSplit) {
		for _, a := range dataToSplit {
			for k, conditions := range s.config.Groups {
				if conditions.Apply(s.logger, a.(map[string]any)) {
					mapOfArrays[k] = append(mapOfArrays[k], a)
				}
			}
		}
	}

	for k, data := range mapOfArrays {
		d, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		result[k] = d
	}

	return mapOfArrays, nil
}

// extractDataToSplit extracts and validates the data to be split based on the specified arrayPath or uses the whole input.
func (s *Splitter) extractDataToSplit(input any) ([]any, error) {
	var dataToSplit any
	if s.arrayPath != "" {
		hierarchy := strings.Split(s.arrayPath, ".")
		var current any
		current = input
		propertyHierarchy := ""
		var value any
		for _, h := range hierarchy {
			var ok bool
			if value, ok = current.(map[string]any); !ok {
				return nil, fmt.Errorf("no json object found at [%s]", propertyHierarchy)
			}
			if propertyHierarchy == "" {
				propertyHierarchy = h
			} else {
				propertyHierarchy = propertyHierarchy + "." + h
			}
			v := value.(map[string]any)
			if d, ok := v[h]; !ok {
				return nil, errors.New("key not found: " + h + "")
			} else {
				value = d
			}
		}
		dataToSplit = value
	} else {
		dataToSplit = input
	}

	isArray := s.isJSONArray(dataToSplit)
	if !isArray {
		return nil, errors.New("input is not an array")
	}
	return dataToSplit.([]any), nil
}

// SplitIo reads data from the provided io.Reader and splits it into groups based on the configuration.
func (s *Splitter) SplitIo(in io.Reader) (map[string][]any, error) {
	data, err := io.ReadAll(in)
	if err != nil {
		return nil, err
	}
	return s.Split(data)
}
