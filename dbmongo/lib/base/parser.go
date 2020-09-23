package base

import (
	"time"

	"github.com/globalsign/mgo/bson"
	"github.com/signaux-faibles/gournal"
)

// Parser fonction de traitement de données en entrée
type Parser func(Cache, *AdminBatch) (chan Tuple, chan Event)

// Tuple unité de donnée à insérer dans un type
type Tuple interface {
	Key() string
	Scope() string
	Type() string
}

// Priority test
type Priority string

// Code test
type Code string

// base.Events base.Event serialisable pour swaggo (TODO: fix this !)
// type base.Events []struct {
// 	ID       bson.ObjectId `json:"-" bson:"_id"`
// 	Date     time.Time     `json:"date" bson:"date"`
// 	Comment  interface{}   `json:"event" bson:"event"`
// 	Priority Priority      `json:"priority" bson:"priority"`
// 	Code     Code          `json:"code" bson:"code"`
// }

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

// GetBSON retourne l'objet base.Event sous une forme sérialisable
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

// InfoReport produit un rapport de niveau Info
func (event Event) InfoReport(report string, tracker gournal.Tracker) {
	event.ReportType = report
	event.Info(tracker.Report(report))
}

// CriticalReport produit un rapport de niveau Critical
func (event Event) CriticalReport(report string, tracker gournal.Tracker) {
	event.ReportType = report
	event.Critical(tracker.Report(report))
}
