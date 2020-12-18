package main

import (
	"errors"
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

	err := runCommand(os.Args[1:])
	if err != nil {
		fmt.Println(err)
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

func runCommand(args []string) error {
	if len(args) < 1 {
		printSupportedCommands()
		return errors.New("Error: You must pass a command")
	}

	// fallback: handle new commands
	newCmd, cmdDef := getNewCommand()
	if newCmd != nil {
		err := newCmd.Validate()
		if err != nil {
			cmdDef.Help(false) // display usage information for this command only
			fmt.Println()
			return err
		}
		connectDb()
		return newCmd.Run()
	}

	// no match
	printSupportedCommands()
	return fmt.Errorf("Unknown command")
}

func printSupportedCommands() {
	fmt.Println("usage: sfdata <command> [--boolean-flag] [--parameter=<value1>,<value2>,...]")
	fmt.Println("")
	fmt.Println("Supported commands:")
	fmt.Println("")

	orderedNewCommandNames := []string{
		"purge",
		"check",
		"pruneEntities",
		"import",
		"validate",
		"compact",
		"reduce",
		"public",
		"etablissements",
		"entreprises",
	}
	commandsMeta := (&cliCommands{}).Metadata()
	for _, cmdName := range orderedNewCommandNames {
		fmt.Printf("   %-16s %s\n", cmdName, commandsMeta[cmdName].Usage)
	}
	fmt.Println("")
}

// Function that uses cosiner/flag to parse CLI args.
func getNewCommand() (command, *cosFlag.FlagSet) {
	var actualArgs = cliCommands{}
	flagSet := cosFlag.NewFlagSet(cosFlag.Flag{})
	flagSet.ParseStruct(&actualArgs, os.Args...)
	// check which command was recognized, based on the fields of cliCommands
	supportedCommands := reflect.ValueOf(actualArgs)
	for i := 0; i < supportedCommands.NumField(); i++ {
		fieldName := supportedCommands.Type().Field(i).Name             // e.g. PruneEntities
		cmdName := strings.ToLower(fieldName[0:1]) + fieldName[1:]      // e.g. pruneEntities
		cmdArgs, ok := supportedCommands.Field(i).Interface().(command) // e.g. pruneEntitiesHandler instance
		if ok != true {
			panic(fmt.Sprintf("Property %v of type cliCommands is not an instance of command", fieldName))
		}
		if cmdArgs.IsEnabled() {
			cmdDef, _ := flagSet.FindSubset(cmdName)
			return cmdArgs, cmdDef
		}
	}
	// no command was recognized in args
	return nil, nil
	// flagSet.Help(false) // display usage information, with list of supported commands
	// os.Exit(1)
}
