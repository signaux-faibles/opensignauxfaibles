package main

import (
	"log/slog"
	"os"
	"runtime/debug"
)

var loglevel *slog.LevelVar

func LogLevel(newLevel slog.LevelVar) {
	loglevel = &newLevel
}

func init() {
	loglevel = new(slog.LevelVar)
	loglevel.Set(slog.LevelDebug)
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     loglevel,
		AddSource: true,
	})

	parentLogger := slog.New(
		handler)
	buildInfo, _ := debug.ReadBuildInfo()
	sha1 := GitCommit
	appLogger := parentLogger.With(
		slog.Group("app", slog.String("sha1", sha1)),
	)
	slog.SetDefault(appLogger)

	slog.Info(
		"initialisation",
		slog.String("go", buildInfo.GoVersion),
		slog.String("path", buildInfo.Path),
		slog.Any("any", buildInfo.Settings),
	)
}
