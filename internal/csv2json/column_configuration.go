package csv2json

import "log/slog"

// IsValid checks if the ColumnConfiguration and its associated Properties are correctly defined with valid values.
func (cc *ColumnConfiguration) IsValid(logger *slog.Logger, key string) bool {
	if len(cc.Property) != 0 && len(cc.Properties) != 0 {
		if logger != nil {
			logger.Warn("column configuration has both a property and properties defined", "key", key)
		}
		return false
	}
	if len(cc.Property) > 0 && len(cc.Type) == 0 {
		if logger != nil {
			logger.Warn("column configuration has a property defined but no type", "key", key)
		}
		return false
	}
	for _, property := range cc.Properties {
		if len(property.Property) == 0 || len(property.Type) == 0 {
			if logger != nil {
				logger.Warn("column configuration has a property defined but no type", "key", key)
			}
			return false
		}
		if property.Condition != nil {
			if property.Condition.Operand1.Type != "column" && property.Condition.Operand1.Type != "value" {
				logger.Error("Operand1.Type must be either column or value", "key", key)
			}
			if property.Condition.Operand2.Type != "column" && property.Condition.Operand2.Type != "value" {
				logger.Error("Operand2.Type must be either column or value", "key", key)
			}
			if property.Condition.Operand1.Type == "column" && property.Condition.Operand1.Value == "" {
				logger.Error("Operand1 has type column but no potential column given", "key", key)
			}
			if property.Condition.Operand2.Type == "column" && property.Condition.Operand2.Value == "" {
				logger.Error("Operand2 has type column but no potential column given", "key", key)
			}
			if property.Condition.Operator != "=" && property.Condition.Operator != "!=" && property.Condition.Operator != ">" && property.Condition.Operator != "<" {
				if logger != nil {
					logger.Warn("column configuration has a property defined but no type", "key", key)
				}
			}
		}
	}
	return true
}
