package main

import (
	"encoding/csv"
	"io"
	"log/slog"
	"os"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/pkg/errors"
)

var loglevel *slog.LevelVar

const HEADER_TO_REMOVE = "NBR_EXPERIENCES_PAIEMENT"

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

	reader := csv.NewReader(inputFile)
	reader.Comma = ';'

	w := csv.NewWriter(os.Stdout)
	defer w.Flush()
	w.Comma = reader.Comma

	slog.Info("manipule les headers")
	columnToRemove := manageHeaders(reader, w)
	slog.Info("suppression d'une colonne", slog.String("status", "start"), slog.Int("index", columnToRemove))

	for {
		record, err := reader.Read()
		if err, ok := err.(*csv.ParseError); ok && err.Err != csv.ErrFieldCount {
			slog.Warn("erreur lors de la lecture", slog.Any("error", err))
			continue
		}
		if err == io.EOF {
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

func manageHeaders(reader *csv.Reader, w *csv.Writer) int {
	columnToRemove := -1
	headers, err := reader.Read()
	if err != nil {
		panic(err)
	}
	for idx, header := range headers {
		if header == HEADER_TO_REMOVE {
			columnToRemove = idx
		}
		headers[idx] = formatHeader(header)
	}
	err = w.Write(removeColumn(headers, columnToRemove))
	if err != nil {
		panic(err)
	}
	slog.Info("colonne à supprimer", slog.Int("idx", columnToRemove))
	return columnToRemove
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
	r = strings.TrimSpace(r)
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
