package urssaf

import (
	"encoding/csv"
	"opensignauxfaibles/lib/base"
	"opensignauxfaibles/lib/engine"
	"os"
	"time"
)

// Type de l'objet
func (delai Delai) Type() base.ParserType {
	return base.Delai
}

// Implémente engine.Parser
type parserDelai struct {
	reader  *csv.Reader
	file    *os.File
	comptes engine.Comptes
	idx     engine.ColMapping
}

func NewParserDelai() *parserDelai {
	return &parserDelai{}
}

func (parser *parserDelai) Type() base.ParserType {
	return base.Delai
}

func (parser *parserDelai) Init(cache *engine.Cache, filter engine.SirenFilter, batch *base.AdminBatch) (err error) {
	parser.comptes, err = engine.GetCompteSiretMapping(*cache, batch, filter, engine.OpenAndReadSiretMapping)
	return err
}

func (parser *parserDelai) Open(filePath base.BatchFile) (err error) {
	parser.file, parser.reader, err = engine.OpenCsvReader(filePath, ';', false)

	if err == nil {
		parser.idx, err = engine.IndexColumnsFromCsvHeader(parser.reader, Delai{})
	}
	return err
}

func (parser *parserDelai) ReadNext(res *engine.ParsedLineResult) error {

	row, err := parser.reader.Read()

	if err != nil {
		return err
	}

	idxRow := parser.idx.IndexRow(row)

	date, err := time.Parse("02/01/2006", idxRow.GetVal("Date_creation"))

	if err != nil {
		res.AddRegularError(err)
		return err

	}

	siret, err := engine.GetSiretFromComptesMapping(idxRow.GetVal("Numero_compte_externe"), &date, parser.comptes)
	if err != nil {
		res.SetFilterError(err)
		return err
	}

	delai := Delai{}
	delai.Siret = siret
	delai.NumeroCompte = idxRow.GetVal("Numero_compte_externe")
	delai.NumeroContentieux = idxRow.GetVal("Numero_structure")
	delai.DateCreation, err = time.Parse("02/01/2006", idxRow.GetVal("Date_creation"))
	res.AddRegularError(err)
	delai.DateEcheance, err = time.Parse("02/01/2006", idxRow.GetVal("Date_echeance"))
	res.AddRegularError(err)
	delai.DureeDelai, err = idxRow.GetInt("Duree_delai")
	delai.Denomination = idxRow.GetVal("Denomination_premiere_ligne")
	delai.Indic6m = idxRow.GetVal("Indic_6M")
	delai.AnneeCreation, err = idxRow.GetInt("Annee_creation")
	res.AddRegularError(err)
	delai.MontantEcheancier, err = idxRow.GetCommaFloat64("Montant_global_echeancier")
	res.AddRegularError(err)
	delai.Stade = idxRow.GetVal("Code_externe_stade")
	delai.Action = idxRow.GetVal("Code_externe_action")

	if len(res.Errors) == 0 {
		res.AddTuple(delai)
	}
	return nil
}

func (parser *parserDelai) Close() error {
	return parser.file.Close()
}
