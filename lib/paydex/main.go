// Le Paydex – index de paiement – mesure statistiquement la régularité de
// paiement d’une entreprise vis-à-vis de ses fournisseurs.
// Il est exprimé en nombre de jours de retard de paiement moyen,
// basé sur trois expériences de paiement minimum
// (provenant de trois fournisseurs distincts).

package paydex

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/signaux-faibles/opensignauxfaibles/lib/base"
	"github.com/signaux-faibles/opensignauxfaibles/lib/marshal"
)

// Paydex décrit le format de chaque entrée de donnée résultant du parsing.
type Paydex struct {
	Siren      string    `col:"SIREN" json:"-" bson:"-"`
	DateValeur time.Time `col:"DATE_VALEUR" json:"date_valeur" bson:"date_valeur"`
	NbJours    int       `col:"NB_JOURS" json:"nb_jours" bson:"nb_jours"`
}

// Key _id de l'objet
func (paydex Paydex) Key() string {
	return paydex.Siren
}

// Scope de l'objet
func (paydex Paydex) Scope() string {
	return "entreprise"
}

// Type de l'objet
func (paydex Paydex) Type() string {
	return "paydex"
}

// ParserPaydex fournit une instance utilisable par ParseFilesFromBatch.
var ParserPaydex = &paydexParser{}

type paydexParser struct {
	file     *os.File
	reader   *csv.Reader
	colIndex marshal.ColMapping
}

func (parser *paydexParser) GetFileType() string {
	return "paydex"
}

func (parser *paydexParser) Init(cache *marshal.Cache, batch *base.AdminBatch) error {
	return nil
}

func (parser *paydexParser) Close() error {
	return parser.file.Close()
}

func (parser *paydexParser) Open(filePath string) (err error) {
	parser.file, parser.reader, err = openPaydexFile(filePath)
	if err == nil {
		parser.colIndex, err = marshal.IndexColumnsFromCsvHeader(parser.reader, Paydex{})
	}
	return err
}

func openPaydexFile(filePath string) (*os.File, *csv.Reader, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return file, nil, err
	}
	reader := csv.NewReader(bufio.NewReader(file))
	reader.Comma = ';'
	return file, reader, err
}

func (parser *paydexParser) ParseLines(parsedLineChan chan marshal.ParsedLineResult) {
	for {
		parsedLine := marshal.ParsedLineResult{}
		row, err := parser.reader.Read()
		if err == io.EOF {
			close(parsedLineChan)
			break
		} else if err != nil {
			parsedLine.AddRegularError(err)
		} else {
			paydex, err := parsePaydexLine(parser.colIndex, row)
			if err != nil {
				parsedLine.AddRegularError(err)
			} else {
				parsedLine.AddTuple(paydex)
			}
		}
		parsedLineChan <- parsedLine
	}
}

func parsePaydexLine(colIndex marshal.ColMapping, row []string) (*Paydex, error) {
	dateValeur, err := time.Parse("02/01/2006", row[colIndex["DATE_VALEUR"]])
	if err != nil {
		return nil, fmt.Errorf("invalid date: %v", row[colIndex["DATE_VALEUR"]])
	}
	nbJours, err := strconv.Atoi(row[colIndex["NB_JOURS"]])
	if err != nil {
		return nil, fmt.Errorf("invalid int: %v", row[colIndex["NB_JOURS"]])
	}
	return &Paydex{
		Siren:      row[colIndex["SIREN"]],
		DateValeur: dateValeur,
		NbJours:    nbJours,
	}, nil
}

// TODO: ajouter détection de colonnes
