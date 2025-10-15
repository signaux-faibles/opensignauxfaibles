package main

import (
	"errors"
	"opensignauxfaibles/lib/base"
	"opensignauxfaibles/lib/engine"
	"opensignauxfaibles/lib/registry"
	"opensignauxfaibles/lib/sinks"
	"os"

	"github.com/cosiner/flag"
)

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

	file := base.NewBatchFile(params.File)
	batch := base.AdminBatch{
		Key:   "parseFile", // dummy batch key
		Files: base.BatchFiles{parserType: []base.BatchFile{file}},
	}

	// stdout csv data output
	dataSinkFactory := sinks.NewStdoutSinkFactory()
	// stdout json report output
	reportSink := &engine.StdoutReportSink{}

	return engine.ImportBatch(
		batch,
		nil, // do not filter any parser
		registry.DefaultParsers,
		engine.NoFilter,
		dataSinkFactory,
		reportSink,
	)
}
