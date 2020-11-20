package marshal

import (
	"fmt"
	"sort"

	"github.com/globalsign/mgo/bson"
	"github.com/signaux-faibles/gournal"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/base"
)

// MaxParsingErrors is the number of parsing errors to report per file.
var MaxParsingErrors = 200

// ParsingTracker permet de collecter puis rapporter des erreurs de parsing.
type ParsingTracker struct {
	gournalTracker gournal.Tracker
}

// Add rapporte une erreur de parsing à la ligne en cours.
func (tracker *ParsingTracker) Add(err base.CriticityError) {
	tracker.gournalTracker.Add(err)
}

// Next informe le Tracker qu'on passe à la ligne suivante.
func (tracker *ParsingTracker) Next() {
	tracker.gournalTracker.Next()
}

// Report génère un rapport de parsing à partir des erreurs rapportées.
func (tracker *ParsingTracker) Report(code string) interface{} {
	return tracker.gournalTracker.Report(code)
}

// NewParsingTracker retourne une instance pour rapporter les erreurs de parsing.
func NewParsingTracker(batchKey string, filePath string) ParsingTracker {
	// fonctions de reporting du moteur
	trackerReports := map[string]gournal.ReportFunction{
		"abstract": reportAbstract,
	}
	context := map[string]string{
		"path":     filePath,
		"batchKey": batchKey,
	}
	return ParsingTracker{
		gournalTracker: gournal.NewTracker(context, trackerReports),
	}
}

func reportAbstract(tracker gournal.Tracker) interface{} {

	var nError = 0
	var nFiltered = 0
	var nFatal = 0

	var fatalErrors = []string{}
	var filterErrors = []string{}
	var errorErrors = []string{}

	// En Golang, l'ordre des clés d'un map n'est pas garanti. (https://blog.golang.org/maps)
	// => On ordonne les erreurs par numéro de cycle, pour permettre la reproductibilité.
	// cf https://github.com/signaux-faibles/opensignauxfaibles/issues/181
	var cycles []int
	for cycle := range tracker.Errors {
		cycles = append(cycles, cycle)
	}
	sort.Ints(cycles)

	// pour chaque cycle qui a au moins une erreur
	for _, cycle := range cycles {
		var lineRejected = false
		// pour chaque erreur du cycle
		cycleErrors := tracker.Errors[cycle]
		for _, err := range cycleErrors {
			switch err.(base.CriticityError).Criticity() {
			case "fatal":
				nFatal++
				if len(fatalErrors) < MaxParsingErrors {
					fatalErrors = append(fatalErrors, fmt.Sprintf("Cycle %d: %v", cycle, err))
				}
			case "error":
				lineRejected = true
				if len(errorErrors) < MaxParsingErrors {
					errorErrors = append(errorErrors, fmt.Sprintf("Cycle %d: %v", cycle, err))
				}
			case "filter":
				nFiltered++
				if len(filterErrors) < MaxParsingErrors {
					filterErrors = append(filterErrors, fmt.Sprintf("Cycle %d: %v", cycle, err))
				}
			}
		}
		if lineRejected {
			nError++
		}
	}

	var nValid int
	if nFatal > 0 {
		nValid = 0
	} else {
		nValid = tracker.Count - nError - nFiltered
	}

	report := fmt.Sprintf(
		"%s: intégration terminée, %d lignes traitées, %d erreurs fatales, %d lignes rejetées, %d lignes filtrées, %d lignes valides",
		tracker.Context["path"],
		tracker.Count,
		nFatal,
		nError,
		nFiltered,
		nValid,
	)

	return bson.M{
		"batchKey":      tracker.Context["batchKey"],
		"summary":       report,
		"linesParsed":   tracker.Count,
		"linesValid":    nValid,
		"linesSkipped":  nFiltered,
		"linesRejected": nError,
		"isFatal":       nFatal > 0,
		"headRejected":  errorErrors,
		"headFatal":     fatalErrors,
	}
}
