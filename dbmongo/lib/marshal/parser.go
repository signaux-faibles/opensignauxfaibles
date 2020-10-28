package marshal

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/signaux-faibles/gournal"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/base"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/sfregexp"
	"github.com/spf13/viper"
)

// Parser associe un type de fichier avec sa fonction de parsing.
type Parser = struct {
	FileType   string
	FileParser ParseFile
}

type filePath = string

// LineParser est une fonction permettant de parser la prochaine ligne d'un fichier
type LineParser func() []Tuple

// ParseFile fonction de traitement de données en entrée
type ParseFile func(filePath, *Cache, *base.AdminBatch, *gournal.Tracker) LineParser

// Tuple unité de donnée à insérer dans un type
type Tuple interface {
	Key() string
	Scope() string
	Type() string
}

func isValid(tuple Tuple) (bool, error) {
	scope := tuple.Scope()
	key := tuple.Key()
	if scope == "entreprise" {
		if !sfregexp.ValidSiren(key) {
			return false, errors.New("siren invalide : " + key)
		}
		return true, nil
	} else if scope == "etablissement" {
		if !sfregexp.ValidSiret(key) {
			return false, errors.New("siret invalide : " + key)
		}
		return true, nil
	}
	return false, errors.New("tuple sans scope")
}

// ParseFilesFromBatch parse tous les fichiers spécifiés dans batch pour un parseur donné.
func ParseFilesFromBatch(cache Cache, batch *base.AdminBatch, parser Parser) (chan Tuple, chan Event) {
	outputChannel := make(chan Tuple)
	eventChannel := make(chan Event)
	event := Event{
		Code:    Code(parser.FileType),
		Channel: eventChannel,
	}
	go func() {
		for _, path := range batch.Files[parser.FileType] {
			tracker := gournal.NewTracker(
				map[string]string{"path": path, "batchKey": batch.ID.Key},
				TrackerReports)
			filePath := viper.GetString("APP_DATA") + path
			runParserWithSirenFilter(parser, filePath, &cache, batch, &tracker, outputChannel)
			event.Info(tracker.Report("abstract"))
		}
		close(outputChannel)
		close(eventChannel)
	}()
	return outputChannel, eventChannel
}

func runParserWithSirenFilter(parser Parser, filePath string, cache *Cache, batch *base.AdminBatch, tracker *gournal.Tracker, outputChannel chan Tuple) {
	filter := GetSirenFilterFromCache(*cache)
	parseLine := parser.FileParser(filePath, cache, batch, tracker)
	if parseLine == nil {
		return
	}
	for {
		tuples := parseLine()
		if tuples == nil {
			break
		}
		for _, tuple := range tuples {
			if _, err := isValid(tuple); err != nil {
				tracker.Add(err)
			} else if filter.Skips(tuple.Key()) {
				tracker.Add(base.NewFilterNotice())
			} else {
				outputChannel <- tuple
			}
		}
		tracker.Next()
	}
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
