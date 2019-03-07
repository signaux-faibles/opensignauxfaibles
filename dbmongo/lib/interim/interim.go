package interim

import (
	"dbmongo/lib/engine"
	"errors"
	"fmt"
	"os"
	"regexp"
	"time"

	"github.com/chrnin/gournal"
	"github.com/kshedden/datareader"

	"github.com/spf13/viper"
)

// Interim Interim – fichier DARES
type Interim struct {
	Siret   string    `json:"siret" bson:"siret"`
	Periode time.Time `json:"periode" bson:"periode"`
	ETP     float64   `json:"etp" bson:"etp"`
}

// Key de l'objet
func (interim Interim) Key() string {
	return interim.Siret
}

// Scope de l'objet
func (interim Interim) Scope() string {
	return "etablissement"
}

// Type de l'objet
func (interim Interim) Type() string {
	return "interim"
}

// Parser fonction qui retourne data et journaux
func Parser(batch engine.AdminBatch) (chan engine.Tuple, chan engine.Event) {
	outputChannel := make(chan engine.Tuple)
	eventChannel := make(chan engine.Event)

	event := engine.Event{
		Code: "parserInterim",
	}

	field := map[string]int{
		"Siret":   0,
		"Periode": 1,
		"ETP":     4,
	}

	go func() {
		for _, path := range batch.Files["interim"] {
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

			reader, err := datareader.NewSAS7BDATReader(file)
			if err != nil {
				tracker.Error(err)
				event.Critical(tracker.Report("fatalError"))
				continue
			}

			row, err := reader.Read(-1)
			if err != nil {
				tracker.Error(err)
				event.Critical(tracker.Report("fatalError"))
				continue
			}

			sirets, missing, err := row[field["Siret"]].AsStringSlice()
			tracker.Error(err)
			periode, _, err := row[field["Periode"]].AsFloat64Slice()
			tracker.Error(err)
			etp, _, err := row[field["ETP"]].AsFloat64Slice()
			tracker.Error(err)
			if tracker.ErrorInCycle() {
				event.Debug(tracker.Report("errors"))
				tracker.Error(errors.New("problème d'accès aux données SAS"))
				event.Critical(tracker.Report("fatalError"))
				continue
			}
			for i := 0; i < len(sirets); i++ {
				if err != nil {
					break
				}
				interim := Interim{}
        validSiret, _ := regexp.MatchString("[0-9]{14}", sirets[i]);
        validSiren :=  (sirets[i][:9] != "000000000")

				if  !missing[i] && validSiret &&  validSiren {
					tracker.Error(err)
					interim.Siret = sirets[i][:14]
					interim.Periode, _ = time.Parse("20060102", fmt.Sprintf("%6.0f", periode[i])+"01")
					interim.ETP = etp[i]
				} else {
					tracker.Error(errors.New("ligne invalide, siret manquant ou invalide: " + sirets[i]))
				}

				if !tracker.ErrorInCycle() {
					outputChannel <- interim
				} else {
					event.Debug(tracker.Report("errors"))
				}

			}
			event.Info(tracker.Report("abstract"))
			file.Close()
		}
		close(eventChannel)
		close(outputChannel)
	}()
	return outputChannel, eventChannel
}
