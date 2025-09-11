package csv2json

import (
	"log/slog"
)

// Applies returns true if the property configuration applies to the given record.
func (pc *PropertyConfiguration) Applies(logger *slog.Logger, recordInfo *RecordWithInformation) bool {
	if pc.Condition == nil {
		return true
	}
	return pc.Condition.Applies(logger, pc.Property, recordInfo)
}
