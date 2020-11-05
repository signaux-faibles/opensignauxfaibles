package reporder

import (
	"encoding/csv"
	"io"
	"os"
	"time"

	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/base"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/marshal"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/misc"
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

// Parser expose le parseur et le type de fichier qu'il supporte.
var Parser = marshal.Parser{FileType: "reporder", FileParser: ParseFile}

// ParseFile permet de lancer le parsing du fichier demandé.
func ParseFile(filePath string, cache *marshal.Cache, batch *base.AdminBatch) (marshal.FileReader, error) {
	file, reader, err := openFile(filePath)
	return reporderReader{
		file:   file,
		reader: reader,
	}, err
}

type reporderReader struct {
	file   *os.File
	reader *csv.Reader
}

func (parser reporderReader) Close() error {
	return parser.file.Close()
}

func openFile(filePath string) (*os.File, *csv.Reader, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return file, nil, err
	}
	reader := csv.NewReader(file)
	reader.Comma = ','
	return file, reader, nil
}

func (parser reporderReader) ParseLines(parsedLineChan chan marshal.ParsedLineResult) {
	for {
		parsedLine := marshal.ParsedLineResult{}
		row, err := parser.reader.Read()
		if err == io.EOF {
			close(parsedLineChan)
			break
		} else if err != nil {
			parsedLine.AddError(base.NewRegularError(err))
		} else {
			parseReporderLine(row, &parsedLine)
			if len(parsedLine.Errors) > 0 {
				parsedLine.Tuples = []marshal.Tuple{}
			}
		}
		parsedLineChan <- parsedLine
	}
}

func parseReporderLine(row []string, parsedLine *marshal.ParsedLineResult) {
	periode, err := time.Parse("2006-01-02", row[1])
	parsedLine.AddError(base.NewRegularError(err))
	randomOrder, err := misc.ParsePFloat(row[2])
	parsedLine.AddError(base.NewRegularError(err))
	parsedLine.AddTuple(RepeatableOrder{
		Siret:       row[0],
		Periode:     periode,
		RandomOrder: randomOrder,
	})
}
