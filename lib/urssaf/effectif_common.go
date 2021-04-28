package urssaf

import (
	"encoding/csv"
	"regexp"
	"time"

	"github.com/signaux-faibles/opensignauxfaibles/lib/marshal"
)

type periodCol struct {
	dateStart time.Time
	colIndex  int
}

func parseEffectifColMapping(reader *csv.Reader, destObject interface{}) (marshal.ColMapping, []periodCol, error) {
	fields, err := reader.Read()
	if err != nil {
		return marshal.ColMapping{}, nil, err
	}
	idx, err := marshal.ValidateAndIndexColumnsFromColTags(marshal.LowercaseFields(fields), destObject)
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
			date, _ := marshal.UrssafToPeriod(field[3:9]) // format: YYQM ou YYYYQM
			periods = append(periods, periodCol{dateStart: date.Start, colIndex: index})
		}
	}
	return periods
}
