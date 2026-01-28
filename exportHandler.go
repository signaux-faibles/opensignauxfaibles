package main

import (
	"fmt"
	"log/slog"
	"opensignauxfaibles/lib/db"
	"opensignauxfaibles/lib/export"

	"github.com/cosiner/flag"
)

var _ commandHandler = exportHandler{}

type exportHandler struct {
	Enable bool   // set to true by cosiner/flag if the user is running this command
	Path   string `names:"--path" env:"EXPORT_DIR" desc:"Directory to export to"`
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

	shouldMigrate := false
	err := db.Init(shouldMigrate)
	if err != nil {
		return fmt.Errorf("error while connecting to db: %w", err)
	}

	return export.CleanViews(params.Path)
}

func (params exportHandler) Validate() error {
	return nil
}
