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
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/misc"

	"github.com/signaux-faibles/gournal"
	"github.com/spf13/viper"
)

// Procol Procédures collectives, extraction URSSAF
type Procol struct {
	DateEffet    time.Time `json:"date_effet" bson:"date_effet"`
	ActionProcol string    `json:"action_procol" bson:"action_procol"`
	StadeProcol  string    `json:"stade_procol" bson:"stade_procol"`
	Siret        string    `json:"-" bson:"-"`
}

// Key _id de l'objet
func (procol Procol) Key() string {
	return procol.Siret
}

// Scope de l'objet
func (procol Procol) Scope() string {
	return "etablissement"
}

// Type de l'objet
func (procol Procol) Type() string {
	return "procol"
}

// Parser transorme le fichier procol en data
func parseProcol(cache engine.Cache, batch *engine.AdminBatch) (chan engine.Tuple, chan engine.Event) {
	outputChannel := make(chan engine.Tuple)
	eventChannel := make(chan engine.Event)

	event := engine.Event{
		Code:    "procolParser",
		Channel: eventChannel,
	}

	go func() {
		for _, path := range batch.Files["procol"] {
			tracker := gournal.NewTracker(
				map[string]string{"path": path},
				engine.TrackerReports)

			file, err := os.Open(viper.GetString("APP_DATA") + path)
			if err != nil {
				tracker.Error(err)
				event.Critical(tracker.Report("fatalError"))
				continue
			}

			reader := csv.NewReader(bufio.NewReader(file))
			reader.Comma = ';'
			reader.LazyQuotes = true
			fields, err := reader.Read()
			if err != nil {
				tracker.Error(err)
				event.Critical(tracker.Report("fatalError"))
				continue
			}

			dateEffetIndex := misc.SliceIndex(len(fields), func(i int) bool { return strings.ToLower(fields[i]) == "dt_effet" })
			actionStadeIndex := misc.SliceIndex(len(fields), func(i int) bool { return strings.ToLower(fields[i]) == "lib_actx_stdx" })
			siretIndex := misc.SliceIndex(len(fields), func(i int) bool { return strings.ToLower(fields[i]) == "siret" })

			if misc.SliceMin(dateEffetIndex, actionStadeIndex, siretIndex) == -1 {
				tracker.Error(errors.New("format de fichier incorrect"))
				event.Critical(tracker.Report("fatalError"))
				continue
			}

			for {

				row, err := reader.Read()

				if err == io.EOF {
					break
				} else if err != nil {
					// Journal(critical, "importProcol", "Erreur de lecture pendant l'import du fichier "+path+". Abandon.")
					close(outputChannel)
				}
				procol := readLineProcol(
					row,
					&tracker,
					dateEffetIndex,
					siretIndex,
					actionStadeIndex,
				)
				if _, err := strconv.Atoi(row[siretIndex]); err == nil && len(row[siretIndex]) == 14 {
					if !tracker.HasErrorInCurrentCycle() {

						outputChannel <- procol
					} else {
						//event.Debug(tracker.Report("errors"))
					}
				}
				tracker.Next()
			}
			event.Info(tracker.Report("abstract"))
			file.Close()
		}
		close(outputChannel)
		close(eventChannel)
	}()
	return outputChannel, eventChannel
}

func readLineProcol(
	row []string,
	tracker *gournal.Tracker,
	dateEffetIndex int,
	siretIndex int,
	actionStadeIndex int,
) Procol {

	procol := Procol{}
	var err error

	dateFormatee := row[dateEffetIndex]
	dateFormatee = dateFormatee[:3] + strings.ToLower(dateFormatee[4:5]) + dateFormatee[6:]
	procol.DateEffet, err = time.Parse("02Jan2006", row[dateEffetIndex])
	tracker.Error(err)
	procol.Siret = row[siretIndex]
	splitted := strings.Split(strings.ToLower(row[actionStadeIndex]), "_")

	for i, v := range splitted {
		r, err := regexp.Compile("liquidation|redressement|sauvegarde")
		tracker.Error(err)
		if match := r.MatchString(v); match {
			procol.ActionProcol = v
			procol.StadeProcol = strings.Join(append(splitted[:i], splitted[i+1:]...), "_")
			break
		}
	}
	return (procol)
}
