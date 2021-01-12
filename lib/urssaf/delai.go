package urssaf

import (
	"bufio"
	"encoding/csv"

	"github.com/signaux-faibles/opensignauxfaibles/lib/base"
	"github.com/signaux-faibles/opensignauxfaibles/lib/marshal"

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
	NumeroCompte      string    `col:"Numero_compte_externe"       json:"numero_compte"      bson:"numero_compte"`
	NumeroContentieux string    `col:"Numero_structure"            json:"numero_contentieux" bson:"numero_contentieux"`
	DateCreation      time.Time `col:"Date_creation"               json:"date_creation"      bson:"date_creation"`
	DateEcheance      time.Time `col:"Date_echeance"               json:"date_echeance"      bson:"date_echeance"`
	DureeDelai        int       `col:"Duree_delai"                 json:"duree_delai"        bson:"duree_delai"`
	Denomination      string    `col:"Denomination_premiere_ligne" json:"denomination"       bson:"denomination"`
	Indic6m           string    `col:"Indic_6M"                    json:"indic_6m"           bson:"indic_6m"`
	AnneeCreation     int       `col:"Annee_creation"              json:"annee_creation"     bson:"annee_creation"`
	MontantEcheancier float64   `col:"Montant_global_echeancier"   json:"montant_echeancier" bson:"montant_echeancier"`
	Stade             string    `col:"Code_externe_stade"          json:"stade"              bson:"stade"`
	Action            string    `col:"Code_externe_action"         json:"action"             bson:"action"`
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
	idx     marshal.ColMapping
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
	parser.file, err = os.Open(filePath)
	if err != nil {
		return err
	}
	parser.reader = csv.NewReader(bufio.NewReader(parser.file))
	parser.reader.Comma = ';'
	header, err := parser.reader.Read()
	if err == nil {
		parser.idx, err = marshal.ValidateAndIndexColumnsFromColTags(header, Delai{})
	}
	return err
}

func (parser *delaiParser) ParseLines(parsedLineChan chan marshal.ParsedLineResult) {
	idx := parser.idx
	for {
		parsedLine := marshal.ParsedLineResult{}
		row, err := parser.reader.Read()
		if err == io.EOF {
			close(parsedLineChan)
			break
		} else if err != nil {
			parsedLine.AddRegularError(err)
		} else {
			date, err := time.Parse("02/01/2006", row[idx["Date_creation"]])
			if err != nil {
				parsedLine.AddRegularError(err)
			} else if siret, err := marshal.GetSiretFromComptesMapping(row[idx["Numero_compte_externe"]], &date, parser.comptes); err == nil {
				parseDelaiLine(row, idx, siret, &parsedLine)
				if len(parsedLine.Errors) > 0 {
					parsedLine.Tuples = []marshal.Tuple{}
				}
			} else {
				parsedLine.SetFilterError(err)
			}
		}
		parsedLineChan <- parsedLine
	}
}

func parseDelaiLine(row []string, idx marshal.ColMapping, siret string, parsedLine *marshal.ParsedLineResult) {
	var err error
	delai := Delai{}
	delai.key = siret
	delai.NumeroCompte = row[idx["Numero_compte_externe"]]
	delai.NumeroContentieux = row[idx["Numero_structure"]]
	delai.DateCreation, err = time.Parse("02/01/2006", row[idx["Date_creation"]])
	parsedLine.AddRegularError(err)
	delai.DateEcheance, err = time.Parse("02/01/2006", row[idx["Date_echeance"]])
	parsedLine.AddRegularError(err)
	delai.DureeDelai, err = strconv.Atoi(row[idx["Duree_delai"]])
	delai.Denomination = row[idx["Denomination_premiere_ligne"]]
	delai.Indic6m = row[idx["Indic_6M"]]
	delai.AnneeCreation, err = strconv.Atoi(row[idx["Annee_creation"]])
	parsedLine.AddRegularError(err)
	delai.MontantEcheancier, err = strconv.ParseFloat(strings.Replace(row[idx["Montant_global_echeancier"]], ",", ".", -1), 64)
	parsedLine.AddRegularError(err)
	delai.Stade = row[idx["Code_externe_stade"]]
	delai.Action = row[idx["Code_externe_action"]]
	parsedLine.AddTuple(delai)
}
