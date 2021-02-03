package marshal

import (
	"fmt"

	"github.com/globalsign/mgo/bson"
)

// MaxParsingErrors is the number of parsing errors to report per file.
var MaxParsingErrors = 200

// ParsingTracker permet de collecter puis rapporter des erreurs de parsing.
type ParsingTracker struct {
	// fields that are included in the report:
	currentLine      int      // Note: line 1 is the first line of data (excluding the header) read from a file
	nbSkippedLines   int      // lines skipped by the perimeter/filter or not found in "comptes" mapping
	nbRejectedLines  int      // lines that have at least one parse error
	firstParseErrors []string // capped by MaxParsingErrors, with line number rendered as string
	fatalErrors      []string // with line number rendered as string
	// private state vars:
	lastSkippedLine        int // to avoid counting 2 lines if 2 "filter" errors are added on a same line
	lastLineWithParseError int // to avoid counting 2 lines if 2 "parser" errors are added on a same line
}

// AddFatalError rapporte une erreur fatale liée au parsing d'un fichier
func (tracker *ParsingTracker) AddFatalError(err error) {
	if err == nil {
		return
	}
	tracker.fatalErrors = append(tracker.fatalErrors, fmt.Sprintf("Fatal: %v", err.Error()))
}

// AddFilterError rapporte le fait que la ligne en cours de parsing est ignorée à cause du filtre/périmètre
func (tracker *ParsingTracker) AddFilterError(err error) {
	if err == nil {
		return
	}
	if tracker.currentLine != tracker.lastSkippedLine {
		tracker.nbSkippedLines++
		tracker.lastSkippedLine = tracker.currentLine
	}
	// fmt.Fprintf(os.Stderr, "Line %d: %v\n", tracker.currentLine, err.Error()) // on ne souhaite pas conserver ces erreurs dans le rapport
}

// AddParseError rapporte une erreur rencontrée sur la ligne en cours de parsing
func (tracker *ParsingTracker) AddParseError(err error) {
	if err == nil {
		return
	}
	if tracker.currentLine != tracker.lastLineWithParseError {
		tracker.nbRejectedLines++
		tracker.lastLineWithParseError = tracker.currentLine
	}
	if len(tracker.firstParseErrors) < MaxParsingErrors {
		tracker.firstParseErrors = append(tracker.firstParseErrors,
			fmt.Sprintf("Line %d: %v", tracker.currentLine, err.Error()))
	}
}

// Next informe le Tracker qu'on va parser la ligne suivante.
func (tracker *ParsingTracker) Next() {
	tracker.currentLine++
}

// Report génère un rapport de parsing à partir des erreurs rapportées.
func (tracker *ParsingTracker) Report(batchKey string, filePath string) bson.M {
	nbParsedLines := tracker.currentLine - 1 // -1 because we started counting at line number 1
	nbValidLines := nbParsedLines - tracker.nbRejectedLines - tracker.nbSkippedLines

	report := fmt.Sprintf(
		"%s: intégration terminée, %d lignes traitées, %d erreurs fatales, %d lignes rejetées, %d lignes filtrées, %d lignes valides",
		filePath,
		nbParsedLines,
		len(tracker.fatalErrors),
		tracker.nbRejectedLines,
		tracker.nbSkippedLines,
		nbValidLines,
	)

	return bson.M{
		"batchKey":      batchKey,
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
func NewParsingTracker() ParsingTracker {
	return ParsingTracker{
		currentLine:            1,
		lastSkippedLine:        -1,
		lastLineWithParseError: -1,
		firstParseErrors:       []string{},
		fatalErrors:            []string{},
	}
}
