package marshal

import (
	"fmt"
	"os"

	"github.com/globalsign/mgo/bson"
)

// MaxParsingErrors is the number of parsing errors to report per file.
var MaxParsingErrors = 200

// ParsingTracker permet de collecter puis rapporter des erreurs de parsing.
type ParsingTracker struct {
	filePath               string
	batchKey               string
	currentLine            int // Note: line 1 is the first line of data (excluding the header) read from a file
	nbSkippedLines         int // lines skipped by the perimeter/filter or not found in "comptes" mapping
	nbRejectedLines        int // lines that have at least one parse error
	lastSkippedLine        int
	lastLineWithParseError int
	firstParseErrors       []string // capped by MaxParsingErrors, with line number rendered as string
	fatalErrors            []string // with line number rendered as string
}

// AddFatalError rapporte une erreur fatale liée au parsing
func (tracker *ParsingTracker) AddFatalError(err error) {
	if err == nil {
		return
	}
	tracker.fatalErrors = append(tracker.fatalErrors, fmt.Sprintf("Fatal: %v", err.Error()))
}

// AddFilterError rapporte le fait que la ligne en cours est ignorée à cause du filtre/périmètre
func (tracker *ParsingTracker) AddFilterError(err error) {
	if tracker.currentLine != tracker.lastSkippedLine {
		tracker.nbSkippedLines++
		tracker.lastSkippedLine = tracker.currentLine
	}
	fmt.Fprintf(os.Stderr, "Line %d: %v\n", tracker.currentLine, err.Error())
}

// AddParseError rapporte une erreur de parsing à la ligne en cours
func (tracker *ParsingTracker) AddParseError(err error) {
	if err == nil {
		return
	}
	if len(tracker.firstParseErrors) < MaxParsingErrors {
		tracker.firstParseErrors = append(tracker.firstParseErrors,
			fmt.Sprintf("Line %d: %v", tracker.currentLine, err.Error()))
	}
	if tracker.currentLine != tracker.lastLineWithParseError {
		tracker.nbRejectedLines++
		tracker.lastLineWithParseError = tracker.currentLine
	}
}

// Next informe le Tracker qu'on passe à la ligne suivante.
func (tracker *ParsingTracker) Next() {
	tracker.currentLine++
}

// Report génère un rapport de parsing à partir des erreurs rapportées.
func (tracker *ParsingTracker) Report(code string) bson.M {
	nbParsedLines := tracker.currentLine - 1 // -1 because we started counting at line number 1
	nbValidLines := nbParsedLines - tracker.nbRejectedLines - tracker.nbSkippedLines

	report := fmt.Sprintf(
		"%s: intégration terminée, %d lignes traitées, %d erreurs fatales, %d lignes rejetées, %d lignes filtrées, %d lignes valides",
		tracker.filePath,
		nbParsedLines,
		len(tracker.fatalErrors),
		tracker.nbRejectedLines,
		tracker.nbSkippedLines,
		nbValidLines,
	)

	return bson.M{
		"batchKey":      tracker.batchKey,
		"summary":       report,
		"linesParsed":   nbParsedLines,
		"linesValid":    nbValidLines,
		"linesSkipped":  tracker.nbSkippedLines,
		"linesRejected": tracker.nbRejectedLines,
		"isFatal":       len(tracker.fatalErrors) > 0,
		"headRejected":  tracker.firstParseErrors,
		"headFatal":     tracker.fatalErrors,
	}
}

// NewParsingTracker retourne une instance pour rapporter les erreurs de parsing.
func NewParsingTracker(batchKey string, filePath string) ParsingTracker {
	return ParsingTracker{
		filePath:               filePath,
		batchKey:               batchKey,
		currentLine:            1,
		lastSkippedLine:        -1,
		lastLineWithParseError: -1,
		firstParseErrors:       []string{},
		fatalErrors:            []string{},
	}
}
