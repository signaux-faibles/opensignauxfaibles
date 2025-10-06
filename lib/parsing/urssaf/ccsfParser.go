package urssaf

import (
	"encoding/csv"
	"errors"
	"os"

	"opensignauxfaibles/lib/base"
	"opensignauxfaibles/lib/engine"
)

type ccsfParser struct {
	file    *os.File
	reader  *csv.Reader
	comptes engine.Comptes
	idx     engine.ColMapping
}

func NewParserCCSF() *ccsfParser {
	return &ccsfParser{}
}

func (parser *ccsfParser) Type() base.ParserType {
	return base.Ccsf
}

func (parser *ccsfParser) Init(cache *engine.Cache, filter engine.SirenFilter, batch *base.AdminBatch) (err error) {
	parser.comptes, err = engine.GetCompteSiretMapping(*cache, batch, filter, engine.OpenAndReadSiretMapping)
	return err
}

func (parser *ccsfParser) Close() error {
	return parser.file.Close()
}

func (parser *ccsfParser) Open(filePath base.BatchFile) (err error) {
	parser.file, parser.reader, err = engine.OpenCsvReader(filePath, ';', false)
	if err == nil {
		parser.idx, err = engine.IndexColumnsFromCsvHeader(parser.reader, CCSF{})
	}
	return err
}

func (parser *ccsfParser) ReadNext(res *engine.ParsedLineResult) error {
	row, err := parser.reader.Read()
	if err != nil {
		return err
	}

	ccsf := CCSF{}
	if len(row) >= 4 {
		idxRow := parser.idx.IndexRow(row)
		ccsf.Action = idxRow.GetVal("Code_externe_action")
		ccsf.Stade = idxRow.GetVal("Code_externe_stade")
		ccsf.DateTraitement, err = engine.UrssafToDate(idxRow.GetVal("Date_de_traitement"))
		res.AddRegularError(err)
		if err != nil {
			return err
		}

		ccsf.key, err = engine.GetSiretFromComptesMapping(idxRow.GetVal("Compte"),
			&ccsf.DateTraitement, parser.comptes)
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
