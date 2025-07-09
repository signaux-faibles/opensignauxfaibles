package urssaf

import (
	"encoding/csv"
	"os"
	"strconv"
	"time"

	"opensignauxfaibles/lib/base"
	"opensignauxfaibles/lib/marshal"
	"opensignauxfaibles/lib/sfregexp"
)

// Effectif Urssaf
type Effectif struct {
	Siret        string    `input:"siret"  json:"-"`
	NumeroCompte string    `input:"compte" json:"numero_compte"`
	Periode      time.Time `             json:"periode"`
	Effectif     int       `             json:"effectif"`
}

func (effectif Effectif) Headers() []string {
	return []string{
		"siret",
		"compte",
		"période",
		"effectif",
	}
}

func (effectif Effectif) Values() []string {
	return []string{
		effectif.Siret,
		effectif.NumeroCompte,
		marshal.TimeToCSV(&effectif.Periode),
		marshal.IntToCSV(&effectif.Effectif),
	}
}

// Key _id de l'objet
func (effectif Effectif) Key() string {
	return effectif.Siret
}

// Scope de l'objet
func (effectif Effectif) Scope() string {
	return "etablissement"
}

// Type de l'objet
func (effectif Effectif) Type() string {
	return "effectif"
}

// ParserEffectif fournit une instance utilisable par ParseFilesFromBatch.
var ParserEffectif = &effectifParser{}

type effectifParser struct {
	file    *os.File
	reader  *csv.Reader
	periods []periodCol
	idx     marshal.ColMapping
}

func (parser *effectifParser) Type() string {
	return "effectif"
}

func (parser *effectifParser) Init(cache *marshal.Cache, batch *base.AdminBatch) error {
	return nil
}

func (parser *effectifParser) Close() error {
	return parser.file.Close()
}

func (parser *effectifParser) Open(filePath base.BatchFile) (err error) {
	parser.file, parser.reader, err = marshal.OpenCsvReader(filePath, ';', false)
	if err == nil {
		parser.idx, parser.periods, err = parseEffectifColMapping(parser.reader, Effectif{})
	}
	return err
}

func (parser *effectifParser) ParseLines(parsedLineChan chan marshal.ParsedLineResult) {
	marshal.ParseLines(parsedLineChan, parser.reader, func(row []string, parsedLine *marshal.ParsedLineResult) {
		parseEffectifLine(row, parser.idx, &parser.periods, parsedLine)
	})
}

func parseEffectifLine(row []string, idx marshal.ColMapping, periods *[]periodCol, parsedLine *marshal.ParsedLineResult) {
	for _, period := range *periods {
		value := row[period.colIndex] // TODO: utiliser idxRow.GetVal(colName) au lieu de row[colIndex] ?
		if value != "" {
			noThousandsSep := sfregexp.RegexpDict["notDigit"].ReplaceAllString(value, "")
			e, err := strconv.Atoi(noThousandsSep)
			parsedLine.AddRegularError(err)
			if e > 0 {
				idxRow := idx.IndexRow(row)
				parsedLine.AddTuple(Effectif{
					Siret:        idxRow.GetVal("siret"),
					NumeroCompte: idxRow.GetVal("compte"),
					Periode:      period.dateStart,
					Effectif:     e,
				})
			}
		}
	}
}
