package engine

import (
	"io"
)

// ParseLines appelle la fonction parseLine() sur chaque ligne du fichier CSV pour transmettre les tuples et/ou erreurs dans parsedLineChan.
func ParseLines(parserInst ParserInst, parsedLineChan chan ParsedLineResult) {
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
