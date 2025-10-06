package apconso

import (
	"encoding/csv"
	"os"
	"time"

	"opensignauxfaibles/lib/base"
	"opensignauxfaibles/lib/engine"
)

type apconsoParser struct {
	file   *os.File
	reader *csv.Reader
	idx    engine.ColMapping
}

func NewParserApconso() *apconsoParser {
	return &apconsoParser{}
}

func (parser *apconsoParser) Type() base.ParserType {
	return base.Apconso
}

func (parser *apconsoParser) Init(_ *engine.Cache, _ engine.SirenFilter, _ *base.AdminBatch) error {
	return nil
}

func (parser *apconsoParser) Open(filePath base.BatchFile) (err error) {
	parser.file, parser.reader, err = engine.OpenCsvReader(filePath, ',', false)
	if err == nil {
		parser.idx, err = engine.IndexColumnsFromCsvHeader(parser.reader, APConso{})
	}
	return err
}

func (parser *apconsoParser) Close() error {
	return parser.file.Close()
}

func (parser *apconsoParser) ReadNext(res *engine.ParsedLineResult) error {
	row, err := parser.reader.Read()
	if err != nil {
		return err
	}

	idxRow := parser.idx.IndexRow(row)
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
