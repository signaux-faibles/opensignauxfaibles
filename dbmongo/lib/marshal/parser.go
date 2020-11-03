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

// ParseFile fonction de traitement de données en entrée
type ParseFile func(filePath, *Cache, *base.AdminBatch) OpenFileResult
type filePath = string

// OpenFileResult permet à runParserWithSirenFilter() de savoir si le fichier à
// parser a bien été ouvert, puis de lancer le parsing des lignes.
type OpenFileResult struct {
	Error      error
	ParseLines func(chan base.ParsedLineResult)
	Close      func() error
}

func isValid(tuple base.Tuple) (bool, error) {
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
func ParseFilesFromBatch(cache Cache, batch *base.AdminBatch, parser Parser) (chan base.Tuple, chan Event) {
	outputChannel := make(chan base.Tuple)
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

func runParserWithSirenFilter(parser Parser, filePath string, cache *Cache, batch *base.AdminBatch, tracker *gournal.Tracker, outputChannel chan base.Tuple) {
	filter := GetSirenFilterFromCache(*cache)
	openFileRes := parser.FileParser(filePath, cache, batch)
	// Note: on ne passe plus le tracker aux parseurs afin de garder ici le controle de la numérotation des lignes où les erreurs sont trouvées
	if openFileRes.Error != nil {
		tracker.Add(base.NewFatalError(openFileRes.Error))
	} else {
		parsedLineChan := make(chan base.ParsedLineResult)
		go openFileRes.ParseLines(parsedLineChan)
		for lineResult := range parsedLineChan {
			for _, err := range lineResult.Errors {
				tracker.Add(err)
			}
			for _, tuple := range lineResult.Tuples {
				if _, err := isValid(tuple); err != nil {
					// Si le siret/siren est invalide, on jette le tuple,
					// et on rapporte une erreur seulement si aucune n'a été
					// rapportée par le parseur.
					if len(lineResult.Errors) == 0 {
						tracker.Add(base.NewRegularError(err))
					}
				} else if filter.Skips(tuple.Key()) {
					tracker.Add(base.NewFilterError(errors.New("ligne filtrée")))
				} else {
					outputChannel <- tuple
				}
			}
			tracker.Next()
		}
	}
	// TODO: s'assurer que openFileRes.Close est bien défini avant de l'appeler ?
	if err := openFileRes.Close(); err != nil {
		tracker.Add(base.NewFatalError(err))
	}
}

// GetJSON sérialise un tuple au format JSON.
func GetJSON(tuple base.Tuple) ([]byte, error) {
	return json.MarshalIndent(tuple, "", "  ")
}

// LogProgress affiche le numéro de ligne en cours de parsing, toutes les 2s.
func LogProgress(lineNumber *int) (stop context.CancelFunc) {
	return base.Cron(time.Second*2, func() {
		fmt.Printf("Reading csv line %d\n", *lineNumber)
	})
}
