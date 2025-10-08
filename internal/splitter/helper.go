package splitter

import (
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/sascha-andres/jsonedit/internal/dataabstraction"
)

func getValue(logger *slog.Logger, jo map[string]any, dataType string, key string) (any, error) {
	logger.Debug("trying to get value for key", "key", key, "data type", dataType)
	hierarchy := strings.Split(key, ".")
	var v map[string]any
	v = jo
	for _, h := range hierarchy {
		var d any
		var ok bool
		if d, ok = v[h]; !ok {
			return nil, errors.New("key not found: " + key + "") // TODO display hierarch
		}
		switch d.(type) {
		case map[string]any:
			v = d.(map[string]any)
		case any:
			return dataabstraction.ConvertToType(dataType, fmt.Sprint(d))
		}
	}
	return nil, nil
}
