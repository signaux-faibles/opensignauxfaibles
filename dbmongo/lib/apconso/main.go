package apconso

import (
	"encoding/csv"
	"io"
	"os"
	"time"

	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/base"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/engine"
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

// Parser produit des lignes de consommation d'activité partielle
func Parser(cache base.Cache, batch *base.AdminBatch) (chan base.Tuple, chan base.Event) {
	outputChannel := make(chan base.Tuple)
	eventChannel := make(chan base.Event)
	event := base.Event{
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

			fields, err := reader.Read()
			if err != nil {
				tracker.Error(err)
				event.Debug(tracker.Report("invalidLine"))
				break
			}
			idxID := misc.SliceIndex(35, func(i int) bool { return fields[i] == "ID_DA" })
			idxSiret := misc.SliceIndex(35, func(i int) bool { return fields[i] == "ETAB_SIRET" })
			idxPeriode := misc.SliceIndex(35, func(i int) bool { return fields[i] == "MOIS" })
			idxHeureConsommee := misc.SliceIndex(35, func(i int) bool { return fields[i] == "HEURES" })
			idxMontants := misc.SliceIndex(35, func(i int) bool { return fields[i] == "MONTANTS" })
			idxEffectifs := misc.SliceIndex(35, func(i int) bool { return fields[i] == "EFFECTIFS" })

			if misc.SliceMin(idxID, idxSiret, idxPeriode, idxHeureConsommee, idxMontants, idxEffectifs) == -1 {
				event.Critical(path + ": entête non conforme, fichier ignoré")
				continue
			}

			for {
				row, err := reader.Read()

				if err == io.EOF {
					break
				} else if err != nil {
					tracker.Error(err)
					break
				}

				if len(row) > 0 {
					apconso := APConso{}

					apconso.ID = row[idxID]
					apconso.Siret = row[idxSiret]
					apconso.Periode, err = time.Parse("01/2006", row[idxPeriode])
					tracker.Error(err)
					apconso.HeureConsommee, err = misc.ParsePFloat(row[idxHeureConsommee])
					tracker.Error(err)
					apconso.Montant, err = misc.ParsePFloat(row[idxMontants])
					tracker.Error(err)
					apconso.Effectif, err = misc.ParsePInt(row[idxEffectifs])
					tracker.Error(err)

					if !tracker.HasErrorInCurrentCycle() && apconso.Siret != "" {
						outputChannel <- apconso
					} else {
						// event.Debug(tracker.Report("errors"))
					}
					tracker.Next()
				}
			}
			event.Info(tracker.Report("abstract"))
		}
	}()

	return outputChannel, eventChannel
}
