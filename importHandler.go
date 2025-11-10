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
		Handles the import and cleaning of input files.

    All imported data is expected to be found in a single directory, that can be defined via "--path" and "--batch" options.

    The imported files are either explicitely provided via a batch configuration file, or by default, inferred from their naming.

    As the import perimeter is usually a fraction of the raw data, the
    pipeline will use either an explicitely provided filter (as part of the
    batch coniguration file), or a filter stored in database (table
    "stg_filter_import"). If an "effectif" file is provided, the filter will
    be updated (or created if none exists).


    The cleaned data is then send to two sinks : one writing CSV files, the other one storing data inside a Postgresql database.

		It is possible to limit execution to certain parsers by specifying the list with the "--parsers" flag.
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
	batchKey, err := engine.NewBatchKey(params.BatchKey)
	if err != nil {
		return err
	}

	// Étape 1
	// On définit d'abord un ensemble de fichiers à importer (batch)
	var batch engine.AdminBatch
	if params.BatchConfigFile != "" {
		// On lit le batch depuis un fichier json
		slog.Info("batch parameter provided, reading batch configuration")
		batch, err = engine.JSONBatchProvider{Path: params.BatchConfigFile}.Get()

	} else {
		// On devine le batch à partir des noms de fichiers
		slog.Info("batch parameter not provided, attempting to determine files to import")
		batch, err = prepareimport.InferBatchProvider{Path: params.Path, BatchKey: batchKey}.Get()
	}

	if err != nil {
		return err
	}

	// Étape 2
	// On définit les parsers à faire tourner
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
		dataSinkFactory = sinks.Combine(
			sinks.NewCSVSinkFactory(batchKey.String()),
			sinks.NewPostgresSinkFactory(db.DB),
		)
		reportSink = engine.NewPostgresReportSink(db.DB)
	} else {
		dataSinkFactory = &engine.DiscardSinkFactory{}
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
	slog.Info("import completed")
	return err
}
