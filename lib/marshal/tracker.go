package marshal

import (
	"fmt"
	"time"
)

// MaxParsingErrors is the number of parsing errors to report per file.
var MaxParsingErrors = 200

// ParsingTracker permet de collecter puis rapporter des erreurs de parsing.
type ParsingTracker struct {
	// fields that are included in the report:
	startDate        time.Time // when the parsingTracker was created - by extension when the process has started
	nbSkippedLines   int64     // lines skipped by the perimeter/filter or not found in "comptes" mapping
	nbRejectedLines  int64     // lines that have at least one parse error
	firstParseErrors []string  // capped by MaxParsingErrors, with line number rendered as string
	fatalErrors      []string  // with line number rendered as string
	// private state vars:
	currentLine            int64 // Note: line 1 is the first line of data (excluding the header) read from a file
	lastSkippedLine        int64 // to avoid counting 2 lines if 2 "filter" errors are added on a same line
	lastLineWithParseError int64 // to avoid counting 2 lines if 2 "parser" errors are added on a same line
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

var gitCommit string

// SetGitCommit spécifie la valeur à stocker dans le CommitHash de chaque
// rapport.
func SetGitCommit(hash string) {
	gitCommit = hash
}

type Report struct {
	Commit        string
	StartDate     time.Time
	Parser        string
	BatchKey      string   `json:"batch_key"`
	HeadFatal     []string `json:"head_fatal"`
	HeadRejected  []string `json:"head_rejected"`
	IsFatal       bool     `json:"is_fatal"`
	LinesParsed   int64    `json:"lines_parsed"`
	LinesRejected int64    `json:"lines_rejected"`
	LinesSkipped  int64    `json:"lines_skipped"`
	LinesValid    int64    `json:"lines_valid"`
	Summary       string   `json:"summary"`
}

// Report génère un rapport de parsing à partir des erreurs rapportées.
func (tracker *ParsingTracker) Report(parser, batchKey, filePath string) Report {
	nbParsedLines := tracker.currentLine - 1 // -1 because we started counting at line number 1
	nbValidLines := nbParsedLines - tracker.nbRejectedLines - tracker.nbSkippedLines

	summary := fmt.Sprintf(
		"%s: intégration terminée, %d lignes traitées, %d erreurs fatales, %d lignes rejetées, %d lignes filtrées, %d lignes valides",
		filePath,
		nbParsedLines,
		len(tracker.fatalErrors),
		tracker.nbRejectedLines,
		tracker.nbSkippedLines,
		nbValidLines,
	)

	return Report{
		Commit:        gitCommit,
		StartDate:     tracker.startDate,
		Parser:        parser,
		BatchKey:      batchKey,
		Summary:       summary,
		LinesParsed:   nbParsedLines,
		LinesValid:    nbValidLines,
		LinesSkipped:  tracker.nbSkippedLines,
		LinesRejected: tracker.nbRejectedLines,
		IsFatal:       len(tracker.fatalErrors) > 0,
		HeadRejected:  tracker.firstParseErrors,
		HeadFatal:     tracker.fatalErrors,
	}
}

// NewParsingTracker retourne une instance pour rapporter les erreurs de parsing.
func NewParsingTracker() ParsingTracker {
	return ParsingTracker{
		startDate:              time.Now(),
		currentLine:            1,
		lastSkippedLine:        -1,
		lastLineWithParseError: -1,
		firstParseErrors:       []string{},
		fatalErrors:            []string{},
	}
}
