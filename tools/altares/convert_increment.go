package main

import (
	"encoding/csv"
	"io"
	"log/slog"
	"os"

	"github.com/pkg/errors"
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

var FORMATS = map[int]func(string) string{
	30: datify,
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

	reader := csv.NewReader(inputFile)
	reader.Comma = ';'

	w := csv.NewWriter(os.Stdout)
	defer w.Flush()
	w.Comma = reader.Comma

	// discard headers
	_, err = reader.Read()
	if err != nil {
		panic(err)
	}

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
		convertor := FORMATS[field]
		if convertor != nil {
			record[field] = convertor(record[field])
		}
		data = append(data, record[field])

	}
	return data
}

func datify(s string) string {
	return s[8:10] + "/" + s[5:7] + "/" + s[0:4]
}
