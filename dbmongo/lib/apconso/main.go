package apconso

import (
	"encoding/csv"
	"errors"
	"io"
	"os"
	"time"

	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/base"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/marshal"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/misc"
)

// APConso Consommation d'activité partielle
type APConso struct {
	ID             string    `json:"id_conso"         bson:"id_conso"`
	Siret          string    `json:"-"                bson:"-"`
	HeureConsommee *float64  `json:"heure_consomme"   bson:"heure_consomme"`
	Montant        *float64  `json:"montant"          bson:"montant"`
	Effectif       *int      `json:"effectif"         bson:"effectif"`
	Periode        time.Time `json:"periode"          bson:"periode"`
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

// Parser expose le parseur et le type de fichier qu'il supporte.
var Parser = marshal.Parser{FileType: "apconso", FileParser: ParseFile}

// ParseFile permet de lancer le parsing du fichier demandé.
func ParseFile(filePath string, cache *marshal.Cache, batch *base.AdminBatch) (marshal.FileReader, error) {
	var idx colMapping
	file, reader, err := openFile(filePath)
	if err == nil {
		idx, err = parseColMapping(reader)
	}
	return apconsoReader{
		file:   file,
		reader: reader,
		idx:    idx,
	}, err
}

type apconsoReader struct {
	file   *os.File
	reader *csv.Reader
	idx    colMapping
}

func (parser apconsoReader) Close() error {
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

type colMapping map[string]int

func parseColMapping(reader *csv.Reader) (colMapping, error) {
	header, err := reader.Read()
	if err != nil {
		return nil, err
	}
	var idx = colMapping{}
	idx["ID"] = misc.SliceIndex(len(header), func(i int) bool { return header[i] == "ID_DA" })
	idx["Siret"] = misc.SliceIndex(len(header), func(i int) bool { return header[i] == "ETAB_SIRET" })
	idx["Periode"] = misc.SliceIndex(len(header), func(i int) bool { return header[i] == "MOIS" })
	idx["HeureConsommee"] = misc.SliceIndex(len(header), func(i int) bool { return header[i] == "HEURES" })
	idx["Montant"] = misc.SliceIndex(len(header), func(i int) bool { return header[i] == "MONTANTS" })
	idx["Effectif"] = misc.SliceIndex(len(header), func(i int) bool { return header[i] == "EFFECTIFS" })
	if misc.SliceMin(idx["ID"], idx["Siret"], idx["Periode"], idx["HeureConsommee"], idx["Montant"], idx["Effectif"]) == -1 {
		return nil, errors.New("entête non conforme, fichier ignoré")
	}
	return idx, nil
}

func (parser apconsoReader) ParseLines(parsedLineChan chan marshal.ParsedLineResult) {
	for {
		parsedLine := marshal.ParsedLineResult{}
		row, err := parser.reader.Read()
		if err == io.EOF {
			close(parsedLineChan)
			break
		} else if err != nil {
			parsedLine.AddError(base.NewRegularError(err))
		} else if len(row) > 0 {
			parseApConsoLine(row, parser.idx, &parsedLine)
			if len(parsedLine.Errors) > 0 {
				parsedLine.Tuples = []marshal.Tuple{}
			}
		}
		parsedLineChan <- parsedLine
	}
}

func parseApConsoLine(row []string, idx colMapping, parsedLine *marshal.ParsedLineResult) {
	apconso := APConso{}
	apconso.ID = row[idx["ID"]]
	apconso.Siret = row[idx["Siret"]]
	var err error
	apconso.Periode, err = time.Parse("01/2006", row[idx["Periode"]])
	parsedLine.AddError(base.NewRegularError(err))
	apconso.HeureConsommee, err = misc.ParsePFloat(row[idx["HeureConsommee"]])
	parsedLine.AddError(base.NewRegularError(err))
	apconso.Montant, err = misc.ParsePFloat(row[idx["Montant"]])
	parsedLine.AddError(base.NewRegularError(err))
	apconso.Effectif, err = misc.ParsePInt(row[idx["Effectif"]])
	parsedLine.AddError(base.NewRegularError(err))
	parsedLine.AddTuple(apconso)
}
