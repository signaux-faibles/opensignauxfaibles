package marshal

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"opensignauxfaibles/lib/base"
	"opensignauxfaibles/lib/sfregexp"
)

// Parser spécifie les fonctions qui doivent être implémentées par chaque parseur de fichier.
type Parser interface {
	Type() string
	Init(cache *Cache, batch *base.AdminBatch) error
	Open(filePath base.BatchFile) error
	ParseLines(parsedLineChan chan ParsedLineResult)
	Close() error
}

// ParsedLineResult est le résultat du parsing d'une ligne.
// Une même ligne peut générer plusieurs tuples.
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

// Tuple spécifie les fonctions que chaque parseur doit implémenter pour ses tuples.
type Tuple interface {
	Key() string   // entité définie par le tuple: numéro SIRET ou SIREN
	Scope() string // type d'entité: "entreprise" ou "etablissement"
	Type() string  // identifiant du parseur qui a extrait ce tuple, ex: "apconso"
}

// ParseFilesFromBatch parse les tuples des fichiers listés dans batch pour le parseur spécifié.
func ParseFilesFromBatch(cache Cache, batch *base.AdminBatch, parser Parser) (chan Tuple, chan Event) {
	outputChannel := make(chan Tuple)
	eventChannel := make(chan Event)
	fileType := parser.Type()
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
	logger := slog.With("batch", batch.ID.Key, "parser", parser.Type(), "filename", path.RelativePath())
	logger.Debug("parsing file")

	tracker := NewParsingTracker()
	fileType := parser.Type()

	err := runParserOnFile(path, parser, batch, cache, &tracker, outputChannel)
	if err != nil {
		tracker.AddFatalError(err)
	}

	logger.Debug("end of file parsing")

	return CreateReportEvent(fileType, tracker.Report(batch.ID.Key, path.RelativePath())) // abstract
}

// runParserOnFile parse les tuples du fichier spécifié, et peut retourner une erreur fatale.
func runParserOnFile(filePath base.BatchFile, parser Parser, batch *base.AdminBatch, cache Cache, tracker *ParsingTracker, outputChannel chan Tuple) error {
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
		tracker.Next()
	}
	return parser.Close()
}

// parseTuplesFromLine extraie les tuples et/ou erreurs depuis une ligne parsée.
func parseTuplesFromLine(lineResult ParsedLineResult, filter *SirenFilter, tracker *ParsingTracker, outputChannel chan Tuple) {
	filterError := lineResult.FilterError
	if filterError != nil {
		tracker.AddFilterError(filterError) // on rapporte le filtrage même si aucun tuple n'est transmis par le parseur
		return
	}
	for _, err := range lineResult.Errors {
		tracker.AddParseError(err)
	}
	for _, tuple := range lineResult.Tuples {
		if _, err := isValid(tuple); err != nil {
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
}

// LogProgress affiche le numéro de ligne en cours de parsing, toutes les 2s.
func LogProgress(lineNumber *int) (stop context.CancelFunc) {
	return base.Cron(time.Minute*1, func() {
		slog.Info("Lis une ligne du fichier csv", slog.Int("line", *lineNumber))
	})
}

// idValid vérifie que la clé (Key) d'un Tuple est valide, selon le type d'entité (Scope) qu'il représente.
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
