package csv2json

import (
	"log/slog"
)

func (pc *PropertyConfiguration) Applies(logger *slog.Logger, named bool, record, header []string) bool {
	if pc.Condition == nil {
		return true
	}
	switch pc.Type {
	case "int":
		return pc.intApplies(logger, named, record, header)
	case "float":
		return pc.floatApplies(logger, named, record, header)
	case "bool":
		return pc.boolApplies(logger, named, record, header)
	case "string":
		return pc.stringApplies(logger, named, record, header)
	default:
		if logger != nil {
			logger.Warn("property configuration has a property defined but no (supported) type, condition check not possible", "key", pc.Property)
		}
		return false
	}
	return true
}

// intApplies evaluates an integer condition based on the specified operator and operand values extracted from the record and header.
func (pc *PropertyConfiguration) intApplies(logger *slog.Logger, named bool, record, header []string) bool {
	// TODO floating point comparison needs a bit more uncertainty
	op1 := pc.Condition.Operand1.getIntValueForApplies(logger, named, record, header)
	op2 := pc.Condition.Operand2.getIntValueForApplies(logger, named, record, header)
	if pc.Condition.Operator == "=" {
		return op1 == op2
	}
	if pc.Condition.Operator == "!=" {
		return op1 != op2
	}
	if pc.Condition.Operator == ">" {
		return op1 > op2
	}
	if pc.Condition.Operator == "<" {
		return op1 < op2
	}
	return false
}

// floatApplies evaluates a float condition using the specified operator and extracted float values from the record and header.
func (pc *PropertyConfiguration) floatApplies(logger *slog.Logger, named bool, record []string, header []string) bool {
	// TODO floating point comparison needs a bit more uncertainty
	op1 := pc.Condition.Operand1.getFloatValueForApplies(logger, named, record, header)
	op2 := pc.Condition.Operand2.getFloatValueForApplies(logger, named, record, header)
	if pc.Condition.Operator == "=" {
		return op1 == op2
	}
	if pc.Condition.Operator == "!=" {
		return op1 != op2
	}
	if pc.Condition.Operator == ">" {
		return op1 > op2
	}
	if pc.Condition.Operator == "<" {
		return op1 < op2
	}
	return false
}

// stringApplies evaluates a string condition using the specified operator and extracted values from the record and header.
func (pc *PropertyConfiguration) stringApplies(logger *slog.Logger, named bool, record []string, header []string) bool {
	op1 := pc.Condition.Operand1.getStringValueForApplies(logger, named, record, header)
	op2 := pc.Condition.Operand2.getStringValueForApplies(logger, named, record, header)
	if pc.Condition.Operator == "=" {
		return op1 == op2
	}
	if pc.Condition.Operator == "!=" {
		return op1 != op2
	}
	if pc.Condition.Operator == ">" {
		return op1 > op2
	}
	if pc.Condition.Operator == "<" {
		return op1 < op2
	}
	return false
}

// boolApplies evaluates a boolean condition based on the specified operator and the extracted value from the record and header.
func (pc *PropertyConfiguration) boolApplies(logger *slog.Logger, named bool, record []string, header []string) bool {
	if pc.Condition.Operator == ">" || pc.Condition.Operator == "<" {
		logger.Error("boolApplies not supported for operators >, <, will always return false")
		return false
	}
	value := pc.Condition.Operand1.getBoolValueForApplies(logger, named, record, header)
	if pc.Condition.Operator == "=" {
		return value
	}
	if pc.Condition.Operator == "!=" {
		return !value
	}
	return false
}
