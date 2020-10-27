package marshal

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/signaux-faibles/gournal"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/base"
	"github.com/spf13/viper"
)

// Parser associe un type de fichier avec sa fonction de parsing.
type Parser = struct {
	FileType   string
	FileParser ParseFile
}

type filePath = string

// ParseFile fonction de traitement de données en entrée
type ParseFile func(filePath, *Cache, *base.AdminBatch, *gournal.Tracker, chan Tuple)

// Tuple unité de donnée à insérer dans un type
type Tuple interface {
	Key() string
	Scope() string
	Type() string
}

// ParseFilesFromBatch parse tous les fichiers spécifiés dans batch pour un parseur donné.
func ParseFilesFromBatch(cache Cache, batch *base.AdminBatch, parser Parser) (chan Tuple, chan Event) {
	filter := GetSirenFilterFromCache(cache)
	outputChannel := make(chan Tuple)
	eventChannel := make(chan Event)
	event := Event{
		Code:    Code(parser.FileType + "_parser"),
		Channel: eventChannel,
	}

	go func() {
		for _, path := range batch.Files[parser.FileType] {
			tracker := gournal.NewTracker(
				map[string]string{"path": path, "batchKey": batch.ID.Key},
				TrackerReports)

			event.Info(path + ": ouverture")
			fullOutputChannel := make(chan Tuple)
			go func() {
				parser.FileParser(viper.GetString("APP_DATA")+path, &cache, batch, &tracker, fullOutputChannel)
				defer close(fullOutputChannel)
			}()
			for tuple := range fullOutputChannel {
				if !filter.Skips(tuple.Key()) {
					outputChannel <- tuple
				} else {
					tracker.Add(base.NewFilterNotice())
				}
			}
			event.Info(tracker.Report("abstract"))
		}
		close(outputChannel)
		close(eventChannel)
	}()

	return outputChannel, eventChannel
}

// GetJSON sérialise un tuple au format JSON.
func GetJSON(tuple Tuple) ([]byte, error) {
	return json.MarshalIndent(tuple, "", "  ")
}

// LogProgress affiche le numéro de ligne en cours de parsing, toutes les 2s.
func LogProgress(lineNumber *int) (stop context.CancelFunc) {
	return base.Cron(time.Second*2, func() {
		fmt.Printf("Reading csv line %d\n", *lineNumber)
	})
}
