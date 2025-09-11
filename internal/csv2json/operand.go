package csv2json

import (
	"log/slog"
	"strconv"
)

// getBoolValueForApplies evaluates and retrieves a boolean value for a condition based on the provided record and header data.
func (op *Operand) getBoolValueForApplies(logger *slog.Logger, named bool, record []string, header []string, headerIndex map[string]int) bool {
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
		if named {
			if idx, ok := headerIndex[op.Value]; ok {
				if idx < len(record) {
					c, err := convertToType("bool", record[idx])
					if err != nil {
						if logger != nil {
							logger.Error("error converting value to bool", "err", err, "value", record[idx])
						}
						return false
					}
					value = c.(bool)
				} else if logger != nil {
					logger.Error("index out of range", "index", idx, "length", len(record))
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
			if i < len(record) {
				c, err := convertToType("bool", record[i])
				if err != nil {
					if logger != nil {
						logger.Error("error converting value to bool", "err", err, "value", record[i])
					}
					return false
				}
				value = c.(bool)
			} else {
				if logger != nil {
					logger.Error("index out of range", "index", i, "length", len(record))
				}
			}
		}
	}
	return value
}

// getFloatValueForApplies retrieves a float value based on the operand configuration and the provided record and header data.
func (op *Operand) getFloatValueForApplies(logger *slog.Logger, named bool, record []string, header []string, headerIndex map[string]int) float64 {
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
		if named {
			if idx, ok := headerIndex[op.Value]; ok {
				if idx < len(record) {
					c, err := convertToType("float", record[idx])
					if err != nil {
						if logger != nil {
							logger.Error("error converting value to float64", "err", err, "value", record[idx])
						}
						return 0.0
					}
					value = c.(float64)
				} else if logger != nil {
					logger.Error("index out of range", "index", idx, "length", len(record))
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
			if i < len(record) {
				c, err := convertToType("float", record[i])
				if err != nil {
					if logger != nil {
						logger.Error("error converting value to float64", "err", err, "value", record[i])
					}
					return 0.0
				}
				value = c.(float64)
			} else {
				if logger != nil {
					logger.Error("index out of range", "index", i, "length", len(record))
				}
			}
		}
	}
	return value
}

// getIntValueForApplies retrieves an integer value based on the operand configuration and the provided record and header data.
func (op *Operand) getIntValueForApplies(logger *slog.Logger, named bool, record []string, header []string, headerIndex map[string]int) int {
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
		if named {
			if idx, ok := headerIndex[op.Value]; ok {
				if idx < len(record) {
					c, err := convertToType("int", record[idx])
					if err != nil {
						if logger != nil {
							logger.Error("error converting value to int", "err", err, "value", record[idx])
						}
						return 0
					}
					value = c.(int)
				} else if logger != nil {
					logger.Error("index out of range", "index", idx, "length", len(record))
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
			if i < len(record) {
				c, err := convertToType("int", record[i])
				if err != nil {
					if logger != nil {
						logger.Error("error converting value to int", "err", err, "value", record[i])
					}
					return 0
				}
				value = c.(int)
			} else {
				if logger != nil {
					logger.Error("index out of range", "index", i, "length", len(record))
				}
			}
		}
	}
	return value
}

// getStringValueForApplies retrieves a string value based on the operand configuration and the provided record and header data.
func (op *Operand) getStringValueForApplies(logger *slog.Logger, named bool, record []string, header []string, headerIndex map[string]int) string {
	value := ""
	if op.Type == "value" {
		value = op.Value
	}
	if op.Type == "column" {
		if named {
			if idx, ok := headerIndex[op.Value]; ok {
				if idx < len(record) {
					value = record[idx]
				} else if logger != nil {
					logger.Error("index out of range", "index", idx, "length", len(record))
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
			if i < len(record) {
				value = record[i]
			} else {
				if logger != nil {
					logger.Error("index out of range", "index", i, "length", len(record))
				}
			}
		}
	}
	return value
}
