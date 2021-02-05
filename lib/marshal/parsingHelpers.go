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
func OpenFileReader(filePath string) (*os.File, io.Reader, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, nil, err
	}
	var fileReader io.Reader
	if strings.HasSuffix(filePath, ".gz") {
		fileReader, err = gzip.NewReader(file)
		if err != nil {
			return file, nil, err
		}
	} else {
		fileReader = bufio.NewReader(file)
	}
	return file, fileReader, err
}
