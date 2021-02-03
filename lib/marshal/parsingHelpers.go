package marshal

import (
	"bufio"
	"compress/gzip"
	"io"
	"os"
	"strings"
)

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
