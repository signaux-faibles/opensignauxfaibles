package main

import (
	"errors"
	"fmt"
	"log/slog"
	"opensignauxfaibles/lib/db"
	"opensignauxfaibles/lib/export"

	"github.com/cosiner/flag"
)

var _ commandHandler = exportHandler{}

type exportHandler struct {
	Enable         bool   // set to true by cosiner/flag if the user is running this command
	Path           string `names:"--path" env:"EXPORT_DIR" desc:"Directory to export to"`
	Schema         string `names:"--schema" desc:"PostgreSQL schema to use (allows running multiple pipelines in parallel on different schemas)"`
	PerimeterOnly  bool   `names:"--perimeter-only" desc:"Export only clean_filter (the perimeter). The detail-filtres.md report written by computePerimeter lives in the same directory."`
}

func (params exportHandler) Documentation() flag.Flag {
	return flag.Flag{
		Usage: "Export DB data for data science",
		Desc: `
    Exports all cleaned views to files.
	`,
	}
}

func (params exportHandler) IsEnabled() bool {
	return params.Enable
}

func (params exportHandler) Run() error {
	slog.Info("executing export command")

	shouldMigrate := true
	err := db.Init(params.Schema, shouldMigrate)
	if err != nil {
		return fmt.Errorf("error while connecting to db: %w", err)
	}

	exporter := export.NewExporter(params.Path, db.DB)
	if params.PerimeterOnly {
		exporter = exporter.WithViews([]string{db.ViewFilter})
	}
	return exporter.CleanViews()
}

func (params exportHandler) Validate() error {
	if params.Schema == "" {
		return errors.New("`schema` parameter is required (use --schema flag)")
	}
	return nil
}
