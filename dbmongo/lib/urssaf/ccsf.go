package urssaf

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/engine"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/marshal"

	"github.com/signaux-faibles/gournal"
	"github.com/spf13/viper"
)

// CCSF information urssaf ccsf
type CCSF struct {
	key            string    `hash:"-"`
	NumeroCompte   string    `json:"-" bson:"-"`
	DateTraitement time.Time `json:"date_traitement" bson:"date_traitement"`
	Stade          string    `json:"stade" bson:"stade"`
	Action         string    `json:"action" json:"action"`
}

// Key _id de l'objet
func (ccsf CCSF) Key() string {
	return ccsf.key
}

// Scope de l'objet
func (ccsf CCSF) Scope() string {
	return "etablissement"
}

// Type de l'objet
func (ccsf CCSF) Type() string {
	return "ccsf"
}

func batchToTime(batch string) (time.Time, error) {
	year, err := strconv.Atoi(batch[0:2])
	if err != nil {
		return time.Time{}, err
	}

	month, err := strconv.Atoi(batch[2:4])
	if err != nil {
		return time.Time{}, err
	}

	date := time.Date(2000+year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	return date, err
}

// Parser produit des lignes CCSF
func parseCCSF(cache engine.Cache, batch *engine.AdminBatch) (chan engine.Tuple, chan engine.Event) {
	outputChannel := make(chan engine.Tuple)
	eventChannel := make(chan engine.Event)

	event := engine.Event{
		Code:    "ccsfParser",
		Channel: eventChannel,
	}

	go func() {

		for _, path := range batch.Files["ccsf"] {
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

			f := map[string]int{
				"NumeroCompte":   2,
				"DateTraitement": 3,
				"Stade":          4,
				"Action":         5,
			}

			for {
				r, err := reader.Read()
				if err == io.EOF {
					break
				} else if err != nil {
					event.Critical(path + "Erreur à la lecture, abandon: " + err.Error())
					continue
				}
				if len(r) >= 4 {
					ccsf := CCSF{}

					ccsf.Action = r[f["Action"]]
					ccsf.Stade = r[f["Stade"]]
					ccsf.DateTraitement, err = urssafToDate(r[f["DateTraitement"]])
					tracker.Error(err)
					if err != nil {
						tracker.Next()
						continue
					}
					ccsf.key, err = marshal.GetSiret(
						r[f["NumeroCompte"]],
						&ccsf.DateTraitement,
						cache,
						batch,
					)
					if err != nil {
						// Compte filtré
						continue
					}
					ccsf.NumeroCompte = r[f["NumeroCompte"]]

					if !tracker.HasErrorInCurrentCycle() {
						outputChannel <- ccsf
					} else {
						//event.Debug(tracker.Report("error"))
					}

				} else {
					tracker.Error(errors.New("Ligne non conforme, moins de 4 champs"))
					event.Warning(tracker.Report("invalidLine"))
				}
				tracker.Next()
			}

			event.Info(tracker.Report("abstract"))

			file.Close()
		}
		close(outputChannel)
		close(eventChannel)
	}()
	return outputChannel, eventChannel
}
