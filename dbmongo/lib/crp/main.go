package crp

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"regexp"

	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/engine"

	"github.com/signaux-faibles/gournal"
	"github.com/spf13/viper"
)

// CRP Ligne de fichier CRP
type CRP struct {
	SiretRenseigne string `json:"siret" bson:"siret"`
	Difficulte     string `json:"difficulte" bson:"difficulte"`
	EtatDossier    string `json:"etat_dossier" bson:"etat_dossier"`
	Actions        string `json:"actions" bson:"actions"`
	Statut         string `json:"statut" bson:"statut"`
	Fichier        string `json:"fichier" bson:"fichier"`
}

// Key id de l'objet
func (crp CRP) Key() string {
	reg, _ := regexp.Compile("[^0-9]+")
	key := reg.ReplaceAllString(crp.SiretRenseigne, "")
	if len(key) <= 9 {
		key = fmt.Sprintf("%09s", key)
	} else {
		key = fmt.Sprintf("%014s", key)[0:9]
	}
	return key
}

// Type de données
func (crp CRP) Type() string {
	return "crp"
}

// Scope de l'objet
func (crp CRP) Scope() string {
	return "entreprise"
}

// Parser produit des lignes de consommation d'activité partielle
func Parser(batch engine.AdminBatch, filter map[string]bool) (chan engine.Tuple, chan engine.Event) {
	outputChannel := make(chan engine.Tuple)
	eventChannel := make(chan engine.Event)
	event := engine.Event{
		Code:    "parserApconso",
		Channel: eventChannel,
	}

	go func() {
		defer close(outputChannel)
		defer close(eventChannel)

		for _, path := range batch.Files["crp"] {
			tracker := gournal.NewTracker(
				map[string]string{"path": path},
				engine.TrackerReports)

			file, err := os.Open(viper.GetString("APP_DATA") + path)
			if err != nil {
				event.Critical(path + ": erreur à l'ouverture du fichier, abandon:" + err.Error())
				continue
			}

			reader := csv.NewReader(bufio.NewReader(file))
			reader.Comma = ','

			event.Info(path + ": ouverture")

			// 2 lignes de titre
			_, err = reader.Read()
			_, err = reader.Read()

			if err != nil {
				tracker.Error(err)
				event.Debug(tracker.Report("invalidLine"))
				break
			}

			for {
				row, err := reader.Read()

				if err == io.EOF {
					break
				} else if err != nil {
					tracker.Error(err)
					event.Debug(tracker.Report("invalidLine"))
					break
				}

				if len(row) > 0 {
					crp := CRP{}
					crp.SiretRenseigne = row[1]
					crp.EtatDossier = row[16]
					crp.Difficulte = row[22]
					crp.Actions = row[24]
					crp.Statut = row[10]
					crp.Fichier = path

					if !tracker.HasErrorInCurrentCycle() && crp.Key() != "000000000" {
						outputChannel <- crp
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
