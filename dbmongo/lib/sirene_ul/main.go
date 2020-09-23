package sireneul

import (
	//"bufio"
	"encoding/csv"
	"errors"
	"io"
	"os"
	"time"

	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/base"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/engine"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/marshal"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/sfregexp"

	"github.com/signaux-faibles/gournal"
	"github.com/spf13/viper"
)

// SireneUL informations sur les entreprises
type SireneUL struct {
	Siren               string     `json:"siren,omitempty"         bson:"siren,omitempty"`
	Nic                 string     `json:"nic,omitempty"           bson:"nic,omitempty"`
	RaisonSociale       string     `json:"raison_sociale"          bson:"raison_sociale"`
	Prenom1UniteLegale  string     `json:"prenom1_unite_legale,omitempty"      bson:"prenom1_unite_legale,omitempty"`
	Prenom2UniteLegale  string     `json:"prenom2_unite_legale,omitempty"      bson:"prenom2_unite_legale,omitempty"`
	Prenom3UniteLegale  string     `json:"prenom3_unite_legale,omitempty"      bson:"prenom3_unite_legale,omitempty"`
	Prenom4UniteLegale  string     `json:"prenom4_unite_legale,omitempty"      bson:"prenom4_unite_legale,omitempty"`
	NomUniteLegale      string     `json:"nom_unite_legale,omitempty"          bson:"nom_unite_legale,omitempty"`
	NomUsageUniteLegale string     `json:"nom_usage_unite_legale,omitempty"     bson:"nom_usage_unite_legale,omitempty"`
	CodeStatutJuridique string     `json:"statut_juridique"        bson:"statut_juridique"`
	Creation            *time.Time `json:"date_creation,omitempty" bson:"date_creation,omitempty"`
}

// Key id de l'objet
func (sirene_ul SireneUL) Key() string {
	return sirene_ul.Siren
}

// Type de données
func (sirene_ul SireneUL) Type() string {
	return "sirene_ul"
}

// Scope de l'objet
func (sirene_ul SireneUL) Scope() string {
	return "entreprise"
}

// Parser produit les données sirene à partir du fichier geosirene
func Parser(cache base.Cache, batch *base.AdminBatch) (chan base.Tuple, chan base.Event) {
	outputChannel := make(chan base.Tuple)
	eventChannel := make(chan base.Event)

	event := base.Event{
		Code:    "sireneULParser",
		Channel: eventChannel,
	}

	go func() {
		for _, path := range batch.Files["sirene_ul"] {
			tracker := gournal.NewTracker(
				map[string]string{"path": path, "batchKey": batch.ID.Key},
				engine.TrackerReports)

			file, err := os.Open(viper.GetString("APP_DATA") + path)
			if err != nil {
				tracker.Error(err)
				tracker.Report("fatalError")
			}
			event.Info(path + ": ouverture")
			reader := csv.NewReader(file)
			reader.Comma = ','
			reader.LazyQuotes = true

			_, _ = reader.Read()

			for {
				row, err := reader.Read()
				if err == io.EOF {
					break
				} else if err != nil {
					tracker.Error(err)
					event.Critical(tracker.Report("fatalError"))
					break
				}

				validSiren := sfregexp.RegexpDict["siren"].MatchString(row[0])
				if !validSiren {
					tracker.Error(errors.New("siren invalide : " + row[0]))
					continue // TODO: exécuter tracker.Next() un fois le TODO ci-dessous traité.
				}
				filter, err := marshal.GetSirenFilter(cache, batch)
				// if filter == nil {
				// 	tracker.Error(errors.New("Veuillez spécifier un fichier filtre SIREN"))
				// 	event.Critical(tracker.Report("fatalError"))
				// 	break
				// }
				if err != nil {
					tracker.Error(err)
					event.Critical(tracker.Report("fatalError"))
					break
				}
				filtered, err := marshal.IsFiltered(row[0], filter)
				if err != nil {
					tracker.Error(err)
				}
				if !filtered {
					sireneul := readLineEtablissement(row, &tracker)
					outputChannel <- sireneul
					tracker.Next() // TODO: garantir que le compteur de lignes
					// correspond au nombre de lignes du fichier. => appeler même si le
					// siren est filtré
				}
			}
			file.Close()
			event.Info(tracker.Report("abstract"))
		}
		close(outputChannel)
		close(eventChannel)
	}()

	return outputChannel, eventChannel
}

func readLineEtablissement(row []string, tracker *gournal.Tracker) SireneUL {
	sireneul := SireneUL{}
	sireneul.Siren = row[0]
	sireneul.RaisonSociale = row[23]
	sireneul.Prenom1UniteLegale = row[6]
	sireneul.Prenom2UniteLegale = row[7]
	sireneul.Prenom3UniteLegale = row[8]
	sireneul.Prenom4UniteLegale = row[9]
	sireneul.NomUniteLegale = row[21]
	sireneul.NomUsageUniteLegale = row[22]
	sireneul.CodeStatutJuridique = row[27]
	creation, err := time.Parse("2006-01-02", row[3])
	if err == nil {
		sireneul.Creation = &creation
	}
	tracker.Error(err)
	return sireneul
}
