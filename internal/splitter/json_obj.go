package splitter

import (
	"errors"
	"log/slog"
	"strings"
)

func (jo *JsonObj) GetValue(logger *slog.Logger, dataType string, key string) (any, error) {
	logger.Debug("trying to get value for key", "key", key, "data type", dataType)
	hierarchy := strings.Split(key, ".")
	var v map[string]any
	v = *jo
	for _, h := range hierarchy {
		var d any
		var ok bool
		if d, ok = v[h]; !ok {
			return nil, errors.New("key not found: " + key + "") // TODO display hierarch
		}
		switch d.(type) {
		case map[string]any:

		}
	}
	return nil, nil
}
