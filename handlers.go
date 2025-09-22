package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"os"
	"sync"

	"github.com/cosiner/flag"
	"github.com/spf13/viper"

	"opensignauxfaibles/lib/base"
	"opensignauxfaibles/lib/engine"
	"opensignauxfaibles/lib/marshal"
	"opensignauxfaibles/lib/parsing"
	prepareimport "opensignauxfaibles/lib/prepare-import"
)

type importBatchHandler struct {
	Enable      bool     // set to true by cosiner/flag if the user is running this command
	BatchKey    string   `names:"--batch" arglist:"batch_key" desc:"Identifiant du batch à importer (ex: 1802, pour Février 2018)"`
	Parsers     []string `names:"--parsers" desc:"Parseurs à employer (ex: apconso, cotisation)"` // TODO: tester la population de ce paramètre
	NoFilter    bool     `names:"--no-filter" desc:"Pour procéder à l'importation même si aucun filtre n'est fourni"`
	BatchConfig string   `names:"--batch-config" desc:"Chemin de définition de l'ensemble des fichiers à importer (batch). À défaut, ces fichiers sont devinés par rapport à leur nommage, dans le répertoire de la variable d'environnement APP_DATA."`
	DryRun      bool     `names:"--dry-run" desc:"Pour parser les fichiers sans créer de fichiers CSV / imports en base"`
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
	batch := base.AdminBatch{}

	if params.BatchConfig != "" {
		// On lit le batch depuis un fichier json
		slog.Info("Batch fourni en paramètre, lecture de la configuration du batch")
		err := engine.Load(&batch, params.BatchConfig)
		if err != nil {
			return fmt.Errorf("impossible de charger la configuration du batch : %w", err)
		}

	} else {
		// On devine le batch à partir des noms de fichiers
		slog.Info("Batch non fourni en paramètre, tentative de déterminer les fichiers à importer")

		batchKey, err := base.NewBatchKey(params.BatchKey)
		if err != nil {
			return err
		}

		batch, err = prepareimport.PrepareImport(viper.GetString("APP_DATA"), batchKey)
		if _, ok := err.(prepareimport.UnsupportedFilesError); ok {
			slog.Warn(fmt.Sprintf("Des fichiers non-identifiés sont présents : %v", err))
		} else if err != nil {
			return fmt.Errorf("une erreur est survenue en préparant l'import : %w", err)
		}

		slog.Info("Batch deviné avec succès")

		batchJSON, _ := json.MarshalIndent(batch, "", "  ")
		if batchJSON != nil {
			slog.Info(string(batchJSON))
		}
	}

	var parserTypes = make([]base.ParserType, 0, len(params.Parsers))
	for _, p := range params.Parsers {
		parserTypes = append(parserTypes, base.ParserType(p))
	}

	var dataSink engine.SinkFactory
	var reportSink engine.ReportSink

	if !params.DryRun {
		dataSink = engine.NewCompositeSinkFactory(
			engine.NewCSVSinkFactory(batch.Key.String()),
			engine.NewPostgresSinkFactory(engine.Db.PostgresDB),
		)
		reportSink = engine.NewPostgresReportSink(engine.Db.PostgresDB)
	} else {
		dataSink = &engine.DiscardSinkFactory{}
		reportSink = &engine.DiscardReportSink{}
	}

	err := engine.ImportBatch(batch, parserTypes, params.NoFilter, dataSink, reportSink)

	if err != nil {
		return err
	}

	printJSON("ok")
	return nil
}

func printJSON(object any) {
	res, _ := json.Marshal(object)
	fmt.Println(string(res))
}

type parseFileHandler struct {
	Enable bool   // set to true by cosiner/flag if the user is running this command
	Parser string `names:"--parser" desc:"Parseur à employer (ex: cotisation)"`
	File   string `names:"--file"   desc:"Nom du fichier à parser. Contrairement à l'import, le chemin du fichier doit être complet et ne tient pas compte de la variable d'environnement APP_DATA"`
}

func (params parseFileHandler) Documentation() flag.Flag {
	return flag.Flag{
		Usage: "Parse un fichier vers la sortie standard",
	}
}

func (params parseFileHandler) IsEnabled() bool {
	return params.Enable
}

func (params parseFileHandler) Validate() error {
	if params.Parser == "" {
		return errors.New("paramètre `parser` obligatoire")
	}
	if params.File == "" {
		return errors.New("paramètre `file` obligatoire")
	}
	if _, err := os.Stat(params.File); err != nil {
		return errors.New("Can't find " + params.File + ": " + err.Error())
	}
	return nil
}

func (params parseFileHandler) Run() error {
	parserType := base.ParserType(params.Parser)
	parsers, err := parsing.ResolveParsers([]base.ParserType{parserType})
	if err != nil {
		return err
	}

	file := base.NewBatchFile(params.File)
	batch := base.AdminBatch{Files: base.BatchFiles{parserType: []base.BatchFile{file}}}
	cache := marshal.NewCache()
	parser := parsers[0]

	// the following code is inspired from marshal.ParseFilesFromBatch()
	outputChannel := make(chan marshal.Tuple)
	reportChannel := make(chan marshal.Report)
	ctx := context.Background()
	go func() {
		reportChannel <- marshal.ParseFile(ctx, file, parser, &batch, cache, outputChannel)
		close(outputChannel)
		close(reportChannel)
	}()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		for tuple := range outputChannel {
			printJSON(tuple) // écriture du tuple dans la sortie standard
		}
	}()

	for e := range reportChannel {
		res, _ := json.MarshalIndent(e, "", "  ")
		log.Println(string(res)) // écriture de l'événement dans stderr
	}

	// Only return once all channels are closed
	wg.Wait()

	return nil
}
