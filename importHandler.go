package main

import (
	"errors"
	"fmt"
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
	BatchKey        string   `names:"--batch" arglist:"batch_key" desc:"Identifiant du batch à importer (ex: 1802, pour Février 2018)"`
	Parsers         []string `names:"--parsers" desc:"Parseurs à employer (ex: apconso, cotisation)"` // TODO: tester la population de ce paramètre
	NoFilter        bool     `names:"--no-filter" desc:"Pour procéder à l'import même si aucun filtre n'est fourni"`
	BatchConfigFile string   `names:"--batch-config" env:"BATCH_CONFIG_FILE" desc:"Chemin de définition de l'ensemble des fichiers à importer (batch). À défaut, ces fichiers sont devinés par rapport à leur nommage, dans le répertoire de la variable d'environnement APP_DATA."`
	DryRun          bool     `names:"--dry-run" desc:"Pour parser les fichiers sans créer de fichiers CSV / imports en base"`
}

func (params importBatchHandler) Documentation() flag.Flag {
	return flag.Flag{
		Usage: "Importe des fichiers de données",
		Desc: `
		Effectue l'import de tous les fichiers du batch donné en paramètre.
		Il est possible de limiter l'exécution à certains parsers en spécifiant la liste dans le flag "--parsers".
		Répond "ok" dans la sortie standard, si le traitement s'est bien déroulé.
	`,
	}
}

func (params importBatchHandler) IsEnabled() bool {
	return params.Enable
}

func (params importBatchHandler) Validate() error {
	if params.BatchKey == "" {
		return errors.New("paramètre `batch` obligatoire")
	}
	return nil
}

// Run importBatchHandler traite les demandes d'import par l'API
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
		slog.Info("Batch fourni en paramètre, lecture de la configuration du batch")
		batch, err = engine.JSONBatchProvider{Path: params.BatchConfigFile}.Get()

	} else {
		// On devine le batch à partir des noms de fichiers
		slog.Info("Batch non fourni en paramètre, tentative de déterminer les fichiers à importer")
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

	// Étape 5 on récupère le périmètre d'import

	// Par défaut, ne filtre pas (retourne engine.NoFilter)
	var filterReader engine.FilterReader = &filter.NoReader{} // No filtering
	var filterWriter engine.FilterWriter                      // No filter update

	if !params.NoFilter {
		// Create filter provider with database dependency
		filterReader = &filter.Reader{Batch: &batch, DB: db.DB}
		if !params.DryRun {
			filterWriter = &filter.DBWriter{DB: db.DB}
		}
	}

	// Étape 5
	// On réalise l'import
	return executeBatchImport(
		batch,
		parserTypes,
		registry.DefaultParsers,
		filterReader,
		filterWriter,
		dataSinkFactory,
		reportSink,
	)
}

// executeBatchImport resolves and updates the filter, and imports data
//
// This function is factored out to facilitate testing the filter state
// changes.
func executeBatchImport(
	batch engine.AdminBatch,
	parserTypes []engine.ParserType,
	registry engine.ParserRegistry,
	filterReader engine.FilterReader,
	filterWriter engine.FilterWriter,
	sinkFactory engine.SinkFactory,
	reportSink engine.ReportSink,
) error {
	// Check if filtering conditions are met
	if err := filter.Check(filterReader, batch.Files); err != nil {
		return fmt.Errorf("filter check failed: %w", err)
	}

	// Update the filter state if needed
	if err := filter.UpdateState(filterWriter, batch.Files); err != nil {
		return fmt.Errorf("filter update failed: %w", err)
	}

	sirenFilter, err := filterReader.Read()
	if err != nil {
		return fmt.Errorf("unable to get filter: %w", err)
	}

	if sirenFilter == nil {
		return errors.New(`
      Le filtre est manquant ou n'a pas été initialisé.
      Lorsque le filtre est manquant, il est nécessaire de l'initialiser via
      l'import d'un fichier 'effectif', ou de placer le fichier filtre à
      importer, préfixé par 'filter_' dans le dossier des données à importer.
      Si vous souhaitez importer sans filtre, utilisez l'option "--no-filter".
      `)
	}
	err = engine.ImportBatch(
		batch,
		parserTypes,
		registry,
		sirenFilter,
		sinkFactory,
		reportSink,
	)
	slog.Info("Import terminé")
	return err
}
