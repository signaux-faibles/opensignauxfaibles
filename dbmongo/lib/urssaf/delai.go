package urssaf

import (
	"bufio"
	"encoding/csv"

	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/engine"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/marshal"

	//"errors"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/signaux-faibles/gournal"
	"github.com/spf13/viper"
)

// Delai tuple fichier ursaff
type Delai struct {
	key               string    `hash:"-"`
	NumeroCompte      string    `json:"numero_compte" bson:"numero_compte"`
	NumeroContentieux string    `json:"numero_contentieux" bson:"numero_contentieux"`
	DateCreation      time.Time `json:"date_creation" bson:"date_creation"`
	DateEcheance      time.Time `json:"date_echeance" bson:"date_echeance"`
	DureeDelai        int       `json:"duree_delai" bson:"duree_delai"`
	Denomination      string    `json:"denomination" bson:"denomination"`
	Indic6m           string    `json:"indic_6m" bson:"indic_6m"`
	AnneeCreation     int       `json:"annee_creation" bson:"annee_creation"`
	MontantEcheancier float64   `json:"montant_echeancier" bson:"montant_echeancier"`
	Stade             string    `json:"stade" bson:"stade"`
	Action            string    `json:"action" bson:"action"`
}

// Key _id de l'objet
func (delai Delai) Key() string {
	return delai.key
}

// Scope de l'objet
func (delai Delai) Scope() string {
	return "etablissement"
}

// Type de l'objet
func (delai Delai) Type() string {
	return "delai"
}

// ParserDelai fonction d'extraction des délais
func ParserDelai(cache engine.Cache, batch *engine.AdminBatch) (chan engine.Tuple, chan engine.Event) {
	outputChannel := make(chan engine.Tuple)
	eventChannel := make(chan engine.Event)

	field := map[string]int{
		"NumeroCompte":      2,
		"NumeroContentieux": 3,
		"DateCreation":      4,
		"DateEcheance":      5,
		"DureeDelai":        6,
		"Denomination":      7,
		"Indic6m":           8,
		"AnneeCreation":     9,
		"MontantEcheancier": 10,
		"Stade":             11,
		"Action":            12,
	}

	event := engine.Event{
		Code:    "delaiParser",
		Channel: eventChannel,
	}

	go func() {

		for _, path := range batch.Files["delai"] {
			tracker := gournal.NewTracker(
				map[string]string{"path": path},
				engine.TrackerReports)

			file, err := os.Open(viper.GetString("APP_DATA") + path)

			if err != nil {
				tracker.Error(err)
				event.Critical(tracker.Report("fatalError"))
				break
			} else {
				event.Info(path + ": ouverture")
			}

			reader := csv.NewReader(bufio.NewReader(file))
			reader.Comma = ';'
			reader.Read()
			for {
				row, err := reader.Read()
				if err == io.EOF {
					break
				} else if err != nil {
					tracker.Error(err)
					event.Debug(tracker.Report("invalidLine"))
					break
				}

				date, err := time.Parse("02/01/2006", row[field["DateCreation"]])
				if err != nil {
					tracker.Error(err)
					continue
				}

				if siret, err := marshal.GetSiret(row[field["NumeroCompte"]], &date, cache, batch); err == nil {
					delai, tracker := readLine(row, field, siret, tracker)
					if !tracker.HasErrorInCurrentCycle() {
						outputChannel <- delai
					}
				} else {
					tracker.Error(engine.NewFilterError(err, "filter"))
					continue
				}

				tracker.Next()
			}
			file.Close()
			event.Debug(tracker.Report("abstract"))
		}
		close(outputChannel)
		close(eventChannel)
	}()

	return outputChannel, eventChannel
}

func readLine(row []string, field map[string]int, siret string, tracker gournal.Tracker) (Delai, gournal.Tracker) {
	var err error
	loc, _ := time.LoadLocation("Europe/Paris")
	delai := Delai{}
	delai.key = siret
	delai.NumeroCompte = row[field["NumeroCompte"]]
	delai.NumeroContentieux = row[field["NumeroContentieux"]]
	delai.DateCreation, err = time.ParseInLocation("02/01/2006", row[field["DateCreation"]], loc)
	tracker.Error(err)
	delai.DateEcheance, err = time.ParseInLocation("02/01/2006", row[field["DateEcheance"]], loc)
	tracker.Error(err)
	delai.DureeDelai, err = strconv.Atoi(row[field["DureeDelai"]])
	delai.Denomination = row[field["Denomination"]]
	delai.Indic6m = row[field["Indic6m"]]
	delai.AnneeCreation, err = strconv.Atoi(row[field["AnneeCreation"]])
	tracker.Error(err)
	delai.MontantEcheancier, err = strconv.ParseFloat(strings.Replace(row[field["MontantEcheancier"]], ",", ".", -1), 64)
	tracker.Error(err)
	delai.Stade = row[field["Stade"]]
	delai.Action = row[field["Action"]]
	return delai, tracker
}
