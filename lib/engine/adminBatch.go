// Package engine is the backbone of the import. It defines the
// interfaces for defining parsers, the filter, and sinks (to which
// data will be sent).
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

// Load loads batch data from JSON
func Load(batch *AdminBatch, reader io.Reader) error {
	batchFileContent, err := io.ReadAll(reader)
	if err != nil {
		return err
	}

	err = json.Unmarshal(batchFileContent, batch)
	return err
}

// ImportBatch runs all parsers on the provided batch
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

	logger.Info("importing raw data...")

	unsupported := checkUnsupportedFiletypes(registry, batchConfig)
	for _, file := range unsupported {
		logger.Warn("unrecognized filetype, skipped", "file", file)
	}

	var g errgroup.Group

	var cache = NewEmptyCache()
	for _, parser := range parsers {
		// We create a parser-specific context. Any error will cancel all
		// parser-related operations
		ctx, cancelParserProcess := context.WithCancelCause(context.Background())
		defer cancelParserProcess(nil)

		if len(batchConfig.Files[parser.Type()]) > 0 {
			logger.Info("parse raw data...", "parser", parser.Type())

			// Parsing files for given parser type
			// outputChannel and eventChannel are populated in background thread
			outputChannel, eventChannel := ParseFilesFromBatch(ctx, cache, &batchConfig, parser, filter)

			// Insert events (parsing logs) into the "Journal" collection
			g.Go(
				func() error {
					eventErr := eventSink.Process(eventChannel)
					return eventErr
				},
			)

			// Stream data to the output sink(s)
			g.Go(
				func() error {
					dataSink, sinkErr := sinkFactory.CreateSink(parser.Type())
					if sinkErr != nil {
						return sinkErr
					}

					sinkErr = dataSink.ProcessOutput(ctx, outputChannel)
					if sinkErr != nil {
						cancelParserProcess(sinkErr)
					}

					return sinkErr
				},
			)
		}
	}
	err = g.Wait() // wait for all events and tuples to be inserted, get the error if any

	if err != nil {
		return err
	}

	logger.Info("raw data import ended")
	logger.Info("inspect the import logs to consult any parsing errors")

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
		return batch, fmt.Errorf("unable to load batch configuration : %w", err)
	}

	return batch, nil
}
