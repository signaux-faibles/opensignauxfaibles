package files

import (
	"time"
)

// FileSummary représente un fichier
type FileSummary struct {
	Name string    `json:"name" bson:"name"`
	Size int64     `json:"size" bson:"size"`
	Date time.Time `json:"date" bson:"date"`
}
