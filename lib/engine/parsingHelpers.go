package engine

import (
	"bufio"
	"compress/gzip"
	"io"
	"os"

	"opensignauxfaibles/lib/base"
)

// OpenFileReader ouvre un fichier potentiellement gzippé et retourne un io.Reader.
func OpenFileReader(batchFile base.BatchFile) (*os.File, io.Reader, error) {
	file, err := os.Open(batchFile.Path())
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

// ParseLines appelle la fonction parseLine() sur chaque ligne du fichier CSV pour transmettre les tuples et/ou erreurs dans parsedLineChan.
func ParseLines(parserInst ParserInstance, parsedLineChan chan ParsedLineResult) {
	defer close(parsedLineChan)

	var lineNumber = 0 // starting with the header

	stopProgressLogger := LogProgress(&lineNumber)
	defer stopProgressLogger()

	for {
		lineNumber++
		parsedLine := ParsedLineResult{}
		err := parserInst.ReadNext(&parsedLine)

		if err == io.EOF {
			break
		} else if err != nil {
			parsedLine.AddRegularError(err)
			parsedLineChan <- parsedLine
			break
		}

		parsedLineChan <- parsedLine
	}
}
