package engine

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/globalsign/mgo/bson"
	"github.com/spf13/viper"
	"golang.org/x/sync/errgroup"

	"opensignauxfaibles/lib/base"
	"opensignauxfaibles/lib/marshal"
	"opensignauxfaibles/lib/parsing"
)

// Load charge les données d'un batch depuis la base de données
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

	startDate := time.Now()
	reportUnsupportedFiletypes(batch)

	var g errgroup.Group

	// Limit maximum concurrent go routines
	g.SetLimit(6)

	for _, parser := range parsers {
		if len(batch.Files[parser.Type()]) > 0 {
			logger.Info("parse raw data", "parser", parser.Type())

			outputChannel, eventChannel := marshal.ParseFilesFromBatch(cache, &batch, parser) // appelle la fonction ParseFile() pour chaque type de fichier

			// Insert events (parsing logs) into the "Journal" collection
			g.Go(
				func() error {
					RelayEvents(eventChannel, "ImportBatch", startDate)
					return nil
				},
			)

			// Stream data to the output sink(s)
			g.Go(
				func() error {
					dataSink, err := sinkFactory.CreateSink(parser.Type())
					if err != nil {
						return err
					}

					return dataSink.ProcessOutput(outputChannel)
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

// CheckBatch checks batch
func CheckBatch(batch base.AdminBatch, parsers []marshal.Parser) (reports []string, err error) {
	if err := CheckBatchPaths(&batch); err != nil {
		return nil, err
	}
	var cache = marshal.NewCache()
	startDate := time.Now()
	for _, parser := range parsers {
		outputChannel, eventChannel := marshal.ParseFilesFromBatch(cache, &batch, parser) // appelle la fonction ParseFile() pour chaque type de fichier
		DiscardTuple(outputChannel)
		parserReports := RelayEvents(eventChannel, "CheckBatch", startDate)
		reports = append(reports, parserReports...)
	}

	return reports, nil
}

func reportUnsupportedFiletypes(batch base.AdminBatch) {
	for fileType := range batch.Files {
		if !parsing.IsSupportedParser(fileType) {
			msg := fmt.Sprintf("Type de fichier non reconnu: %v", fileType)
			log.Println(msg) // notification dans la sortie d'erreurs
			event := marshal.CreateReportEvent(fileType, bson.M{
				"batchKey": batch.ID.Key,
				"summary":  msg,
			})
			event.ReportType = "ImportBatch_error"
			mainMessageChannel <- event
		}
	}
}
