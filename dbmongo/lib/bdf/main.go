package bdf

import (
	"dbmongo/lib/engine"
	"dbmongo/lib/misc"
	"strings"
	"time"

	"github.com/signaux-faibles/gournal"
	"github.com/spf13/viper"
	"github.com/tealeg/xlsx"
)

// BDF Information Banque de France
type BDF struct {
	Siren               string    `json:"siren" bson:"siren"`
	Annee               *int      `json:"annee_bdf" bson:"annee_bdf"`
	ArreteBilan         time.Time `json:"arrete_bilan_bdf" bson:"arrete_bilan_bdf"`
	RaisonSociale       string    `json:"raison_sociale" bson:"raison_sociale"`
	Secteur             string    `json:"secteur" bson:"secteur"`
	PoidsFrng           *float64  `json:"poids_frng" bson:"poids_frng"`
	TauxMarge           *float64  `json:"taux_marge" bson:"taux_marge"`
	DelaiFournisseur    *float64  `json:"delai_fournisseur" bson:"delai_fournisseur"`
	DetteFiscale        *float64  `json:"dette_fiscale" bson:"dette_fiscale"`
	FinancierCourtTerme *float64  `json:"financier_court_terme" bson:"financier_court_terme"`
	FraisFinancier      *float64  `json:"frais_financier" bson:"frais_financier"`
}

// Key id de l'objet
func (bdf BDF) Key() string {
	return bdf.Siren
}

// Type de données
func (bdf BDF) Type() string {
	return "bdf"
}

// Scope de l'objet
func (bdf BDF) Scope() string {
	return "entreprise"
}

// Parser produit les datas BDF à partir des fichiers source
func Parser(batch engine.AdminBatch, filter map[string]bool) (chan engine.Tuple, chan engine.Event) {
	outputChannel := make(chan engine.Tuple)
	eventChannel := make(chan engine.Event)
	event := engine.Event{
		Code:    "bdfParser",
		Channel: eventChannel,
	}

	go func() {
		for _, path := range batch.Files["bdf"] {
			tracker := gournal.NewTracker(
				map[string]string{"path": path},
				engine.TrackerReports)

			xlFile, err := xlsx.OpenFile(viper.GetString("APP_DATA") + path)
			if err != nil {
				tracker.Error(err)
				event.Critical(tracker.Report("fataError"))
				continue
			} else {
				event.Debug(path + ": ouverture " + path)

				for _, sheet := range xlFile.Sheets {
					for _, row := range sheet.Rows[1:] {
						bdf := BDF{}
						bdf.Siren = strings.Replace(row.Cells[0].Value, " ", "", -1)
						bdf.Annee, err = misc.ParsePInt(row.Cells[1].Value)
						tracker.Error(err)
						bdf.ArreteBilan, err = time.Parse("2006-01-02", row.Cells[2].Value)
						tracker.Error(err)
						bdf.RaisonSociale = row.Cells[3].Value
						bdf.Secteur = row.Cells[6].Value
						if len(row.Cells) > 7 {
							bdf.PoidsFrng, err = misc.ParsePFloat(row.Cells[7].Value)
							tracker.Error(err)
						} else {
							bdf.PoidsFrng = nil
						}
						if len(row.Cells) > 8 {
							bdf.TauxMarge, err = misc.ParsePFloat(row.Cells[8].Value)
							tracker.Error(err)
						} else {
							bdf.TauxMarge = nil
						}
						if len(row.Cells) > 9 {
							bdf.DelaiFournisseur, err = misc.ParsePFloat(row.Cells[9].Value)
							tracker.Error(err)
						} else {
							bdf.DelaiFournisseur = nil
						}
						if len(row.Cells) > 10 {
							bdf.DetteFiscale, err = misc.ParsePFloat(row.Cells[10].Value)
							tracker.Error(err)
						} else {
							bdf.DetteFiscale = nil
						}
						if len(row.Cells) > 11 {
							bdf.FinancierCourtTerme, err = misc.ParsePFloat(row.Cells[11].Value)
							tracker.Error(err)
						} else {
							bdf.FinancierCourtTerme = nil
						}
						if len(row.Cells) > 12 {
							bdf.FraisFinancier, err = misc.ParsePFloat(row.Cells[12].Value)
							tracker.Error(err)
						} else {
							bdf.FraisFinancier = nil
						}

						if !tracker.ErrorInCycle() {
							outputChannel <- bdf
						} else {
							event.Debug(tracker.Report("error"))
						}
						tracker.Next()
					}
				}
				event.Info(tracker.Report("abstract"))
			}
		}

		close(outputChannel)
		close(eventChannel)

	}()
	return outputChannel, eventChannel
}
