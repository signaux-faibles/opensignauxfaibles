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

// Events Event serialisable pour swaggo (TODO: fix this !)
// type Events []struct {
// 	ID       bson.ObjectId `json:"-" bson:"_id"`
// 	Date     time.Time     `json:"date" bson:"date"`
// 	Comment  interface{}   `json:"event" bson:"event"`
// 	Priority Priority      `json:"priority" bson:"priority"`
// 	Code     Code          `json:"code" bson:"code"`
// }

// Event est un objet de journal
// swagger:ignore
type Event struct {
	ActualEvent actualEvent
	Channel     chan Event
}

type actualEvent struct {
	ID         bson.ObjectId `json:"id" bson:"_id"`
	Date       time.Time     `json:"date" bson:"date"`
	StartDate  time.Time     `json:"startDate" bson:"startDate"`
	Comment    interface{}   `json:"event" bson:"event"`
	Priority   Priority      `json:"priority" bson:"priority"`
	Code       Code          `json:"parserCode" bson:"parserCode"`
	ReportType string        `json:"reportType" bson:"reportType"`
}

// GetBSON retourne l'objet Event sous une forme sérialisable
func (event Event) GetBSON() (interface{}, error) {
	return event.ActualEvent, nil
}

// CreateEvent initialise un évènement avec les valeurs par défaut.
func CreateEvent() (event Event) {
	actualEvent := actualEvent{
		ID:       bson.NewObjectId(),
		Date:     time.Now(),
		Priority: Priority("info"),
	}
	return Event{
		ActualEvent: actualEvent,
	}
}

func (event Event) throw(comment interface{}, logLevel string) {
	event.ActualEvent.ID = bson.NewObjectId()
	event.ActualEvent.Date = time.Now()
	event.ActualEvent.Comment = comment
	if event.ActualEvent.Code == "" {
		event.ActualEvent.Code = Code("unknown")
	}
	event.ActualEvent.Priority = Priority("info")
	event.Channel <- event
}

// Info produit un évènement de niveau Info
func (event Event) Info(comment interface{}) {
	event.throw(comment, "info")
}

// ParseReport permet d'accéder aux propriétés d'un rapport de parsing.
func (event Event) ParseReport() (map[string]interface{}, error) {
	var jsonDocument map[string]interface{}
	temporaryBytes, err := bson.MarshalJSON(event.ActualEvent.Comment)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(temporaryBytes, &jsonDocument)
	return jsonDocument, err
}
