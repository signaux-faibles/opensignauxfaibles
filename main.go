package main

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	cosFlag "github.com/cosiner/flag"

	"github.com/signaux-faibles/opensignauxfaibles/lib/engine"

	"github.com/signaux-faibles/opensignauxfaibles/lib/naf"
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
	cmdHandlerWithArgs := parseCommandFromArgs()
	// exit if no command was recognized in args
	if cmdHandlerWithArgs == nil {
		fmt.Printf("Commande non reconnue. Utilisez %v --help pour lister les commandes.\n", strings.Join(os.Args, " "))
		os.Exit(1)
		return
	}
	// validate command parameters
	if err := cmdHandlerWithArgs.Validate(); err != nil {
		fmt.Printf("Erreur: %v. Utilisez %v --help pour consulter la documentation.", err, strings.Join(os.Args, " "))
		os.Exit(2)
	}
	// execute the command
	connectDb()
	if err := cmdHandlerWithArgs.Run(); err != nil {
		fmt.Printf("\nErreur: %v\n", err)
		os.Exit(3)
	}
	engine.FlushEventQueue()
}

// Ask cosiner/flag to parse arguments
func parseCommandFromArgs() commandHandler {
	var actualArgs = cliCommands{}
	actualArgs.populateFromArgs()
	for _, cmdHandlerWithArgs := range actualArgs.index() {
		if cmdHandlerWithArgs.IsEnabled() {
			return cmdHandlerWithArgs
		}
	}
	return nil
}

// Interface that each command should implement
type commandHandler interface {
	Documentation() cosFlag.Flag // returns documentation to display in the CLI
	IsEnabled() bool             // returns true when the user invokes this command from the CLI
	Validate() error             // returns an error if some command parameters don't meet expectations
	Run() error                  // executes the command and return an error if it fails
}

// List of command handlers that cosiner/flag should recognize in CLI arguments.
// Each entry will be populated with parameters parsed from command line arguments.
// Each entry must implement the commandHandler interface.
type cliCommands struct {
	Purge             purgeBatchHandler
	Check             checkBatchHandler
	PruneEntities     pruneEntitiesHandler
	Import            importBatchHandler
	PurgeNotCompacted purgeNotCompactedHandler
	Validate          validateHandler
	Compact           compactHandler
	Reduce            reduceHandler
	Public            publicHandler
	Etablissements    exportEtablissementsHandler
	Entreprises       exportEntreprisesHandler
}

func (cmds *cliCommands) populateFromArgs() {
	flagSet := cosFlag.NewFlagSet(cosFlag.Flag{})
	_ = flagSet.ParseStruct(cmds, os.Args...) // may panic with "unexpected non-flag value: unknown_command"
}

// Metadata returns the documentation that will be displayed by cosiner/flag
// if the user invokes "--help", or if some parameters are invalid.
func (cmds *cliCommands) Metadata() map[string]cosFlag.Flag {
	commandMetadata := map[string]cosFlag.Flag{}
	// we use reflection to get the documentation of each prop from cliCommands
	for cmdName, cmdArgs := range cmds.index() {
		commandMetadata[cmdName] = cmdArgs.Documentation()
	}
	return commandMetadata
}

// List and index the commandHandler entries, using reflection.
func (cmds *cliCommands) index() map[string]commandHandler {
	commandByName := map[string]commandHandler{}
	supportedCommands := reflect.ValueOf(*cmds)
	for i := 0; i < supportedCommands.NumField(); i++ {
		fieldName := supportedCommands.Type().Field(i).Name                    // e.g. "PruneEntities"
		cmdName := strings.ToLower(fieldName[0:1]) + fieldName[1:]             // e.g. "pruneEntities"
		cmdArgs, ok := supportedCommands.Field(i).Interface().(commandHandler) // e.g. pruneEntitiesHandler instance
		if ok != true {
			panic(fmt.Sprintf("Property %v of type cliCommands is not an instance of commandHandler", fieldName))
		}
		commandByName[cmdName] = cmdArgs
	}
	return commandByName
}
