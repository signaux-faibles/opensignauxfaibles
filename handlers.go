package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/cosiner/flag"

	"opensignauxfaibles/lib/base"
	"opensignauxfaibles/lib/engine"
	"opensignauxfaibles/lib/marshal"
	"opensignauxfaibles/lib/parsing"
)

type importBatchHandler struct {
	Enable   bool     // set to true by cosiner/flag if the user is running this command
	BatchKey string   `names:"--batch" arglist:"batch_key" desc:"Identifiant du batch à importer (ex: 1802, pour Février 2018)"`
	Parsers  []string `names:"--parsers" desc:"Parseurs à employer (ex: altares,cotisation)"` // TODO: tester la population de ce paramètre
	NoFilter bool     `names:"--no-filter" desc:"Pour procéder à l'importation même si aucun filtre n'est fourni"`
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

	err := engine.Load(&batch, params.BatchKey)
	if err != nil {
		return errors.New("Impossible de charger la configuration du batch: " + err.Error())
	}

	parsers, err := parsing.ResolveParsers(params.Parsers)
	if err != nil {
		return err
	}

	dataSink := engine.NewCompositeSinkFactory(
		engine.NewCSVSinkFactory(batch.ID.Key),
		engine.NewPostgresSinkFactory(engine.Db.PostgresDB),
	)
	eventSink := engine.NewPostgresEventSink(engine.Db.PostgresDB)

	err = engine.ImportBatch(batch, parsers, params.NoFilter, dataSink, eventSink)

	if err != nil {
		return err
	}

	printJSON("ok")
	return nil
}

type checkBatchHandler struct {
	Enable   bool     // set to true by cosiner/flag if the user is running this command
	BatchKey string   `names:"--batch" arglist:"batch_key" desc:"Identifiant du batch à vérifier (ex: 1802, pour Février 2018)"`
	Parsers  []string `names:"--parsers" desc:"Parseurs à employer (ex: altares,cotisation)"`
}

func (params checkBatchHandler) Documentation() flag.Flag {
	return flag.Flag{
		Usage: "Vérifie la validité d'un batch avant son importation",
		Desc: `
		Vérifie la validité du batch sur le point d'être importé et des fichiers qui le constituent.
		Il est possible de limiter l'exécution à certains parsers en spécifiant la liste dans le flag "--parsers".
		Répond avec un propriété JSON "reports" qui contient les rapports textuels de parsing de chaque fichier.
	`,
	}
}

func (params checkBatchHandler) IsEnabled() bool {
	return params.Enable
}

func (params checkBatchHandler) Validate() error {
	if params.BatchKey == "" {
		return errors.New("paramètre `batch` obligatoire")
	}
	return nil
}

func (params checkBatchHandler) Run() error {
	batch := base.AdminBatch{}
	err := engine.Load(&batch, params.BatchKey)
	if err != nil {
		return errors.New("Batch inexistant: " + err.Error())
	}

	parsers, err := parsing.ResolveParsers(params.Parsers)
	if err != nil {
		return err
	}

	eventSink := engine.NewPostgresEventSink(engine.Db.PostgresDB)
	err = engine.CheckBatch(batch, parsers, eventSink)
	if err != nil {
		return fmt.Errorf("erreurs détectées: %v", err)
	}
	return nil
}

func printJSON(object interface{}) {
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
	parsers, err := parsing.ResolveParsers([]string{params.Parser})
	if err != nil {
		return err
	}

	file := base.NewBatchFileWithBasePath(params.File, "")
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

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		for tuple := range outputChannel {
			printJSON(tuple) // écriture du tuple dans la sortie standard
		}
	}()

	for e := range eventChannel {
		res, _ := json.MarshalIndent(e, "", "  ")
		log.Println(string(res)) // écriture de l'événement dans stderr
	}

	// Only return once all channels are closed
	wg.Wait()

	return nil
}
