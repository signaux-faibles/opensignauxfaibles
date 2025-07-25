package apconso

import (
	"encoding/csv"
	"os"
	"time"

	"opensignauxfaibles/lib/base"
	"opensignauxfaibles/lib/marshal"
)

// APConso Consommation d'activité partielle
type APConso struct {
	ID             string    `input:"ID_DA"      json:"id_conso"       sql:"id_demande"       csv:"ID"`
	Siret          string    `input:"ETAB_SIRET" json:"-"              sql:"siret"            csv:"Siret"`
	HeureConsommee *float64  `input:"HEURES"     json:"heure_consomme" sql:"heures_consommees"csv:"HeureConsommee"`
	Montant        *float64  `input:"MONTANTS"   json:"montant"        sql:"montant"          csv:"Montant"`
	Effectif       *int      `input:"EFFECTIFS"  json:"effectif"       sql:"effectif"         csv:"Effectif"`
	Periode        time.Time `input:"MOIS"       json:"periode"        sql:"periode"          csv:"Periode"`
}

// Key id de l'objet
func (apconso APConso) Key() string {
	return apconso.Siret
}

// Type de données
func (apconso APConso) Type() string {
	return "apconso"
}

// Scope de l'objet
func (apconso APConso) Scope() string {
	return "etablissement"
}

// Parser fournit une instance utilisable par ParseFilesFromBatch.
var Parser = &apconsoParser{}

type apconsoParser struct {
	file   *os.File
	reader *csv.Reader
	idx    marshal.ColMapping
}

func (parser *apconsoParser) Type() string {
	return "apconso"
}

func (parser *apconsoParser) Init(_ *marshal.Cache, _ *base.AdminBatch) error {
	return nil
}

func (parser *apconsoParser) Open(filePath base.BatchFile) (err error) {
	parser.file, parser.reader, err = marshal.OpenCsvReader(filePath, ',', false)
	if err == nil {
		parser.idx, err = marshal.IndexColumnsFromCsvHeader(parser.reader, APConso{})
	}
	return err
}

func (parser *apconsoParser) Close() error {
	return parser.file.Close()
}

func (parser *apconsoParser) ParseLines(parsedLineChan chan marshal.ParsedLineResult) {
	marshal.ParseLines(parsedLineChan, parser.reader, func(row []string, parsedLine *marshal.ParsedLineResult) {
		parseApConsoLine(row, parser.idx, parsedLine)
	})
}

func parseApConsoLine(row []string, idx marshal.ColMapping, parsedLine *marshal.ParsedLineResult) {
	idxRow := idx.IndexRow(row)
	apconso := APConso{}
	apconso.ID = idxRow.GetVal("ID_DA")
	apconso.Siret = idxRow.GetVal("ETAB_SIRET")
	var err error
	apconso.Periode, err = time.Parse("2006-01-02", idxRow.GetVal("MOIS"))
	parsedLine.AddRegularError(err)
	apconso.HeureConsommee, err = idxRow.GetFloat64("HEURES")
	parsedLine.AddRegularError(err)
	apconso.Montant, err = idxRow.GetFloat64("MONTANTS")
	parsedLine.AddRegularError(err)
	apconso.Effectif, err = idxRow.GetIntFromFloat("EFFECTIFS")
	parsedLine.AddRegularError(err)
	if len(parsedLine.Errors) == 0 {
		parsedLine.AddTuple(apconso)
	}
}
