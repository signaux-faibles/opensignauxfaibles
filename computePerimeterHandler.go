package main

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/cosiner/flag"
	"github.com/spf13/viper"

	"opensignauxfaibles/lib/db"
	"opensignauxfaibles/lib/engine"
	"opensignauxfaibles/lib/filter"
	prepareimport "opensignauxfaibles/lib/prepare-import"
)

type computePerimeterHandler struct {
	Enable          bool   // set to true by cosiner/flag if the user is running this command
	Path            string `names:"--path" env:"APP_DATA" desc:"Directory where raw data can be found. If the batch is not explicitly defined via \"--batch-config\", then it is expected to be in a subfolder named after the batchkey provided with \"--batch\""`
	BatchKey        string `names:"--batch" arglist:"batch_key" desc:"Batch identifier (e.g., 1802 for February 2018)"`
	BatchConfigFile string `names:"--batch-config" env:"BATCH_CONFIG_FILE" desc:"Path to batch definition file. If not provided, files are inferred from their naming in the data directory."`
}

func (params computePerimeterHandler) Documentation() flag.Flag {
	return flag.Flag{
		Usage: "Compute import perimeter from effectif_ent file",
		Desc: `
    Computes the SIREN import perimeter from an effectif_ent file and stores it
    in the database (table "stg_filter_import").

    This command should be run before importing data, so that the import can
    filter input data based on the computed perimeter.

    The effectif_ent file is resolved from the batch configuration (explicit or
    inferred from filenames in the data directory).
		`,
	}
}

func (params computePerimeterHandler) IsEnabled() bool {
	return params.Enable
}

func (params computePerimeterHandler) Validate() error {
	if params.BatchKey == "" {
		return errors.New("`batch` parameter is required")
	}
	return nil
}

func (params computePerimeterHandler) Run() error {
	slog.Info("executing computePerimeter command")

	shouldMigrate := true
	err := db.Init(shouldMigrate)
	if err != nil {
		return fmt.Errorf("error while connecting to db: %w", err)
	}
	defer db.DB.Close()

	batchKey, err := engine.NewBatchKey(params.BatchKey)
	if err != nil {
		return err
	}

	viper.Set("batch", params.BatchKey)

	var batch engine.AdminBatch
	if params.BatchConfigFile != "" {
		slog.Info("--batch-config provided, reading batch configuration")
		batch, err = engine.JSONBatchProvider{Path: params.BatchConfigFile}.Get()
	} else {
		slog.Info("no --batch-config provided, attempting to infer files from filenames")
		batch, err = prepareimport.InferBatchProvider{Path: params.Path, BatchKey: batchKey}.Get()
	}
	if err != nil {
		return err
	}

	effectifEntFile := batch.Files.GetEffectifEntFile()
	if effectifEntFile == nil {
		return errors.New("no effectif_ent file found in batch: an effectif_ent file is required to compute the perimeter")
	}

	slog.Info("computing perimeter from effectif_ent file...")
	sirenFilter, err := filter.CreateFilter(effectifEntFile, filter.DefaultNbMois, filter.DefaultMinEffectif)
	if err != nil {
		return fmt.Errorf("failed to compute perimeter: %w", err)
	}

	writer := &filter.DBWriter{DB: db.DB}
	if err := writer.Write(sirenFilter); err != nil {
		return fmt.Errorf("failed to write perimeter to database: %w", err)
	}

	slog.Info("perimeter computed and saved successfully", "n_sirens", len(sirenFilter.All()))
	return nil
}
