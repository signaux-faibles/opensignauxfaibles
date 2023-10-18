package engine

import (
	"encoding/csv"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/spf13/viper"

	"opensignauxfaibles/lib/marshal"
)

const exportPath = "/export/csv"

var csvFiles = map[string]*os.File{}

func InsertIntoCSVs() chan *Value {
	importing.Add(1)
	source := make(chan *Value, 10)
	go func() {
		defer importing.Done()
		for v := range source {
			writeBatchesToCSV(v.Value.Batch)
		}
	}()
	return source
}

// FlushImportedData finalise l'insertion des données dans ImportedData.
func FlushImportedData(channel chan *Value) {
	closeCSVs()
	importing.Wait()
}

func closeCSVs() {
	for _, file := range csvFiles {
		err := file.Close()
		if err != nil {
			slog.Error(
				"erreur pendant la fermeture du fichier",
				slog.Any("error", err),
				slog.String("filename", file.Name()),
			)
		}
	}
}

func writeBatchesToCSV(batchs map[string]Batch) {
	for k, v := range batchs {
		writeBatchToCSV(k, v)
	}
}

func writeBatchToCSV(key string, batch Batch) {
	for _, tuples := range batch {
		writeLinesToCSV(key, tuples)
	}
}

func writeLinesToCSV(key string, tuples map[string]marshal.Tuple) {
	for _, tuple := range tuples {
		logger := slog.Default().With(slog.Any("tuple", tuple))
		csvWriter := openFile(key, tuple)
		err := csvWriter.Write(tuple.Values())
		if err != nil {
			logger.Error("erreur pendant l'écriture du tuple en csvWriter")
		}
		csvWriter.Flush()
	}
}

func openFile(key string, tuple marshal.Tuple) *csv.Writer {
	logger := slog.Default().With(slog.Any("tuple", tuple))
	file, found := csvFiles[tuple.Type()]
	if found {
		return csv.NewWriter(file)
	}
	var err error
	fullFilename := prepareFilename(key, tuple.Type())

	logger = logger.With(slog.String("filename", fullFilename))
	file, err = os.OpenFile(fullFilename, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0644)
	csvFiles[tuple.Type()] = file
	if err != nil {
		logger.Error(
			"erreur pendant l'ouverture du fichier",
			slog.Any("error", err),
		)
		panic(err)
	}
	writer := csv.NewWriter(file)
	headers := tuple.Headers()
	logger.Info(
		"write headers",
		slog.Any("headers", headers),
	)
	err = writer.Write(headers)
	if err != nil {
		logger.Error(
			"erreur pendant l'écriture des headers",
			slog.Any("error", err),
		)
	}
	writer.Flush()
	return writer
}

func prepareFilename(key string, s string) string {
	rootPath := exportPath
	filename := string(s) + ".csv"
	if viper.IsSet("export.path") {
		rootPath = viper.GetString("export.path")
	}
	exportPath := filepath.Join(rootPath, key)
	createExportFolder(exportPath)
	return filepath.Join(exportPath, filename)
}

func createExportFolder(path string) {
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		slog.Error(
			"erreur pendant la création du répertoire d'export",
			slog.String("path", path),
			slog.Any("error", err),
		)
		panic(errors.Wrap(err, "erreur pendant la création du répertoire d'export"))
	}
}
