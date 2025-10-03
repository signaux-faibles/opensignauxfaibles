package urssaf

import (
	"encoding/csv"
	"regexp"
	"time"

	"opensignauxfaibles/lib/engine"
)

type periodCol struct {
	dateStart time.Time
	colIndex  int
}

func parseEffectifColMapping(reader *csv.Reader, destObject interface{}) (engine.ColMapping, []periodCol, error) {
	fields, err := reader.Read()
	if err != nil {
		return engine.ColMapping{}, nil, err
	}
	idx, err := engine.ValidateAndIndexColumnsFromInputTags(engine.LowercaseFields(fields), destObject)
	// Dans quels champs lire l'effectifEnt
	periods := parseEffectifPeriod(fields)
	return idx, periods, err
}

// ParseEffectifPeriod extrait les périodes depuis une liste de noms de colonnes csv.
func parseEffectifPeriod(fields []string) []periodCol {
	periods := []periodCol{}
	re, _ := regexp.Compile("^eff")
	for index, field := range fields {
		if re.MatchString(field) {
			dateStart, _, _ := engine.UrssafToPeriod(field[3:9]) // format: YYQM ou YYYYQM
			periods = append(periods, periodCol{dateStart: dateStart, colIndex: index})
		}
	}
	return periods
}
