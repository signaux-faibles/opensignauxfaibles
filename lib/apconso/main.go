package apconso

import (
	"encoding/csv"
	"os"
	"time"

	"github.com/signaux-faibles/opensignauxfaibles/lib/base"
	"github.com/signaux-faibles/opensignauxfaibles/lib/marshal"
)

// APConso Consommation d'activité partielle
type APConso struct {
	ID             string    `col:"ID_DA"      json:"id_conso"       bson:"id_conso"`
	Siret          string    `col:"ETAB_SIRET" json:"-"              bson:"-"`
	HeureConsommee *float64  `col:"HEURES"     json:"heure_consomme" bson:"heure_consomme"`
	Montant        *float64  `col:"MONTANTS"   json:"montant"        bson:"montant"`
	Effectif       *int      `col:"EFFECTIFS"  json:"effectif"       bson:"effectif"`
	Periode        time.Time `col:"MOIS"       json:"periode"        bson:"periode"`
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

func (parser *apconsoParser) GetFileType() string {
	return "apconso"
}

func (parser *apconsoParser) Init(cache *marshal.Cache, batch *base.AdminBatch) error {
	return nil
}

func (parser *apconsoParser) Open(filePath string) (err error) {
	parser.file, parser.reader, err = marshal.OpenCsvReader(base.BatchFile(filePath), ',', false)
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
	apconso.Periode, err = time.Parse("01/2006", idxRow.GetVal("MOIS"))
	parsedLine.AddRegularError(err)
	apconso.HeureConsommee, err = idxRow.GetFloat64("HEURES")
	parsedLine.AddRegularError(err)
	apconso.Montant, err = idxRow.GetFloat64("MONTANTS")
	parsedLine.AddRegularError(err)
	apconso.Effectif, err = idxRow.GetInt("EFFECTIFS")
	parsedLine.AddRegularError(err)
	if len(parsedLine.Errors) == 0 {
		parsedLine.AddTuple(apconso)
	}
}
