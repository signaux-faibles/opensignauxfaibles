package apconso

import (
	"encoding/csv"
	"io"
	"os"
	"time"

	"github.com/signaux-faibles/opensignauxfaibles/lib/base"
	"github.com/signaux-faibles/opensignauxfaibles/lib/marshal"
	"github.com/signaux-faibles/opensignauxfaibles/lib/misc"
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
	parser.file, parser.reader, err = openFile(filePath)
	if err == nil {
		parser.idx, err = parseColMapping(parser.reader)
	}
	return err
}

func (parser *apconsoParser) Close() error {
	return parser.file.Close()
}

func openFile(filePath string) (*os.File, *csv.Reader, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return file, nil, err
	}
	reader := csv.NewReader(file)
	reader.Comma = ','
	return file, reader, nil
}

func parseColMapping(reader *csv.Reader) (marshal.ColMapping, error) {
	header, err := reader.Read()
	if err != nil {
		return nil, err
	}
	return marshal.ValidateAndIndexColumnsFromColTags(header, APConso{})
}

func (parser *apconsoParser) ParseLines(parsedLineChan chan marshal.ParsedLineResult) {
	for {
		parsedLine := marshal.ParsedLineResult{}
		row, err := parser.reader.Read()
		if err == io.EOF {
			close(parsedLineChan)
			break
		} else if err != nil {
			parsedLine.AddRegularError(err)
		} else if len(row) > 0 {
			parseApConsoLine(row, parser.idx, &parsedLine)
			if len(parsedLine.Errors) > 0 {
				parsedLine.Tuples = []marshal.Tuple{}
			}
		}
		parsedLineChan <- parsedLine
	}
}

func parseApConsoLine(row []string, idx marshal.ColMapping, parsedLine *marshal.ParsedLineResult) {
	apconso := APConso{}
	apconso.ID = row[idx["ID_DA"]]
	apconso.Siret = row[idx["ETAB_SIRET"]]
	var err error
	apconso.Periode, err = time.Parse("01/2006", row[idx["MOIS"]])
	parsedLine.AddRegularError(err)
	apconso.HeureConsommee, err = misc.ParsePFloat(row[idx["HEURES"]])
	parsedLine.AddRegularError(err)
	apconso.Montant, err = misc.ParsePFloat(row[idx["MONTANTS"]])
	parsedLine.AddRegularError(err)
	apconso.Effectif, err = misc.ParsePInt(row[idx["EFFECTIFS"]])
	parsedLine.AddRegularError(err)
	parsedLine.AddTuple(apconso)
}
