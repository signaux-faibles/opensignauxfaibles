package main

import (
	"errors"
	"log/slog"
	"os"
	"runtime/debug"
	"strings"

	"github.com/spf13/viper"
)

var loglevel *slog.LevelVar

func initLogger() {
	loglevel = new(slog.LevelVar)
	loglevel.Set(slog.LevelInfo)
	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: loglevel,
		//AddSource: true,
	})

	parentLogger := slog.New(handler)
	buildInfo, _ := debug.ReadBuildInfo()
	sha1 := GitCommit
	appLogger := parentLogger.With(
		slog.Group("app", slog.String("sha1", sha1)),
	)
	slog.SetDefault(appLogger)

	level, err := parseLogLevel(viper.GetString("log.level"))
	if err != nil {
		slog.Warn("erreur de log level", slog.Any("error", err))
	}
	loglevel.Set(level)

	slog.Info(
		"initialisation",
		slog.String("go", buildInfo.GoVersion),
		slog.String("path", buildInfo.Path),
		slog.Any("any", buildInfo.Settings),
		slog.Any("level", loglevel),
	)
}

func parseLogLevel(logLevel string) (slog.Level, error) {
	switch strings.ToUpper(logLevel) {
	case "DEBUG":
		return slog.LevelDebug, nil
	case "INFO":
		return slog.LevelInfo, nil
	case "WARN":
		return slog.LevelWarn, nil
	case "ERROR":
		return slog.LevelError, nil
	default:
		return slog.LevelInfo, errors.New("log level inconnu : '" + logLevel + "'")
	}
}

func ConfigureLogLevel(logLevel string) {
	var level, err = parseLogLevel(logLevel)
	if err != nil {
		slog.Warn("Erreur de configuration sur le loglevel", slog.String("cause", err.Error()))
		return
	}
	loglevel.Set(level)
}
