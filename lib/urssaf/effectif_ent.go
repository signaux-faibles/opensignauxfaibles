package urssaf

import (
	"bufio"
	"encoding/csv"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/signaux-faibles/opensignauxfaibles/lib/base"
	"github.com/signaux-faibles/opensignauxfaibles/lib/marshal"
	"github.com/signaux-faibles/opensignauxfaibles/lib/misc"
	"github.com/signaux-faibles/opensignauxfaibles/lib/sfregexp"
)

// EffectifEnt Urssaf
type EffectifEnt struct {
	Siren       string    `json:"-" bson:"-"`
	Periode     time.Time `json:"periode" bson:"periode"`
	EffectifEnt int       `json:"effectif" bson:"effectif"`
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
	idx     colMapping
}

func (parser *effectifEntParser) GetFileType() string {
	return "effectif_ent"
}

func (parser *effectifEntParser) Close() error {
	return parser.file.Close()
}

func (parser *effectifEntParser) Init(cache *marshal.Cache, batch *base.AdminBatch) error {
	return nil
}

func (parser *effectifEntParser) Open(filePath string) (err error) {
	parser.file, err = os.Open(filePath)
	if err == nil {
		parser.reader = csv.NewReader(bufio.NewReader(parser.file))
		parser.reader.Comma = ';'
	}
	if err == nil {
		parser.idx, parser.periods, err = parseEffectifEntColMapping(parser.reader)
	}
	return err
}

func (parser *effectifEntParser) ParseLines(parsedLineChan chan marshal.ParsedLineResult) {
	for {
		parsedLine := marshal.ParsedLineResult{}
		row, err := parser.reader.Read()
		if err == io.EOF {
			close(parsedLineChan)
			break
		} else if err != nil {
			parsedLine.AddRegularError(err)
		} else {
			parseEffectifEntLine(row, parser.idx, &parser.periods, &parsedLine)
		}
		parsedLineChan <- parsedLine
	}
}

func parseEffectifEntColMapping(reader *csv.Reader) (colMapping, []periodCol, error) {
	fields, err := reader.Read()
	if err != nil {
		return nil, nil, err
	}
	var idx = colMapping{
		"siren": misc.SliceIndex(len(fields), func(i int) bool { return strings.ToLower(fields[i]) == "siren" }),
	}
	// Dans quels champs lire l'effectifEnt
	periods := parseEffectifPeriod(fields)
	return idx, periods, nil
}

type periodCol struct {
	dateStart time.Time
	colIndex  int
}

// ParseEffectifPeriod extrait les périodes depuis une liste de noms de colonnes csv.
func parseEffectifPeriod(fields []string) []periodCol {
	periods := []periodCol{}
	re, _ := regexp.Compile("^eff")
	for index, field := range fields {
		if re.MatchString(field) {
			date, _ := marshal.UrssafToPeriod(field[3:9])
			periods = append(periods, periodCol{dateStart: date.Start, colIndex: index})
		}
	}
	return periods
}

func parseEffectifEntLine(row []string, idx colMapping, periods *[]periodCol, parsedLine *marshal.ParsedLineResult) {
	for _, period := range *periods {
		value := row[period.colIndex]
		if value != "" {
			noThousandsSep := sfregexp.RegexpDict["notDigit"].ReplaceAllString(value, "")
			s, err := strconv.ParseFloat(noThousandsSep, 64)
			parsedLine.AddRegularError(err)
			e := int(s)
			if e > 0 {
				parsedLine.AddTuple(EffectifEnt{
					Siren:       row[idx["siren"]],
					Periode:     period.dateStart,
					EffectifEnt: e,
				})
			}
		}
	}
}