package engine

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/globalsign/mgo/bson"
	"github.com/signaux-faibles/gournal"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/base"
)

// Db connecteur exportable
var Db DB

// MaxParsingErrors is the number of parsing errors that are needed to
// interrupt a parser
var MaxParsingErrors = 200

func reportAbstract(tracker gournal.Tracker) interface{} {

	var nError = 0
	var nFiltered = 0
	var nFatal = 0

	var fatalErrors = []string{}
	var filterErrors = []string{}
	var errorErrors = []string{}
	for cycle := range tracker.Errors {
		for _, err := range tracker.Errors[cycle] {
			switch c := err.(type) {
			case base.CriticityError:
				if c.Criticity() == "fatal" {
					nFatal++
					if len(fatalErrors) < MaxParsingErrors {
						fatalErrors = append(fatalErrors, fmt.Sprintf("Cycle %d: %v", cycle, err))
					}
				}
				if c.Criticity() == "error" {
					nError++
					if len(errorErrors) < MaxParsingErrors {
						errorErrors = append(errorErrors, fmt.Sprintf("Cycle %d: %v", cycle, err))
					}
				}
				if c.Criticity() == "filter" {
					nFiltered++
					if len(filterErrors) < MaxParsingErrors {
						filterErrors = append(filterErrors, fmt.Sprintf("Cycle %d: %v", cycle, err))
					}
				}
			default:
				nFatal++
				if len(fatalErrors) < MaxParsingErrors {
					fatalErrors = append(fatalErrors, fmt.Sprintf("Cycle %d: %v", cycle, err))
					fmt.Printf("Cycle %d: %v", cycle, err)
				}
			}
		}
	}

	// En Golang, l'ordre des clés d'un map n'est pas garanti. (https://blog.golang.org/maps)
	// => On ordonne les erreurs pour permettre la reproductibilité.
	// cf https://github.com/signaux-faibles/opensignauxfaibles/issues/181
	sort.Strings(fatalErrors)
	sort.Strings(filterErrors)
	sort.Strings(errorErrors)

	nValid := tracker.Count - nFatal - nError - nFiltered
	report := fmt.Sprintf(
		"%s: intégration terminée, %d lignes traitées, %d erreures fatales, %d rejets, %d lignes filtrées, %d lignes valides",
		tracker.Context["path"],
		tracker.Count,
		nFatal,
		nError,
		nFiltered,
		nValid,
	)

	return bson.M{
		"batchKey":    tracker.Context["batchKey"],
		"report":      report,
		"total":       tracker.Count,
		"valid":       nValid,
		"filtered":    nFiltered,
		"error":       nError,
		"headFilters": filterErrors,
		"headErrors":  errorErrors,
		"headFatal":   fatalErrors,
	}
}

func reportCycleErrors(tracker gournal.Tracker) interface{} {
	return bson.M{
		"report":      tracker.Context["path"] + ": ligne " + strconv.Itoa(tracker.Count) + " ignorée",
		"errorReport": tracker.ErrorsInCurrentCycle(),
	}
}

func reportFatalError(tracker gournal.Tracker) interface{} {
	report := "Erreur fatale, abandon"
	if errs, ok := tracker.Errors[tracker.Count]; ok {
		report = report + ": " + errs[len(errs)-1].Error()
	}
	return bson.M{
		"report": report,
	}
}

// ShouldBreak returns true when there are to much errors regarding maxErrors
func ShouldBreak(tracker gournal.Tracker, maxErrors int) bool {
	l := 0
	hasError := false
	for _, errs := range tracker.Errors {
		for _, e := range errs {
			switch c := e.(type) {
			case base.CriticityError:
				if c.Criticity() == "fatal" {
					hasError = true
				}
				if c.Criticity() == "error" {
					hasError = true
				}
				if c.Criticity() == "filter" {
				}
			default:
				hasError = true
			}
		}
		if hasError {
			l++
		}
	}
	return l > maxErrors
}

// TrackerReports contient les fonctions de reporting du moteur
var TrackerReports = map[string]gournal.ReportFunction{
	"abstract":   reportAbstract,
	"errors":     reportCycleErrors,
	"fatalError": reportFatalError,
}
