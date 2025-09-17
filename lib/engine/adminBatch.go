package engine

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"os"

	"github.com/spf13/viper"
	"golang.org/x/sync/errgroup"

	"opensignauxfaibles/lib/base"
	"opensignauxfaibles/lib/marshal"
	"opensignauxfaibles/lib/parsing"
)

// Load charge les données d'un batch depuis le fichier de configuration
func Load(batch *base.AdminBatch, batchKey string) error {
	batchFileContent, err := os.ReadFile(viper.GetString("BATCH_CONFIG_FILE"))
	if err != nil {
		return err
	}

	err = json.Unmarshal(batchFileContent, batch)
	return err
}

// ImportBatch lance tous les parsers sur le batch fourni
func ImportBatch(
	batch base.AdminBatch,
	parsers []marshal.Parser,
	skipFilter bool,
	sinkFactory SinkFactory,
	eventSink ReportSink,
) error {

	logger := slog.With("batch", batch.ID.Key)
	logger.Info("starting raw data import")

	var cache = marshal.NewCache()

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
			if _, err := os.Stat(batchFile.FilePath()); err != nil {
				ErrorString += batchFile.FilePath() + " is missing (" + err.Error() + ").\n"
			}
		}
	}
	if ErrorString != "" {
		return errors.New(ErrorString)
	}
	return nil

}

// CheckBatch checks batch, discard all data but logs events
func CheckBatch(
	batch base.AdminBatch,
	parsers []marshal.Parser,
	eventSink ReportSink,
) error {
	ctx := context.Background()
	if err := CheckBatchPaths(&batch); err != nil {
		return err
	}
	var cache = marshal.NewCache()
	for _, parser := range parsers {
		logger := slog.With("batch", batch.ID.Key, "parser", parser.Type())
		outputChannel, eventChannel := marshal.ParseFilesFromBatch(ctx, cache, &batch, parser)

		DiscardTuple(outputChannel)
		for report := range eventChannel {
			if report.LinesRejected > 0 {
				logger.Error(report.Summary)
			} else {
				logger.Info(report.Summary)
			}
		}
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
