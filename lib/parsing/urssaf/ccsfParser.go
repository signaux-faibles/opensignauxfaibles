package urssaf

import (
	"errors"
	"io"

	"opensignauxfaibles/lib/base"
	"opensignauxfaibles/lib/engine"
	"opensignauxfaibles/lib/parsing"
)

type CCSFParser struct{}

func NewCCSFParser() engine.Parser {
	return &CCSFParser{}
}

func (parser *CCSFParser) Type() base.ParserType { return base.Ccsf }
func (parser *CCSFParser) New(r io.Reader) engine.ParserInst {
	return &UrssafParserInst{
		parsing.CsvParserInst{
			Reader:        r,
			RowParser:     &ccsfRowParser{},
			Comma:         ';',
			LazyQuotes:    false,
			CaseSensitive: false,
			DestTuple:     CCSF{},
		},
	}
}

type ccsfRowParser struct {
	UrssafRowParser
}

func (rp *ccsfRowParser) ParseRow(row []string, res *engine.ParsedLineResult, idx parsing.ColIndex) error {

	var err error
	ccsf := CCSF{}
	if len(row) >= 4 {
		idxRow := idx.IndexRow(row)
		ccsf.Action = idxRow.GetVal("Code_externe_action")
		ccsf.Stade = idxRow.GetVal("Code_externe_stade")
		ccsf.DateTraitement, err = UrssafToDate(idxRow.GetVal("Date_de_traitement"))
		res.AddRegularError(err)

		if err != nil {
			return err
		}

		ccsf.key, err = rp.GetComptes().GetSiret(
			idxRow.GetVal("Compte"),
			&ccsf.DateTraitement,
		)

		if err != nil {
			// Compte filtré
			res.SetFilterError(err)
			return err
		}
		ccsf.NumeroCompte = idxRow.GetVal("Compte")
	} else {
		res.AddRegularError(errors.New("ligne non conforme, moins de 4 champs"))
	}
	if len(res.Errors) == 0 {
		res.AddTuple(ccsf)
	}
	return nil
}
