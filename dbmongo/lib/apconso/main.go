package apconso

import (
	"encoding/csv"
	"errors"
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

// APConso Consommation d'activité partielle
type APConso struct {
	ID             string    `json:"id_conso"         bson:"id_conso"`
	Siret          string    `json:"-"                bson:"-"`
	HeureConsommee *float64  `json:"heure_consomme"   bson:"heure_consomme"`
	Montant        *float64  `json:"montant"          bson:"montant"`
	Effectif       *int      `json:"effectif"         bson:"effectif"`
	Periode        time.Time `json:"periode"          bson:"periode"`
}

// Key id de l'objet
func (apconso APConso) Key() string {
	return apconso.Siret
}

// Type de données
func (apconso APConso) Type() string {
	return "apconso"
}

// Scope de l'objet
func (apconso APConso) Scope() string {
	return "etablissement"
}

type colMapping map[string]int

// Parser produit des lignes de consommation d'activité partielle
func Parser(cache marshal.Cache, batch *base.AdminBatch) (chan marshal.Tuple, chan marshal.Event) {
	outputChannel := make(chan marshal.Tuple)
	eventChannel := make(chan marshal.Event)
	event := marshal.Event{
		Code:    "parserApconso",
		Channel: eventChannel,
	}

	go func() {
		defer close(outputChannel)
		defer close(eventChannel)

		for _, path := range batch.Files["apconso"] {
			tracker := gournal.NewTracker(
				map[string]string{"path": path, "batchKey": batch.ID.Key},
				engine.TrackerReports)

			file, err := os.Open(viper.GetString("APP_DATA") + path)
			if err != nil {
				event.Critical(path + ": erreur à l'ouverture du fichier, abandon:" + err.Error())
				continue
			}

			reader := csv.NewReader(file)
			reader.Comma = ','

			event.Info(path + ": ouverture")
			parseApConsoFile(reader, &tracker, outputChannel)
			event.Info(tracker.Report("abstract"))
		}
	}()

	return outputChannel, eventChannel
}

func parseApConsoFile(reader *csv.Reader, tracker *gournal.Tracker, outputChannel chan marshal.Tuple) {
	fields, err := reader.Read()
	if err != nil {
		tracker.Add(err)
		return
	}
	var idx = colMapping{}
	idx["ID"] = misc.SliceIndex(35, func(i int) bool { return fields[i] == "ID_DA" })
	idx["Siret"] = misc.SliceIndex(35, func(i int) bool { return fields[i] == "ETAB_SIRET" })
	idx["Periode"] = misc.SliceIndex(35, func(i int) bool { return fields[i] == "MOIS" })
	idx["HeureConsommee"] = misc.SliceIndex(35, func(i int) bool { return fields[i] == "HEURES" })
	idx["Montants"] = misc.SliceIndex(35, func(i int) bool { return fields[i] == "MONTANTS" })
	idx["Effectifs"] = misc.SliceIndex(35, func(i int) bool { return fields[i] == "EFFECTIFS" })

	if misc.SliceMin(idx["ID"], idx["Siret"], idx["Periode"], idx["HeureConsommee"], idx["Montants"], idx["Effectifs"]) == -1 {
		tracker.Add(errors.New("entête non conforme, fichier ignoré"))
		return
	}

	for {
		row, err := reader.Read()

		if err == io.EOF {
			break
		} else if err != nil {
			tracker.Add(err)
			break
		}

		// TODO: filtrer et/ou valider siret ?

		if len(row) > 0 {
			apconso := parseApConsoLine(row, tracker, idx)

			if !tracker.HasErrorInCurrentCycle() && apconso.Siret != "" {
				outputChannel <- apconso
			} /* else {
				event.Debug(tracker.Report("errors"))
			} */

			tracker.Next() // TODO: executer même si len(row) === 0 ?
		}

	}
}

func parseApConsoLine(row []string, tracker *gournal.Tracker, idx colMapping) APConso {
	apconso := APConso{}
	apconso.ID = row[idx["ID"]]
	apconso.Siret = row[idx["Siret"]]
	var err error
	apconso.Periode, err = time.Parse("01/2006", row[idx["Periode"]])
	tracker.Add(err)
	apconso.HeureConsommee, err = misc.ParsePFloat(row[idx["HeureConsommee"]])
	tracker.Add(err)
	apconso.Montant, err = misc.ParsePFloat(row[idx["Montants"]])
	tracker.Add(err)
	apconso.Effectif, err = misc.ParsePInt(row[idx["Effectifs"]])
	tracker.Add(err)
	return apconso
}
