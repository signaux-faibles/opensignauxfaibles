package urssaf

import (
	"bufio"
	"encoding/csv"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/base"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/marshal"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/misc"
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

// ParserCotisation expose le parseur et le type de fichier qu'il supporte.
var ParserCotisation = marshal.Parser{FileType: "cotisation", FileParser: ParseCotisationFile}

// ParseCotisationFile extrait les tuples depuis le fichier demandé et génère un rapport Gournal.
func ParseCotisationFile(filePath string, cache *marshal.Cache, batch *base.AdminBatch) (marshal.ParsedLineChan, error) {
	comptes, err := marshal.GetCompteSiretMapping(*cache, batch, marshal.OpenAndReadSiretMapping)
	if err != nil {
		return nil, err
	}
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	// defer file.Close() // TODO: à réactiver
	reader := csv.NewReader(bufio.NewReader(file))
	reader.Comma = ';'
	reader.LazyQuotes = true

	// ligne de titre
	reader.Read()

	var idx = colMapping{
		"NumeroCompte": 2,
		"Periode":      3,
		"Encaisse":     5,
		"Du":           6,
	}

	parsedLineChan := make(marshal.ParsedLineChan)
	go func() {
		for {
			parsedLine := base.ParsedLineResult{}
			row, err := reader.Read()
			if err == io.EOF {
				close(parsedLineChan)
				break
			} else if err != nil {
				parsedLine.AddError(err)
			} else {
				parseCotisationLine(row, &comptes, idx, &parsedLine)
				if len(parsedLine.Errors) > 0 {
					parsedLine.Tuples = []base.Tuple{}
				}
			}
			parsedLineChan <- parsedLine
		}
	}()
	return parsedLineChan, nil
}

func parseCotisationLine(row []string, comptes *marshal.Comptes, idx colMapping, parsedLine *base.ParsedLineResult) {
	cotisation := Cotisation{}

	periode, err := marshal.UrssafToPeriod(row[idx["Periode"]])
	date := periode.Start
	parsedLine.AddError(err)

	siret, err := marshal.GetSiretFromComptesMapping(row[idx["NumeroCompte"]], &date, *comptes)
	if err != nil {
		parsedLine.AddError(base.NewFilterError(err))
	} else {
		cotisation.key = siret
		cotisation.NumeroCompte = row[idx["NumeroCompte"]]
		cotisation.Periode, err = marshal.UrssafToPeriod(row[idx["Periode"]])
		parsedLine.AddError(err)
		cotisation.Encaisse, err = strconv.ParseFloat(strings.Replace(row[idx["Encaisse"]], ",", ".", -1), 64)
		parsedLine.AddError(err)
		cotisation.Du, err = strconv.ParseFloat(strings.Replace(row[idx["Du"]], ",", ".", -1), 64)
		parsedLine.AddError(err)
	}
	parsedLine.AddTuple(cotisation)
}
