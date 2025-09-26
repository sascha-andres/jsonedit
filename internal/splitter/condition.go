package splitter

import (
	"log/slog"
)

// Apply checks whether all conditions apply to the specified record and header.
func (cs *Conditions) Apply(logger *slog.Logger, property string, recordInfo *RecordWithInformation) bool {
	if cs == nil || len(*cs) == 0 {
		return false
	}
	for _, condition := range *cs {
		if !condition.Applies(logger, property, recordInfo) {
			return false
		}
	}
	return true
}

// Applies checks whether the condition applies to the specified record and header.
func (c *Condition) Applies(logger *slog.Logger, property string, recordInfo *RecordWithInformation) bool {
	switch c.Type {
	case "int":
		return c.intApplies(logger, recordInfo)
	case "float":
		return c.floatApplies(logger, recordInfo)
	case "bool":
		return c.boolApplies(logger, recordInfo)
	case "string":
		return c.stringApplies(logger, recordInfo)
	default:
		if logger != nil {
			logger.Warn("property configuration has a property defined but no (supported) type, condition check not possible", "key", property)
		}
		return false
	}
}

// intApplies evaluates an integer condition based on the specified operator and operand values extracted from the record and header.
func (c *Condition) intApplies(logger *slog.Logger, recordInfo *RecordWithInformation) bool {
	op1 := c.Operand1.getIntValueForApplies(logger, recordInfo)
	op2 := c.Operand2.getIntValueForApplies(logger, recordInfo)
	if c.Operator == "=" {
		return op1 == op2
	}
	if c.Operator == "!=" {
		return op1 != op2
	}
	if c.Operator == ">" {
		return op1 > op2
	}
	if c.Operator == "<" {
		return op1 < op2
	}
	return false
}

// floatApplies evaluates a float condition using the specified operator and extracted float values from the record and header.
func (c *Condition) floatApplies(logger *slog.Logger, recordInfo *RecordWithInformation) bool {
	// TODO floating point comparison needs a bit more uncertainty
	op1 := c.Operand1.getFloatValueForApplies(logger, recordInfo)
	op2 := c.Operand2.getFloatValueForApplies(logger, recordInfo)
	if c.Operator == "=" {
		return op1 == op2
	}
	if c.Operator == "!=" {
		return op1 != op2
	}
	if c.Operator == ">" {
		return op1 > op2
	}
	if c.Operator == "<" {
		return op1 < op2
	}
	return false
}

// stringApplies evaluates a string condition using the specified operator and extracted values from the record and header.
func (c *Condition) stringApplies(logger *slog.Logger, recordInfo *RecordWithInformation) bool {
	op1 := c.Operand1.getStringValueForApplies(logger, recordInfo)
	op2 := c.Operand2.getStringValueForApplies(logger, recordInfo)
	if c.Operator == "=" {
		return op1 == op2
	}
	if c.Operator == "!=" {
		return op1 != op2
	}
	if c.Operator == ">" {
		return op1 > op2
	}
	if c.Operator == "<" {
		return op1 < op2
	}
	return false
}

// boolApplies evaluates a boolean condition based on the specified operator and the extracted value from the record and header.
func (c *Condition) boolApplies(logger *slog.Logger, recordInfo *RecordWithInformation) bool {
	if c.Operator == ">" || c.Operator == "<" {
		logger.Error("boolApplies not supported for operators >, <, will always return false")
		return false
	}
	value := c.Operand1.getBoolValueForApplies(logger, recordInfo)
	if c.Operator == "=" {
		return value
	}
	if c.Operator == "!=" {
		return !value
	}
	return false
}
