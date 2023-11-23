package main

import (
	"io"
	"log/slog"
	"os"

	"opensignauxfaibles/tools/altares/pkg/altares"
	"opensignauxfaibles/tools/altares/pkg/utils"
)

var loglevel *slog.LevelVar

func main() {
	inputs, o := readArgs()
	output, err := os.Create(o)
	utils.ManageError(err, "erreur à la création du fichier de sortie")
	slog.Debug("fichier de sortie créé", slog.String("filename", output.Name()))
	defer utils.CloseIt(output, "fermeture du fichier de sortie : "+os.Args[1])
	convertAndConcat(inputs, output)
}

func convertAndConcat(altaresFiles []string, outputCsv io.Writer) {
	slog.Debug("démarrage de la conversion et de la concaténation", slog.Any("inputs", altaresFiles))
	altares.ConvertStock(altaresFiles[0], outputCsv)
	if len(altaresFiles) == 1 {
		slog.Info("terminé, pas de fichier incrément")
	}
	for _, filename := range altaresFiles[1:] {
		altares.ConvertIncrement(filename, outputCsv)
	}
}

func readArgs() (inputs []string, output string) {
	slog.Debug("lecture des arguments", slog.String("status", "start"), slog.Any("all", os.Args))
	if len(os.Args) <= 2 {
		slog.Warn("rien à faire, car pas de fichiers altares ou pas de fichier source")
		os.Exit(0)
	}
	output = os.Args[len(os.Args)-1]
	inputs = os.Args[1 : len(os.Args)-1]
	if len(inputs) == 0 {
		slog.Warn("rien à faire, car pas de fichiers altares")
		os.Exit(0)
	}
	slog.Debug("lecture des arguments", slog.String("status", "end"), slog.String("output", output), slog.Any("inputs", inputs))
	return inputs, output
}

func init() {
	loglevel = new(slog.LevelVar)
	loglevel.Set(slog.LevelDebug)
	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: loglevel,
	})

	logger := slog.New(
		handler)
	slog.SetDefault(logger)
}
