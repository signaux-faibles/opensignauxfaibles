package urssaf

import (
	"errors"
	"io"

	"opensignauxfaibles/lib/engine"
	"opensignauxfaibles/lib/parsing"
)

type CCSFParser struct{}

func NewCCSFParser() engine.Parser {
	return &CCSFParser{}
}

func (parser *CCSFParser) Type() engine.ParserType { return engine.Ccsf }
func (parser *CCSFParser) New(r io.Reader) engine.ParserInst {
	return &parsing.CsvParserInst{
		Reader:        r,
		RowParser:     &ccsfRowParser{},
		Comma:         ';',
		LazyQuotes:    false,
		CaseSensitive: false,
		DestTuple:     CCSF{},
	}
}

type ccsfRowParser struct{}

func (rp *ccsfRowParser) ParseRow(row []string, res *engine.ParsedLineResult, idx parsing.ColIndex) {

	var err error
	ccsf := CCSF{}
	if len(row) >= 4 {
		idxRow := idx.IndexRow(row)
		ccsf.Action = idxRow.GetVal("Code_externe_action")
		ccsf.Stade = idxRow.GetVal("Code_externe_stade")
		ccsf.DateTraitement, err = UrssafToDate(idxRow.GetVal("Date_de_traitement"))
		res.AddRegularError(err)

		if err != nil {
			return
		}

		ccsf.Siret = idxRow.GetVal("Siret")
		ccsf.NumeroCompte = idxRow.GetVal("Compte")
	} else {
		res.AddRegularError(errors.New("invalid line, fewer than 4 fields"))
	}
	if len(res.Errors) == 0 {
		res.AddTuple(ccsf)
	}
}
