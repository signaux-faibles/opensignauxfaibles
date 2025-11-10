package main

import (
	"errors"
	"log/slog"

	"github.com/cosiner/flag"

	"opensignauxfaibles/lib/db"
	"opensignauxfaibles/lib/engine"
	"opensignauxfaibles/lib/filter"
	prepareimport "opensignauxfaibles/lib/prepare-import"
	"opensignauxfaibles/lib/registry"
	"opensignauxfaibles/lib/sinks"
)

type importBatchHandler struct {
	Enable          bool     // set to true by cosiner/flag if the user is running this command
	Path            string   `names:"--path" env:"APP_DATA" desc:"Directory where raw data can be found. If the batch is not explicitly defined via \"--batch-config\", then it is expected to be in a subfolder named after the batchkey provided with \"--batch\""`
	BatchKey        string   `names:"--batch" arglist:"batch_key" desc:"Batch identifier to import (e.g., 1802 for February 2018)"`
	Parsers         []string `names:"--parsers" desc:"Parsers to use (e.g., apconso, cotisation). Consult documentation with --help for full list of available parsers."` // TODO: tester la population de ce paramètre
	NoFilter        bool     `names:"--no-filter" desc:"Proceed with import without filtering input data, and without updating the filter stored in DB."`
	BatchConfigFile string   `names:"--batch-config" env:"BATCH_CONFIG_FILE" desc:"Path to batch definition file. If not provided, files are inferred from their naming in the data directory (defined by \"APP_DATA\" environment variable or --path option."`
	DryRun          bool     `names:"--dry-run" desc:"Parse files without creating CSV files / database imports. Import report is printed to stdout."`
}

func (params importBatchHandler) Documentation() flag.Flag {
	return flag.Flag{
		Usage: "Import data files",
		Desc: `
    Handles the import and cleaning of input data files for a specific batch.

    DATA LOCATION:
    All imported data is expected to be relative to a base directory defined via
    "--path". If a "--batch-config" option is provided, the relative paths of
    the data files are explicitely given. Otherwise, data files are expected
    to be in a directory with the same name as the "--batch" option.

    BATCH CONFIGURATION:
    The imported files are either explicitly provided via a batch configuration file
    (--batch-config), or inferred from their naming in the data directory.


    FILTERING:
    As the import perimeter is usually a fraction of the raw data, the pipeline uses
    a SIREN filter with the following priority:
      1. Explicitly provided filter (in batch configuration or "filter_..." file)
      2. Filter stored in database (table "stg_filter_import")

    If an "effectif" file is provided, the filter will be updated (or created if none
    exists). If no filter is available and no effectif file is provided, the import
    will fail, unless the "--no-filter" flag is provided.

    OUTPUT:
    The cleaned data is sent to two sinks:
      - CSV files (written to disk)
      - PostgreSQL database (table inserts)

    Import logs can be consulted in the "import_logs" table in PostgreSQL.

    DRY RUN:
    The "--dry-run" flag will discard the data instead of sending it to sinks.
    Import logs are printed to stdout, and no write operations are performed i
    the database (though a filter can still be read from the database if available).

    SELECTIVE PARSING:
    Limit imports to specific parsers using the "--parsers" flag. Note: This only
    limits data import, but does not affect filter updates if an effectif file is
    provided. Use an explicit batch configuration to completely ignore certain files.
	`,
	}
}

func (params importBatchHandler) IsEnabled() bool {
	return params.Enable
}

func (params importBatchHandler) Validate() error {
	if params.BatchKey == "" {
		return errors.New("`batch` parameter is required")
	}
	return nil
}

// Run importBatchHandler processes import requests from the API
// on peut demander l'exécution de tous les parsers sans fournir d'option
// ou demander l'exécution de parsers particuliers en fournissant une liste de leurs codes.
func (params importBatchHandler) Run() error {
	slog.Info("executing import command")

	batchKey, err := engine.NewBatchKey(params.BatchKey)
	if err != nil {
		return err
	}

	// Étape 1
	// On définit d'abord un ensemble de fichiers à importer (batch)
	var batch engine.AdminBatch
	if params.BatchConfigFile != "" {
		// On lit le batch depuis un fichier json
		slog.Info("--batch-config provided, reading batch configuration")
		batch, err = engine.JSONBatchProvider{Path: params.BatchConfigFile}.Get()

	} else {
		// On devine le batch à partir des noms de fichiers
		slog.Info("no --batch-config provided, attempting to infer files to import from filenames")
		batch, err = prepareimport.InferBatchProvider{Path: params.Path, BatchKey: batchKey}.Get()
	}

	if err != nil {
		return err
	}

	// Étape 2
	// On définit les parsers à faire tourner
	if len(params.Parsers) >= 1 {
		slog.Info("import restricted to provided parsers", "parsers", params.Parsers)
	}

	var parserTypes = make([]engine.ParserType, 0, len(params.Parsers))
	for _, p := range params.Parsers {
		parserTypes = append(parserTypes, engine.ParserType(p))
	}

	// Étape 3
	// On définit la destination des données parsées et des rapports de
	// validation
	var dataSinkFactory engine.SinkFactory
	var reportSink engine.ReportSink

	if !params.DryRun {
		slog.Info("data will be written as CSV and to Postgresql tables")
		dataSinkFactory = sinks.Combine(
			sinks.NewCSVSinkFactory(batchKey.String()),
			sinks.NewPostgresSinkFactory(db.DB),
		)

		slog.Info("import logs will be written to Postgresql")
		reportSink = engine.NewPostgresReportSink(db.DB)
	} else {
		slog.Info("dry-run mode: no writing to the database")

		slog.Info("data will be discarded")
		dataSinkFactory = &engine.DiscardSinkFactory{}
		slog.Info("import logs will be written to stdout")
		reportSink = &engine.StdoutReportSink{}
	}

	// Étape 5: Configure filter resolution strategy
	var filterResolver engine.FilterResolver

	if params.NoFilter {
		filterResolver = &filter.NoFilterResolver{}
	} else {
		reader := &filter.StandardReader{Batch: &batch, DB: db.DB}
		var writer filter.Writer
		if !params.DryRun {
			writer = &filter.DBWriter{DB: db.DB}
		}
		filterResolver = &filter.StandardFilterResolver{
			Reader: reader,
			Writer: writer,
		}
	}

	// Étape 6: Execute import
	return executeBatchImport(
		batch,
		parserTypes,
		registry.DefaultParsers,
		filterResolver,
		dataSinkFactory,
		reportSink,
	)
}

// executeBatchImport resolves the filter and imports data
//
// This function is factored out to facilitate testing the filter state
// changes.
func executeBatchImport(
	batch engine.AdminBatch,
	parserTypes []engine.ParserType,
	registry engine.ParserRegistry,
	filterResolver engine.FilterResolver,
	sinkFactory engine.SinkFactory,
	reportSink engine.ReportSink,
) error {
	// Resolve the filter (check, update, and read the filter)
	sirenFilter, err := filterResolver.Resolve(batch.Files)
	if err != nil {
		slog.Error("filter resolution failed", "error", err)
		return err
	}

	if sirenFilter == nil {
		return errors.New(`
      The filter is missing or has not been initialized.
      When the filter is missing, it must be initialized by importing an 'effectif' file,
      or by placing a filter file (prefixed with 'filter_') in the data import directory.
      If you wish to import without a filter, use the "--no-filter" option.
      `)
	}

	// Import with the resolved filter
	err = engine.ImportBatch(
		batch,
		parserTypes,
		registry,
		sirenFilter,
		sinkFactory,
		reportSink,
	)
	return err
}
