package effectif

import (
	"io"
	"strconv"

	"opensignauxfaibles/lib/engine"
	"opensignauxfaibles/lib/parsing"
	"opensignauxfaibles/lib/sfregexp"
)

type EffectifParser struct{}

func NewEffectifParser() engine.Parser {
	return &EffectifParser{}
}

func (parser *EffectifParser) Type() engine.ParserType { return engine.Effectif }
func (parser *EffectifParser) New(r io.Reader) engine.ParserInst {
	return &EffectifParserInst{
		parsing.CsvParserInst{
			Reader:     r,
			RowParser:  &effectifRowParser{},
			Comma:      ';',
			LazyQuotes: false,
			DestTuple:  Effectif{},
		},
	}

}

type effectifRowParser struct {
	periods []periodCol
}

// Used at `Init`
func (rp *effectifRowParser) setPeriods(periods []periodCol) {
	rp.periods = periods
}

func (rp *effectifRowParser) ParseRow(row []string, res *engine.ParsedLineResult, idx parsing.ColIndex) {
	for _, period := range rp.periods {
		value := row[period.colIndex]

		if value != "" {
			noThousandsSep := sfregexp.RegexpDict["notDigit"].ReplaceAllString(value, "")
			e, err := strconv.Atoi(noThousandsSep)
			res.AddRegularError(err)
			if e >= 0 {
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
}
