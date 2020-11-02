package urssaf

import (
	"bufio"
	"encoding/csv"
	"errors"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/base"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/marshal"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/misc"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/sfregexp"
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

// ParserEffectifEnt expose le parseur et le type de fichier qu'il supporte.
var ParserEffectifEnt = marshal.Parser{FileType: "effectif_ent", FileParser: ParseEffectifEntFile}

// ParseEffectifEntFile permet de lancer le parsing du fichier demandé.
func ParseEffectifEntFile(filePath string, cache *marshal.Cache, batch *base.AdminBatch) marshal.OpenFileResult {
	var idx colMapping
	var periods []periodCol
	closeFct, reader, err := openEffectifEntFile(filePath)
	if err == nil {
		idx, periods, err = parseEffectifEntColMapping(reader)
	}
	return marshal.OpenFileResult{
		Error: err,
		ParseLines: func(parsedLineChan chan base.ParsedLineResult) {
			parseEffectifEntLines(reader, idx, periods, parsedLineChan)
		},
		Close: closeFct,
	}
}

func openEffectifEntFile(filePath string) (func() error, *csv.Reader, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return file.Close, nil, err
	}
	reader := csv.NewReader(bufio.NewReader(file))
	reader.Comma = ';'
	return file.Close, reader, err
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

func parseEffectifEntLines(reader *csv.Reader, idx colMapping, periods []periodCol, parsedLineChan chan base.ParsedLineResult) {
	for {
		parsedLine := base.ParsedLineResult{}
		row, err := reader.Read()
		if err == io.EOF {
			close(parsedLineChan)
			break
		} else if err != nil {
			parsedLine.AddError(err)
		} else {
			parseEffectifEntLine(row, idx, periods, &parsedLine)
		}
		parsedLineChan <- parsedLine
	}
}

func parseEffectifEntLine(row []string, idx colMapping, periods []periodCol, parsedLine *base.ParsedLineResult) {
	siren := row[idx["siren"]]
	if !sfregexp.ValidSiren(siren) {
		parsedLine.AddError(errors.New("Format de siren incorrect : " + siren)) // TODO: remove validation
	} else {
		for _, period := range periods {
			value := row[period.colIndex]
			if value != "" {
				noThousandsSep := sfregexp.RegexpDict["notDigit"].ReplaceAllString(value, "")
				s, err := strconv.ParseFloat(noThousandsSep, 64)
				parsedLine.AddError(err)
				e := int(s)
				if e > 0 {
					parsedLine.AddTuple(EffectifEnt{
						Siren:       siren,
						Periode:     period.dateStart,
						EffectifEnt: e,
					})
				}
			}
		}
	}
}
