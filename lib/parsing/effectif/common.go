// Package effectif fournit les parseurs liés à l'extraction des données
// d'effectif
package effectif

import (
	"regexp"
	"time"

	"opensignauxfaibles/lib/engine"
	"opensignauxfaibles/lib/parsing"
	"opensignauxfaibles/lib/parsing/urssaf"
)

// EffectifParserInst is a variant of parsing.CsvParsirInst
// where the index is defined dynamically
type EffectifParserInst struct {
	parsing.CsvParserInst
}

func (p *EffectifParserInst) Init(cache *engine.Cache, filter engine.SirenFilter, batch *engine.AdminBatch) error {

	err := p.CsvParserInst.Init(cache, filter, batch)
	if err != nil {
		return err
	}

	if rowParser, ok := p.RowParser.(interface{ setPeriods([]periodCol) }); ok {
		periods := parseEffectifPeriod(p.Header())
		rowParser.setPeriods(periods)
	}

	return nil
}

type periodCol struct {
	dateStart time.Time
	colIndex  int
}

// ParseEffectifPeriod extrait les périodes depuis une liste de noms de colonnes csv.
func parseEffectifPeriod(fields []string) []periodCol {
	periods := []periodCol{}
	re, _ := regexp.Compile("^eff")
	for index, field := range fields {
		if re.MatchString(field) {
			dateStart, _, _ := urssaf.UrssafToPeriod(field[3:9]) // format: YYQM ou YYYYQM
			periods = append(periods, periodCol{dateStart: dateStart, colIndex: index})
		}
	}
	return periods
}
