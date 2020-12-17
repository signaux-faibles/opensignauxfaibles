package marshal

import (
	"encoding/json"
	"time"

	"github.com/globalsign/mgo/bson"
)

// Priority test
type Priority string

// Code test
type Code string

// EventInChannel envelope un objet de journal avec sa destination
type EventInChannel struct {
	ActualEvent Event
	Channel     chan Event
}

// Event est un objet de journal
// swagger:ignore
type Event struct {
	ID         bson.ObjectId `json:"-" bson:"_id"`
	Date       time.Time     `json:"date" bson:"date"`
	StartDate  time.Time     `json:"startDate" bson:"startDate"`
	Comment    interface{}   `json:"event" bson:"event"`
	Priority   Priority      `json:"priority" bson:"priority"`
	Code       Code          `json:"parserCode" bson:"parserCode"`
	ReportType string        `json:"report_type" bson:"reportType"`
}

// CreateEvent initialise un évènement avec les valeurs par défaut.
func CreateEvent() (event Event) {
	return Event{
		ID:       bson.NewObjectId(),
		Date:     time.Now(),
		Priority: Priority("info"),
	}
}

func (event EventInChannel) throw(comment interface{}, logLevel string) {
	event.ActualEvent.ID = bson.NewObjectId()
	event.ActualEvent.Date = time.Now()
	event.ActualEvent.Comment = comment
	if event.ActualEvent.Code == "" {
		event.ActualEvent.Code = Code("unknown")
	}
	event.ActualEvent.Priority = Priority("info")
	event.Channel <- event.ActualEvent
}

// Info produit un évènement de niveau Info
func (event EventInChannel) Info(comment interface{}) {
	event.throw(comment, "info")
}

// ParseReport permet d'accéder aux propriétés d'un rapport de parsing.
func (event Event) ParseReport() (map[string]interface{}, error) {
	var jsonDocument map[string]interface{}
	temporaryBytes, err := bson.MarshalJSON(event.Comment)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(temporaryBytes, &jsonDocument)
	return jsonDocument, err
}
