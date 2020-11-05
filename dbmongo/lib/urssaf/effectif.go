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

	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/base"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/marshal"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/misc"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/sfregexp"
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

// ParserEffectif expose le parseur et le type de fichier qu'il supporte.
var ParserEffectif = marshal.Parser{FileType: "effectif", FileParser: ParseEffectifFile}

// ParseEffectifFile permet de lancer le parsing du fichier demandé.
func ParseEffectifFile(filePath string, cache *marshal.Cache, batch *base.AdminBatch) (marshal.FileReader, error) {
	var idx colMapping
	var periods []periodCol
	file, reader, err := openEffectifFile(filePath)
	if err == nil {
		idx, periods, err = parseEffectifColMapping(reader)
	}
	return effectifReader{
		file:    file,
		reader:  reader,
		periods: &periods,
		idx:     idx,
	}, err
}

type effectifReader struct {
	file    *os.File
	reader  *csv.Reader
	periods *[]periodCol
	idx     colMapping
}

func (parser effectifReader) Close() error {
	return parser.file.Close()
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

func (parser effectifReader) ParseLines(parsedLineChan chan marshal.ParsedLineResult) {
	for {
		parsedLine := marshal.ParsedLineResult{}
		row, err := parser.reader.Read()
		if err == io.EOF {
			close(parsedLineChan)
			break
		} else if err != nil {
			parsedLine.AddError(base.NewRegularError(err))
		} else {
			parseEffectifLine(row, parser.idx, parser.periods, &parsedLine)
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
			parsedLine.AddError(base.NewRegularError(err))
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
