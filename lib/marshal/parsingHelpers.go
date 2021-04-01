package marshal

import (
	"bufio"
	"compress/gzip"
	"encoding/csv"
	"io"
	"os"

	"github.com/signaux-faibles/opensignauxfaibles/lib/base"
)

// OpenCsvReader ouvre un fichier CSV potentiellement gzippé et retourne un csv.Reader.
func OpenCsvReader(batchFile base.BatchFile, comma rune, lazyQuotes bool) (*os.File, *csv.Reader, error) {
	file, fileReader, err := OpenFileReader(batchFile)
	if err != nil {
		return file, nil, err
	}
	reader := csv.NewReader(fileReader)
	reader.Comma = comma
	reader.LazyQuotes = lazyQuotes
	return file, reader, err
}

// OpenFileReader ouvre un fichier potentiellement gzippé et retourne un io.Reader.
func OpenFileReader(batchFile base.BatchFile) (*os.File, io.Reader, error) {
	file, err := os.Open(batchFile.FilePath())
	if err != nil {
		return nil, nil, err
	}
	var fileReader io.Reader
	if batchFile.IsCompressed() {
		fileReader, err = gzip.NewReader(file)
		if err != nil {
			return file, nil, err
		}
	} else {
		fileReader = bufio.NewReader(file)
	}
	return file, fileReader, err
}
