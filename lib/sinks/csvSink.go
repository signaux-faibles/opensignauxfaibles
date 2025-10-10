package sinks

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"time"

	"github.com/spf13/viper"

	"opensignauxfaibles/lib/base"
	"opensignauxfaibles/lib/engine"
	"opensignauxfaibles/lib/parsing"
)

const DefaultExportPath = "/export/csv"

type CSVSinkFactory struct {
	relativeDirPath string
}

func NewCSVSinkFactory(relativeDirPath string) engine.SinkFactory {
	return &CSVSinkFactory{relativeDirPath}
}

func (f *CSVSinkFactory) CreateSink(parserType base.ParserType) (engine.DataSink, error) {

	// resolve filename given parserType
	rootDir := DefaultExportPath
	if viper.IsSet("export.path") {
		rootDir = viper.GetString("export.path")
	}

	exportPath := filepath.Join(rootDir, f.relativeDirPath)
	filename := filepath.Join(exportPath, string(parserType)+".csv")

	slog.Debug(
		"Résolution du nom de fichier d'export",
		slog.String("key", f.relativeDirPath),
		slog.String("type", string(parserType)),
		slog.String("filename", filename),
	)

	return &CSVSink{filename, nil}, nil
}

// CSVSink writes the output to a csv file. It implements `DataSink`
// If writer is nil, it will stream into the "file".
// Otherwise, it will stream to the io.Writer (mainly used for tests)
type CSVSink struct {
	file   string
	writer io.Writer
}

// ProcessOutput creates and writes all incoming data to a csv file.
// If the file already exists it will be overwritten
//
// All incoming tuples are expected to be of same type.
//
// The filename is derived from the tuple type with extension ".csv".
// The directory is derived from the CSVSink's directory
// path, relative to the export root directory ("export.path"
// configuration, or by default the `DefaultExportPath` constant)
func (s *CSVSink) ProcessOutput(ctx context.Context, ch chan engine.Tuple) error {
	logger := slog.With("sink", "csv", "file", s.file)
	logger.Debug("stream data to CSV file")

	var w *csv.Writer

	if s.writer != nil {
		logger.Debug("a writer has been provided with the CSVSink, it has precedence over any file provided")
		w = csv.NewWriter(s.writer)
	} else {
		logger.Debug("set up file writer, create file and directory if needed", "output_file", s.file)

		file, err := createFile(s.file)

		if err != nil {
			return fmt.Errorf("an error occurred while creating an output CSV file: %v", err)
		}
		defer file.Close()

		w = csv.NewWriter(file)
	}

	logger.Debug("data writing")

	nWritten := 0

	headersWritten := false
	for tuple := range ch {
		if !headersWritten {
			w.Write(ExtractCSVHeaders(tuple))
			headersWritten = true
		}
		w.Write(ExtractCSVRow(tuple))
		nWritten++
	}
	w.Flush()

	logger.Debug("output streaming to CSV file ended successfully", "n_written", nWritten)
	return nil
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

// ExtractCSVHeaders extrait les en-têtes csv via le tag "csv"
func ExtractCSVHeaders(tuple engine.Tuple) (header []string) {
	return parsing.ExtractFieldsByTags(tuple, "csv")
}

// ExtractCSVRow returns the tuple values, in same order as the header, and converted to strings
func ExtractCSVRow(tuple engine.Tuple) (values []string) {
	rawValues := parsing.ExtractValuesByTags(tuple, "csv")
	for _, v := range rawValues {
		values = append(values, valueToString(v))
	}
	return values
}

func valueToString(v reflect.Value) string {
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return ""
	}

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.String:
		return v.String()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(v.Uint(), 10)
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(v.Float(), 'f', -1, 64)
	case reflect.Bool:
		return strconv.FormatBool(v.Bool())
	case reflect.Struct:
		if v.Type() == reflect.TypeOf(time.Time{}) {
			t := v.Interface().(time.Time)
			return t.Format(time.DateOnly)
		}
		return ""
	default:
		// Fallback to string representation
		return v.String()
	}
}
