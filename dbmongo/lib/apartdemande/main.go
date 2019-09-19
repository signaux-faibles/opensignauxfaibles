package apartdemande

import (
	"bufio"
	"dbmongo/lib/engine"
	"dbmongo/lib/misc"
	"encoding/csv"
	"io"
	"os"
	"strings"
	"time"

	"github.com/signaux-faibles/gournal"
	"github.com/spf13/viper"
)

// Periode Période de temps avec un début et une fin

// APDemande Demande d'activité partielle
type APDemande struct {
	ID                 string       `json:"id_demande" bson:"id_demande"`
	Siret              string       `json:"-" bson:"-"`
	EffectifEntreprise *int         `json:"effectif_entreprise" bson:"effectif_entreprise"`
	Effectif           *int         `json:"effectif" bson:"effectif"`
	DateStatut         time.Time    `json:"date_statut" bson:"date_statut"`
	Periode            misc.Periode `json:"periode" bson:"periode"`
	HTA                *float64     `json:"hta" bson:"hta"`
	MTA                *float64     `json:"mta" bson:"mta"`
	EffectifAutorise   *int         `json:"effectif_autorise" bson:"effectif_autorise"`
	MotifRecoursSE     *int         `json:"motif_recours_se" bson:"motif_recours_se"`
	HeureConsommee     *float64     `json:"heure_consommee" bson:"heure_consommee"`
	MontantConsomme    *float64     `json:"montant_consommee" bson:"montant_consommee"`
	EffectifConsomme   *int         `json:"effectif_consomme" bson:"effectif_consomme"`
}

// Key id de l'objet
func (apdemande APDemande) Key() string {
	return apdemande.Siret
}

// Type de données
func (apdemande APDemande) Type() string {
	return "apdemande"
}

// Scope de l'objet
func (apdemande APDemande) Scope() string {
	return "etablissement"
}

// Parser produit les lignes
func Parser(batch engine.AdminBatch, filter map[string]bool) (chan engine.Tuple, chan engine.Event) {
	outputChannel := make(chan engine.Tuple)
	eventChannel := make(chan engine.Event)
	event := engine.Event{
		Code:    "parserApdemande",
		Channel: eventChannel,
	}

	go func() {
		defer close(outputChannel)
		defer close(eventChannel)

		for _, path := range batch.Files["apdemande"] {
			tracker := gournal.NewTracker(
				map[string]string{"path": path},
				engine.TrackerReports)

			file, err := os.Open(viper.GetString("APP_DATA") + path)
			if err != nil {
				event.Critical(path + ": erreur à l'ouverture du fichier: " + err.Error())
				return
			}
			reader := csv.NewReader(bufio.NewReader(file))
			reader.Comma = ','
			reader.LazyQuotes = true

			event.Info(path + ": ouverture")

			header, err := reader.Read()
			if err != nil {
				tracker.Error(err)
				event.Debug(tracker.Report("invalidLine"))
				break
			}

			f := make(map[string]int)
			for idx, field := range header {
				f[field] = idx
			}
			fields := []string{
				"ID_DA",
				"ETAB_SIRET",
				"EFF_ENT",
				"EFF_ETAB",
				"DATE_STATUT",
				"DATE_DEB",
				"DATE_FIN",
				"HTA",
				"EFF_AUTO",
				"MOTIF_RECOURS_SE",
				"S_HEURE_CONSOM_TOT",
				"S_EFF_CONSOM_TOT",
			}

			for _, field := range fields {
				if _, found := f[field]; !found {
					event.Critical("Import du fichier " + path + ". " + field + " non trouvé. Abandon.")
					continue
				}
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

				if row[f["ETAB_SIRET"]] != "" {
					apdemande := APDemande{}
					apdemande.ID = row[f["ID_DA"]]
					apdemande.Siret = row[f["ETAB_SIRET"]]
					apdemande.EffectifEntreprise, err = misc.ParsePInt(row[f["EFF_ENT"]])
					tracker.Error(err)
					apdemande.Effectif, err = misc.ParsePInt(row[f["EFF_ETAB"]])
					tracker.Error(err)
					apdemande.DateStatut, err = time.Parse("02/01/2006", row[f["DATE_STATUT"]])
					tracker.Error(err)
					apdemande.Periode = misc.Periode{}
					apdemande.Periode.Start, err = time.Parse("02/01/2006", row[f["DATE_DEB"]])
					tracker.Error(err)
					apdemande.Periode.End, err = time.Parse("02/01/2006", row[f["DATE_FIN"]])
					tracker.Error(err)
					apdemande.HTA, err = misc.ParsePFloat(row[f["HTA"]])
					tracker.Error(err)
					apdemande.MTA, err = misc.ParsePFloat(strings.ReplaceAll(row[f["MTA"]], ",", "."))
					tracker.Error(err)
					apdemande.EffectifAutorise, err = misc.ParsePInt(row[f["EFF_AUTO"]])
					tracker.Error(err)
					apdemande.MotifRecoursSE, err = misc.ParsePInt(row[f["MOTIF_RECOURS_SE"]])
					tracker.Error(err)
					apdemande.HeureConsommee, err = misc.ParsePFloat(row[f["S_HEURE_CONSOM_TOT"]])
					tracker.Error(err)
					apdemande.EffectifConsomme, err = misc.ParsePInt(row[f["S_EFF_CONSOM_TOT"]])
					tracker.Error(err)
					apdemande.MontantConsomme, err = misc.ParsePFloat(strings.ReplaceAll(row[f["S_MONTANT_CONSOM_TOT"]], ",", "."))
					tracker.Error(err)
					if !tracker.ErrorInCycle() {
						outputChannel <- apdemande
					} else {
						// event.Debug(tracker.Report("error"))
					}
				} else {
					event.Debug(tracker.Report("invalidLine"))
				}
				tracker.Next()
			}
			event.Info(tracker.Report("abstract"))
		}
	}()

	return outputChannel, eventChannel
}
