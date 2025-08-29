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
	}
	return true
}
