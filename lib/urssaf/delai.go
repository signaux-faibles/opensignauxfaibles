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
	key               string    `hash:"-"`
	NumeroCompte      string    `col:"Numero_compte_externe"       json:"numero_compte"      bson:"numero_compte"`
	NumeroContentieux string    `col:"Numero_structure"            json:"numero_contentieux" bson:"numero_contentieux"`
	DateCreation      time.Time `col:"Date_creation"               json:"date_creation"      bson:"date_creation"`
	DateEcheance      time.Time `col:"Date_echeance"               json:"date_echeance"      bson:"date_echeance"`
	DureeDelai        *int      `col:"Duree_delai"                 json:"duree_delai"        bson:"duree_delai"`
	Denomination      string    `col:"Denomination_premiere_ligne" json:"denomination"       bson:"denomination"`
	Indic6m           string    `col:"Indic_6M"                    json:"indic_6m"           bson:"indic_6m"`
	AnneeCreation     *int      `col:"Annee_creation"              json:"annee_creation"     bson:"annee_creation"`
	MontantEcheancier *float64  `col:"Montant_global_echeancier"   json:"montant_echeancier" bson:"montant_echeancier"`
	Stade             string    `col:"Code_externe_stade"          json:"stade"              bson:"stade"`
	Action            string    `col:"Code_externe_action"         json:"action"             bson:"action"`
}

func (delai Delai) Headers() []string {
	r := []string{"key"}
	r = append(r, marshal.ExtractColTags(delai)...)
	return r
}

func (delai Delai) Values() []string {
	r := []string{
		delai.key,
		delai.NumeroCompte,
		delai.NumeroContentieux,
		marshal.TimeToCSV(&delai.DateCreation),
		marshal.TimeToCSV(&delai.DateEcheance),
		marshal.IntToCSV(delai.DureeDelai),
		delai.Denomination,
		delai.Indic6m,
		marshal.IntToCSV(delai.AnneeCreation),
		marshal.FloatToCSV(delai.MontantEcheancier),
		delai.Stade,
		delai.Action,
	}
	return r
}

// Key _id de l'objet
func (delai Delai) Key() string {
	return delai.key
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

func (parser *delaiParser) GetFileType() string {
	return "delai"
}

func (parser *delaiParser) Close() error {
	return parser.file.Close()
}

func (parser *delaiParser) Init(cache *marshal.Cache, batch *base.AdminBatch) (err error) {
	parser.comptes, err = marshal.GetCompteSiretMapping(*cache, batch, marshal.OpenAndReadSiretMapping)
	return err
}

func (parser *delaiParser) Open(filePath string) (err error) {
	parser.file, parser.reader, err = marshal.OpenCsvReader(base.BatchFile(filePath), ';', false)
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
	delai.key = siret
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
