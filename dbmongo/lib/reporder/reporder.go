package reporder

import (
	"encoding/csv"
	"io"
	"os"
	"time"

	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/base"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/engine"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/marshal"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/misc"

	"github.com/signaux-faibles/gournal"

	"github.com/spf13/viper"
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

// Parser fonction qui retourne data et journaux
func Parser(cache base.Cache, batch *base.AdminBatch) (chan base.Tuple, chan base.Event) {
	outputChannel := make(chan base.Tuple)
	eventChannel := make(chan base.Event)

	event := base.Event{
		Code:    "parserRepeatableOrder",
		Channel: eventChannel,
	}

	filter := marshal.GetSirenFilterFromCache(cache)

	go func() {
		for _, path := range batch.Files["reporder"] {
			tracker := gournal.NewTracker(
				map[string]string{"path": path, "batchKey": batch.ID.Key},
				engine.TrackerReports)
			// get current file name

			file, err := os.Open(viper.GetString("APP_DATA") + path)
			if err != nil {
				tracker.Error(err)
				event.Critical(tracker.Report("fatalError"))
				continue
			}

			event.Info(path + ": ouverture")

			reader := csv.NewReader(file)
			reader.Comma = ','

			if err != nil {
				tracker.Error(err)
				event.Critical(tracker.Report("fatalError"))
				continue
			}

			for {
				row, err := reader.Read()
				if err == io.EOF {
					file.Close()
					break
				} else if err != nil {
					file.Close()
					event.Critical(path + ": abandon suite à un problème de lecture du fichier: " + err.Error())
					break
				}

				periode, err := time.Parse("2006-01-02", row[1])
				tracker.Error(err)
				randomOrder, err := misc.ParsePFloat(row[2])
				tracker.Error(err)

				reporder := RepeatableOrder{
					Siret:       row[0],
					Periode:     periode,
					RandomOrder: randomOrder,
				}
				filtered, err := marshal.IsFiltered(reporder.Siret[0:9], filter)
				if err != nil {
					tracker.Error(err)
				}
				if !tracker.HasErrorInCurrentCycle() && !filtered {
					outputChannel <- reporder
					tracker.Next()
				}
			}
			event.Info(tracker.Report("abstract"))
		}
		close(eventChannel)
		close(outputChannel)
	}()
	return outputChannel, eventChannel
}
