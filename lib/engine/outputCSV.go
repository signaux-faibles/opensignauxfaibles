package engine

import (
	"encoding/csv"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/spf13/viper"

	"opensignauxfaibles/lib/marshal"
)

const DefaultExportPath = "/export/csv"

// CSVOutputStreamer writes the output to csv files. It implements `OutputStreamer`
// If writer is nil, it will stream into csv files in the "relativeDirPath"
// directory.
// Otherwise, it will stream to the io.Writer.
type CSVOutputStreamer struct {
	relativeDirPath string
	parserType      string
	writer          io.Writer
}

// NewCSVOutputStreamer creates a streamer that will write CSV files into a
// given directory.
func NewCSVOutputStreamer(relativeDirPath, parserType string) OutputStreamer {
	out := CSVOutputStreamer{relativeDirPath, parserType, nil}
	return out
}

// Stream creates and writes all incoming data to a csv file.
// If the file already exists it will be overwritten
//
// All incoming tuples are expected to be of same type.
//
// The filename is derived from the tuple type with extension ".csv".
// The directory is derived from the CSVOutputStreamer's directory
// path, relative to the export root directory ("export.path"
// configuration, or by default the `DefaultExportPath` constant)
func (out CSVOutputStreamer) Stream(ch chan marshal.Tuple) error {
	logger := slog.With("parser", out.parserType, "streamer", "csv")
	logger.Debug("stream data to CSV file")

	var w *csv.Writer

	if out.writer != nil {
		logger.Debug("a writer has been provided with the CSVOutputStreamer")
		w = csv.NewWriter(out.writer)
	} else {

		filePath := resolveFilePath(out.relativeDirPath, out.parserType)
		logger.Debug("set up file writer, create file and directory if needed", "output_file", filePath)

		file, err := createFile(filePath)

		if err != nil {
			return fmt.Errorf("an error occurred while creating an output CSV file: %v", err)
		}
		defer file.Close()

		w = csv.NewWriter(file)
	}

	logger.Debug("data writing")

	nWritten := 0

	headersWritten := false
	for tuple := range ch { // remaining data
		if !headersWritten {
			w.Write(tuple.Headers())
			headersWritten = true
		}
		w.Write(tuple.Values())
		nWritten++
	}
	w.Flush()

	logger.Debug("output streaming to CSV file ended successfully", "n_written", nWritten)

	return nil
}

// resolveFilePath returns the file path for a given tuple type and batch key
func resolveFilePath(relativePath string, filetype string) string {
	rootDir := DefaultExportPath
	if viper.IsSet("export.path") {
		rootDir = viper.GetString("export.path")
	}

	filename := filetype + ".csv"

	exportPath := filepath.Join(rootDir, relativePath)
	filename = filepath.Join(exportPath, filename)

	slog.Debug(
		"Résolution du nom de fichier d'export",
		slog.String("key", relativePath),
		slog.String("type", filetype),
		slog.String("filename", filename),
	)

	return filename
}

// createFile creates all necessary directories and the file given at path
func createFile(filePath string) (*os.File, error) {
	dir := filepath.Dir(filePath)

	err := os.MkdirAll(dir, 0755) // No error if already exists
	if err != nil {
		return nil, err
	}

	file, err := os.Create(filePath) // this will truncate if it already exists
	if err != nil {
		return nil, err
	}

	return file, nil
}
