package urssaf

import (
	"bufio"
	"encoding/csv"
	"errors"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/signaux-faibles/opensignauxfaibles/lib/base"
	"github.com/signaux-faibles/opensignauxfaibles/lib/marshal"
	"github.com/signaux-faibles/opensignauxfaibles/lib/misc"
	"github.com/signaux-faibles/opensignauxfaibles/lib/sfregexp"
)

// Effectif Urssaf
type Effectif struct {
	Siret        string    `json:"-" bson:"-"`
	NumeroCompte string    `json:"numero_compte" bson:"numero_compte"`
	Periode      time.Time `json:"periode" bson:"periode"`
	Effectif     int       `json:"effectif" bson:"effectif"`
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
	idx     colMapping
}

func (parser *effectifParser) GetFileType() string {
	return "effectif"
}

func (parser *effectifParser) Init(cache *marshal.Cache, batch *base.AdminBatch) error {
	return nil
}

func (parser *effectifParser) Close() error {
	return parser.file.Close()
}

func (parser *effectifParser) Open(filePath string) (err error) {
	parser.file, parser.reader, err = openEffectifFile(filePath)
	if err == nil {
		parser.idx, parser.periods, err = parseEffectifColMapping(parser.reader)
	}
	return err
}

func openEffectifFile(filePath string) (*os.File, *csv.Reader, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return file, nil, err
	}
	reader := csv.NewReader(bufio.NewReader(file))
	reader.Comma = ';'
	return file, reader, err
}

func parseEffectifColMapping(reader *csv.Reader) (colMapping, []periodCol, error) {
	fields, err := reader.Read()
	if err != nil {
		return nil, nil, err
	}

	var idx = colMapping{
		"siret":  misc.SliceIndex(len(fields), func(i int) bool { return strings.ToLower(fields[i]) == "siret" }),
		"compte": misc.SliceIndex(len(fields), func(i int) bool { return strings.ToLower(fields[i]) == "compte" }),
	}

	if misc.SliceMin(idx["siret"], idx["compte"]) == -1 {
		return nil, nil, errors.New("erreur à l'analyse du fichier, abandon, l'un " +
			"des champs obligatoires n'a pu etre trouve:" +
			" siretIndex = " + strconv.Itoa(idx["siret"]) +
			", compteIndex = " + strconv.Itoa(idx["compte"]))
	}

	// Dans quels champs lire l'effectif
	periods := parseEffectifPeriod(fields)
	return idx, periods, err
}

func (parser *effectifParser) ParseLines(parsedLineChan chan marshal.ParsedLineResult) {
	for {
		parsedLine := marshal.ParsedLineResult{}
		row, err := parser.reader.Read()
		if err == io.EOF {
			close(parsedLineChan)
			break
		} else if err != nil {
			parsedLine.AddRegularError(err)
		} else {
			parseEffectifLine(row, parser.idx, &parser.periods, &parsedLine)
		}
		parsedLineChan <- parsedLine
	}
}

func parseEffectifLine(row []string, idx colMapping, periods *[]periodCol, parsedLine *marshal.ParsedLineResult) {
	for _, period := range *periods {
		value := row[period.colIndex]
		if value != "" {
			noThousandsSep := sfregexp.RegexpDict["notDigit"].ReplaceAllString(value, "")
			e, err := strconv.Atoi(noThousandsSep)
			parsedLine.AddRegularError(err)
			if e > 0 {
				parsedLine.AddTuple(Effectif{
					Siret:        row[idx["siret"]],
					NumeroCompte: row[idx["compte"]],
					Periode:      period.dateStart,
					Effectif:     e,
				})
			}
		}
	}
}