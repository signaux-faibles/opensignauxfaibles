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
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/sfregexp"

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

type periodCol struct {
	dateStart time.Time
	colIndex  int
}

// ParseEffectifEntPeriod Transforme un tableau de périodes telles qu'écrites dans l'entête du tableau d'effectifEnt urssaf en date de début
func parseEffectifEntPeriod(fields []string) []periodCol {
	periods := []periodCol{}
	re, _ := regexp.Compile("^eff")
	for index, field := range fields {
		if re.MatchString(field) {
			date, _ := marshal.UrssafToPeriod(field[3:9])
			periods = append(periods, periodCol{dateStart: date.Start, colIndex: index})
		}
	}
	return periods
}

// ParserEffectifEnt retourne un channel fournissant des données extraites
func ParserEffectifEnt(cache marshal.Cache, batch *base.AdminBatch) (chan marshal.Tuple, chan marshal.Event) {
	outputChannel := make(chan marshal.Tuple)
	eventChannel := make(chan marshal.Event)
	event := marshal.Event{
		Code:    "effectifEntParser",
		Channel: eventChannel,
	}
	filter := marshal.GetSirenFilterFromCache(cache)
	go func() {
		for _, path := range batch.Files["effectif_ent"] {

			file, err := os.Open(viper.GetString("APP_DATA") + path)
			if err != nil {
				event.Critical(path + ": erreur à l'ouverture du fichier, abandon: " + err.Error())
			} else {
				tracker := gournal.NewTracker(
					map[string]string{"path": path, "batchKey": batch.ID.Key},
					engine.TrackerReports)

				event.Info(path + ": ouverture")
				reader := csv.NewReader(bufio.NewReader(file))
				reader.Comma = ';'

				parseEffectifEntFile(reader, filter, &tracker, outputChannel)
				file.Close()
				event.Debug(tracker.Report("abstract"))
			}
		}
		close(outputChannel)
		close(eventChannel)
	}()
	return outputChannel, eventChannel
}

func parseEffectifEntFile(reader *csv.Reader, filter map[string]bool, tracker *gournal.Tracker, outputChannel chan marshal.Tuple) {
	fields, err := reader.Read()
	if err != nil {
		tracker.Add(err)
		return
	}

	var idx = colMapping{
		"siren": misc.SliceIndex(len(fields), func(i int) bool { return strings.ToLower(fields[i]) == "siren" }),
	}

	// Dans quels champs lire l'effectifEnt
	periods := parseEffectifEntPeriod(fields)

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			tracker.Add(err)
			break
		}

		effectifs := parseEffectifEntLine(periods, row, idx, filter, tracker)
		for _, eff := range effectifs {
			outputChannel <- eff
		}
		tracker.Next()
	}
}

func parseEffectifEntLine(periods []periodCol, row []string, idx colMapping, filter map[string]bool, tracker *gournal.Tracker) []EffectifEnt {
	var effectifs = []EffectifEnt{}
	siren := row[idx["siren"]]
	filtered, err := marshal.IsFiltered(siren, filter)
	tracker.Add(err)
	if len(siren) != 9 {
		tracker.Add(errors.New("Format de siren incorrect : " + siren))
	} else if !filtered {
		for _, period := range periods {
			value := row[period.colIndex]
			if value != "" {
				noThousandsSep := sfregexp.RegexpDict["notDigit"].ReplaceAllString(value, "")
				s, err := strconv.ParseFloat(noThousandsSep, 64)
				tracker.Add(err)
				e := int(s)
				if e > 0 {
					effectifs = append(effectifs, EffectifEnt{
						Siren:       siren,
						Periode:     period.dateStart,
						EffectifEnt: e,
					})
				}
			}
		}
	}
	return effectifs
}
