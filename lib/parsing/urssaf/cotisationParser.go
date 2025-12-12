package urssaf

import (
	"io"

	"opensignauxfaibles/lib/engine"
	"opensignauxfaibles/lib/parsing"
)

type CotisationParser struct{}

func NewCotisationParser() engine.Parser {
	return &CotisationParser{}
}

func (parser *CotisationParser) Type() engine.ParserType { return engine.Cotisation }
func (parser *CotisationParser) New(r io.Reader) engine.ParserInst {
	return &parsing.CsvParserInst{
		Reader:        r,
		RowParser:     &cotisationRowParser{},
		Comma:         ';',
		LazyQuotes:    true,
		CaseSensitive: false,
		DestTuple:     Cotisation{},
	}
}

type cotisationRowParser struct{}

func (rp *cotisationRowParser) ParseRow(row []string, res *engine.ParsedLineResult, idx parsing.ColIndex) {

	idxRow := idx.IndexRow(row)

	cotisation := Cotisation{}

	periodeDebut, periodeFin, err := UrssafToPeriod(idxRow.GetVal("periode"))
	res.AddRegularError(err)

	cotisation.Siret = idxRow.GetVal("Siret")
	cotisation.NumeroCompte = idxRow.GetVal("Compte")
	cotisation.PeriodeDebut = periodeDebut
	cotisation.PeriodeFin = periodeFin
	cotisation.Encaisse, err = idxRow.GetCommaFloat64("enc_direct")
	res.AddRegularError(err)
	cotisation.Du, err = idxRow.GetCommaFloat64("cotis_due")
	res.AddRegularError(err)

	if len(res.Errors) == 0 {
		res.AddTuple(cotisation)
	}
}
