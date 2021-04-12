package main

import (
	"encoding/json"
	"errors"
	"log"
	"os"

	flag "github.com/cosiner/flag"
	"github.com/signaux-faibles/opensignauxfaibles/lib/base"
	"github.com/signaux-faibles/opensignauxfaibles/lib/marshal"
	"github.com/signaux-faibles/opensignauxfaibles/lib/parsing"
)

type parseFileHandler struct {
	Enable bool   // set to true by cosiner/flag if the user is running this command
	Parser string `names:"--parser" desc:"Parseur à employer (ex: cotisation)"`
	File   string `names:"--file"   desc:"Nom du fichier à parser"`
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

	parsers, err := parsing.ResolveParsers([]string{params.Parser})
	if err != nil {
		return err
	}

	file := base.BatchFile(params.File)
	batch := base.AdminBatch{Files: base.BatchFiles{params.Parser: []base.BatchFile{file}}}
	cache := marshal.NewCache()
	parser := parsers[0]

	// the following code is inspired from marshal.ParseFilesFromBatch()
	outputChannel := make(chan marshal.Tuple)
	eventChannel := make(chan marshal.Event)
	go func() {
		eventChannel <- marshal.ParseFile(file, parser, &batch, cache, outputChannel)
		close(outputChannel)
		close(eventChannel)
	}()

	go func() {
		for tuple := range outputChannel {
			printJSON(tuple)
		}
	}()

	for e := range eventChannel {
		res, _ := json.MarshalIndent(e, "", "  ")
		log.Println(string(res))
	}
	return nil
}
