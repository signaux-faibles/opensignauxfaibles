package utils

import (
	"io"
	"log/slog"
)

func CloseIt(o io.Closer, s string) {
	err := o.Close()
	slog.Debug(s)
	ManageError(err, "erreur "+s)
}
