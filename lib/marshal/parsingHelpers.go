package marshal

import (
	"bufio"
	"compress/gzip"
	"encoding/csv"
	"io"
	"os"
	"strings"
)

// OpenCsvReader ouvre un fichier CSV potentiellement gzippé et retourne un csv.Reader.
func OpenCsvReader(filePath string, comma rune, lazyQuotes bool) (*os.File, *csv.Reader, error) {
	file, fileReader, err := OpenFileReader(filePath)
	if err != nil {
		return file, nil, err
	}
	reader := csv.NewReader(fileReader)
	reader.Comma = comma
	reader.LazyQuotes = lazyQuotes
	return file, reader, err
}

// OpenFileReader ouvre un fichier potentiellement gzippé et retourne un io.Reader.
// Un fichier gzippé est caractérisé par une des propriétés suivantes:
// - une extension ".gz"
// - ou la présence du préfixe "gzip:"
func OpenFileReader(filePath string) (*os.File, io.Reader, error) {
	isCompressed := strings.HasSuffix(filePath, ".gz")
	if strings.HasPrefix(filePath, "gzip:") {
		isCompressed = true
		filePath = strings.Replace(filePath, "gzip:", "", 1)
	}
	file, err := os.Open(filePath)
	if err != nil {
		return nil, nil, err
	}
	var fileReader io.Reader
	if isCompressed {
		fileReader, err = gzip.NewReader(file)
		if err != nil {
			return file, nil, err
		}
	} else {
		fileReader = bufio.NewReader(file)
	}
	return file, fileReader, err
}
