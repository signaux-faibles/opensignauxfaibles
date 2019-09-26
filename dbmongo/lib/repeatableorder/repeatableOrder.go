package repeatableorder

import (
	"dbmongo/lib/engine"
	"dbmongo/lib/misc"
	"encoding/csv"
	"io"
	"os"
	"time"

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
	return "repOrder"
}

// Parser fonction qui retourne data et journaux
func Parser(batch engine.AdminBatch, filter map[string]bool) (chan engine.Tuple, chan engine.Event) {
	outputChannel := make(chan engine.Tuple)
	eventChannel := make(chan engine.Event)

	event := engine.Event{
		Code:    "parserRepeatableOrder",
		Channel: eventChannel,
	}

	go func() {
		for _, path := range batch.Files["repOrder"] {
			tracker := gournal.NewTracker(
				map[string]string{"path": path},
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

				repOrder := RepeatableOrder{
					Siret:       row[0],
					Periode:     periode,
					RandomOrder: randomOrder,
				}

				if !tracker.ErrorInCycle() && filter[repOrder.Siret[0:9]] {
					outputChannel <- repOrder
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
