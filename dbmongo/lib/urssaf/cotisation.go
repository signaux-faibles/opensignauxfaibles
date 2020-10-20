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

	"github.com/signaux-faibles/gournal"
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

// ParserCotisation transforme les fichiers en données à intégrer
func ParserCotisation(cache marshal.Cache, batch *base.AdminBatch) (chan marshal.Tuple, chan marshal.Event) {
	return marshal.ParseFilesFromBatch(cache, batch, marshal.Parser{FileType: "cotisation", FileParser: ParseCotisationFile})
}

// ParseCotisationFile extrait les tuples depuis le fichier demandé et génère un rapport Gournal.
func ParseCotisationFile(filePath string, cache *marshal.Cache, batch *base.AdminBatch, tracker *gournal.Tracker, outputChannel chan marshal.Tuple) {
	comptes, err := marshal.GetCompteSiretMapping(*cache, batch, marshal.OpenAndReadSiretMapping)
	if err != nil {
		tracker.Add(err)
		return
	}
	file, err := os.Open(filePath)
	if err != nil {
		tracker.Add(err)
		return
	}
	defer file.Close()
	reader := csv.NewReader(bufio.NewReader(file))
	reader.Comma = ';'
	reader.LazyQuotes = true
	parseCotisationFile(reader, &comptes, tracker, outputChannel)
}

func parseCotisationFile(reader *csv.Reader, comptes *marshal.Comptes, tracker *gournal.Tracker, outputChannel chan marshal.Tuple) {
	// ligne de titre
	reader.Read()

	var idx = colMapping{
		"NumeroCompte": 2,
		"Periode":      3,
		"Encaisse":     5,
		"Du":           6,
	}

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			tracker.Add(err)
		} else {
			cotisation := parseCotisationLine(row, tracker, comptes, idx)
			if !tracker.HasErrorInCurrentCycle() {
				outputChannel <- cotisation
			}
		}
		tracker.Next()
	}
}

func parseCotisationLine(row []string, tracker *gournal.Tracker, comptes *marshal.Comptes, idx colMapping) Cotisation {
	cotisation := Cotisation{}

	periode, err := marshal.UrssafToPeriod(row[idx["Periode"]])
	date := periode.Start
	tracker.Add(err)

	siret, err := marshal.GetSiretFromComptesMapping(row[idx["NumeroCompte"]], &date, *comptes)
	if err != nil {
		tracker.Add(base.NewFilterError(err))
	} else {
		cotisation.key = siret
		cotisation.NumeroCompte = row[idx["NumeroCompte"]]
		cotisation.Periode, err = marshal.UrssafToPeriod(row[idx["Periode"]])
		tracker.Add(err)
		cotisation.Encaisse, err = strconv.ParseFloat(strings.Replace(row[idx["Encaisse"]], ",", ".", -1), 64)
		tracker.Add(err)
		cotisation.Du, err = strconv.ParseFloat(strings.Replace(row[idx["Du"]], ",", ".", -1), 64)
		tracker.Add(err)
	}
	return cotisation
}
