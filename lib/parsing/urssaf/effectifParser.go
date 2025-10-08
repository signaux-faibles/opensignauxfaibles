package urssaf

import (
	"io"
	"strconv"

	"opensignauxfaibles/lib/base"
	"opensignauxfaibles/lib/engine"
	"opensignauxfaibles/lib/parsing"
	"opensignauxfaibles/lib/sfregexp"
)

type EffectifParser struct{}

func (parser *EffectifParser) Type() base.ParserType { return base.Effectif }
func (parser *EffectifParser) New(r io.Reader) engine.ParserInst {
	return &EffectifParserInst{
		parsing.CsvParserInst{
			Reader:     r,
			RowParser:  &effectifRowParser{},
			Comma:      ';',
			LazyQuotes: false,
			DestTuple:  Effectif{},
		},
		// proper idx defined at `Init`
		engine.ColMapping{},
	}

}

type effectifRowParser struct {
	UrssafRowParser
	periods []periodCol
}

// Used at `Init`
func (rp *effectifRowParser) setPeriods(periods []periodCol) {
	rp.periods = periods
}

func (rp *effectifRowParser) ParseRow(row []string, res *engine.ParsedLineResult, idx engine.ColMapping) error {
	for _, period := range rp.periods {
		value := row[period.colIndex]

		if value != "" {
			noThousandsSep := sfregexp.RegexpDict["notDigit"].ReplaceAllString(value, "")
			e, err := strconv.Atoi(noThousandsSep)
			res.AddRegularError(err)
			if e > 0 {
				idxRow := idx.IndexRow(row)
				res.AddTuple(Effectif{
					Siret:        idxRow.GetVal("siret"),
					NumeroCompte: idxRow.GetVal("compte"),
					Periode:      period.dateStart,
					Effectif:     e,
				})
			}
		}
	}
	return nil
}
