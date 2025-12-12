package engine

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"time"

	"opensignauxfaibles/lib/sfregexp"
)

// Parser crée une instance de parser `ParserInst` avec le contenu spécifique
// du `io.Reader`
type Parser interface {
	New(io.Reader) ParserInst
	Type() ParserType
}

// ParserInst extrait des données structurées d'un contenu, lu via un
// `io.Reader`
type ParserInst interface {
	io.Reader

	Init(filter SirenFilter, batch *AdminBatch) error

	// ReadNext extracts tuples from next line.
	// Returns an error only if reading further is not possible: any error
	// definitively interrupts the parsing.
	// Should return io.EOF when there is nothing left to parse.
	ReadNext(*ParsedLineResult) error
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
	Key() string      // entité définie par le tuple: numéro SIRET ou SIREN
	Scope() Scope     // type d'entité: "entreprise" ou "etablissement"
	Type() ParserType // identifiant du parseur qui a extrait ce tuple, ex: "apconso"
}

// ParseFilesFromBatch parse les tuples des fichiers listés dans batch pour le parseur spécifié.
// Retourne les channels des données et des rapports, qui seront populés en
// arrière plan, puis fermés quand toute la donnée aura été traitée.
func ParseFilesFromBatch(
	ctx context.Context,
	batch *AdminBatch,
	parser Parser,
	filter SirenFilter,
) (chan Tuple, chan Report) {

	outputChannel := make(chan Tuple)
	reportChannel := make(chan Report)
	fileType := parser.Type()

	go func() {
		for _, path := range batch.Files[fileType] {
			reportChannel <- parseFileWithReport(ctx, path, parser, batch, outputChannel, filter)
		}
		close(outputChannel)
		close(reportChannel)
	}()
	return outputChannel, reportChannel
}

// parseFileWithReport parses tuples on a single file, and returns a parsing report.
// This function is responsible for logging and error handling, but delegates actual
// parsing down the line.
func parseFileWithReport(
	ctx context.Context,
	batchFile BatchFile,
	parser Parser,
	batch *AdminBatch,
	outputChannel chan Tuple,
	filter SirenFilter,
) Report {
	logger := slog.With("batch", batch.Key, "parser", parser.Type(), "filename", batchFile.Path())
	logger.Debug("parsing file")

	// One tracker per file
	tracker := NewParsingTracker()

	err := runParserOnFile(ctx, batchFile, parser, batch, &tracker,
		outputChannel, filter)

	if err != nil {
		slog.Error("fatal error while parsing file", "parser", parser.Type(), "file", batchFile.Path(), "error", err)
		tracker.AddFatalError(err)
	}

	logger.Debug("end of file parsing")

	return tracker.Report(parser.Type(), batch.Key, batchFile.Path())
}

// runParserOnFile parses tuples on a single file. An error is returned if
// parsing cannot continue.
func runParserOnFile(
	ctx context.Context,
	batchFile BatchFile,
	parser Parser,
	batch *AdminBatch,
	tracker *ParsingTracker,
	outputChannel chan Tuple,
	filter SirenFilter,
) error {

	file, err := batchFile.Open()
	if err != nil {
		return err
	}
	defer file.Close()

	parserInst := parser.New(bufio.NewReader(file))

	if err = parserInst.Init(filter, batch); err != nil {
		return err
	}

	parsedLineChan := make(chan ParsedLineResult)
	errChan := make(chan error)

	go func() {
		errChan <- parseLines(parserInst, parsedLineChan, batchFile.Filename())
	}()

	for lineResult := range parsedLineChan {
		err = processParsedLineResult(ctx, lineResult, filter, tracker, outputChannel)
		if err != nil {
			// Fatal error
			return err
		}

		tracker.Next()
	}

	// Fatal error if any
	return (<-errChan)
}

// parseLines appelle la fonction parseLine() sur chaque ligne du fichier CSV pour transmettre les tuples et/ou erreurs dans parsedLineChan.
//
// "filename" is used for logging purposes.
func parseLines(parserInst ParserInst, parsedLineChan chan ParsedLineResult,
	filename string) error {
	defer close(parsedLineChan)

	var lineNumber = 0 // starting with the header

	stopProgressLogger := LogProgress(&lineNumber, filename)
	defer stopProgressLogger()

	for {
		lineNumber++
		parsedLine := ParsedLineResult{}
		err := parserInst.ReadNext(&parsedLine)

		if err == io.EOF {
			return nil
		} else if err != nil {
			return err
		}

		parsedLineChan <- parsedLine
	}
}

// processParsedLineResult extraie les tuples et/ou erreurs depuis une ligne parsée.
// Return an error only if parsing cannot proceed. Otherwise, track errors
// with the ParsingTracker.
func processParsedLineResult(ctx context.Context, lineResult ParsedLineResult, filter SirenFilter, tracker *ParsingTracker, outputChannel chan Tuple) error {
	// tracking lines filtered by parser even if no data is transmitted
	filterError := lineResult.FilterError
	if filterError != nil {
		tracker.AddFilterError(filterError)
		return nil
	}

	// Report parsing errors
	for _, err := range lineResult.Errors {
		tracker.AddParseError(err)
	}

	for _, tuple := range lineResult.Tuples {
		if _, err := hasValidKey(tuple); err != nil {
			// in addition, invalid siret/siren is reported iff no other parsing
			// error occurred
			if len(lineResult.Errors) == 0 {
				tracker.AddParseError(err)
			}
		} else if filter.ShouldSkip(tuple.Key()) {
			tracker.AddFilterError(errors.New("(filtered)"))

		} else {
			select {
			case <-ctx.Done():
				return fmt.Errorf("Parser interrupted by cancelled context with following cause: %v", context.Cause(ctx))
			case outputChannel <- tuple:
			}
		}
	}
	return nil
}

// LogProgress affiche le numéro de ligne en cours de parsing, toutes les 2s.
func LogProgress(lineNumber *int, filename string) (stop context.CancelFunc) {
	return Cron(time.Minute*1, func() {
		slog.Info(
			"Lis une ligne du fichier csv",
			slog.String("filename", filename),
			slog.Int("line", *lineNumber),
		)
	})
}

// hasValidKey vérifie que la clé (Key) d'un Tuple est valide, selon le type d'entité (Scope) qu'il représente.
func hasValidKey(tuple Tuple) (bool, error) {
	scope := tuple.Scope()
	key := tuple.Key()
	switch scope {
	case ScopeEntreprise:
		if !sfregexp.ValidSiren(key) {
			return false, errors.New("siren invalide : " + key)
		}
		return true, nil

	case ScopeEtablissement:
		if !sfregexp.ValidSiret(key) {
			return false, errors.New("siret invalide : " + key)
		}
		return true, nil
	}

	return false, errors.New("tuple sans scope")
}
