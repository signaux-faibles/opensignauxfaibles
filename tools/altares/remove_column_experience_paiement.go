package main

import (
	"encoding/csv"
	"io"
	"log/slog"
	"os"
	"strings"

	"github.com/iancoleman/strcase"
)

var loglevel *slog.LevelVar

func main() {
	slog.Debug("c'est parti")
	columnToRemove := 7

	slog.Info("suppression d'une colonne", slog.String("status", "start"), slog.Int("index", columnToRemove))

	inputFile, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	reader := csv.NewReader(inputFile)
	reader.Comma = ';'

	w := csv.NewWriter(os.Stdout)
	w.Comma = reader.Comma

	slog.Info("manipule les headers")
	headers, err := reader.Read()
	if err != nil {
		panic(err)
	}
	for idx, header := range headers {
		headers[idx] = formatHeader(header)
	}
	err = w.Write(removeColumn(headers, columnToRemove))
	if err != nil {
		panic(err)
	}

	//x	headers[idx] = strcase.ToSnakeWithIgnore(value, ".,")
	//}
	//removeColumn(headers)
	//slog.Info("transforme les headers", slog.Any("headers", headers))
	//err = w.Write(headers)
	//if err != nil {
	//	panic(err)
	//}

	slog.Info("manipule les headers")
	for {
		record, err := reader.Read()
		if err, ok := err.(*csv.ParseError); ok && err.Err != csv.ErrFieldCount {
			continue
		}
		if err == io.EOF {
			w.Flush()
			slog.Info("suppression d'une colonne", slog.String("status", "end"), slog.Int("index", columnToRemove))
			return
		}
		newRecord := removeColumn(record, columnToRemove)
		err = w.Write(newRecord)
		if err != nil {
			panic(err)
		}
	}
}

func removeColumn(record []string, remove int) []string {
	var r []string
	for idx, value := range record {
		if idx != remove {
			r = append(r, value)
		}
	}
	return r
}

func formatHeader(input string) string {
	r := strcase.ToSnakeWithIgnore(input, ".,")
	r = strings.Replace(r, "etat", "état", 1)
	r = strings.Replace(r, "etudies", "étudiés", 1)
	return r
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
