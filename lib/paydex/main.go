package paydex

import (
	"time"
)

// Paydex décrit le format de chaque entrée de donnée résultant du parsing.
type Paydex struct {
	Siren   string    `json:"siren" bson:"siren"`
	Periode time.Time `json:"periode" bson:"periode"`
	Jours   int       `json:"paydex_jours" bson:"paydex_jours"`
}

func parsePaydexLine(row []string) Paydex {
	return Paydex{
		Siren:   "000000001",
		Periode: time.Date(2018, 12, 01, 00, 00, 00, 0, time.UTC),
		Jours:   2,
	}
}

