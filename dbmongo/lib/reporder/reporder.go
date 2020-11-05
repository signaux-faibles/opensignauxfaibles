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
func ParseFile(filePath string, cache *marshal.Cache, batch *base.AdminBatch) marshal.OpenFileResult {
	closeFct, reader, err := openFile(filePath)
	return marshal.OpenFileResult{
		Error: err,
		ParseLines: func(parsedLineChan chan marshal.ParsedLineResult) {
			parseLines(reader, parsedLineChan)
		},
		Close: closeFct,
	}
}

func openFile(filePath string) (func() error, *csv.Reader, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return file.Close, nil, err
	}
	reader := csv.NewReader(file)
	reader.Comma = ','
	return file.Close, reader, nil
}

func parseLines(reader *csv.Reader, parsedLineChan chan marshal.ParsedLineResult) {
	for {
		parsedLine := marshal.ParsedLineResult{}
		row, err := reader.Read()
		if err == io.EOF {
			close(parsedLineChan)
			break
		} else if err != nil {
			parsedLine.AddRegularError(err)
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
	parsedLine.AddRegularError(err)
	randomOrder, err := misc.ParsePFloat(row[2])
	parsedLine.AddRegularError(err)
	parsedLine.AddTuple(RepeatableOrder{
		Siret:       row[0],
		Periode:     periode,
		RandomOrder: randomOrder,
	})
}
