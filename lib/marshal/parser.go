package marshal

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/signaux-faibles/opensignauxfaibles/lib/base"
	"github.com/signaux-faibles/opensignauxfaibles/lib/sfregexp"
	"github.com/spf13/viper"
)

// Parser fournit les fonctions de parsing d'un type de fichier donné.
type Parser interface {
	GetFileType() string
	Init(cache *Cache, batch *base.AdminBatch) error
	Open(filePath string) error
	ParseLines(parsedLineChan chan ParsedLineResult)
	Close() error
}

// ParsedLineResult est le résultat du parsing d'une ligne.
type ParsedLineResult struct {
	Tuples      []Tuple
	Errors      []error
	FilterError error
}

// AddTuple permet au parseur d'ajouter un tuple extrait depuis la ligne en cours.
func (res *ParsedLineResult) AddTuple(tuple Tuple) {
	if tuple != nil {
		res.Tuples = append(res.Tuples, tuple)
	}
}

// AddRegularError permet au parseur de rapporter une erreur d'extraction.
func (res *ParsedLineResult) AddRegularError(err error) {
	if err != nil {
		res.Errors = append(res.Errors, err)
	}
}

// SetFilterError permet au parseur de rapporter que la ligne doit été filtrée.
func (res *ParsedLineResult) SetFilterError(err error) {
	if err != nil {
		res.FilterError = err
	}
}

// Tuple unité de donnée à insérer dans un type
type Tuple interface {
	Key() string
	Scope() string
	Type() string
}

// ParseFilesFromBatch parse les tuples des fichiers listés dans batch pour le parseur spécifié.
func ParseFilesFromBatch(cache Cache, batch *base.AdminBatch, parser Parser) (chan Tuple, chan Event) {
	outputChannel := make(chan Tuple)
	eventChannel := make(chan Event)
	fileType := parser.GetFileType()
	go func() {
		for _, path := range batch.Files[fileType] {
			eventChannel <- ParseFile(path, parser, batch, cache, outputChannel)
		}
		close(outputChannel)
		close(eventChannel)
	}()
	return outputChannel, eventChannel
}

// ParseFile parse les tuples du fichier spécifié puis retourne un rapport de journal.
func ParseFile(path base.BatchFile, parser Parser, batch *base.AdminBatch, cache Cache, outputChannel chan Tuple) Event {
	tracker := NewParsingTracker()
	fileType := parser.GetFileType()
	filePath := path.Prefix() + viper.GetString("APP_DATA") + path.FilePath()
	err := runParserOnFile(filePath, parser, batch, cache, &tracker, outputChannel)
	if err != nil {
		tracker.AddFatalError(err)
	}
	return CreateReportEvent(fileType, tracker.Report(batch.ID.Key, path.FilePath())) // abstract
}

// runParserOnFile parse les tuples du fichier spécifié, et peut retourner une erreur fatale.
func runParserOnFile(filePath string, parser Parser, batch *base.AdminBatch, cache Cache, tracker *ParsingTracker, outputChannel chan Tuple) error {
	filter := GetSirenFilterFromCache(cache)
	if err := parser.Init(&cache, batch); err != nil {
		return err
	}
	if err := parser.Open(filePath); err != nil {
		return err
	}
	parsedLineChan := make(chan ParsedLineResult)
	go parser.ParseLines(parsedLineChan)
	for lineResult := range parsedLineChan {
		parseTuplesFromLine(lineResult, &filter, tracker, outputChannel)
	}
	return parser.Close()
}

// parseTuplesFromLine extraie les tuples et/ou erreurs depuis une ligne parsée.
func parseTuplesFromLine(lineResult ParsedLineResult, filter *SirenFilter, tracker *ParsingTracker, outputChannel chan Tuple) {
	filterError := lineResult.FilterError
	if filterError != nil {
		tracker.AddFilterError(filterError) // on rapporte le filtrage même si aucun tuple n'est transmis par le parseur
	}
	for _, err := range lineResult.Errors {
		tracker.AddParseError(err)
	}
	for _, tuple := range lineResult.Tuples {
		if filterError != nil {
			continue // l'erreur de filtrage a déjà été rapportée => on se contente de passer au tuple suivant
		} else if _, err := isValid(tuple); err != nil {
			// On rapporte une erreur de siret/siren invalide seulement si aucune autre error n'a été rapportée par le parseur
			if len(lineResult.Errors) == 0 {
				tracker.AddParseError(err)
			}
		} else if filter.Skips(tuple.Key()) {
			tracker.AddFilterError(errors.New("(filtered)"))
		} else {
			outputChannel <- tuple
		}
	}
	tracker.Next()
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
