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
	err := runCommand()
	if err != nil {
		fmt.Printf("\nErreur: %v\n", err)
		os.Exit(1)
	}
	engine.FlushEventQueue()
}

// Interface that each command should implement
type command interface {
	Documentation() cosFlag.Flag // returns documentation to display in the CLI
	IsEnabled() bool             // returns true when the user invokes this command from the CLI
	Validate() error             // returns an error if some command parameters don't meet expectations
	Run() error                  // executes the command and return an error if it fails
}

// List of command handlers that cosiner/flag should recognize in CLI arguments
type cliCommands struct {
	Purge          purgeBatchHandler
	Check          checkBatchHandler
	PruneEntities  pruneEntitiesHandler
	Import         importBatchHandler
	Validate       validateHandler
	Compact        compactHandler
	Reduce         reduceHandler
	Public         publicHandler
	Etablissements exportEtablissementsHandler
	Entreprises    exportEntreprisesHandler
}

// Metadata returns the documentation that will be displayed by cosiner/flag
// if the user invokes "--help", or if some parameters are invalid.
func (cmds *cliCommands) Metadata() map[string]cosFlag.Flag {
	// TODO: use reflection to generate that list automatically
	return map[string]cosFlag.Flag{
		"purge":          cmds.Purge.Documentation(),
		"check":          cmds.Check.Documentation(),
		"pruneEntities":  cmds.PruneEntities.Documentation(),
		"import":         cmds.Import.Documentation(),
		"validate":       cmds.Validate.Documentation(),
		"compact":        cmds.Compact.Documentation(),
		"reduce":         cmds.Reduce.Documentation(),
		"public":         cmds.Public.Documentation(),
		"etablissements": cmds.Etablissements.Documentation(),
		"entreprises":    cmds.Entreprises.Documentation(),
	}
}

// "purgeNotCompacted": {"TODO - summary", func(args []string) error {
// 	return purgeNotCompactedHandler() // TODO: écrire rapport dans Journal ?
// }},

// Detect and run the command specified in CLI args.
// Returns any validation or execution error.
func runCommand() error {
	// ask cosiner/flag to parse arguments
	var actualArgs = cliCommands{}
	flagSet := cosFlag.NewFlagSet(cosFlag.Flag{})
	if err := flagSet.ParseStruct(&actualArgs, os.Args...); err != nil {
		return err // Note: parsing may exit instead of reporting "unexpected non-flag value: unknown_command"
	}
	// find and execute the command, if any
	cmdName, cmdHandler := getCommand(actualArgs)
	if cmdHandler != nil {
		err := cmdHandler.Validate() // validate command parameters
		if err != nil {
			cmdDef, _ := flagSet.FindSubset(cmdName)
			cmdDef.Help(false) // display usage information for this command
			return err
		}
		connectDb()
		return cmdHandler.Run()
	}
	// no command was recognized in args
	flagSet.Help(false) // display usage information, with list of supported commands
	return fmt.Errorf("Commande non reconnue")
}

// Find which command was recognized from CLI args, based on the fields of cliCommands.
func getCommand(actualArgs cliCommands) (string, command) {
	supportedCommands := reflect.ValueOf(actualArgs)
	for i := 0; i < supportedCommands.NumField(); i++ {
		fieldName := supportedCommands.Type().Field(i).Name             // e.g. PruneEntities
		cmdName := strings.ToLower(fieldName[0:1]) + fieldName[1:]      // e.g. pruneEntities
		cmdArgs, ok := supportedCommands.Field(i).Interface().(command) // e.g. pruneEntitiesHandler instance
		if ok != true {
			panic(fmt.Sprintf("Property %v of type cliCommands is not an instance of command", fieldName))
		}
		if cmdArgs.IsEnabled() {
			return cmdName, cmdArgs
		}
	}
	return "", nil
}
