package apartdemande

import (
	"dbmongo/lib/engine"
	"dbmongo/lib/misc"
	"time"

	"github.com/chrnin/gournal"
	"github.com/spf13/viper"
	"github.com/tealeg/xlsx"
)

// Periode Période de temps avec un début et une fin

// APDemande Demande d'activité partielle
type APDemande struct {
	key                string
	scope              string
	datatype           string
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
	return apdemande.key
}

// Type de données
func (apdemande APDemande) Type() string {
	return apdemande.datatype
}

// Scope de l'objet
func (apdemande APDemande) Scope() string {
	return apdemande.scope
}

// Parser produit les lignes
func Parser(batch engine.AdminBatch) (chan engine.Tuple, chan engine.Event) {
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

			xlFile, err := xlsx.OpenFile(viper.GetString("APP_DATA") + path)
			if err != nil {
				event.Critical(path + ": erreur à l'ouverture du fichier: " + err.Error())
				return
			}
			event.Info(path + ": ouverture")

			sheet := xlFile.Sheets[0]
			f := make(map[string]int)
			for idx, cell := range sheet.Rows[0].Cells {
				f[cell.Value] = idx
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
			minLength := 0
			for _, field := range fields {
				if i, err := f[field]; err {
					minLength = misc.Max(minLength, i)
				} else {
					event.Critical("Import du fichier " + path + ". " + field + " non trouvé. Abandon.")
					continue
				}
			}
			for _, row := range sheet.Rows[1:] {
				if len(row.Cells) >= minLength && row.Cells[f["ETAB_SIRET"]].Value != "" {
					apdemande := APDemande{}
					apdemande.ID = row.Cells[f["ID_DA"]].Value
					apdemande.Siret = row.Cells[f["ETAB_SIRET"]].Value
					apdemande.EffectifEntreprise, err = misc.ParsePInt(row.Cells[f["EFF_ENT"]].Value)
					tracker.Error(err)
					apdemande.Effectif, err = misc.ParsePInt(row.Cells[f["EFF_ETAB"]].Value)
					tracker.Error(err)
					apdemande.DateStatut, err = misc.ExcelToTime(row.Cells[f["DATE_STATUT"]].Value)
					tracker.Error(err)
					apdemande.Periode = misc.Periode{}
					apdemande.Periode.Start, err = misc.ExcelToTime(row.Cells[f["DATE_DEB"]].Value)
					tracker.Error(err)
					apdemande.Periode.End, err = misc.ExcelToTime(row.Cells[f["DATE_FIN"]].Value)
					tracker.Error(err)
					apdemande.HTA, err = misc.ParsePFloat(row.Cells[f["HTA"]].Value)
					tracker.Error(err)
					apdemande.MTA, err = misc.ParsePFloat(row.Cells[f["MTA"]].Value)
					tracker.Error(err)
					apdemande.EffectifAutorise, err = misc.ParsePInt(row.Cells[f["EFF_AUTO"]].Value)
					tracker.Error(err)
					apdemande.MotifRecoursSE, err = misc.ParsePInt(row.Cells[f["MOTIF_RECOURS_SE"]].Value)
					tracker.Error(err)
					apdemande.HeureConsommee, err = misc.ParsePFloat(row.Cells[f["S_HEURE_CONSOM_TOT"]].Value)
					tracker.Error(err)
					apdemande.EffectifConsomme, err = misc.ParsePInt(row.Cells[f["S_EFF_CONSOM_TOT"]].Value)
					tracker.Error(err)
					apdemande.MontantConsomme, err = misc.ParsePFloat(row.Cells[f["S_MONTANT_CONSOM_TOT"]].Value)
					tracker.Error(err)
					if !tracker.ErrorInCycle() {
						outputChannel <- apdemande
					} else {
						event.Debug(tracker.Report("error"))
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
