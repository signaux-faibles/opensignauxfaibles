package urssaf

import (
	"io"
	"strconv"

	"opensignauxfaibles/lib/base"
	"opensignauxfaibles/lib/engine"
	"opensignauxfaibles/lib/parsing"
	"opensignauxfaibles/lib/sfregexp"
)

type EffectifEntParser struct{}

func (parser *EffectifEntParser) Type() base.ParserType { return base.EffectifEnt }
func (parser *EffectifEntParser) New(r io.Reader) engine.ParserInst {
	return &EffectifParserInst{
		parsing.CsvParserInst{
			Reader:     r,
			RowParser:  &effectifEntRowParser{},
			Comma:      ';',
			LazyQuotes: false,
			DestTuple:  EffectifEnt{},
		},
		// proper idx defined at `Init`
		engine.ColMapping{},
	}

}

type effectifEntRowParser struct {
	UrssafRowParser
	periods []periodCol
}

// Used at `Init`
func (rp *effectifEntRowParser) setPeriods(periods []periodCol) {
	rp.periods = periods
}

func (rp *effectifEntRowParser) ParseRow(row []string, res *engine.ParsedLineResult, idx engine.ColMapping) error {

	for _, period := range rp.periods {
		value := row[period.colIndex]

		if value != "" {
			noThousandsSep := sfregexp.RegexpDict["notDigit"].ReplaceAllString(value, "")
			s, err := strconv.ParseFloat(noThousandsSep, 64)
			res.AddRegularError(err)
			e := int(s)
			if e > 0 {
				idxRow := idx.IndexRow(row)
				res.AddTuple(EffectifEnt{
					Siren:       idxRow.GetVal("siren"),
					Periode:     period.dateStart,
					EffectifEnt: e,
				})
			}
		}
	}
	return nil
}
