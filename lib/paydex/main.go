package paydex

import (
	"bufio"
	"encoding/csv"
	"io"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/signaux-faibles/opensignauxfaibles/lib/base"
	"github.com/signaux-faibles/opensignauxfaibles/lib/marshal"
)

// Paydex décrit le format de chaque entrée de donnée résultant du parsing.
type Paydex struct {
	Siren   string    `json:"siren" bson:"siren"`
	Periode time.Time `json:"periode" bson:"periode"`
	Jours   int       `json:"paydex_jours" bson:"paydex_jours"`
}

// Key _id de l'objet
func (paydex Paydex) Key() string {
	return paydex.Siren
}

// Scope de l'objet
func (paydex Paydex) Scope() string {
	return "entreprise"
}

// Type de l'objet
func (paydex Paydex) Type() string {
	return "paydex"
}

// ParserPaydex fournit une instance utilisable par ParseFilesFromBatch.
var ParserPaydex = &paydexParser{}

type paydexParser struct {
	file   *os.File
	reader *csv.Reader
}

func (parser *paydexParser) GetFileType() string {
	return "paydex"
}

func (parser *paydexParser) Init(cache *marshal.Cache, batch *base.AdminBatch) error {
	return nil
}

func (parser *paydexParser) Close() error {
	return parser.file.Close()
}

func (parser *paydexParser) Open(filePath string) (err error) {
	parser.file, parser.reader, err = openPaydexFile(filePath)
	return err
}

func openPaydexFile(filePath string) (*os.File, *csv.Reader, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return file, nil, err
	}
	reader := csv.NewReader(bufio.NewReader(file))
	reader.Comma = ';'
	return file, reader, err
}

func (parser *paydexParser) ParseLines(parsedLineChan chan marshal.ParsedLineResult) {
	parser.reader.Read() // parse header
	for {
		parsedLine := marshal.ParsedLineResult{}
		row, err := parser.reader.Read()
		if err == io.EOF {
			close(parsedLineChan)
			break
		} else if err != nil {
			parsedLine.AddRegularError(err)
		} else {
			parsedLine.AddTuple(parsePaydexLine(row))
		}
		parsedLineChan <- parsedLine
	}
}

func parsePaydexLine(row []string) Paydex {
	periode, err := time.Parse("02/01/2006", row[3])
	if err != nil {
		log.Fatalf("invalid date: %v", row[3])
	}
	jours, err := strconv.Atoi(row[1])
	if err != nil {
		log.Fatalf("invalid date: %v", row[3])
	}
	return Paydex{
		Siren:   row[0],
		Periode: beginningOfMonth(periode),
		Jours:   jours,
	}
}

func beginningOfMonth(date time.Time) time.Time {
	return date.AddDate(0, 0, -date.Day()+1)
}

