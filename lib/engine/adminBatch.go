package engine

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"

	"golang.org/x/sync/errgroup"

	"opensignauxfaibles/lib/base"
	"opensignauxfaibles/lib/marshal"
	"opensignauxfaibles/lib/parsing"
)

// Load charge les données d'un batch depuis du JSON
func Load(batch *base.AdminBatch, reader io.Reader) error {
	batchFileContent, err := io.ReadAll(reader)
	if err != nil {
		return err
	}

	err = json.Unmarshal(batchFileContent, batch)
	return err
}

// ImportBatch lance tous les parsers sur le batch fourni
// Si "sinkFactory" implémente également l'interface `Finalizer`, alors la
// fonction `Finalize` est déclenchée à la fin de tous les imports.
func ImportBatch(
	batchProvider base.BatchProvider,
	parserTypes []base.ParserType,
	skipFilter bool,
	sinkFactory SinkFactory,
	eventSink ReportSink,
) error {

	batch, err := batchProvider.Get()
	if err != nil {
		return err
	}

	parsers, err := parsing.ResolveParsers(parserTypes)
	if err != nil {
		return err
	}

	logger := slog.With("batch", batch.Key)
	logger.Info("starting raw data import")

	var cache = marshal.NewEmptyCache()

	filter, err := marshal.GetSirenFilter(cache, &batch)
	if err != nil {
		return err
	}
	if !skipFilter && filter == nil {
		return errors.New("veuillez inclure un filtre")
	}

	unsupported := checkUnsupportedFiletypes(batch)
	for _, file := range unsupported {
		logger.Warn("Type de fichier non reconnu", "file", file)
	}

	var g errgroup.Group

	for _, parser := range parsers {
		// We create a parser-specific context. Any error will cancelParserProcess all
		// parser-related operations
		ctx, cancelParserProcess := context.WithCancelCause(context.Background())
		defer cancelParserProcess(nil)

		if len(batch.Files[parser.Type()]) > 0 {
			logger.Info("parse raw data", "parser", parser.Type())

			outputChannel, eventChannel := marshal.ParseFilesFromBatch(ctx, cache, &batch, parser) // appelle la fonction ParseFile() pour chaque type de fichier

			// Insert events (parsing logs) into the "Journal" collection
			g.Go(
				func() error {
					err := eventSink.Process(eventChannel)
					return err
				},
			)

			// Stream data to the output sink(s)
			g.Go(
				func() error {
					dataSink, err := sinkFactory.CreateSink(parser.Type())
					if err != nil {
						return err
					}

					err = dataSink.ProcessOutput(ctx, outputChannel)
					if err != nil {
						cancelParserProcess(err)
					}

					return err
				},
			)
		}
	}
	err = g.Wait() // wait for all events and tuples to be inserted, get the error if any

	if finalizer, ok := sinkFactory.(Finalizer); ok {
		// If any finalizer function is detected (e.g. update materialized views
		// for postgresql sink factory), run it now
		logger.Debug("Fonction de finalisation de l'import détectée.")
		err = finalizer.Finalize()
		if err != nil {
			return fmt.Errorf("une erreur est survenue lors de l'exécution de la fonction de finalisation : %w", err)
		}
		logger.Info("Finalisation de l'import effectuée avec succès")
	}

	if err != nil {
		return err
	}
	logger.Info("raw data parsing ended")
	logger.Info("inspect the \"Journal\" database to consult parsing errors.")

	return nil
}

// CheckBatchPaths checks if the filepaths of batch.Files exist
func CheckBatchPaths(batch *base.AdminBatch) error {
	var ErrorString string
	for _, filepaths := range batch.Files {
		for _, batchFile := range filepaths {
			if _, err := os.Stat(batchFile.Path()); err != nil {
				ErrorString += batchFile.Path() + " is missing (" + err.Error() + ").\n"
			}
		}
	}
	if ErrorString != "" {
		return errors.New(ErrorString)
	}
	return nil

}

func checkUnsupportedFiletypes(batch base.AdminBatch) []base.ParserType {
	var errFileTypes []base.ParserType
	for fileType := range batch.Files {
		if !parsing.IsSupportedParser(fileType) {
			errFileTypes = append(errFileTypes, fileType)
		}
	}
	return errFileTypes
}

// JSONBatchProvider reads an admin batch from a JSON file.
// Implements base.BatchProvider interface.
type JSONBatchProvider struct {
	Path string
}

func (p JSONBatchProvider) Get() (base.AdminBatch, error) {
	fileReader, err := os.Open(p.Path)

	var batch = base.AdminBatch{}

	if err == nil {
		err = Load(&batch, fileReader)
	}

	if err != nil {
		return batch, fmt.Errorf("impossible de charger la configuration du batch : %w", err)
	}

	return batch, nil
}
