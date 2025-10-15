// Package engine est la colonne vertebrale de l'import. Elle définit les
// interfaces pour définer les parsers, le filtre, les sinks (vers lesquels
// la donnée sera envoyée).
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
)

// Load charge les données d'un batch depuis du JSON
func Load(batch *AdminBatch, reader io.Reader) error {
	batchFileContent, err := io.ReadAll(reader)
	if err != nil {
		return err
	}

	err = json.Unmarshal(batchFileContent, batch)
	return err
}

// ImportBatch lance tous les parsers sur le batch fourni
func ImportBatch(
	batchConfig AdminBatch,
	parserTypes []ParserType,
	registry ParserRegistry,
	filter SirenFilter,
	sinkFactory SinkFactory,
	eventSink ReportSink,
) error {

	parsers, err := ResolveParsers(registry, parserTypes)
	if err != nil {
		return err
	}

	logger := slog.With("batch", batchConfig.Key)
	logger.Info("starting raw data import")

	unsupported := checkUnsupportedFiletypes(registry, batchConfig)
	for _, file := range unsupported {
		logger.Warn("Type de fichier non reconnu", "file", file)
	}

	var g errgroup.Group

	var cache = NewEmptyCache()
	for _, parser := range parsers {
		// We create a parser-specific context. Any error will cancelParserProcess all
		// parser-related operations
		ctx, cancelParserProcess := context.WithCancelCause(context.Background())
		defer cancelParserProcess(nil)

		if len(batchConfig.Files[parser.Type()]) > 0 {
			logger.Info("parse raw data", "parser", parser.Type())

			outputChannel, eventChannel := ParseFilesFromBatch(ctx, cache, &batchConfig, parser, filter)

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

	if err != nil {
		return err
	}
	logger.Info("raw data parsing ended")
	logger.Info("inspect the \"Journal\" database to consult parsing errors.")

	return nil
}

// CheckBatchPaths checks if the filepaths of batch.Files exist
func CheckBatchPaths(batch *AdminBatch) error {
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

func checkUnsupportedFiletypes(registry ParserRegistry, batch AdminBatch) []ParserType {
	var errFileTypes []ParserType
	for parserType := range batch.Files {
		if parserType != Filter && registry.Resolve(parserType) == nil {
			errFileTypes = append(errFileTypes, parserType)
		}
	}
	return errFileTypes
}

// JSONBatchProvider reads an admin batch from a JSON file.
// Implements BatchProvider interface.
type JSONBatchProvider struct {
	Path string
}

func (p JSONBatchProvider) Get() (AdminBatch, error) {
	fileReader, err := os.Open(p.Path)

	var batch = AdminBatch{}

	if err == nil {
		err = Load(&batch, fileReader)
	}

	if err != nil {
		return batch, fmt.Errorf("impossible de charger la configuration du batch : %w", err)
	}

	return batch, nil
}
