package urssaf

import (
	"bufio"
	"encoding/csv"

	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/base"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/marshal"

	//"errors"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/signaux-faibles/gournal"
)

// Delai tuple fichier ursaff
type Delai struct {
	key               string    `hash:"-"`
	NumeroCompte      string    `json:"numero_compte" bson:"numero_compte"`
	NumeroContentieux string    `json:"numero_contentieux" bson:"numero_contentieux"`
	DateCreation      time.Time `json:"date_creation" bson:"date_creation"`
	DateEcheance      time.Time `json:"date_echeance" bson:"date_echeance"`
	DureeDelai        int       `json:"duree_delai" bson:"duree_delai"`
	Denomination      string    `json:"denomination" bson:"denomination"`
	Indic6m           string    `json:"indic_6m" bson:"indic_6m"`
	AnneeCreation     int       `json:"annee_creation" bson:"annee_creation"`
	MontantEcheancier float64   `json:"montant_echeancier" bson:"montant_echeancier"`
	Stade             string    `json:"stade" bson:"stade"`
	Action            string    `json:"action" bson:"action"`
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

// ParserDelai expose le parseur et le type de fichier qu'il supporte.
var ParserDelai = marshal.Parser{FileType: "delai", FileParser: ParseDelaiFile}

// ParseDelaiFile extrait les tuples depuis le fichier demandé et génère un rapport Gournal.
func ParseDelaiFile(filePath string, cache *marshal.Cache, batch *base.AdminBatch, tracker *gournal.Tracker) marshal.ParsedLineChan {
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

	var idx = colMapping{
		"NumeroCompte":      2,
		"NumeroContentieux": 3,
		"DateCreation":      4,
		"DateEcheance":      5,
		"DureeDelai":        6,
		"Denomination":      7,
		"Indic6m":           8,
		"AnneeCreation":     9,
		"MontantEcheancier": 10,
		"Stade":             11,
		"Action":            12,
	}

	reader.Read()

	parsedLineChan := make(marshal.ParsedLineChan)
	go func() {
		for {
			parsedLine := marshal.ParsedLineResult{}
			row, err := reader.Read()
			if err == io.EOF {
				close(parsedLineChan)
				break
			} else if err != nil {
				parsedLine.AddError(err)
			} else {
				date, err := time.Parse("02/01/2006", row[idx["DateCreation"]])
				if err != nil {
					parsedLine.AddError(err)
				} else if siret, err := marshal.GetSiretFromComptesMapping(row[idx["NumeroCompte"]], &date, comptes); err == nil {
					parseDelaiLine(row, idx, siret, &parsedLine)
					if tracker.HasErrorInCurrentCycle() {
						parsedLine.Tuples = []marshal.Tuple{}
					}
				} else {
					parsedLine.AddError(base.NewFilterError(err))
				}
			}
			parsedLineChan <- parsedLine
		}
	}()
	return parsedLineChan
}

func parseDelaiLine(row []string, idx colMapping, siret string, parsedLine *marshal.ParsedLineResult) {
	var err error
	loc, _ := time.LoadLocation("Europe/Paris")
	delai := Delai{}
	delai.key = siret
	delai.NumeroCompte = row[idx["NumeroCompte"]]
	delai.NumeroContentieux = row[idx["NumeroContentieux"]]
	delai.DateCreation, err = time.ParseInLocation("02/01/2006", row[idx["DateCreation"]], loc)
	parsedLine.AddError(err)
	delai.DateEcheance, err = time.ParseInLocation("02/01/2006", row[idx["DateEcheance"]], loc)
	parsedLine.AddError(err)
	delai.DureeDelai, err = strconv.Atoi(row[idx["DureeDelai"]])
	delai.Denomination = row[idx["Denomination"]]
	delai.Indic6m = row[idx["Indic6m"]]
	delai.AnneeCreation, err = strconv.Atoi(row[idx["AnneeCreation"]])
	parsedLine.AddError(err)
	delai.MontantEcheancier, err = strconv.ParseFloat(strings.Replace(row[idx["MontantEcheancier"]], ",", ".", -1), 64)
	parsedLine.AddError(err)
	delai.Stade = row[idx["Stade"]]
	delai.Action = row[idx["Action"]]
	parsedLine.AddTuple(delai)
}
