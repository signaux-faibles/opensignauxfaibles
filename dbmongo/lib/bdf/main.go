package bdf

import (
	"bufio"
	"encoding/csv"
	"io"
	"os"
	"strings"
	"time"

	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/engine"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/misc"
	"github.com/spf13/viper"

	"github.com/signaux-faibles/gournal"
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
func Parser(cache engine.Cache, batch *engine.AdminBatch) (chan engine.Tuple, chan engine.Event) {
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

			file, err := os.Open(viper.GetString("APP_DATA") + path)
			if err != nil {
				tracker.Error(err)
				event.Critical(tracker.Report("fataError"))
				continue
			}

			reader := csv.NewReader(bufio.NewReader(file))
			reader.Comma = ';'
			reader.LazyQuotes = true
			event.Info(path + ": ouverture " + path)

			for {
				row, err := reader.Read()
				if err == io.EOF {
					break
				} else if err != nil {
					tracker.Error(err)
					break
				}
				bdf := BDF{}
				bdf.Siren = strings.Replace(row[0], " ", "", -1)
				bdf.Annee, err = misc.ParsePInt(row[1])
				tracker.Error(err)
				var arrete = row[2]
				arrete = strings.Replace(arrete, "janv", "-01-", -1)
				arrete = strings.Replace(arrete, "JAN", "-01-", -1)
				arrete = strings.Replace(arrete, "févr", "-02-", -1)
				arrete = strings.Replace(arrete, "FEB", "-02-", -1)
				arrete = strings.Replace(arrete, "mars", "-03-", -1)
				arrete = strings.Replace(arrete, "MAR", "-03-", -1)
				arrete = strings.Replace(arrete, "avr", "-04-", -1)
				arrete = strings.Replace(arrete, "APR", "-04-", -1)
				arrete = strings.Replace(arrete, "mai", "-05-", -1)
				arrete = strings.Replace(arrete, "MAY", "-05-", -1)
				arrete = strings.Replace(arrete, "juin", "-06-", -1)
				arrete = strings.Replace(arrete, "JUN", "-06-", -1)
				arrete = strings.Replace(arrete, "juil", "-07-", -1)
				arrete = strings.Replace(arrete, "JUL", "-07-", -1)
				arrete = strings.Replace(arrete, "août", "-08-", -1)
				arrete = strings.Replace(arrete, "AUG", "-08-", -1)
				arrete = strings.Replace(arrete, "sept", "-09-", -1)
				arrete = strings.Replace(arrete, "SEP", "-09-", -1)
				arrete = strings.Replace(arrete, "oct", "-10-", -1)
				arrete = strings.Replace(arrete, "OCT", "-10-", -1)
				arrete = strings.Replace(arrete, "nov", "-11-", -1)
				arrete = strings.Replace(arrete, "NOV", "-11-", -1)
				arrete = strings.Replace(arrete, "déc", "-12-", -1)
				arrete = strings.Replace(arrete, "DEC", "-12-", -1)
				bdf.ArreteBilan, err = time.Parse("02-01-2006", arrete)
				tracker.Error(err)
				bdf.RaisonSociale = row[3]
				bdf.Secteur = row[6]
				if len(row) > 7 {
					bdf.PoidsFrng, err = misc.ParsePFloat(row[7])
					tracker.Error(err)
				} else {
					bdf.PoidsFrng = nil
				}
				if len(row) > 8 {
					bdf.TauxMarge, err = misc.ParsePFloat(row[8])
					tracker.Error(err)
				} else {
					bdf.TauxMarge = nil
				}
				if len(row) > 9 {
					bdf.DelaiFournisseur, err = misc.ParsePFloat(row[9])
					tracker.Error(err)
				} else {
					bdf.DelaiFournisseur = nil
				}
				if len(row) > 10 {
					bdf.DetteFiscale, err = misc.ParsePFloat(row[10])
					tracker.Error(err)
				} else {
					bdf.DetteFiscale = nil
				}
				if len(row) > 11 {
					bdf.FinancierCourtTerme, err = misc.ParsePFloat(row[11])
					tracker.Error(err)
				} else {
					bdf.FinancierCourtTerme = nil
				}
				if len(row) > 12 {
					bdf.FraisFinancier, err = misc.ParsePFloat(row[12])
					tracker.Error(err)
				} else {
					bdf.FraisFinancier = nil
				}

				if !tracker.HasErrorInCurrentCycle() {
					outputChannel <- bdf
				}
				tracker.Next()
			}
			event.Info(tracker.Report("abstract"))
		}

		close(outputChannel)
		close(eventChannel)

	}()
	return outputChannel, eventChannel
}
