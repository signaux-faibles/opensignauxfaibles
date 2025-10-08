package urssaf

import (
	"io"

	"opensignauxfaibles/lib/base"
	"opensignauxfaibles/lib/engine"
	"opensignauxfaibles/lib/parsing"
)

type CotisationParser struct{}

func NewCotisationParser() engine.Parser {
	return &CotisationParser{}
}

func (parser *CotisationParser) Type() base.ParserType { return base.Cotisation }
func (parser *CotisationParser) New(r io.Reader) engine.ParserInst {
	return &UrssafParserInst{
		parsing.CsvParserInst{
			Reader:     r,
			RowParser:  &cotisationRowParser{},
			Comma:      ';',
			LazyQuotes: true,
			DestTuple:  Cotisation{},
		},
	}
}

type cotisationRowParser struct {
	UrssafRowParser
}

func (rp *cotisationRowParser) ParseRow(row []string, res *engine.ParsedLineResult, idx parsing.ColIndex) error {

	idxRow := idx.IndexRow(row)

	cotisation := Cotisation{}

	periodeDebut, periodeFin, err := UrssafToPeriod(idxRow.GetVal("periode"))
	res.AddRegularError(err)

	siret, err := rp.GetComptes().GetSiret(idxRow.GetVal("Compte"), &periodeDebut)
	if err != nil {
		res.SetFilterError(err)
	} else {
		cotisation.Siret = siret
		cotisation.NumeroCompte = idxRow.GetVal("Compte")
		cotisation.PeriodeDebut = periodeDebut
		cotisation.PeriodeFin = periodeFin
		cotisation.Encaisse, err = idxRow.GetCommaFloat64("enc_direct")
		res.AddRegularError(err)
		cotisation.Du, err = idxRow.GetCommaFloat64("cotis_due")
		res.AddRegularError(err)
	}
	if len(res.Errors) == 0 {
		res.AddTuple(cotisation)
	}
	return nil
}
