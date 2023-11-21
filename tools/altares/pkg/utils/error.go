package utils

import (
	"log/slog"
)

func ManageError(err error, message string) {
	if err == nil {
		return
	}
	slog.Error(message, slog.Any("error", err))
	panic(err)
}
