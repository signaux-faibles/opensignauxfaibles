package urssaf

import (
	"encoding/csv"
	"os"
	"time"

	"opensignauxfaibles/lib/base"
	"opensignauxfaibles/lib/marshal"
)

// Delai tuple fichier ursaff
type Delai struct {
	Siret             string    `                                    json:"-"                  sql:"siret"                csv:"siret"`
	NumeroCompte      string    `input:"Numero_compte_externe"       json:"numero_compte"                                 csv:"numéro_compte"`
	NumeroContentieux string    `input:"Numero_structure"            json:"numero_contentieux"                            csv:"numéro_contentieux"`
	DateCreation      time.Time `input:"Date_creation"               json:"date_creation"      sql:"date_creation"        csv:"date_création"`
	DateEcheance      time.Time `input:"Date_echeance"               json:"date_echeance"      sql:"date_echeance"        csv:"date_échéance"`
	DureeDelai        *int      `input:"Duree_delai"                 json:"duree_delai"        sql:"duree_delai"          csv:"durée_délai"`
	Denomination      string    `input:"Denomination_premiere_ligne" json:"denomination"                                  csv:"dénomination"`
	Indic6m           string    `input:"Indic_6M"                    json:"indic_6m"                                      csv:"indic_6mois"`
	AnneeCreation     *int      `input:"Annee_creation"              json:"annee_creation"                                csv:"année_création"`
	MontantEcheancier *float64  `input:"Montant_global_echeancier"   json:"montant_echeancier" sql:"montant_echeancier"   csv:"montant_échéancier"`
	Stade             string    `input:"Code_externe_stade"          json:"stade"              sql:"stade"                csv:"stade"`
	Action            string    `input:"Code_externe_action"         json:"action"             sql:"action"               csv:"action"`
}

// Key _id de l'objet
func (delai Delai) Key() string {
	return delai.Siret
}

// Scope de l'objet
func (delai Delai) Scope() string {
	return "etablissement"
}

// Type de l'objet
func (delai Delai) Type() string {
	return "delai"
}

// ParserDelai fournit une instance utilisable par ParseFilesFromBatch.
var ParserDelai = &delaiParser{}

type delaiParser struct {
	file    *os.File
	reader  *csv.Reader
	comptes marshal.Comptes
	idx     marshal.ColMapping
}

func (parser *delaiParser) Type() string {
	return "delai"
}

func (parser *delaiParser) Close() error {
	return parser.file.Close()
}

func (parser *delaiParser) Init(cache *marshal.Cache, batch *base.AdminBatch) (err error) {
	parser.comptes, err = marshal.GetCompteSiretMapping(*cache, batch, marshal.OpenAndReadSiretMapping)
	return err
}

func (parser *delaiParser) Open(filePath base.BatchFile) (err error) {
	parser.file, parser.reader, err = marshal.OpenCsvReader(filePath, ';', false)
	if err == nil {
		parser.idx, err = marshal.IndexColumnsFromCsvHeader(parser.reader, Delai{})
	}
	return err
}

func (parser *delaiParser) ParseLines(parsedLineChan chan marshal.ParsedLineResult) {
	marshal.ParseLines(parsedLineChan, parser.reader, func(row []string, parsedLine *marshal.ParsedLineResult) {
		parser.parseDelaiLine(row, parsedLine)
	})
}

func (parser *delaiParser) parseDelaiLine(row []string, parsedLine *marshal.ParsedLineResult) {
	idxRow := parser.idx.IndexRow(row)
	date, err := time.Parse("02/01/2006", idxRow.GetVal("Date_creation"))
	if err != nil {
		parsedLine.AddRegularError(err)
	} else if siret, err := marshal.GetSiretFromComptesMapping(idxRow.GetVal("Numero_compte_externe"), &date, parser.comptes); err == nil {
		parseDelaiLine(idxRow, siret, parsedLine)
	} else {
		parsedLine.SetFilterError(err)
	}
}

func parseDelaiLine(idxRow marshal.IndexedRow, siret string, parsedLine *marshal.ParsedLineResult) {
	var err error
	delai := Delai{}
	delai.Siret = siret
	delai.NumeroCompte = idxRow.GetVal("Numero_compte_externe")
	delai.NumeroContentieux = idxRow.GetVal("Numero_structure")
	delai.DateCreation, err = time.Parse("02/01/2006", idxRow.GetVal("Date_creation"))
	parsedLine.AddRegularError(err)
	delai.DateEcheance, err = time.Parse("02/01/2006", idxRow.GetVal("Date_echeance"))
	parsedLine.AddRegularError(err)
	delai.DureeDelai, err = idxRow.GetInt("Duree_delai")
	delai.Denomination = idxRow.GetVal("Denomination_premiere_ligne")
	delai.Indic6m = idxRow.GetVal("Indic_6M")
	delai.AnneeCreation, err = idxRow.GetInt("Annee_creation")
	parsedLine.AddRegularError(err)
	delai.MontantEcheancier, err = idxRow.GetCommaFloat64("Montant_global_echeancier")
	parsedLine.AddRegularError(err)
	delai.Stade = idxRow.GetVal("Code_externe_stade")
	delai.Action = idxRow.GetVal("Code_externe_action")
	if len(parsedLine.Errors) == 0 {
		parsedLine.AddTuple(delai)
	}
}
