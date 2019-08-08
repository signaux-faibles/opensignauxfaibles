package urssaf

import (
	"bufio"
	"dbmongo/lib/engine"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/signaux-faibles/gournal"
	"github.com/spf13/viper"
)

// DPAE Déclaration préalabre à l'embauche
type DPAE struct {
	Siret    string    `json:"-" bson:"-"`
	Date     time.Time `json:"date" bson:"date"`
	CDI      float64   `json:"cdi" bson:"cdi"`
	CDDLong  float64   `json:"cdd_long" bson:"cdd_long"`
	CDDCourt float64   `json:"cdd_court" bson:"cdd_court"`
}

// Key _id de l'objet
func (dpae DPAE) Key() string {
	return dpae.Siret
}

// Scope de l'objet
func (dpae DPAE) Scope() string {
	return "etablissement"
}

// Type de l'objet
func (dpae DPAE) Type() string {
	return "dpae"
}

// Parser produit les datas DPAE
func parseDPAE(batch engine.AdminBatch, mapping Comptes) (chan engine.Tuple, chan engine.Event) {
	outputChannel := make(chan engine.Tuple)
	eventChannel := make(chan engine.Event)
	event := engine.Event{
		Code:    "dpaeParser",
		Channel: eventChannel,
	}

	go func() {
		for _, path := range batch.Files["dpae"] {
			tracker := gournal.NewTracker(
				map[string]string{"path": path},
				engine.TrackerReports)

			file, err := os.Open(viper.GetString("APP_DATA") + path)
			if err != nil {
				fmt.Println("Error", err)
			}

			reader := csv.NewReader(bufio.NewReader(file))
			reader.Comma = ';'
			reader.Read()
			for {
				row, err := reader.Read()
				if err == io.EOF {
					break
				} else if err != nil {
					event.Critical(path + ": erreur à la lecture du fichier, abandon: " + err.Error())
					file.Close()
					continue
				}

				date, err := time.Parse("20060102", row[1]+row[2]+"01")
				tracker.Error(err)

				dpae := DPAE{
					Siret: row[0],
					Date:  date,
				}
				dpae.CDI, err = strconv.ParseFloat(row[3], 64)
				tracker.Error(err)
				dpae.CDDLong, err = strconv.ParseFloat(row[4], 64)
				tracker.Error(err)
				dpae.CDDCourt, err = strconv.ParseFloat(row[5], 64)
				tracker.Error(err)

				if !tracker.ErrorInCycle() {
					outputChannel <- dpae
				} else {
					//event.Debug(tracker.Report("errors"))
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
