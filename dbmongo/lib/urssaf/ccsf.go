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

			comptes, err := marshal.GetCompteSiretMapping(cache, batch, marshal.OpenAndReadSiretMapping)
			if err != nil {
				tracker.Add(err)
				return
			}

			reader := csv.NewReader(bufio.NewReader(file))
			reader.Comma = ';'
			reader.Read()

			parseCcsfFile(reader, &comptes, &tracker, outputChannel)
			event.Info(tracker.Report("abstract"))

			file.Close()
		}
		close(outputChannel)
		close(eventChannel)
	}()
	return outputChannel, eventChannel
}

var idx = map[string]int{
	"NumeroCompte":   2,
	"DateTraitement": 3,
	"Stade":          4,
	"Action":         5,
}

func parseCcsfFile(reader *csv.Reader, comptes *marshal.Comptes, tracker *gournal.Tracker, outputChannel chan marshal.Tuple) {
	for {
		r, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			tracker.Add(err)
			continue
		}

		ccsf := parseCcsfLine(r, tracker, idx, comptes)
		if !tracker.HasErrorInCurrentCycle() {
			outputChannel <- ccsf
		}
		tracker.Next()
	}
}

func parseCcsfLine(r []string, tracker *gournal.Tracker, idx colMapping, comptes *marshal.Comptes) CCSF {
	var err error
	ccsf := CCSF{}
	if len(r) >= 4 {
		ccsf.Action = r[idx["Action"]]
		ccsf.Stade = r[idx["Stade"]]
		ccsf.DateTraitement, err = marshal.UrssafToDate(r[idx["DateTraitement"]])
		tracker.Add(err)
		if err != nil {
			return ccsf
		}

		ccsf.key, err = marshal.GetSiretFromComptesMapping(r[idx["NumeroCompte"]], &ccsf.DateTraitement, *comptes)
		if err != nil {
			// Compte filtré
			tracker.Add(base.NewFilterError(err))
			return ccsf
		}
		ccsf.NumeroCompte = r[idx["NumeroCompte"]]

	} else {
		tracker.Add(errors.New("Ligne non conforme, moins de 4 champs"))
	}
	return ccsf
}
