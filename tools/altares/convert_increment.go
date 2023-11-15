package main

import (
	"encoding/csv"
	"io"
	"log/slog"
	"os"

	"github.com/pkg/errors"
	"golang.org/x/text/encoding/charmap"
)

var FIELDS = []int{
	1,
	18,
	21,
	21,
	22,
	23,
	24,
	26,
	27,
	30,
}

func main() {
	inputFile, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer func() {
		closeErr := inputFile.Close()
		if closeErr != nil {
			panic(errors.Wrap(closeErr, "erreur à la fermeture du fichier"))
		}
	}()

	fromISO8859_15toUTF8 := charmap.ISO8859_15.NewDecoder()
	convertReader := fromISO8859_15toUTF8.Reader(inputFile)
	reader := csv.NewReader(convertReader)
	reader.Comma = ';'

	w := csv.NewWriter(os.Stdout)
	defer w.Flush()

	// discard headers
	headers, err := reader.Read()
	if err != nil {
		panic(err)
	}
	slog.Debug("description des headers", slog.Any("headers", headers))

	for {
		record, err := reader.Read()
		if err, ok := err.(*csv.ParseError); ok && err.Err != csv.ErrFieldCount {
			continue
		}
		if err == io.EOF {
			return
		}
		out := selectFields(record)

		if out != nil {
			err = w.Write(out)
			if err != nil {
				panic(err)
			}
		}
	}
}

func selectFields(record []string) []string {
	var data []string
	for _, field := range FIELDS {
		if field > len(record)-1 {
			slog.Error(
				"erreur de longueur de ligne",
				slog.Any("record", record),
			)
			return nil
		}
		data = append(data, record[field])
	}
	return data
}

func init() {
	loglevel := new(slog.LevelVar)
	loglevel.Set(slog.LevelDebug)
	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: loglevel,
	})

	logger := slog.New(
		handler)
	slog.SetDefault(logger)
}
