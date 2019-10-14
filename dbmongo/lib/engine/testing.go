package engine

import (
	"sync"

	"github.com/globalsign/mgo/bson"
)

// DiscardEvents supprime les évènements
func DiscardEvents(events chan Event) {
	go func() {
		for range events {
		}
	}()
}

// DiscardTuple supprime les évènements
func DiscardTuple(tuples chan Tuple) {
	go func() {
		for range tuples {
		}
	}()
}

// AnalyseEvents extracts information from events. Wait til waitgroup is done before doing
// anything with the output
func AnalyseEvents(events chan Event, wg *sync.WaitGroup) (*int, *int, *int, *int, bool) {

	type abstractData struct {
		Report   string
		Total    int
		Valid    int
		Filtered int
		Error    int
	}

	var s abstractData
	var fatal bool
	go func(a *abstractData) {
		if wg != nil {
			defer wg.Done()
		}
		for event := range events {
			if event.ReportType == "abstract" {
				comment := event.Comment.(bson.M)
				bsonBytes, _ := bson.Marshal(comment)
				err := bson.Unmarshal(bsonBytes, a)
				if err != nil {
					panic("Could not unmarshal abstract report data: " + err.Error())
				}
			}
			if event.Priority == "critical" {
				fatal = true
			}
		}
	}(&s)
	return &s.Total, &s.Valid, &s.Filtered, &s.Error, fatal
}
