package csv2json

import (
	"log/slog"
)

// getBoolValueForApplies evaluates and retrieves a boolean value for a condition based on the provided record and header data.
func (op *Operand) getBoolValueForApplies(logger *slog.Logger, recordInfo *RecordWithInformation) bool {
	value := false
	if op.Type == "value" {
		c, err := convertToType("bool", op.Value)
		if err != nil {
			if logger != nil {
				logger.Error("error converting value to bool", "err", err, "value", op.Value)
			}
			return false
		}
		value = c.(bool)
	}
	if op.Type == "column" {
		val, err := recordInfo.GetValue(logger, "bool", op.Value)
		if err != nil {
			if logger != nil {
				logger.Error("error converting value to bool", "err", err)
			}
			return false
		}
		value = val.(bool)
	}
	return value
}

// getFloatValueForApplies retrieves a float value based on the operand configuration and the provided record and header data.
func (op *Operand) getFloatValueForApplies(logger *slog.Logger, recordInfo *RecordWithInformation) float64 {
	value := 0.0
	if op.Type == "value" {
		c, err := convertToType("float", op.Value)
		if err != nil {
			if logger != nil {
				logger.Error("error converting value to float64", "err", err, "value", op.Value)
			}
			return 0.0
		}
		value = c.(float64)
	}
	if op.Type == "column" {
		val, err := recordInfo.GetValue(logger, "float", op.Value)
		if err != nil {
			if logger != nil {
				logger.Error("error converting value to float", "err", err)
			}
			return 0
		}
		value = val.(float64)
	}
	return value
}

// getIntValueForApplies retrieves an integer value based on the operand configuration and the provided record and header data.
func (op *Operand) getIntValueForApplies(logger *slog.Logger, recordInfo *RecordWithInformation) int {
	value := 0
	if op.Type == "value" {
		c, err := convertToType("int", op.Value)
		if err != nil {
			if logger != nil {
				logger.Error("error converting value to int", "err", err, "value", op.Value)
			}
			return 0
		}
		value = c.(int)
	}
	if op.Type == "column" {
		val, err := recordInfo.GetValue(logger, "int", op.Value)
		if err != nil {
			if logger != nil {
				logger.Error("error converting value to int", "err", err)
			}
			return 0
		}
		value = val.(int)
	}
	return value
}

// getStringValueForApplies retrieves a string value based on the operand configuration and the provided record and header data.
func (op *Operand) getStringValueForApplies(logger *slog.Logger, recordInfo *RecordWithInformation) string {
	value := ""
	if op.Type == "value" {
		value = op.Value
	}
	if op.Type == "column" {
		val, err := recordInfo.GetValue(logger, "string", op.Value)
		if err != nil {
			if logger != nil {
				logger.Error("error converting value to string", "err", err)
			}
			return ""
		}
		value = val.(string)
	}
	return value
}
