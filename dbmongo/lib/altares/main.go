package altares

import (
	"encoding/csv"
	"io"
	"os"
	"time"

	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/base"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/engine"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/marshal"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/misc"

	"github.com/signaux-faibles/gournal"
	"github.com/spf13/viper"
)

// Altares Extrait du récapitulatif altarès
type Altares struct {
	DateEffet     time.Time `json:"date_effet" bson:"date_effet"`
	DateParution  time.Time `json:"date_parution" bson:"date_parution"`
	CodeJournal   string    `json:"code_journal" bson:"code_journal"`
	CodeEvenement string    `json:"code_evenement" bson:"code_evenement"`
	Siret         string    `json:"-" bson:"-"`
}

// Key id de l'objet
func (altares Altares) Key() string {
	return altares.Siret
}

// Type de données
func (altares Altares) Type() string {
	return "altares"
}

// Scope de l'objet
func (altares Altares) Scope() string {
	return "etablissement"
}

func sliceIndex(limit int, predicate func(i int) bool) int {
	for i := 0; i < limit; i++ {
		if predicate(i) {
			return i
		}
	}
	return -1
}

// Parser  Altares
func Parser(cache marshal.Cache, batch *base.AdminBatch) (chan marshal.Tuple, chan marshal.Event) {
	// TODO: appliquer filtre, après appel à marshal.GetSirenFilterFromCache(cache)
	outputChannel := make(chan marshal.Tuple)
	eventChannel := make(chan marshal.Event)

	event := marshal.Event{
		Code:    "altaresParser",
		Channel: eventChannel,
	}

	go func() {
		for _, path := range batch.Files["altares"] {
			tracker := gournal.NewTracker(
				map[string]string{"path": path, "batchKey": batch.ID.Key},
				engine.TrackerReports)

			file, err := os.Open(viper.GetString("APP_DATA") + path)
			if err != nil {
				event.Critical(path + ": erreur à l'ouverture: " + err.Error())
				continue
			}
			event.Info(path + ": ouverture")
			reader := csv.NewReader(file)
			reader.Comma = ','
			reader.LazyQuotes = true
			// ligne de titre
			fields, err := reader.Read()

			dateEffetIndex := sliceIndex(len(fields), func(i int) bool { return fields[i] == "Date d'effet" })
			dateParutionIndex := sliceIndex(len(fields), func(i int) bool { return fields[i] == "Date parution" })
			codeJournalIndex := sliceIndex(len(fields), func(i int) bool { return fields[i] == "Code du journal" })
			codeEvenementIndex := sliceIndex(len(fields), func(i int) bool { return fields[i] == "Code de la nature de l'événement" })
			siretIndex := sliceIndex(len(fields), func(i int) bool { return fields[i] == "Siret" })

			if misc.SliceMin(dateEffetIndex, dateParutionIndex, codeJournalIndex, codeEvenementIndex, siretIndex) == -1 {
				event.Critical(path + ": entête non conforme, fichier ignoré")
				file.Close()
				continue
			}

			for {
				row, err := reader.Read()
				if err == io.EOF {
					file.Close()
					break
				} else if err != nil {
					file.Close()
					event.Critical(path + ": abandon suite à un problème de lecture du fichier: " + err.Error())
					break
				}

				dateEffet, err := time.Parse("2006-01-02", row[dateEffetIndex])
				tracker.Add(err)
				dateParution, err := time.Parse("2006-01-02", row[dateParutionIndex])
				tracker.Add(err)

				altares := Altares{
					Siret:         row[siretIndex],
					DateEffet:     dateEffet,
					DateParution:  dateParution,
					CodeJournal:   row[codeJournalIndex],
					CodeEvenement: row[codeEvenementIndex],
				}
				if !tracker.HasErrorInCurrentCycle() {
					outputChannel <- altares
				} else {
					event.Debug(tracker.Report("errors"))
				}
				tracker.Next()
			}
			event.Info(tracker.Report("abstract"))
		}
		close(eventChannel)
		close(outputChannel)
	}()
	return outputChannel, eventChannel
}
