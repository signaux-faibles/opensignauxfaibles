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

// An OutputHandler directs a stream of output data to the desired sink
type OutputHandler interface {
	Stream(ch chan marshal.Tuple) error
}

// CSVOutputHandler writes the output to CSVs. Implements OutputHandler
type CSVOutputHandler struct {
	directory  string
	chanToCSVs chan *Data
}

func NewOutputHandler(batchKey BatchKey) *CSVOutputHandler {
	ch := InsertIntoCSVs()
	out := CSVOutputHandler{batchKey, ch}
	return &out
}

func (out *CSVOutputHandler) Stream(ch chan marshal.Tuple) error {
	for tuple := range ch {

		value := Data{
			Scope: tuple.Scope(),
			Key:   tuple.Key(),
			Batch: map[BatchKey]Batch{
				out.directory: {
					tuple.Type(): tuple,
				},
			},
		}
		out.chanToCSVs <- &value
	}
	return nil
}

func (out *CSVOutputHandler) Close() {
	close(out.chanToCSVs)
	closeCSVs()
}

func InsertIntoCSVs() chan *Data {
	data := make(chan *Data, 10)

	go func() {
		for v := range data {
			writeBatchesToCSV(v.Batch)
		}
	}()

	return data
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

func writeBatchesToCSV(batchs map[BatchKey]Batch) {
	for key, batch := range batchs {
		for _, tuples := range batch {
			writeLinesToCSV(key, tuples)
		}
	}
}

func writeLinesToCSV(key BatchKey, tuples map[string]marshal.Tuple) {
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
	file, found := csvFiles[tuple.Type()]
	if found {
		return csv.NewWriter(file)
	}
	var err error
	fullFilename := prepareFilename(key, tuple.Type())

	logger := slog.Default().With(slog.String("filename", fullFilename))
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
		"écrit les headers",
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

func prepareFilename(key string, filetype string) string {
	rootPath := exportPath
	filename := filetype + ".csv"
	if viper.IsSet("export.path") {
		rootPath = viper.GetString("export.path")
	}
	exportPath := filepath.Join(rootPath, key)
	createExportFolder(exportPath)
	filename = filepath.Join(exportPath, filename)
	slog.Debug(
		"le nom de fichier est généré",
		slog.String("key", key),
		slog.String("type", filetype),
		slog.String("filename", filename),
	)
	return filename
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
