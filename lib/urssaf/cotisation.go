package urssaf

import (
	"bufio"
	"encoding/csv"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/signaux-faibles/opensignauxfaibles/lib/base"
	"github.com/signaux-faibles/opensignauxfaibles/lib/marshal"
	"github.com/signaux-faibles/opensignauxfaibles/lib/misc"
)

// Cotisation Objet cotisation
type Cotisation struct {
	key          string       `hash:"-"`
	NumeroCompte string       `json:"numero_compte" bson:"numero_compte"`
	Periode      misc.Periode `json:"period" bson:"periode"`
	Encaisse     float64      `json:"encaisse" bson:"encaisse"`
	Du           float64      `json:"du" bson:"du"`
}

// Key _id de l'objet
func (cotisation Cotisation) Key() string {
	return cotisation.key
}

// Scope de l'objet
func (cotisation Cotisation) Scope() string {
	return "etablissement"
}

// Type de l'objet
func (cotisation Cotisation) Type() string {
	return "cotisation"
}

// ParserCotisation fournit une instance utilisable par ParseFilesFromBatch.
var ParserCotisation = &cotisationParser{}

type cotisationParser struct {
	file    *os.File
	reader  *csv.Reader
	comptes marshal.Comptes
}

func (parser *cotisationParser) GetFileType() string {
	return "cotisation"
}

func (parser *cotisationParser) Close() error {
	return parser.file.Close()
}

func (parser *cotisationParser) Init(cache *marshal.Cache, batch *base.AdminBatch) (err error) {
	parser.comptes, err = marshal.GetCompteSiretMapping(*cache, batch, marshal.OpenAndReadSiretMapping)
	return err
}

func (parser *cotisationParser) Open(filePath string) (err error) {
	parser.file, err = os.Open(filePath)
	if err != nil {
		return err
	}
	parser.reader = csv.NewReader(bufio.NewReader(parser.file))
	parser.reader.Comma = ';'
	parser.reader.LazyQuotes = true
	_, err = parser.reader.Read() // Sauter l'en-tête
	return err
}

var idxCotisation = colMapping{
	"NumeroCompte": 2,
	"Periode":      3,
	"Encaisse":     5,
	"Du":           6,
}

func (parser *cotisationParser) ParseLines(parsedLineChan chan marshal.ParsedLineResult) {
	for {
		parsedLine := marshal.ParsedLineResult{}
		row, err := parser.reader.Read()
		if err == io.EOF {
			close(parsedLineChan)
			break
		} else if err != nil {
			parsedLine.AddRegularError(err)
		} else {
			parseCotisationLine(row, &parser.comptes, &parsedLine)
			if len(parsedLine.Errors) > 0 {
				parsedLine.Tuples = []marshal.Tuple{}
			}
		}
		parsedLineChan <- parsedLine
	}
}

func parseCotisationLine(row []string, comptes *marshal.Comptes, parsedLine *marshal.ParsedLineResult) {
	idx := idxCotisation
	cotisation := Cotisation{}

	periode, err := marshal.UrssafToPeriod(row[idx["Periode"]])
	date := periode.Start
	parsedLine.AddRegularError(err)

	siret, err := marshal.GetSiretFromComptesMapping(row[idx["NumeroCompte"]], &date, *comptes)
	if err != nil {
		parsedLine.SetFilterError(err)
	} else {
		cotisation.key = siret
		cotisation.NumeroCompte = row[idx["NumeroCompte"]]
		cotisation.Periode, err = marshal.UrssafToPeriod(row[idx["Periode"]])
		parsedLine.AddRegularError(err)
		cotisation.Encaisse, err = strconv.ParseFloat(strings.Replace(row[idx["Encaisse"]], ",", ".", -1), 64)
		parsedLine.AddRegularError(err)
		cotisation.Du, err = strconv.ParseFloat(strings.Replace(row[idx["Du"]], ",", ".", -1), 64)
		parsedLine.AddRegularError(err)
	}
	parsedLine.AddTuple(cotisation)
}
