package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"opensignauxfaibles/lib/base"
	"opensignauxfaibles/lib/engine"
	"opensignauxfaibles/lib/parsing"
	"os"
	"sync"

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
	parsers, err := engine.ResolveParsers(parsing.DefaultParsers, []base.ParserType{parserType})
	if err != nil {
		return err
	}

	file := base.NewBatchFile(params.File)
	batch := base.AdminBatch{Files: base.BatchFiles{parserType: []base.BatchFile{file}}}
	cache := engine.NewEmptyCache()
	parser := parsers[0]

	// the following code is inspired from engine.ParseFilesFromBatch()
	outputChannel := make(chan engine.Tuple)
	reportChannel := make(chan engine.Report)
	ctx := context.Background()
	go func() {
		reportChannel <- engine.ParseFile(ctx, file, parser, &batch, cache, outputChannel, nil)
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
