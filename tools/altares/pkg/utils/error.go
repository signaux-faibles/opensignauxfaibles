package utils

import (
	"context"
	"log/slog"
)

func ManageError(err error, message string, args ...slog.Attr) {
	if err == nil {
		return
	}
	allArgs := []slog.Attr{slog.Any("error", err)}
	for _, arg := range args {
		allArgs = append(allArgs, arg)
	}
	slog.LogAttrs(context.Background(), slog.LevelError, message, allArgs...)
	panic(err)
}
