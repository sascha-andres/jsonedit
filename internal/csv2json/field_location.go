package csv2json

// IsValid checks if the FieldLocation has a valid value, returning true for RecordLocation or DocumentLocation, false otherwise.
func (fl *FieldLocation) IsValid() bool {
	switch *fl {
	case RecordLocation:
		return true
	case DocumentLocation:
		return true
	}
	return false
}
