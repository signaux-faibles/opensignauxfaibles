package urssaf

import (
	"bufio"
	"encoding/csv"
	"errors"
	"io"
	"os"
	"time"

	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/base"
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
	Action         string    `json:"action" bson:"action"`
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

// ParserCCSF produit des lignes CCSF
func ParserCCSF(cache marshal.Cache, batch *base.AdminBatch) (chan marshal.Tuple, chan marshal.Event) {
	outputChannel := make(chan marshal.Tuple)
	eventChannel := make(chan marshal.Event)

	event := marshal.Event{
		Code:    "ccsfParser",
		Channel: eventChannel,
	}

	go func() {

		for _, path := range batch.Files["ccsf"] {
			tracker := gournal.NewTracker(
				map[string]string{"path": path, "batchKey": batch.ID.Key},
				engine.TrackerReports)

			file, err := os.Open(viper.GetString("APP_DATA") + path)
			if err != nil {
				tracker.Add(err)
				event.Critical(tracker.Report("fatalError"))
				continue
			} else {
				event.Info(path + ": ouverture")
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
					ccsf.DateTraitement, err = marshal.UrssafToDate(r[f["DateTraitement"]])
					tracker.Add(err)
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
						tracker.Add(base.NewFilterError(err))
						continue
					}
					ccsf.NumeroCompte = r[f["NumeroCompte"]]

					if !tracker.HasErrorInCurrentCycle() {
						outputChannel <- ccsf
					} else {
						//event.Debug(tracker.Report("error"))
					}

				} else {
					tracker.Add(errors.New("Ligne non conforme, moins de 4 champs"))
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
