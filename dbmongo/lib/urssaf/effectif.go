package urssaf

import (
	"bufio"
	"encoding/csv"
	"errors"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/base"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/engine"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/marshal"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/misc"

	"github.com/signaux-faibles/gournal"
	//"github.com/globalsign/mgo/bson"
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

// ParserEffectif retourne un channel fournissant des données extraites
func ParserEffectif(cache base.Cache, batch *base.AdminBatch) (chan base.Tuple, chan base.Event) {
	outputChannel := make(chan base.Tuple)
	eventChannel := make(chan base.Event)
	event := base.Event{
		Code:    "effectifParser",
		Channel: eventChannel,
	}
	go func() {
		for _, path := range batch.Files["effectif"] {
			tracker := gournal.NewTracker(
				map[string]string{"path": path, "batchKey": batch.ID.Key},
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

			siretIndex := misc.SliceIndex(len(fields), func(i int) bool { return strings.ToLower(fields[i]) == "siret" })
			compteIndex := misc.SliceIndex(len(fields), func(i int) bool { return strings.ToLower(fields[i]) == "compte" })
			if misc.SliceMin(siretIndex, compteIndex) == -1 {
				event.Critical(path + ": erreur à l'analyse du fichier, abandon, l'un " +
					"des champs obligatoires n'a pu etre trouve:" +
					" siretIndex = " + strconv.Itoa(siretIndex) +
					", compteIndex = " + strconv.Itoa(compteIndex))
				continue
			}
			// Dans quels champs lire l'effectif
			re, _ := regexp.Compile("^eff")
			var effectifFields []string
			var effectifIndexes []int
			for ind, field := range fields {
				if re.MatchString(field) {
					effectifFields = append(effectifFields, field)
					effectifIndexes = append(effectifIndexes, ind)
				}
			}

			periods, err := parseEffectifPeriod(effectifFields)
			if err != nil {
				event.Critical(path + ": erreur a l'analyse du fichier, abandon: " + err.Error())
				continue
			}

			for {
				row, err := reader.Read()
				if err == io.EOF {
					break
				} else if err != nil {
					tracker.Error(err)
					event.Critical(tracker.Report("fatalError"))
					break
				}

				notDigit := regexp.MustCompile("[^0-9]")

				filter, err := marshal.GetSirenFilter(cache, batch)
				// if filter == nil {
				// 	tracker.Error(errors.New("Veuillez spécifier un fichier filtre SIREN"))
				// 	event.Critical(tracker.Report("fatalError"))
				// 	break
				// }

				siret := row[siretIndex]
				filtered, err := marshal.IsFiltered(siret, filter)
				tracker.Error(err)
				if len(siret) == 14 && !filtered {
					for i, j := range effectifIndexes {
						if row[j] != "" {
							noThousandsSep := notDigit.ReplaceAllString(row[j], "")
							e, err := strconv.Atoi(noThousandsSep)
							tracker.Error(err)
							if e > 0 {
								eff := Effectif{
									Siret:        siret,
									NumeroCompte: row[compteIndex],
									Periode:      periods[i],
									Effectif:     e}
								outputChannel <- eff
							}
						}
					}
				}
				if engine.ShouldBreak(tracker, engine.MaxParsingErrors) {
					tracker.Error(base.NewCriticError(errors.New("Parser interrompu: trop d'erreurs"), "fatal"))
					event.Critical(tracker.Report("fatalError"))
					break
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
