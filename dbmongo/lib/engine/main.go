package engine

import (
	"fmt"
	"strconv"
	"time"

	"github.com/globalsign/mgo/bson"
	"github.com/signaux-faibles/gournal"
)

//go:generate go run js/loadJS.go

// Db connecteur exportable
var Db DB

// Priority test
type Priority string

// Code test
type Code string

// MaxParsingErrors is the number of parsing errors that are needed to
// interrupt a parser
var MaxParsingErrors = 200

// Event est un objet de journal
// swagger:ignore
type Event struct {
	ID         bson.ObjectId `json:"-" bson:"_id"`
	Date       time.Time     `json:"date" bson:"date"`
	Comment    interface{}   `json:"event" bson:"event"`
	Priority   Priority      `json:"priority" bson:"priority"`
	Code       Code          `json:"code" bson:"code"`
	ReportType string        `json:"report_type" bson:"record_type"`
	Channel    chan Event    `json:"-"`
}

// Events Event serialisable pour swaggo (TODO: fix this !)
// type Events []struct {
// 	ID       bson.ObjectId `json:"-" bson:"_id"`
// 	Date     time.Time     `json:"date" bson:"date"`
// 	Comment  interface{}   `json:"event" bson:"event"`
// 	Priority Priority      `json:"priority" bson:"priority"`
// 	Code     Code          `json:"code" bson:"code"`
// }

// GetBSON retourne l'objet Event sous une forme sérialisable
func (event Event) GetBSON() (interface{}, error) {
	var tmp struct {
		ID       bson.ObjectId `json:"id" bson:"_id"`
		Date     time.Time     `json:"date" bson:"date"`
		Comment  interface{}   `json:"event" bson:"event"`
		Priority Priority      `json:"priority" bson:"priority"`
		Code     Code          `json:"code" bson:"code"`
	}
	tmp.ID = event.ID
	tmp.Date = event.Date
	tmp.Comment = event.Comment
	tmp.Priority = event.Priority
	tmp.Code = event.Code
	return tmp, nil
}

// Debug test
var Debug = Priority("debug")

// Info test
var Info = Priority("info")

// Warning test
var Warning = Priority("warning")

// Critical test
var Critical = Priority("critical")

var unknownCode = Code("unknown")

func (event Event) throw(comment interface{}, logLevel string) {
	event.ID = bson.NewObjectId()
	event.Date = time.Now()
	event.Comment = comment
	if event.Code == "" {
		event.Code = unknownCode
	}
	switch logLevel {
	case "debug":
		event.Priority = Debug
	case "info":
		event.Priority = Info
	case "warning":
		event.Priority = Warning
	case "critical":
		event.Priority = Critical
	default:
		panic("Wrong use of throw function")
	}
	event.Channel <- event
}

// Debug produit un évènement de niveau Debug
func (event Event) Debug(comment interface{}) {
	event.throw(comment, "debug")
}

// Info produit un évènement de niveau Info
func (event Event) Info(comment interface{}) {
	event.throw(comment, "info")
}

// Warning produit un évènement de niveau Warning
func (event Event) Warning(comment interface{}) {
	event.throw(comment, "warning")
}

// Critical produit un évènement de niveau Critical
func (event Event) Critical(comment interface{}) {
	event.throw(comment, "critical")
}

// DebugReport produit un rapport de niveau Debug
func (event Event) DebugReport(report string, tracker gournal.Tracker) {
	event.ReportType = report
	event.Debug(tracker.Report(report))
}

// InfoReport produit un rapport de niveau Info
func (event Event) InfoReport(report string, tracker gournal.Tracker) {
	event.ReportType = report
	event.Info(tracker.Report(report))
}

// WarningReport produit un rapport de niveau Warning
func (event Event) WarningReport(report string, tracker gournal.Tracker) {
	event.ReportType = report
	event.Warning(tracker.Report(report))
}

// CriticalReport produit un rapport de niveau Critical
func (event Event) CriticalReport(report string, tracker gournal.Tracker) {
	event.ReportType = report
	event.Critical(tracker.Report(report))
}

func reportAbstract(tracker gournal.Tracker) interface{} {

	var nError = 0
	var nFiltered = 0
	var nFatal = 0

	var fatalErrors = []string{}
	var filterErrors = []string{}
	var errorErrors = []string{}
	for ind := range tracker.Errors {
		var hasError = false
		var hasFilter = false
		var hasFatal = false
		for _, e := range tracker.Errors[ind] {
			switch c := e.(type) {
			case CriticityError:
				if c.Criticity() == "fatal" {
					hasFatal = true
					if len(fatalErrors) < MaxParsingErrors {
						fatalErrors = append(fatalErrors, fmt.Sprintf("Cycle %d: %v", ind, e))
					}
				}
				if c.Criticity() == "error" {
					hasError = true
					if len(errorErrors) < MaxParsingErrors {
						errorErrors = append(errorErrors, fmt.Sprintf("Cycle %d: %v", ind, e))
					}
				}
				if c.Criticity() == "filter" {
					hasFilter = true
					if len(filterErrors) < MaxParsingErrors {
						filterErrors = append(filterErrors, fmt.Sprintf("Cycle %d: %v", ind, e))
					}
				}
			default:
				hasFatal = true
				if len(fatalErrors) < MaxParsingErrors {
					fatalErrors = append(fatalErrors, fmt.Sprintf("Cycle %d: %v", ind, e))
					fmt.Printf("Cycle %d: %v", ind, e)
				}
			}
		}
		if hasFatal {
			nFatal = nFatal + 1
		} else if hasError {
			nError = nError + 1
		} else if hasFilter {
			nFiltered = nFiltered + 1
		}
	}
	nValid := tracker.Count + 1 - nFatal - nError - nFiltered
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

func ShouldBreak(tracker gournal.Tracker, maxErrors int) bool {
	l := 0
	hasError := false
	for _, errs := range tracker.Errors {
		for _, e := range errs {
			switch c := e.(type) {
			case CriticityError:
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
			l += 1
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
