package bdf

import (
	"bufio"
	"encoding/csv"
	"errors"
	"io"
	"os"
	"strings"
	"time"

	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/base"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/marshal"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/misc"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/sfregexp"
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
func Parser(cache marshal.Cache, batch *base.AdminBatch) (chan marshal.Tuple, chan marshal.Event) {
	outputChannel := make(chan marshal.Tuple)
	eventChannel := make(chan marshal.Event)
	event := marshal.Event{
		Code:    "bdfParser",
		Channel: eventChannel,
	}

	go func() {
		for _, path := range batch.Files["bdf"] {
			tracker := gournal.NewTracker(
				map[string]string{"path": path, "batchKey": batch.ID.Key},
				marshal.TrackerReports)

			event.Info(path + ": ouverture " + path)
			ParseFile(viper.GetString("APP_DATA")+path, &cache, batch, &tracker, outputChannel)
			event.Info(tracker.Report("abstract"))
		}

		close(outputChannel)
		close(eventChannel)

	}()
	return outputChannel, eventChannel
}

// ParseFile extrait les tuples depuis le fichier demandé et génère un rapport Gournal.
func ParseFile(filePath string, cache *marshal.Cache, batch *base.AdminBatch, tracker *gournal.Tracker, outputChannel chan marshal.Tuple) {
	filter := marshal.GetSirenFilterFromCache(*cache)
	file, err := os.Open(filePath)
	if err != nil {
		tracker.Add(err)
		return
	}
	defer file.Close()
	reader := csv.NewReader(bufio.NewReader(file))
	reader.Comma = ';'
	reader.LazyQuotes = true
	parseBdfFile(reader, filter, tracker, outputChannel)
}

var field = map[string]int{
	"siren":               0,
	"année":               1,
	"arrêtéBilan":         2,
	"raisonSociale":       3,
	"secteur":             6,
	"poidsFrng":           7,
	"tauxMarge":           8,
	"delaiFournisseur":    9,
	"detteFiscale":        10,
	"financierCourtTerme": 11,
	"fraisFinancier":      12,
}

func parseBdfLine(row []string, tracker *gournal.Tracker, filter map[string]bool) BDF {
	bdf := BDF{}
	bdf.Siren = strings.Replace(row[field["siren"]], " ", "", -1)

	validSiren := sfregexp.RegexpDict["siren"].MatchString(bdf.Siren)
	if !validSiren {
		tracker.Add(errors.New("siren invalide : " + bdf.Siren))
		return BDF{}
	}

	filtered, err := marshal.IsFiltered(bdf.Siren, filter)
	tracker.Add(err)
	if filtered {
		tracker.Add(base.NewFilterNotice())
		return BDF{}
	}

	bdf.Annee, err = misc.ParsePInt(row[field["année"]])
	tracker.Add(err)
	var arrete = row[field["arrêtéBilan"]]
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
	tracker.Add(err)
	bdf.RaisonSociale = row[field["raisonSociale"]]
	bdf.Secteur = row[field["secteur"]]
	if len(row) > field["poidsFrng"] {
		bdf.PoidsFrng, err = misc.ParsePFloat(row[field["poidsFrng"]])
		tracker.Add(err)
	} else {
		bdf.PoidsFrng = nil
	}
	if len(row) > field["tauxMarge"] {
		bdf.TauxMarge, err = misc.ParsePFloat(row[field["tauxMarge"]])
		tracker.Add(err)
	} else {
		bdf.TauxMarge = nil
	}
	if len(row) > field["delaiFournisseur"] {
		bdf.DelaiFournisseur, err = misc.ParsePFloat(row[field["delaiFournisseur"]])
		tracker.Add(err)
	} else {
		bdf.DelaiFournisseur = nil
	}
	if len(row) > field["detteFiscale"] {
		bdf.DetteFiscale, err = misc.ParsePFloat(row[field["detteFiscale"]])
		tracker.Add(err)
	} else {
		bdf.DetteFiscale = nil
	}
	if len(row) > field["financierCourtTerme"] {
		bdf.FinancierCourtTerme, err = misc.ParsePFloat(row[field["financierCourtTerme"]])
		tracker.Add(err)
	} else {
		bdf.FinancierCourtTerme = nil
	}
	if len(row) > field["fraisFinancier"] {
		bdf.FraisFinancier, err = misc.ParsePFloat(row[field["fraisFinancier"]])
		tracker.Add(err)
	} else {
		bdf.FraisFinancier = nil
	}

	return bdf
}

func parseBdfFile(reader *csv.Reader, filter map[string]bool, tracker *gournal.Tracker, outputChannel chan marshal.Tuple) {
	// Lecture en-tête
	_, err := reader.Read()
	tracker.Add(err)

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			tracker.Add(err)
		} else {
			bdf := parseBdfLine(row, tracker, filter)
			var errorInCurrentCycle = tracker.HasErrorInCurrentCycle()
			if !errorInCurrentCycle {
				outputChannel <- bdf
			}
		}
		tracker.Next()
	}
}
