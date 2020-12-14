package urssaf

import (
	"bufio"
	"encoding/csv"
	"errors"
	"io"
	"os"
	"time"

	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/base"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/marshal"
)

// CCSF information urssaf ccsf
type CCSF struct {
	key            string    `hash:"-"`
	NumeroCompte   string    `json:"-" bson:"-"`
	DateTraitement time.Time `json:"date_traitement" bson:"date_traitement"`
	Stade          string    `json:"stade" bson:"stade"`
	Action         string    `json:"action" bson:"action"`
}

// Key _id de l'objet
func (ccsf CCSF) Key() string {
	return ccsf.key
}

// Scope de l'objet
func (ccsf CCSF) Scope() string {
	return "etablissement"
}

// Type de l'objet
func (ccsf CCSF) Type() string {
	return "ccsf"
}

// ParserCCSF fournit une instance utilisable par ParseFilesFromBatch.
var ParserCCSF = &ccsfParser{}

type ccsfParser struct {
	file    *os.File
	reader  *csv.Reader
	comptes marshal.Comptes
}

func (parser *ccsfParser) GetFileType() string {
	return "ccsf"
}

func (parser *ccsfParser) Init(cache *marshal.Cache, batch *base.AdminBatch) (err error) {
	parser.comptes, err = marshal.GetCompteSiretMapping(*cache, batch, marshal.OpenAndReadSiretMapping)
	return err
}

func (parser *ccsfParser) Close() error {
	return parser.file.Close()
}

func (parser *ccsfParser) Open(filePath string) (err error) {
	parser.file, err = os.Open(filePath)
	if err != nil {
		return err
	}
	parser.reader = csv.NewReader(bufio.NewReader(parser.file))
	parser.reader.Comma = ';'
	_, err = parser.reader.Read() // Sauter l'en-tête
	return err
}

var idxCcsf = colMapping{
	"NumeroCompte":   2,
	"DateTraitement": 3,
	"Stade":          4,
	"Action":         5,
}

func (parser *ccsfParser) ParseLines(parsedLineChan chan marshal.ParsedLineResult) {
	for {
		parsedLine := marshal.ParsedLineResult{}
		row, err := parser.reader.Read()
		if err == io.EOF {
			close(parsedLineChan)
			break
		} else if err != nil {
			parsedLine.AddRegularError(err)
		} else {
			parseCcsfLine(row, &parser.comptes, &parsedLine)
			if len(parsedLine.Errors) > 0 {
				parsedLine.Tuples = []marshal.Tuple{}
			}
		}
		parsedLineChan <- parsedLine
	}
}

func parseCcsfLine(row []string, comptes *marshal.Comptes, parsedLine *marshal.ParsedLineResult) {
	idx := idxCcsf
	var err error
	ccsf := CCSF{}
	if len(row) >= 4 {
		ccsf.Action = row[idx["Action"]]
		ccsf.Stade = row[idx["Stade"]]
		ccsf.DateTraitement, err = marshal.UrssafToDate(row[idx["DateTraitement"]])
		parsedLine.AddRegularError(err)
		if err != nil {
			return
		}

		ccsf.key, err = marshal.GetSiretFromComptesMapping(row[idx["NumeroCompte"]], &ccsf.DateTraitement, *comptes)
		if err != nil {
			// Compte filtré
			parsedLine.SetFilterError(err)
			return
		}
		ccsf.NumeroCompte = row[idx["NumeroCompte"]]

	} else {
		parsedLine.AddRegularError(errors.New("Ligne non conforme, moins de 4 champs"))
	}
	parsedLine.AddTuple(ccsf)
}
