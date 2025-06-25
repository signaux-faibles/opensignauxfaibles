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

// EffectifEnt Urssaf
type EffectifEnt struct {
	Siren       string    `col:"siren" json:"-"        bson:"-"`
	Periode     time.Time `            json:"periode"  bson:"periode"`
	EffectifEnt int       `            json:"effectif" bson:"effectif"`
}

func (effectifEnt EffectifEnt) Headers() []string {
	return []string{
		"siren",
		"période",
		"effectif_entreprise",
	}
}

func (effectifEnt EffectifEnt) Values() []string {
	return []string{
		effectifEnt.Siren,
		marshal.TimeToCSV(&effectifEnt.Periode),
		marshal.IntToCSV(&effectifEnt.EffectifEnt),
	}
}

// Key _id de l'objet
func (effectifEnt EffectifEnt) Key() string {
	return effectifEnt.Siren
}

// Scope de l'objet
func (effectifEnt EffectifEnt) Scope() string {
	return "entreprise"
}

// Type de l'objet
func (effectifEnt EffectifEnt) Type() string {
	return "effectif_ent"
}

// ParserEffectifEnt fournit une instance utilisable par ParseFilesFromBatch.
var ParserEffectifEnt = &effectifEntParser{}

type effectifEntParser struct {
	file    *os.File
	reader  *csv.Reader
	periods []periodCol
	idx     marshal.ColMapping
}

func (parser *effectifEntParser) GetFileType() string {
	return "effectif_ent"
}

func (parser *effectifEntParser) Close() error {
	return parser.file.Close()
}

func (parser *effectifEntParser) Init(_ *marshal.Cache, _ *base.AdminBatch) error {
	return nil
}

func (parser *effectifEntParser) Open(filePath base.BatchFile) (err error) {
	parser.file, parser.reader, err = marshal.OpenCsvReader(filePath, ';', false)
	if err == nil {
		parser.idx, parser.periods, err = parseEffectifColMapping(parser.reader, EffectifEnt{})
	}
	return err
}

func (parser *effectifEntParser) ParseLines(parsedLineChan chan marshal.ParsedLineResult) {
	marshal.ParseLines(parsedLineChan, parser.reader, func(row []string, parsedLine *marshal.ParsedLineResult) {
		parseEffectifEntLine(row, parser.idx, &parser.periods, parsedLine)
	})
}

func parseEffectifEntLine(row []string, idx marshal.ColMapping, periods *[]periodCol, parsedLine *marshal.ParsedLineResult) {
	for _, period := range *periods {
		value := row[period.colIndex] // TODO: utiliser idxRow.GetVal(colName) au lieu de row[colIndex] ?
		if value != "" {
			noThousandsSep := sfregexp.RegexpDict["notDigit"].ReplaceAllString(value, "")
			s, err := strconv.ParseFloat(noThousandsSep, 64)
			parsedLine.AddRegularError(err)
			e := int(s)
			if e > 0 {
				idxRow := idx.IndexRow(row)
				parsedLine.AddTuple(EffectifEnt{
					Siren:       idxRow.GetVal("siren"),
					Periode:     period.dateStart,
					EffectifEnt: e,
				})
			}
		}
	}
}
