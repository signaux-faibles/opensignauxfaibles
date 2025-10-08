package apconso

import (
	"io"
	"time"

	"opensignauxfaibles/lib/base"
	"opensignauxfaibles/lib/engine"
	"opensignauxfaibles/lib/parsing"
)

type ApconsoParser struct{}

func NewApconsoParser() engine.Parser {
	return &ApconsoParser{}
}

func (p *ApconsoParser) Type() base.ParserType { return base.Apconso }
func (p *ApconsoParser) New(r io.Reader) engine.ParserInst {
	return &parsing.CsvParserInst{
		Reader:     r,
		RowParser:  &apconsoRowParser{},
		Comma:      ',',
		LazyQuotes: false,
		DestTuple:  APConso{},
	}
}

type apconsoRowParser struct{}

func (rp *apconsoRowParser) ParseRow(row []string, res *engine.ParsedLineResult, idx engine.ColMapping) error {
	var err error

	idxRow := idx.IndexRow(row)
	apconso := APConso{}
	apconso.ID = idxRow.GetVal("ID_DA")
	apconso.Siret = idxRow.GetVal("ETAB_SIRET")

	apconso.Periode, err = time.Parse("2006-01-02", idxRow.GetVal("MOIS"))
	res.AddRegularError(err)
	apconso.HeureConsommee, err = idxRow.GetFloat64("HEURES")
	res.AddRegularError(err)
	apconso.Montant, err = idxRow.GetFloat64("MONTANTS")
	res.AddRegularError(err)
	apconso.Effectif, err = idxRow.GetIntFromFloat("EFFECTIFS")
	res.AddRegularError(err)

	if len(res.Errors) == 0 {
		res.AddTuple(apconso)
	}
	return nil
}
