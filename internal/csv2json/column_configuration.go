package csv2json

// IsValid checks if the ColumnConfiguration and its associated Properties are correctly defined with valid values.
func (cc *ColumnConfiguration) IsValid() bool {
	if len(cc.Property) != 0 && len(cc.Properties) != 0 {
		return false
	}
	if len(cc.Property) > 0 && len(cc.Type) == 0 {
		return false
	}
	for _, property := range cc.Properties {
		if len(property.Property) == 0 || len(property.Type) == 0 {
			return false
		}
	}
	return true
}
