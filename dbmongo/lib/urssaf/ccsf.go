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

// ParserCCSF expose le parseur et le type de fichier qu'il supporte.
var ParserCCSF = marshal.Parser{FileType: "ccsf", FileParser: ParseCcsfFile}

// ParseCcsfFile permet de lancer le parsing du fichier demandé.
func ParseCcsfFile(filePath string, cache *marshal.Cache, batch *base.AdminBatch) marshal.OpenFileResult {
	var comptes marshal.Comptes
	closeFct, reader, err := openCcsfFile(filePath)
	if err != nil {
		comptes, err = marshal.GetCompteSiretMapping(*cache, batch, marshal.OpenAndReadSiretMapping)
	}
	return marshal.OpenFileResult{
		Error: err,
		ParseLines: func(parsedLineChan chan base.ParsedLineResult) {
			parseCcsfLines(reader, &comptes, parsedLineChan)
		},
		Close: closeFct,
	}
}

func openCcsfFile(filePath string) (func() error, *csv.Reader, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return file.Close, nil, err
	}
	reader := csv.NewReader(bufio.NewReader(file))
	reader.Comma = ';'
	_, err = reader.Read() // Sauter l'en-tête
	return file.Close, reader, err
}

var idxCcsf = colMapping{
	"NumeroCompte":   2,
	"DateTraitement": 3,
	"Stade":          4,
	"Action":         5,
}

func parseCcsfLines(reader *csv.Reader, comptes *marshal.Comptes, parsedLineChan chan base.ParsedLineResult) {
	for {
		parsedLine := base.ParsedLineResult{}
		row, err := reader.Read()
		if err == io.EOF {
			close(parsedLineChan)
			break
		} else if err != nil {
			parsedLine.AddError(err)
		} else {
			parseCcsfLine(row, comptes, &parsedLine)
			if len(parsedLine.Errors) > 0 {
				parsedLine.Tuples = []base.Tuple{}
			}
		}
		parsedLineChan <- parsedLine
	}
}

func parseCcsfLine(row []string, comptes *marshal.Comptes, parsedLine *base.ParsedLineResult) {
	idx := idxCcsf
	var err error
	ccsf := CCSF{}
	if len(row) >= 4 {
		ccsf.Action = row[idx["Action"]]
		ccsf.Stade = row[idx["Stade"]]
		ccsf.DateTraitement, err = marshal.UrssafToDate(row[idx["DateTraitement"]])
		parsedLine.AddError(err)
		if err != nil {
			return
		}

		ccsf.key, err = marshal.GetSiretFromComptesMapping(row[idx["NumeroCompte"]], &ccsf.DateTraitement, *comptes)
		if err != nil {
			// Compte filtré
			parsedLine.AddError(base.NewFilterError(err))
			return
		}
		ccsf.NumeroCompte = row[idx["NumeroCompte"]]

	} else {
		parsedLine.AddError(errors.New("Ligne non conforme, moins de 4 champs"))
	}
	parsedLine.AddTuple(ccsf)
}
