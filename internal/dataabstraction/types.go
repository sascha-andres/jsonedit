package dataabstraction

import "log/slog"

type ValueGetter interface {
	GetValue(logger *slog.Logger, t, column string) (any, error)
}
