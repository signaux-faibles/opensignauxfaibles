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

	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/engine"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/marshal"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/misc"

	"github.com/signaux-faibles/gournal"
	//"github.com/globalsign/mgo/bson"
	"github.com/spf13/viper"
)

// EffectifEnt Urssaf
type EffectifEnt struct {
	Siren       string    `json:"-" bson:"-"`
	Periode     time.Time `json:"periode" bson:"periode"`
	EffectifEnt int       `json:"effectif" bson:"effectif"`
}

// Key _id de l'objet
func (effectifEnt EffectifEnt) Key() string {
	return effectifEnt.Siren
}

// Scope de l'objet
func (effectifEnt EffectifEnt) Scope() string {
	return "entreprise"
}

// Type de l'objet
func (effectifEnt EffectifEnt) Type() string {
	return "effectif_ent"
}

// ParseEffectifEntPeriod Transforme un tableau de périodes telles qu'écrites dans l'entête du tableau d'effectifEnt urssaf en date de début
func parseEffectifEntPeriod(effectifEntPeriods []string) ([]time.Time, error) {
	periods := []time.Time{}
	for _, period := range effectifEntPeriods {
		urssaf := period[3:9]
		date, _ := urssafToPeriod(urssaf)
		periods = append(periods, date.Start)
	}

	return periods, nil
}

// ParserEffectifEnt retourne un channel fournissant des données extraites
func ParserEffectifEnt(cache engine.Cache, batch *engine.AdminBatch) (chan engine.Tuple, chan engine.Event) {
	outputChannel := make(chan engine.Tuple)
	eventChannel := make(chan engine.Event)
	event := engine.Event{
		Code:    "effectifEntParser",
		Channel: eventChannel,
	}
	go func() {
		for _, path := range batch.Files["effectif_ent"] {
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

			sirenIndex := misc.SliceIndex(len(fields), func(i int) bool { return strings.ToLower(fields[i]) == "siren" })

			// Dans quels champs lire l'effectifEnt
			re, _ := regexp.Compile("^eff")
			var effectifEntFields []string
			var effectifEntIndexes []int
			for ind, field := range fields {
				if re.MatchString(field) {
					effectifEntFields = append(effectifEntFields, field)
					effectifEntIndexes = append(effectifEntIndexes, ind)
				}
			}

			periods, err := parseEffectifEntPeriod(effectifEntFields)
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

				siren := row[sirenIndex]
				filtered, err := marshal.IsFiltered(siren, cache, batch)
				tracker.Error(err)
				notDigit := regexp.MustCompile("[^0-9]")
				if len(siren) != 9 {
					tracker.Error(errors.New("Format de siren incorrect : " + siren))
				} else if !filtered {
					for i, j := range effectifEntIndexes {
						if row[j] != "" {
							noThousandsSep := notDigit.ReplaceAllString(row[j], "")
							s, err := strconv.ParseFloat(noThousandsSep, 64)
							tracker.Error(err)
							e := int(s)
							if e > 0 {
								eff := EffectifEnt{
									Siren:       siren,
									Periode:     periods[i],
									EffectifEnt: e,
								}

								outputChannel <- eff
							}
						}
					}
				}

				if engine.ShouldBreak(tracker, engine.MaxParsingErrors) {
					tracker.Error(engine.NewCriticError(errors.New("Parser interrompu: trop d'erreurs"), "fatal"))
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
