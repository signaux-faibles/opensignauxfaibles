package reporder

import (
	"encoding/csv"
	"io"
	"os"
	"time"

	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/base"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/marshal"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/misc"

	"github.com/signaux-faibles/gournal"
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

// ParseFile extrait les tuples depuis un fichier Reporder et génère un rapport Gournal.
func ParseFile(filePath string, cache *marshal.Cache, batch *base.AdminBatch, tracker *gournal.Tracker) marshal.ParsedLineChan {
	file, err := os.Open(filePath)
	if err != nil {
		tracker.Add(err)
		return nil
	}
	// defer file.Close() // TODO: à réactiver

	reader := csv.NewReader(file)
	reader.Comma = ','

	parsedLineChan := make(marshal.ParsedLineChan)
	go func() {
		for {
			tuples := []marshal.Tuple{}
			row, err := reader.Read()
			if err == io.EOF {
				close(parsedLineChan)
				break
			} else if err != nil {
				tracker.Add(err)
			} else {
				reporder := parseReporderLine(row, tracker)
				if !tracker.HasErrorInCurrentCycle() {
					tuples = []marshal.Tuple{reporder}
				}
			}
			parsedLineChan <- marshal.ParsedLineResult{Tuples: tuples, Errors: []marshal.ParseError{}}
		}
	}()
	return parsedLineChan
}

func parseReporderLine(row []string, tracker *gournal.Tracker) RepeatableOrder {
	periode, err := time.Parse("2006-01-02", row[1])
	tracker.Add(err)
	randomOrder, err := misc.ParsePFloat(row[2])
	tracker.Add(err)
	return RepeatableOrder{
		Siret:       row[0],
		Periode:     periode,
		RandomOrder: randomOrder,
	}
}
