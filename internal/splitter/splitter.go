package splitter

import (
	"log/slog"
)

// OptionFunc defines a function type used to configure a Splitter by modifying its internal properties.
type OptionFunc func(*Splitter) error

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
