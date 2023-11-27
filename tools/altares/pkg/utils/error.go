package utils

import (
	"log/slog"
)

func ManageError(err error, message string, args ...slog.Attr) {
	if err == nil {
		return
	}
	group := slog.Group("args", args)
	slog.Error(message, slog.Any("error", err), group)
	panic(err)
}
