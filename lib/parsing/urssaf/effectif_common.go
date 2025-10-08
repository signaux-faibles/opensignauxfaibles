package urssaf

import (
	"regexp"
	"time"

	"opensignauxfaibles/lib/base"
	"opensignauxfaibles/lib/engine"
	"opensignauxfaibles/lib/parsing"
)

// EffectifParserInst is a variant of parsing.CsvParsirInst
// where the index is defined dynamically
type EffectifParserInst struct {
	parsing.CsvParserInst

	idx engine.ColMapping
}

func (p *EffectifParserInst) Init(cache *engine.Cache, filter engine.SirenFilter,
	batch *base.AdminBatch) error {
	err := p.CsvParserInst.Init(cache, filter, batch)
	if err != nil {
		return err
	}

	idx, periods, err := parseEffectifColMapping(p.Header(), EffectifEnt{})
	p.idx = idx

	if rowParser, ok := p.CsvParserInst.RowParser.(interface{ setPeriods([]periodCol) }); ok {
		rowParser.setPeriods(periods)
	}

	return nil
}

type periodCol struct {
	dateStart time.Time
	colIndex  int
}

func parseEffectifColMapping(fields []string, destObject interface{}) (engine.ColMapping, []periodCol, error) {
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
