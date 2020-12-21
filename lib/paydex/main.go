package paydex

import (
	"log"
	"strconv"
	"time"
)

// Paydex décrit le format de chaque entrée de donnée résultant du parsing.
type Paydex struct {
	Siren   string    `json:"siren" bson:"siren"`
	Periode time.Time `json:"periode" bson:"periode"`
	Jours   int       `json:"paydex_jours" bson:"paydex_jours"`
}

func parsePaydexLine(row []string) Paydex {
	periode, err := time.Parse("02/01/2006", row[3])
	if err != nil {
		log.Fatalf("invalid date: %v", row[3])
	}
	jours, err := strconv.Atoi(row[1])
	if err != nil {
		log.Fatalf("invalid date: %v", row[3])
	}
	return Paydex{
		Siren:   row[0],
		Periode: beginningOfMonth(periode),
		Jours:   jours,
	}
}

func beginningOfMonth(date time.Time) time.Time {
	return date.AddDate(0, 0, -date.Day()+1)
}

	}
}

