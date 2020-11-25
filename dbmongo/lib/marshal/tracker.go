package marshal

import (
	"fmt"
	"os"

	"github.com/globalsign/mgo/bson"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/base"
)

// MaxParsingErrors is the number of parsing errors to report per file.
var MaxParsingErrors = 200

type parseError struct {
	line int
	err  error
}

// ParsingTracker permet de collecter puis rapporter des erreurs de parsing.
type ParsingTracker struct {
	filePath       string
	batchKey       string
	currentLine    int
	nbSkippedLines int
	fatalErrors    []error
	parseErrors    []parseError
}

// Add rapporte une erreur de parsing à la ligne en cours.
func (tracker *ParsingTracker) Add(err base.CriticityError) {
	if err.Criticity() == "fatal" {
		tracker.fatalErrors = append(tracker.fatalErrors, err)
	} else if err.Criticity() == "filter" {
		// TODO: make sure that we never add more than 1 filter error per line
		tracker.nbSkippedLines++
		fmt.Fprintf(os.Stderr, "Line %d: %v\n", tracker.currentLine, err.Error())
	} else {
		tracker.parseErrors = append(tracker.parseErrors, parseError{
			line: tracker.currentLine,
			err:  err,
		})
	}
}

// Next informe le Tracker qu'on passe à la ligne suivante.
func (tracker *ParsingTracker) Next() {
	tracker.currentLine++
}

// Report génère un rapport de parsing à partir des erreurs rapportées.
func (tracker *ParsingTracker) Report(code string) interface{} {
	var nbRejectedLines = 0

	var headFatal = []string{}
	for _, err := range tracker.fatalErrors {
		if len(headFatal) < MaxParsingErrors {
			rendered := fmt.Sprintf("Line %d: %v", tracker.currentLine, err.Error())
			headFatal = append(headFatal, rendered)
		}
	}

	var headRejected = []string{}
	var lastLineWithError = -1
	for _, err := range tracker.parseErrors {
		if err.line != lastLineWithError {
			nbRejectedLines++
			lastLineWithError = err.line
		}
		if len(headRejected) < MaxParsingErrors {
			rendered := fmt.Sprintf("Line %d: %v", tracker.currentLine, err.err.Error())
			headRejected = append(headRejected, rendered)
		}
	}

	nbParsedLines := tracker.currentLine - 1
	nbValidLines := nbParsedLines - nbRejectedLines - tracker.nbSkippedLines

	report := fmt.Sprintf(
		"%s: intégration terminée, %d lignes traitées, %d erreurs fatales, %d lignes rejetées, %d lignes filtrées, %d lignes valides",
		tracker.filePath,
		nbParsedLines,
		len(tracker.fatalErrors),
		nbRejectedLines,
		tracker.nbSkippedLines,
		nbValidLines,
	)

	return bson.M{
		"batchKey":      tracker.batchKey,
		"summary":       report,
		"linesParsed":   nbParsedLines,
		"linesValid":    nbValidLines,
		"linesSkipped":  tracker.nbSkippedLines,
		"linesRejected": nbRejectedLines,
		"isFatal":       len(tracker.fatalErrors) > 0,
		"headRejected":  headRejected,
		"headFatal":     headFatal,
	}
}

// NewParsingTracker retourne une instance pour rapporter les erreurs de parsing.
func NewParsingTracker(batchKey string, filePath string) ParsingTracker {
	return ParsingTracker{
		filePath:    filePath,
		batchKey:    batchKey,
		currentLine: 1,
	}
}
