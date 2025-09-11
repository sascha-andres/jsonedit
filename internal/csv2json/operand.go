package csv2json

import (
	"log/slog"
	"strconv"
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
		if recordInfo.HeaderIndex != nil {
			if idx, ok := recordInfo.HeaderIndex[op.Value]; ok {
				if idx < len(recordInfo.Record) {
					c, err := convertToType("bool", recordInfo.Record[idx])
					if err != nil {
						if logger != nil {
							logger.Error("error converting value to bool", "err", err, "value", recordInfo.Record[idx])
						}
						return false
					}
					value = c.(bool)
				} else if logger != nil {
					logger.Error("index out of range", "index", idx, "length", len(recordInfo.Record))
				}
			} else if logger != nil {
				logger.Error("header not found", "name", op.Value)
			}
		} else {
			i, err := strconv.Atoi(op.Value)
			if err != nil {
				if logger != nil {
					logger.Error("error converting value to int for index", "err", err, "value", op.Value)
				}
			}
			if i < len(recordInfo.Record) {
				c, err := convertToType("bool", recordInfo.Record[i])
				if err != nil {
					if logger != nil {
						logger.Error("error converting value to bool", "err", err, "value", recordInfo.Record[i])
					}
					return false
				}
				value = c.(bool)
			} else {
				if logger != nil {
					logger.Error("index out of range", "index", i, "length", len(recordInfo.Record))
				}
			}
		}
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
				logger.Error("error converting value to bool", "err", err, "value", op.Value)
			}
			return 0.0
		}
		value = c.(float64)
	}
	if op.Type == "column" {
		if recordInfo.HeaderIndex != nil {
			if idx, ok := recordInfo.HeaderIndex[op.Value]; ok {
				if idx < len(recordInfo.Record) {
					c, err := convertToType("float", recordInfo.Record[idx])
					if err != nil {
						if logger != nil {
							logger.Error("error converting value to float64", "err", err, "value", recordInfo.Record[idx])
						}
						return 0.0
					}
					value = c.(float64)
				} else if logger != nil {
					logger.Error("index out of range", "index", idx, "length", len(recordInfo.Record))
				}
			} else if logger != nil {
				logger.Error("header not found", "name", op.Value)
			}
		} else {
			i, err := strconv.Atoi(op.Value)
			if err != nil {
				if logger != nil {
					logger.Error("error converting value to int for index", "err", err, "value", op.Value)
				}
			}
			if i < len(recordInfo.Record) {
				c, err := convertToType("float", recordInfo.Record[i])
				if err != nil {
					if logger != nil {
						logger.Error("error converting value to float64", "err", err, "value", recordInfo.Record[i])
					}
					return 0.0
				}
				value = c.(float64)
			} else {
				if logger != nil {
					logger.Error("index out of range", "index", i, "length", len(recordInfo.Record))
				}
			}
		}
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
		if recordInfo.HeaderIndex != nil {
			if idx, ok := recordInfo.HeaderIndex[op.Value]; ok {
				if idx < len(recordInfo.Record) {
					c, err := convertToType("int", recordInfo.Record[idx])
					if err != nil {
						if logger != nil {
							logger.Error("error converting value to int", "err", err, "value", recordInfo.Record[idx])
						}
						return 0
					}
					value = c.(int)
				} else if logger != nil {
					logger.Error("index out of range", "index", idx, "length", len(recordInfo.Record))
				}
			} else if logger != nil {
				logger.Error("header not found", "name", op.Value)
			}
		} else {
			i, err := strconv.Atoi(op.Value)
			if err != nil {
				if logger != nil {
					logger.Error("error converting value to int for index", "err", err, "value", op.Value)
				}
			}
			if i < len(recordInfo.Record) {
				c, err := convertToType("int", recordInfo.Record[i])
				if err != nil {
					if logger != nil {
						logger.Error("error converting value to int", "err", err, "value", recordInfo.Record[i])
					}
					return 0
				}
				value = c.(int)
			} else {
				if logger != nil {
					logger.Error("index out of range", "index", i, "length", len(recordInfo.Record))
				}
			}
		}
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
		if recordInfo.HeaderIndex != nil {
			if idx, ok := recordInfo.HeaderIndex[op.Value]; ok {
				if idx < len(recordInfo.Record) {
					value = recordInfo.Record[idx]
				} else if logger != nil {
					logger.Error("index out of range", "index", idx, "length", len(recordInfo.Record))
				}
			} else if logger != nil {
				logger.Error("header not found", "name", op.Value)
			}
		} else {
			i, err := strconv.Atoi(op.Value)
			if err != nil {
				if logger != nil {
					logger.Error("error converting value to int for index", "err", err, "value", op.Value)
				}
			}
			if i < len(recordInfo.Record) {
				value = recordInfo.Record[i]
			} else {
				if logger != nil {
					logger.Error("index out of range", "index", i, "length", len(recordInfo.Record))
				}
			}
		}
	}
	return value
}
