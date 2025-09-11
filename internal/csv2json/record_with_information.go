package csv2json

import (
	"errors"
	"log/slog"
	"strconv"
)

// GetValue retrieves a value from the record based on the specified type and column, converting it to the target type t.
// Logs errors for invalid headers or indices using the provided logger.
// Returns the converted value or an error if the operation fails.
func (rwi *RecordWithInformation) GetValue(logger *slog.Logger, t, column string) (any, error) {
	var idx int
	var err error
	if rwi.HeaderIndex != nil {
		if _, ok := rwi.HeaderIndex[column]; !ok {
			if logger != nil {
				logger.Error("header not found", "name", column)
			}
			return nil, errors.New("header not found: " + column + "")
		}
		return convertToType(t, rwi.Record[rwi.HeaderIndex[column]])
	} else {
		idx, err = strconv.Atoi(column)
		if err != nil {
			return nil, err
		}
	}
	if idx >= len(rwi.Record) {
		if logger != nil {
			logger.Error("index out of range", "index", idx, "length", len(rwi.Record))
		}
		return nil, errors.New("index out of range: " + column + "")
	}
	return convertToType(t, rwi.Record[idx])
}
