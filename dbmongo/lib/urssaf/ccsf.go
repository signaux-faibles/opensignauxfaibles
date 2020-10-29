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

// ParseCcsfFile extrait les tuples depuis le fichier demandé et génère un rapport Gournal.
func ParseCcsfFile(filePath string, cache *marshal.Cache, batch *base.AdminBatch) (marshal.ParsedLineChan, error) {
	comptes, err := marshal.GetCompteSiretMapping(*cache, batch, marshal.OpenAndReadSiretMapping)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	// defer file.Close() // TODO: à réactiver

	reader := csv.NewReader(bufio.NewReader(file))
	reader.Comma = ';'

	reader.Read() // en-tête du fichier

	var idx = colMapping{
		"NumeroCompte":   2,
		"DateTraitement": 3,
		"Stade":          4,
		"Action":         5,
	}

	parsedLineChan := make(marshal.ParsedLineChan)
	go func() {
		for {
			parsedLine := base.ParsedLineResult{}
			row, err := reader.Read()
			if err == io.EOF {
				close(parsedLineChan)
				break
			} else if err != nil {
				parsedLine.AddError(err)
			} else {
				parseCcsfLine(row, idx, &comptes, &parsedLine)
				if len(parsedLine.Errors) > 0 {
					parsedLine.Tuples = []base.Tuple{}
				}
			}
			parsedLineChan <- parsedLine
		}
	}()
	return parsedLineChan, nil
}

func parseCcsfLine(row []string, idx colMapping, comptes *marshal.Comptes, parsedLine *base.ParsedLineResult) {
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
