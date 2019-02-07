package urssaf

import (
	"bufio"
	"dbmongo/lib/engine"
	"dbmongo/lib/misc"
	"encoding/csv"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/chrnin/gournal"
	"github.com/globalsign/mgo/bson"
	"github.com/spf13/viper"
)

// Effectif Urssaf
type Effectif struct {
	Siret        string    `json:"-" bson:"-"`
	NumeroCompte string    `json:"numero_compte" bson:"numero_compte"`
	Periode      time.Time `json:"periode" bson:"periode"`
	Effectif     int       `json:"effectif" bson:"effectif"`
}

// Key _id de l'objet
func (effectif Effectif) Key() string {
	return effectif.Siret
}

// Scope de l'objet
func (effectif Effectif) Scope() string {
	return "etablissement"
}

// Type de l'objet
func (effectif Effectif) Type() string {
	return "effectif"
}

// ParseEffectifPeriod Transforme un tableau de périodes telles qu'écrites dans l'entête du tableau d'effectif urssaf en date de début
func parseEffectifPeriod(effectifPeriods []string) ([]time.Time, error) {
	periods := []time.Time{}
	for _, period := range effectifPeriods {
		urssaf := period[3:9]
		date, _ := urssafToPeriod(urssaf)
		periods = append(periods, date.Start)
	}

	return periods, nil
}

// Parser retourne un channel fournissant des données extraites
func parseEffectif(batch engine.AdminBatch, mapping map[string]string) (chan engine.Tuple, chan engine.Event) {
	outputChannel := make(chan engine.Tuple)
	eventChannel := make(chan engine.Event)
	event := engine.Event{
		Code:    "effectifParser",
		Channel: eventChannel,
	}
	go func() {
		for _, path := range batch.Files["effectif"] {
			tracker := gournal.NewTracker(
				map[string]string{"path": path},
				engine.TrackerReports)

			file, err := os.Open(viper.GetString("APP_DATA") + path)
			if err != nil {
				event.Critical(path + ": erreur à l'ouverture du fichier, abandon: " + err.Error())
				continue
			} else {
				event.Info(path + ": ouverture")
			}

			reader := csv.NewReader(bufio.NewReader(file))
			reader.Comma = ';'
			fields, err := reader.Read()

			if err != nil {
				event.Critical(path + ": erreur à la lecture du fichier, abandon: " + err.Error())
				continue
			}

			siretIndex := misc.SliceIndex(len(fields), func(i int) bool { return fields[i] == "SIRET" })
			compteIndex := misc.SliceIndex(len(fields), func(i int) bool { return fields[i] == "compte" })
			boundaryIndex := misc.SliceIndex(len(fields), func(i int) bool { return fields[i] == "ape_ins" })
			if misc.SliceMin(siretIndex, compteIndex, boundaryIndex) == -1 {
				event.Critical(path + ": erreur à l'analyse du fichier, abandon: " + err.Error())
				continue
			}

			periods, err := parseEffectifPeriod(fields[0:boundaryIndex])
			if err != nil {
				event.Critical(path + ": erreur à l'analyse du fichier, abandon: " + err.Error())
				continue
			}

			for {
				row, err := reader.Read()
				if err == io.EOF {
					if tracker.Errors != nil {
						event.Warning(bson.M{
							"errorReport": tracker.Errors,
						})
					}
					event.Info(tracker.Report("abstract"))
					file.Close()
					break
				} else if err != nil {
					tracker.Error(err)
					event.Critical(tracker.Report("fatalError"))
					break
				}

				i := 0
				if len(row[siretIndex]) == 14 {
					for i < boundaryIndex {
						e, err := strconv.Atoi(row[i])
						tracker.Error(err)
						if e > 0 {
							eff := Effectif{
								Siret:        row[siretIndex],
								NumeroCompte: row[compteIndex],
								Periode:      periods[i],
								Effectif:     e}
							outputChannel <- eff
						}
						i++
					}
				}
				tracker.Next()
			}
			file.Close()
		}
		close(outputChannel)
		close(eventChannel)
	}()
	return outputChannel, eventChannel
}
