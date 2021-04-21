package reporder

import (
	"encoding/csv"
	"os"
	"time"

	"github.com/signaux-faibles/opensignauxfaibles/lib/base"
	"github.com/signaux-faibles/opensignauxfaibles/lib/marshal"
	"github.com/signaux-faibles/opensignauxfaibles/lib/misc"
)

// RepeatableOrder random number
type RepeatableOrder struct {
	Siret       string    `json:"siret"          bson:"siret"`
	Periode     time.Time `json:"periode"        bson:"periode"`
	RandomOrder *float64  `json:"random_order"   bson:"random_order"`
}

// Key de l'objet
func (rep RepeatableOrder) Key() string {
	return rep.Siret
}

// Scope de l'objet
func (rep RepeatableOrder) Scope() string {
	return "etablissement"
}

// Type de l'objet
func (rep RepeatableOrder) Type() string {
	return "reporder"
}

// Parser fournit une instance utilisable par ParseFilesFromBatch.
var Parser = &reporderParser{}

type reporderParser struct {
	file   *os.File
	reader *csv.Reader
}

func (parser *reporderParser) GetFileType() string {
	return "reporder"
}

func (parser *reporderParser) Init(cache *marshal.Cache, batch *base.AdminBatch) error {
	return nil
}

func (parser *reporderParser) Open(filePath string) (err error) {
	parser.file, parser.reader, err = marshal.OpenCsvReader(base.BatchFile(filePath), ',', false)
	return err
}

func (parser *reporderParser) Close() error {
	return parser.file.Close()
}

func (parser *reporderParser) ParseLines(parsedLineChan chan marshal.ParsedLineResult) {
	marshal.ParseLines(parsedLineChan, parser.reader, func(row []string, parsedLine *marshal.ParsedLineResult) {
		parseReporderLine(row, parsedLine)
	})
}

func parseReporderLine(row []string, parsedLine *marshal.ParsedLineResult) {
	periode, err := time.Parse("2006-01-02", row[1])
	parsedLine.AddRegularError(err)
	randomOrder, err := misc.ParsePFloat(row[2])
	parsedLine.AddRegularError(err)
	parsedLine.AddTuple(RepeatableOrder{
		Siret:       row[0],
		Periode:     periode,
		RandomOrder: randomOrder,
	})
}
