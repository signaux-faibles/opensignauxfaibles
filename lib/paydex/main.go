// Le Paydex – index de paiement – mesure statistiquement la régularité de
// paiement d’une entreprise vis-à-vis de ses fournisseurs.
// Il est exprimé en nombre de jours de retard de paiement moyen,
// basé sur trois expériences de paiement minimum
// (provenant de trois fournisseurs distincts).

package paydex

import (
	"encoding/csv"
	"fmt"
	"os"
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
	parser.file, parser.reader, err = marshal.OpenCsvReader(base.BatchFile(filePath), ';', false)
	if err == nil {
		parser.colIndex, err = marshal.IndexColumnsFromCsvHeader(parser.reader, Paydex{})
	}
	return err
}

func (parser *paydexParser) ParseLines(parsedLineChan chan marshal.ParsedLineResult) {
	marshal.ParseLines(parsedLineChan, parser.reader, func(row []string, parsedLine *marshal.ParsedLineResult) {
		parser.parseLine(row, parsedLine)
	})
}

func (parser *paydexParser) parseLine(row []string, parsedLine *marshal.ParsedLineResult) {
	paydex, err := parsePaydexLine(parser.colIndex, row)
	if err != nil {
		parsedLine.AddRegularError(err)
	} else {
		parsedLine.AddTuple(paydex)
	}
}

func parsePaydexLine(idx marshal.ColMapping, row []string) (*Paydex, error) {
	idxRow := idx.IndexRow(row)
	dateValeur, err := time.Parse("02/01/2006", idxRow.GetVal("DATE_VALEUR"))
	if err != nil {
		return nil, fmt.Errorf("invalid date: %v", idxRow.GetVal("DATE_VALEUR"))
	}
	nbJours, err := idxRow.GetInt("NB_JOURS")
	if err != nil {
		return nil, fmt.Errorf("invalid int: %v", idxRow.GetVal("NB_JOURS"))
	}
	return &Paydex{
		Siren:      idxRow.GetVal("SIREN"),
		DateValeur: dateValeur,
		NbJours:    *nbJours,
	}, nil
}
