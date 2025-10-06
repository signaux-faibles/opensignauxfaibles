package urssaf

import (
	"encoding/csv"
	"os"

	"opensignauxfaibles/lib/base"
	"opensignauxfaibles/lib/engine"
)

func NewParserCotisation() *parserCotisation {
	return &parserCotisation{}
}

// parserCotisation implements engine.Parser
type parserCotisation struct {
	file    *os.File
	reader  *csv.Reader
	comptes engine.Comptes
	idx     engine.ColMapping
}

func (parser *parserCotisation) Type() base.ParserType {
	return base.Cotisation
}

func (parser *parserCotisation) Init(cache *engine.Cache, filter engine.SirenFilter, batch *base.AdminBatch) (err error) {
	parser.comptes, err = engine.GetCompteSiretMapping(*cache, batch, filter, engine.OpenAndReadSiretMapping)
	return err
}

func (parser *parserCotisation) Open(filePath base.BatchFile) (err error) {
	parser.file, parser.reader, err = engine.OpenCsvReader(filePath, ';', true)
	if err == nil {
		parser.idx, err = engine.IndexColumnsFromCsvHeader(parser.reader, Cotisation{})
	}
	return err
}

func (parser *parserCotisation) ReadNext(res *engine.ParsedLineResult) error {
	row, err := parser.reader.Read()
	if err != nil {
		return err
	}

	idxRow := parser.idx.IndexRow(row)

	cotisation := Cotisation{}

	periodeDebut, periodeFin, err := engine.UrssafToPeriod(idxRow.GetVal("periode"))
	res.AddRegularError(err)

	siret, err := engine.GetSiretFromComptesMapping(idxRow.GetVal("Compte"), &periodeDebut, parser.comptes)
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

func (parser *parserCotisation) Close() error {
	return parser.file.Close()
}
