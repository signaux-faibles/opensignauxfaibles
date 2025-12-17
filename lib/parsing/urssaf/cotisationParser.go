package urssaf

import (
	"fmt"
	"io"

	"opensignauxfaibles/lib/engine"
	"opensignauxfaibles/lib/parsing"
	"opensignauxfaibles/lib/sfregexp"
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

	siret := idxRow.GetVal("Siret")

	if sfregexp.RegexpDict["acossInternal"].MatchString(siret) {
		res.SetFilterError(fmt.Errorf("acoss internal id: %s", siret))
		return
	}

	cotisation.Siret = siret
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
