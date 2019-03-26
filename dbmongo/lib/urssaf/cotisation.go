package urssaf

import (
	"bufio"
	"dbmongo/lib/engine"
	"encoding/csv"
	"io"
	"os"
	"strconv"
	"strings"
  "time"

	"github.com/chrnin/gournal"
	"github.com/spf13/viper"
)

// Cotisation Objet cotisation
type Cotisation struct {
	key          string  `hash:"-"`
	NumeroCompte string  `json:"numero_compte" bson:"numero_compte"`
	PeriodeDebit string  `json:"periode_debit" bson:"periode_debit"`
	Periode      Periode `json:"period" bson:"periode"`
	Recouvrement float64 `json:"recouvrement" bson:"recouvrement"`
	Encaisse     float64 `json:"encaisse" bson:"encaisse"`
	Du           float64 `json:"du" bson:"du"`
	Ecriture     string  `json:"ecriture" bson:"ecriture"`
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
func parseCotisation(batch engine.AdminBatch, mapping Comptes) (chan engine.Tuple, chan engine.Event) {
	outputChannel := make(chan engine.Tuple)
	eventChannel := make(chan engine.Event)

	field := map[string]int{
		"NumeroCompte": 0,
		"PeriodeDebit": 1,
		"Periode":      4,
		"Recouvrement": 2,
		"Encaisse":     3,
		"Du":           5,
		"Ecriture":     6,
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
          date, err := urssafToDate(row[field["Periode"]])
          if err != nil { date = time.Now() }

          if siret, err := mapping.GetSiret(row[field["NumeroCompte"]], date); err == nil {
						cotisation := Cotisation{}
						cotisation.key = siret
						cotisation.NumeroCompte = row[field["NumeroCompte"]]
						cotisation.Periode, err = urssafToPeriod(row[field["Periode"]])
						tracker.Error(err)
						cotisation.PeriodeDebit = row[field["PeriodeDebit"]]
						cotisation.Recouvrement, err = strconv.ParseFloat(strings.Replace(row[field["Recouvrement"]], ",", ".", -1), 64)
						tracker.Error(err)
						cotisation.Encaisse, err = strconv.ParseFloat(strings.Replace(row[field["Encaisse"]], ",", ".", -1), 64)
						tracker.Error(err)
						cotisation.Du, err = strconv.ParseFloat(strings.Replace(row[field["Du"]], ",", ".", -1), 64)
						tracker.Error(err)
						cotisation.Ecriture = row[field["Ecriture"]]

						if !tracker.ErrorInCycle() {
							outputChannel <- cotisation
						} else {
							//event.Debug(tracker.Report("errors"))
						}
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
