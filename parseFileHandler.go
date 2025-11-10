package main

import (
	"errors"
	"opensignauxfaibles/lib/engine"
	"opensignauxfaibles/lib/registry"
	"opensignauxfaibles/lib/sinks"
	"os"

	"github.com/cosiner/flag"
)

type parseFileHandler struct {
	Enable bool   // set to true by cosiner/flag if the user is running this command
	Parser string `names:"--parser" desc:"Parser to use (ex: cotisation)"`
	File   string `names:"--file"   desc:"Path to the file to parse. Unlike import, the file path must be complete and does not take into account the APP_DATA environment variable.`
}

func (params parseFileHandler) Documentation() flag.Flag {
	return flag.Flag{
		Usage: "Parses a file to stdout",
	}
}

func (params parseFileHandler) IsEnabled() bool {
	return params.Enable
}

func (params parseFileHandler) Validate() error {
	if params.Parser == "" {
		return errors.New("required `parser` parameter")
	}
	if params.File == "" {
		return errors.New("required `file` parameter")
	}
	if _, err := os.Stat(params.File); err != nil {
		return errors.New("Can't find " + params.File + ": " + err.Error())
	}
	return nil
}

func (params parseFileHandler) Run() error {
	parserType := engine.ParserType(params.Parser)

	file := engine.NewBatchFile(params.File)
	batch := engine.AdminBatch{
		Key:   "parseFile", // dummy batch key
		Files: engine.BatchFiles{parserType: []engine.BatchFile{file}},
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
