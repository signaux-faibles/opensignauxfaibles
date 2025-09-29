package createfilter

import (
	"errors"
	"strconv"
	"time"
)

// Periode est un type temporel avec un début et une fin employé dans les types et
// fonctions opensignauxfaibles manipulant des périodes temporelles. La date de fin
// est exclue de la période.
type Periode struct {
	Start time.Time `json:"start" bson:"start"`
	End   time.Time `json:"end" bson:"end"`
}

// UrssafToPeriod convertit le format de période urssaf en type Periode. On trouve ces
// périodes formatées en 4 ou 6 caractère (YYQM ou YYYYQM).
// si YY < 50 alors YYYY = 20YY sinon YYYY = 19YY.
// si QM == 62 alors période annuelle sur YYYY.
// si M == 0 alors période trimestrielle sur le trimestre Q de YYYY.
// si 0 < M < 4 alors mois M du trimestre Q.
func UrssafToPeriod(urssaf string) (Periode, error) {
	period := Periode{}

	if len(urssaf) == 4 {
		if urssaf[0:2] < "50" {
			urssaf = "20" + urssaf
		} else {
			urssaf = "19" + urssaf
		}
	}

	if len(urssaf) != 6 {
		return period, errors.New("valeur non autorisée")
	}

	year, err := strconv.Atoi(urssaf[0:4])
	if err != nil {
		return Periode{}, errors.New("valeur non autorisée")
	}

	if urssaf[4:6] == "62" {
		period.Start = time.Date(year, time.Month(1), 1, 0, 0, 0, 0, time.UTC)
		period.End = time.Date(year+1, time.Month(1), 1, 0, 0, 0, 0, time.UTC)
	} else {
		quarter, err := strconv.Atoi(urssaf[4:5])
		if err != nil {
			return period, err
		}
		monthOfQuarter, err := strconv.Atoi(urssaf[5:6])
		if err != nil {
			return period, err
		}
		if monthOfQuarter == 0 {
			period.Start = time.Date(year, time.Month((quarter-1)*3+1), 1, 0, 0, 0, 0, time.UTC)
			period.End = time.Date(year, time.Month((quarter-1)*3+4), 1, 0, 0, 0, 0, time.UTC)
		} else {
			period.Start = time.Date(year, time.Month((quarter-1)*3+monthOfQuarter), 1, 0, 0, 0, 0, time.UTC)
			period.End = time.Date(year, time.Month((quarter-1)*3+monthOfQuarter+1), 1, 0, 0, 0, 0, time.UTC)
		}
	}
	return period, nil
}

func effectifColNameToDate(colName string) (time.Time, error) {
	if colName[0:3] != "eff" {
		return time.Time{}, errors.New("this column is not a valid period: " + colName)
	}
	urssafPeriode := colName[3:9]
	periode, err := UrssafToPeriod(urssafPeriode)
	return periode.Start, err
}
