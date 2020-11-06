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

// ParserDelai fournit une instance utilisable par ParseFilesFromBatch.
var ParserDelai = &delaiParser{}

type delaiParser struct {
	file    *os.File
	reader  *csv.Reader
	comptes marshal.Comptes
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
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	reader := csv.NewReader(bufio.NewReader(file))
	reader.Comma = ';'
	_, err = reader.Read() // Sauter l'en-tête
	return err
}

var idxDelai = colMapping{
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

func (parser *delaiParser) ParseLines(parsedLineChan chan marshal.ParsedLineResult) {
	idx := idxDelai
	for {
		parsedLine := marshal.ParsedLineResult{}
		row, err := parser.reader.Read()
		if err == io.EOF {
			close(parsedLineChan)
			break
		} else if err != nil {
			parsedLine.AddError(base.NewRegularError(err))
		} else {
			date, err := time.Parse("02/01/2006", row[idx["DateCreation"]])
			if err != nil {
				parsedLine.AddError(base.NewRegularError(err))
			} else if siret, err := marshal.GetSiretFromComptesMapping(row[idx["NumeroCompte"]], &date, parser.comptes); err == nil {
				parseDelaiLine(row, idx, siret, &parsedLine)
				if len(parsedLine.Errors) > 0 {
					parsedLine.Tuples = []marshal.Tuple{}
				}
			} else {
				parsedLine.AddError(base.NewFilterError(err))
			}
		}
		parsedLineChan <- parsedLine
	}
}

func parseDelaiLine(row []string, idx colMapping, siret string, parsedLine *marshal.ParsedLineResult) {
	var err error
	loc, _ := time.LoadLocation("Europe/Paris")
	delai := Delai{}
	delai.key = siret
	delai.NumeroCompte = row[idx["NumeroCompte"]]
	delai.NumeroContentieux = row[idx["NumeroContentieux"]]
	delai.DateCreation, err = time.ParseInLocation("02/01/2006", row[idx["DateCreation"]], loc)
	parsedLine.AddError(base.NewRegularError(err))
	delai.DateEcheance, err = time.ParseInLocation("02/01/2006", row[idx["DateEcheance"]], loc)
	parsedLine.AddError(base.NewRegularError(err))
	delai.DureeDelai, err = strconv.Atoi(row[idx["DureeDelai"]])
	delai.Denomination = row[idx["Denomination"]]
	delai.Indic6m = row[idx["Indic6m"]]
	delai.AnneeCreation, err = strconv.Atoi(row[idx["AnneeCreation"]])
	parsedLine.AddError(base.NewRegularError(err))
	delai.MontantEcheancier, err = strconv.ParseFloat(strings.Replace(row[idx["MontantEcheancier"]], ",", ".", -1), 64)
	parsedLine.AddError(base.NewRegularError(err))
	delai.Stade = row[idx["Stade"]]
	delai.Action = row[idx["Action"]]
	parsedLine.AddTuple(delai)
}
