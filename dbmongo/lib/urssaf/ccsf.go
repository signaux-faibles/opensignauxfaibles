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

	"github.com/signaux-faibles/gournal"
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
func ParseCcsfFile(filePath string, cache *marshal.Cache, batch *base.AdminBatch, tracker *gournal.Tracker) marshal.TupleGenerator {
	comptes, err := marshal.GetCompteSiretMapping(*cache, batch, marshal.OpenAndReadSiretMapping)
	if err != nil {
		tracker.Add(err)
		return nil
	}

	file, err := os.Open(filePath)
	if err != nil {
		tracker.Add(err)
		return nil
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

	tupleGenerator := make(marshal.TupleGenerator)
	go func() {
		for {
			tuples := []marshal.Tuple{}

			row, err := reader.Read()
			if err == io.EOF {
				close(tupleGenerator)
				break
			} else if err != nil {
				tracker.Add(err)
			} else {
				ccsf := parseCcsfLine(row, tracker, &comptes, idx)
				if !tracker.HasErrorInCurrentCycle() {
					tuples = []marshal.Tuple{ccsf}
				}
			}
			tupleGenerator <- tuples
		}
	}()
	return tupleGenerator
}

func parseCcsfLine(row []string, tracker *gournal.Tracker, comptes *marshal.Comptes, idx colMapping) CCSF {
	var err error
	ccsf := CCSF{}
	if len(row) >= 4 {
		ccsf.Action = row[idx["Action"]]
		ccsf.Stade = row[idx["Stade"]]
		ccsf.DateTraitement, err = marshal.UrssafToDate(row[idx["DateTraitement"]])
		tracker.Add(err)
		if err != nil {
			return ccsf
		}

		ccsf.key, err = marshal.GetSiretFromComptesMapping(row[idx["NumeroCompte"]], &ccsf.DateTraitement, *comptes)
		if err != nil {
			// Compte filtré
			tracker.Add(base.NewFilterError(err))
			return ccsf
		}
		ccsf.NumeroCompte = row[idx["NumeroCompte"]]

	} else {
		tracker.Add(errors.New("Ligne non conforme, moins de 4 champs"))
	}
	return ccsf
}
