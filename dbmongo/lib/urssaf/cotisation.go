package urssaf

import (
	"bufio"
	"encoding/csv"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/base"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/engine"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/marshal"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/misc"

	"github.com/signaux-faibles/gournal"
	"github.com/spf13/viper"
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
	outputChannel := make(chan marshal.Tuple)
	eventChannel := make(chan marshal.Event)

	go func() {
		event := marshal.Event{
			Code:    "cotisationParser",
			Channel: eventChannel,
		}

		for _, path := range batch.Files["cotisation"] {
			tracker := gournal.NewTracker(
				map[string]string{"path": path, "batchKey": batch.ID.Key},
				engine.TrackerReports)

			file, err := os.Open(viper.GetString("APP_DATA") + path)
			if err != nil {
				tracker.Add(err)
				event.Critical(tracker.Report("fatalError"))
				break
			} else {
				event.Info(path + ": ouverture")
			}
			defer file.Close()

			comptes, err := marshal.GetCompteSiretMapping(cache, batch, marshal.OpenAndReadSiretMapping)
			if err != nil {
				tracker.Add(err)
			} else {
				reader := csv.NewReader(bufio.NewReader(file))
				reader.Comma = ';'
				reader.LazyQuotes = true
				parseCotisationFile(reader, &comptes, &tracker, outputChannel)
			}

			event.Info(tracker.Report("abstract"))
			file.Close()
		}
		close(eventChannel)
		close(outputChannel)
	}()
	return outputChannel, eventChannel
}

func parseCotisationFile(reader *csv.Reader, comptes *marshal.Comptes, tracker *gournal.Tracker, outputChannel chan marshal.Tuple) {
	// ligne de titre
	reader.Read()

	var field = colMapping{
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
			break
		} else {
			cotisation := parseCotisationLine(row, tracker, comptes, field)
			if !tracker.HasErrorInCurrentCycle() {
				outputChannel <- cotisation
			}
		}
		tracker.Next()
	}
}

func parseCotisationLine(row []string, tracker *gournal.Tracker, comptes *marshal.Comptes, field colMapping) Cotisation {
	cotisation := Cotisation{}

	periode, err := marshal.UrssafToPeriod(row[field["Periode"]])
	date := periode.Start
	tracker.Add(err)

	siret, err := marshal.GetSiretFromComptesMapping(row[field["NumeroCompte"]], &date, *comptes)
	if err != nil {
		tracker.Add(base.NewFilterError(err))
	} else {
		cotisation.key = siret
		cotisation.NumeroCompte = row[field["NumeroCompte"]]
		cotisation.Periode, err = marshal.UrssafToPeriod(row[field["Periode"]])
		tracker.Add(err)
		cotisation.Encaisse, err = strconv.ParseFloat(strings.Replace(row[field["Encaisse"]], ",", ".", -1), 64)
		tracker.Add(err)
		cotisation.Du, err = strconv.ParseFloat(strings.Replace(row[field["Du"]], ",", ".", -1), 64)
		tracker.Add(err)
	}
	return cotisation
}
