package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

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

	if err := runCommand(os.Args[1:]); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var cmds = []commandDefinition{
	// {
	// 	Name: "purgeBatch",
	// 	Run: func([]string) error {
	// 		var params purgeBatchParams
	// 		return purgeBatchHandler(params)
	// 	},
	// },
	{
		Name: "hello",
		Run: func(args []string) error {
			params := helloParams{}
			flag.StringVar(&params.FirstName, "first", "", "First name")
			flag.StringVar(&params.LastName, "last", "", "First name")
			flag.CommandLine.Parse(args)
			return helloHandler(params)
		},
	},
}

type commandDefinition struct {
	Name string
	Run  func([]string) error
}

func runCommand(args []string) error {
	if len(args) < 1 {
		printSupportedCommands()
		return errors.New("Error: You must pass a command")
	}

	command := os.Args[1]

	for _, cmd := range cmds {
		if cmd.Name == command {
			return cmd.Run(os.Args[2:])
		}
	}

	printSupportedCommands()
	return fmt.Errorf("Unknown command: %s", command)
}

func printSupportedCommands() {
	fmt.Println("Supported commands:")
	for _, cmd := range cmds {
		fmt.Printf(" - %s\n", cmd.Name)
	}
}

type helloParams struct {
	FirstName string `names:"-f, --first-name" usage:"hello -f Adrien"`
	LastName  string `names:"-l, --last-name"`
}

func helloHandler(params helloParams) error {
	if params.FirstName == "" {
		return errors.New("firstname is not defined")
	}
	fmt.Printf("Hello %s %s\n", params.FirstName, params.LastName)
	return nil
}
