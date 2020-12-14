package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/signaux-faibles/opensignauxfaibles/lib/engine"

	"github.com/signaux-faibles/opensignauxfaibles/lib/naf"

	_ "github.com/signaux-faibles/opensignauxfaibles/docs"
)

func connectDb() {
	engine.Db = engine.InitDB()
	go engine.MessageSocketAddClient()

	var err error
	naf.Naf, err = naf.LoadNAF()
	if err != nil {
		panic(err)
	}
}

// main Fonction Principale
func main() {

	// connectDb()
	// api.POST("/data/batch/purge", purgeBatchHandler)      // [x] écrit dans Journal
	// api.POST("/data/check", checkBatchHandler)            // [x] écrit dans Journal
	// api.POST("/data/import", importBatchHandler)          // [x] écrit dans Journal
	// api.POST("/data/validate", validateHandler)           // [x] écrit dans Journal
	// api.POST("/data/compact", compactHandler)             // [x] écrit dans Journal
	// api.POST("/data/reduce", reduceHandler)               // [x] écrit dans Journal
	// api.POST("/data/public", publicHandler)               // [x] écrit dans Journal
	// api.POST("/data/pruneEntities", pruneEntitiesHandler) // [x] écrit dans Journal
	// api.GET("/data/purgeNotCompacted", purgeNotCompactedHandler)
	// api.GET("/data/etablissements", exportEtablissementsHandler)
	// api.GET("/data/entreprises", exportEntreprisesHandler)

	err := runCommand(os.Args[1:])
	time.Sleep(2 * time.Second) // TODO: trouver un meilleur moyen d'assurer que les données ont fini d'être enregistrées en db

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var cmds = map[string]commandDefinition{
	// {
	// 	Name: "purgeBatch",
	// 	Run: func([]string) error {
	// 		var params purgeBatchParams
	// 		return purgeBatchHandler(params)
	// 	},
	// },
	"check": func(args []string) error {
		var parsers string
		params := checkBatchParams{}
		flag.StringVar(&params.BatchKey, "batch", "", "Batch identifier")
		flag.StringVar(&parsers, "parsers", "", "List of parsers")
		flag.CommandLine.Parse(args)
		params.Parsers = strings.Split(parsers, ",")
		connectDb()
		return checkBatchHandler(params)
	},
	"pruneEntities": func(args []string) error {
		params := pruneEntitiesParams{}
		flag.StringVar(&params.BatchKey, "batch", "", "Batch identifier")
		flag.BoolVar(&params.Delete, "delete", false, "Delete entities")
		flag.CommandLine.Parse(args)
		connectDb()
		return pruneEntitiesHandler(params)
	},
	"compact": func(args []string) error {
		params := compactParams{}
		flag.StringVar(&params.FromBatchKey, "since-batch", "", "Batch identifier")
		flag.CommandLine.Parse(args)
		connectDb()
		return compactHandler(params)
	},
	"reduce": func(args []string) error {
		params := reduceParams{} // TODO: also populate other parameters
		flag.StringVar(&params.BatchKey, "until-batch", "", "Batch identifier")
		flag.StringVar(&params.Key, "key", "", "SIRET or SIREN to focus on")
		flag.CommandLine.Parse(args)
		connectDb()
		return reduceHandler(params)
	},
	"public": func(args []string) error {
		params := publicParams{} // TODO: also populate other parameters
		flag.StringVar(&params.BatchKey, "until-batch", "", "Batch identifier")
		flag.StringVar(&params.Key, "key", "", "SIRET or SIREN to focus on")
		flag.CommandLine.Parse(args)
		connectDb()
		return publicHandler(params)
	},
}

type commandDefinition func(args []string) error

func runCommand(args []string) error {
	if len(args) < 1 {
		printSupportedCommands()
		return errors.New("Error: You must pass a command")
	}

	command := os.Args[1]
	commandFct := cmds[command]
	if commandFct != nil {
		return commandFct(os.Args[2:])
	}

	printSupportedCommands()
	return fmt.Errorf("Unknown command: %s", command)
}

func printSupportedCommands() {
	fmt.Println("Supported commands:")
	for cmd := range cmds {
		fmt.Printf(" - %s\n", cmd)
	}
}
