package base

import (
	"time"

	"github.com/globalsign/mgo/bson"
)

// Parser fonction de traitement de données en entrée
type Parser func(Cache, *AdminBatch) (chan Tuple, chan Event)

// Tuple unité de donnée à insérer dans un type
type Tuple interface {
	Key() string
	Scope() string
	Type() string
}

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
