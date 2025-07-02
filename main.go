package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"reflect"
	"strings"

	cosFlag "github.com/cosiner/flag"
	"github.com/spf13/viper"

	"opensignauxfaibles/lib/engine"
	"opensignauxfaibles/lib/marshal"
)

// GitCommit est le hash du dernier commit à inclure dans le binaire.
var GitCommit string // (populé lors de la compilation, par `make build`)

func connectDB() {
	engine.Db = engine.InitDB()
	go engine.InitEventQueue()
}

// main Fonction Principale
func main() {
	initConfig()
	initLogger()
	exitCode := runCLI(os.Args...)
	os.Exit(exitCode)
}

func runCLI(args ...string) int {
	cmdHandlerWithArgs := parseCommandFromArgs(args)
	slog.Debug("run cli", slog.Any("args", args))
	useDB := os.Getenv("NO_DB") != "1"
	// exit if no command was recognized in args
	if cmdHandlerWithArgs == nil {
		fmt.Printf("Commande non reconnue. Utilisez %v --help pour lister les commandes.\n", strings.Join(args, " "))
		return 1
	}
	// validate command parameters
	if err := cmdHandlerWithArgs.Validate(); err != nil {
		fmt.Printf("Erreur: %v. Utilisez %v --help pour consulter la documentation.", err, strings.Join(args, " "))
		return 2
	}
	// execute the command
	if useDB {
		connectDB()
		defer engine.FlushEventQueue()
	}
	if err := cmdHandlerWithArgs.Run(); err != nil {
		fmt.Printf("\nErreur: %v\n", err)
		return 3
	}
	return 0
}

func initConfig() {
	viper.SetConfigType("toml")
	viper.SetConfigName("config") // => will look for config.toml in the following paths:
	viper.AddConfigPath("/etc/opensignauxfaibles")
	viper.AddConfigPath("$HOME/.opensignauxfaibles")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		slog.Error("erreur à la lecture de la config toml", slog.Any("error", err))
	}
	marshal.SetGitCommit(GitCommit)
}

// Ask cosiner/flag to parse arguments
func parseCommandFromArgs(args []string) commandHandler {
	var actualArgs = cliCommands{}
	actualArgs.populateFromArgs(args)
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
	ParseFile parseFileHandler
	Check     checkBatchHandler
	Import    importBatchHandler
}

func (cmds *cliCommands) populateFromArgs(args []string) {
	flagSet := cosFlag.NewFlagSet(cosFlag.Flag{})
	_ = flagSet.ParseStruct(cmds, args...) // may panic with "unexpected non-flag value: unknown_command"
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
		if !ok {
			log.Fatal(fmt.Sprintf("Property %v of type cliCommands is not an instance of commandHandler", fieldName))
		}
		commandByName[cmdName] = cmdArgs
	}
	return commandByName
}
