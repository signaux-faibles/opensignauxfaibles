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

// Event est un objet de journal
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
func CreateEvent() Event {
	return Event{
		ID:       bson.NewObjectId(),
		Date:     time.Now(),
		Priority: Priority("info"),
		Code:     Code("unknown"),
	}
}

// CreateReportEvent initialise un évènement contenant un rapport de parsing.
func CreateReportEvent(fileType string, report interface{}) Event {
	event := CreateEvent()
	event.Code = Code(fileType)
	event.Comment = report
	return event
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
