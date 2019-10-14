package urssaf

import (
	"bufio"
	"encoding/csv"
	"io"
	"opensignauxfaibles/dbmongo/lib/engine"
	"opensignauxfaibles/dbmongo/lib/marshal"
	"opensignauxfaibles/dbmongo/lib/misc"
	"os"
	"strconv"
	"strings"

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

// ParseCotisation transforme les fichiers en données à intégrer
func parseCotisation(cache engine.Cache, batch *engine.AdminBatch) (chan engine.Tuple, chan engine.Event) {
	outputChannel := make(chan engine.Tuple)
	eventChannel := make(chan engine.Event)

	field := map[string]int{
		"NumeroCompte": 2,
		"Periode":      3,
		"Encaisse":     5,
		"Du":           6,
	}

	go func() {
		event := engine.Event{
			Code:    "cotisationParser",
			Channel: eventChannel,
		}

		for _, path := range batch.Files["cotisation"] {
			tracker := gournal.NewTracker(
				map[string]string{"path": path},
				engine.TrackerReports)

			file, err := os.Open(viper.GetString("APP_DATA") + path)
			defer file.Close()

			if err != nil {
				tracker.Error(err)
				event.Critical(tracker.Report("fatalError"))
				break
			}

			reader := csv.NewReader(bufio.NewReader(file))
			reader.Comma = ';'
			reader.LazyQuotes = true
			// ligne de titre
			reader.Read()

			for {
				row, err := reader.Read()
				if err == io.EOF {
					break
				} else if err != nil {
				} else {
					periode, err := urssafToPeriod(row[field["Periode"]])
					date := periode.Start
					tracker.Error(err)
					// if err != nil { date = time.Now() }

					if siret, err := marshal.GetSiret(row[field["NumeroCompte"]], &date, cache, batch); err == nil {
						cotisation := Cotisation{}
						cotisation.key = siret
						cotisation.NumeroCompte = row[field["NumeroCompte"]]
						cotisation.Periode, err = urssafToPeriod(row[field["Periode"]])
						tracker.Error(err)
						cotisation.Encaisse, err = strconv.ParseFloat(strings.Replace(row[field["Encaisse"]], ",", ".", -1), 64)
						tracker.Error(err)
						cotisation.Du, err = strconv.ParseFloat(strings.Replace(row[field["Du"]], ",", ".", -1), 64)
						tracker.Error(err)

						if !tracker.HasErrorInCurrentCycle() {
							outputChannel <- cotisation
						}
					} else {
						continue
					}
				}
				tracker.Next()
			}
			event.Info(tracker.Report("abstract"))
			file.Close()
		}
		close(eventChannel)
		close(outputChannel)
	}()
	return outputChannel, eventChannel
}
