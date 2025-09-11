package csv2json

import (
	"log/slog"
)

// Apply checks whether all conditions apply to the specified record and header.
func (cs *Conditions) Apply(logger *slog.Logger, property string, named bool, record, header []string) bool {
	if cs == nil || len(*cs) == 0 {
		return false
	}
	for _, condition := range *cs {
		if !condition.Applies(logger, property, named, record, header) {
			return false
		}
	}
	return true
}

// Applies checks whether the condition applies to the specified record and header.
func (c *Condition) Applies(logger *slog.Logger, property string, named bool, record, header []string) bool {
	switch c.Type {
	case "int":
		return c.intApplies(logger, named, record, header)
	case "float":
		return c.floatApplies(logger, named, record, header)
	case "bool":
		return c.boolApplies(logger, named, record, header)
	case "string":
		return c.stringApplies(logger, named, record, header)
	default:
		if logger != nil {
			logger.Warn("property configuration has a property defined but no (supported) type, condition check not possible", "key", property)
		}
		return false
	}
}

// intApplies evaluates an integer condition based on the specified operator and operand values extracted from the record and header.
func (c *Condition) intApplies(logger *slog.Logger, named bool, record, header []string) bool {
	op1 := c.Operand1.getIntValueForApplies(logger, named, record, header)
	op2 := c.Operand2.getIntValueForApplies(logger, named, record, header)
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
func (c *Condition) floatApplies(logger *slog.Logger, named bool, record []string, header []string) bool {
	// TODO floating point comparison needs a bit more uncertainty
	op1 := c.Operand1.getFloatValueForApplies(logger, named, record, header)
	op2 := c.Operand2.getFloatValueForApplies(logger, named, record, header)
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
func (c *Condition) stringApplies(logger *slog.Logger, named bool, record []string, header []string) bool {
	op1 := c.Operand1.getStringValueForApplies(logger, named, record, header)
	op2 := c.Operand2.getStringValueForApplies(logger, named, record, header)
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
func (c *Condition) boolApplies(logger *slog.Logger, named bool, record []string, header []string) bool {
	if c.Operator == ">" || c.Operator == "<" {
		logger.Error("boolApplies not supported for operators >, <, will always return false")
		return false
	}
	value := c.Operand1.getBoolValueForApplies(logger, named, record, header)
	if c.Operator == "=" {
		return value
	}
	if c.Operator == "!=" {
		return !value
	}
	return false
}
