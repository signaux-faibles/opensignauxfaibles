package urssaf

import (
	"bufio"
	"dbmongo/lib/engine"
	"encoding/csv"
	"errors"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/chrnin/gournal"
	"github.com/spf13/viper"
)

// Delai tuple fichier ursaff
type Delai struct {
	key               string    `hash:"-"`
	NumeroCompte      string    `json:"numero_compte" bson:"numero_compte"`
	NumeroContentieux string    `json:"numero_contentieux" bson:"numero_contentieux"`
	DateCreation      time.Time `json:"date_creation" bson:"date_creation"`
	DateEcheanche     time.Time `json:"date_echeance" bson:"date_echeance"`
	DureeDelai        int       `json:"duree_delai" bson:"duree_delai"`
	Denomination      string    `json:"denomination" bson:"denomination"`
	Indic6m           string    `json:"indic_6m" bson:"indic_6m"`
	AnneeCreation     int       `json:"annee_creation" bson:"annee_creation"`
	MontantEcheancier float64   `json:"montant_echeancier" bson:"montant_echeancier"`
	NumeroStructure   string    `json:"numero_structure" bson:"numero_structure"`
	Stade             string    `json:"stade" bson:"stade"`
	Action            string    `json:"action" bson:"action"`
}

// Key _id de l'objet
func (delai Delai) Key() string {
	return delai.key
}

// Scope de l'objet
func (delai Delai) Scope() string {
	return "etablissemnt"
}

// Type de l'objet
func (delai Delai) Type() string {
	return "delai"
}

// Parser fonction d'extraction des délais
func parseDelai(batch engine.AdminBatch, mapping Comptes) (chan engine.Tuple, chan engine.Event) {
	outputChannel := make(chan engine.Tuple)
	eventChannel := make(chan engine.Event)

	field := map[string]int{
		"NumeroCompte":      0,
		"NumeroContentieux": 1,
		"DateCreation":      2,
		"DateEcheanche":     3,
		"DureeDelai":        4,
		"Denomination":      5,
		"Indic6m":           6,
		"AnneeCreation":     7,
		"MontantEcheancier": 8,
		"NumeroStructure":   9,
		"Stade":             10,
		"Action":            11,
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
					} else {
            date, err := time.Parse("2006-01-02", row[field["DateCreation"]])
            if err != nil { date = time.Now() }
						if siret, err := mapping.GetSiret(row[field["NumeroCompte"]], date); err == nil {
							delai := Delai{}
							delai.key = siret
							delai.NumeroCompte = row[field["NumeroCompte"]]
							delai.NumeroContentieux = row[field["NumeroContentieux"]]
							delai.DateCreation, err = time.Parse("2006-01-02", row[field["DateCreation"]])
							tracker.Error(err)
							delai.DateEcheanche, err = time.Parse("2006-01-02", row[field["DateEcheanche"]])
							tracker.Error(err)
							delai.DureeDelai, err = strconv.Atoi(row[field["DureeDelai"]])
							delai.Denomination = row[field["Denomination"]]
							delai.Indic6m = row[field["Indic6m"]]
							delai.AnneeCreation, err = strconv.Atoi(row[field["AnneeCreation"]])
							tracker.Error(err)
							delai.MontantEcheancier, err = strconv.ParseFloat(strings.Replace(row[field["MontantEcheancier"]], ",", ".", -1), 64)
							tracker.Error(err)
							delai.NumeroStructure = row[field["NumeroStructure"]]
							delai.Stade = row[field["Stade"]]
							delai.Action = row[field["Action"]]
							if !tracker.ErrorInCycle() {
								outputChannel <- delai
							} else {
								event.Debug(tracker.Report("errors"))
							}
						} else {
              tracker.Error(errors.New("Compte absent du mapping : " + row[field["NumeroCompte"]]))
							event.Debug(tracker.Report("invalidLine"))
						}
					}
					tracker.Next()
				}
				file.Close()
				event.Debug(tracker.Report("abstract"))
			}
		}
		close(outputChannel)
		close(eventChannel)
	}()

	return outputChannel, eventChannel
}
