package main

import (
	"errors"
	"flag"
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
	Purge         purgeBatchHandler
	Check         checkBatchHandler
	PruneEntities pruneEntitiesHandler
	Import        importBatchHandler
}

// Metadata returns the documentation that will be displayed by cosiner/flag
// if the user invokes "--help", or if some parameters are invalid.
func (cmds *cliCommands) Metadata() map[string]cosFlag.Flag {
	return map[string]cosFlag.Flag{
		"purge":         cmds.Purge.Documentation(),
		"check":         cmds.Check.Documentation(),
		"pruneEntities": cmds.PruneEntities.Documentation(),
		"import":        cmds.Import.Documentation(),
	}
}

// Structure to define commands that will be migrated over cosiner/flag's format.
type legacyCommandDefinition struct {
	name    string
	summary string
	run     func(args []string) error
}

// List of commands that will be migrated over cosiner/flag's format.
var legacyCommandDefs = []*legacyCommandDefinition{
	{
		"validate",
		"Liste les entrées de données invalides",
		/**
		Vérifie la validité des entrées de données contenues dans les documents de la collection RawData ou ImportedData.
		Répond en listant dans la sortie standard les entrées invalides au format JSON.
		*/
		func(args []string) error {
			params := validateParams{}
			flag.StringVar(&params.Collection, "collection", "", "Nom de la collection à valider: RawData ou ImportedData")
			flag.CommandLine.Parse(args)
			connectDb()
			return validateHandler(params) // [x] écrit dans Journal
		}}, {
		"compact",
		"Compacte la base de données",
		/**
		Ce traitement permet le compactage de la base de données.
		Ce compactage a pour effet de réduire tous les objets en clé uniques comportant dans la même arborescence toutes les données en rapport avec ces clés.
		Ce traitement est nécessaire avant l'usage des commandes `reduce` et `public`, après chaque import de données.
		Répond "ok" dans la sortie standard, si le traitement s'est bien déroulé.
		*/
		func(args []string) error {
			params := compactParams{}
			flag.StringVar(&params.FromBatchKey, "since-batch", "", "Identifiant du batch à partir duquel compacter (ex: `1802`, pour Février 2018)")
			flag.CommandLine.Parse(args)
			connectDb()
			return compactHandler(params) // [x] écrit dans Journal
		}}, {
		// "purgeNotCompacted": {"TODO - summary", func(args []string) error {
		// 	return purgeNotCompactedHandler() // TODO: écrire rapport dans Journal ?
		// }},
		"reduce",
		"Calcule les variables destinées à la prédiction",
		/**
		Alimente la collection Features en calculant les variables avec le traitement mapreduce demandé dans la propriété `features`.
		Le traitement remplace les objets similaires en sortie du calcul dans la collection Features, les objets non concernés par le traitement ne seront ainsi pas remplacés, de sorte que si un seul siret est demandé le calcul ne remplacera qu'un seul objet.
		Ces traitements ne prennent en compte que les objets déjà compactés.
		Répond "ok" dans la sortie standard, si le traitement s'est bien déroulé.
		*/
		func(args []string) error {
			params := reduceParams{} // TODO: also populate other parameters
			flag.StringVar(&params.BatchKey, "until-batch", "", "Identifiant du batch jusqu'auquel calculer (ex: `1802`, pour Février 2018)")
			flag.StringVar(&params.Key, "key", "", "Numéro SIRET or SIREN d'une entité à calculer exclusivement")
			flag.CommandLine.Parse(args)
			connectDb()
			return reduceHandler(params) // [x] écrit dans Journal
		}}, {
		"public",
		"Génère les données destinées au site web",
		/**
		Alimente la collection Public avec les objets calculés pour le batch cité en paramètre, à partir de la collection RawData.
		Le traitement prend en paramètre la clé du batch (obligatoire) et un SIREN (optionnel). Lorsque le SIREN n'est pas précisé, tous les objets lié au batch sont traités, à conditions qu'ils soient dans le périmètre de scoring "algo2".
		Cette collection sera ensuite accédée par les utilisateurs pour consulter les données des entreprises.
		Des niveaux d'accéditation fins (ligne ou colonne) pour la consultation de ces données peuvent être mis en oeuvre.
		Ces filtrages sont effectués grace à la notion de scope. Les objets et les utilisateurs disposent d'un ensemble de tags et les objets partageant au moins un tag avec les utilisateurs peuvent être consultés par ceux-ci.
		Ces tags sont exploités pour traiter la notion de région (ligne) mais aussi les permissions (colonne).
		Répond "ok" dans la sortie standard, si le traitement s'est bien déroulé.
		*/
		func(args []string) error {
			params := publicParams{}
			flag.StringVar(&params.BatchKey, "until-batch", "", "Identifiant du batch jusqu'auquel calculer (ex: `1802`, pour Février 2018)")
			flag.StringVar(&params.Key, "key", "", "Numéro SIRET or SIREN d'une entité à calculer exclusivement")
			flag.CommandLine.Parse(args)
			connectDb()
			return publicHandler(params) // [x] écrit dans Journal
		}}, {
		"etablissements",
		"Exporte la liste des établissements",
		/**
		Exporte la liste des établissements depuis la collection Public.
		Répond dans la sortie standard une ligne JSON par établissement.
		*/
		func(args []string) error {
			params := exportParams{}
			flag.StringVar(&params.Key, "key", "", "Numéro SIREN à utiliser pour filtrer les résultats.")
			flag.CommandLine.Parse(args)
			connectDb()
			return exportEtablissementsHandler(params) // TODO: écrire rapport dans Journal ?
		}}, {
		"entreprises",
		"Exporte la liste des entreprises",
		/**
		Exporte la liste des entreprises depuis la collection Public.
		Répond dans la sortie standard une ligne JSON par entreprise.
		*/
		func(args []string) error {
			params := exportParams{}
			flag.StringVar(&params.Key, "key", "", "Numéro SIREN à utiliser pour filtrer les résultats.")
			flag.CommandLine.Parse(args)
			connectDb()
			return exportEntreprisesHandler(params) // TODO: écrire rapport dans Journal ?
		}},
}

func runCommand(args []string) error {
	if len(args) < 1 {
		printSupportedCommands()
		return errors.New("Error: You must pass a command")
	}

	// handle legacy commands
	var legacyCmds = map[string]*legacyCommandDefinition{}
	for _, commandDef := range legacyCommandDefs {
		legacyCmds[commandDef.name] = commandDef
	}
	command := os.Args[1]
	commandDef := legacyCmds[command]
	if commandDef != nil {
		return commandDef.run(os.Args[2:])
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
	return fmt.Errorf("Unknown command: %s", command)
}

func printSupportedCommands() {
	fmt.Println("usage: sfdata <command> [--boolean-flag] [--parameter=<value1>,<value2>,...]")
	fmt.Println("")
	fmt.Println("Supported commands:")
	fmt.Println("")
	commandsMeta := (&cliCommands{}).Metadata()
	for cmdName, cmdMeta := range commandsMeta {
		fmt.Printf("   %-16s %s\n", cmdName, cmdMeta.Usage)
	}
	for _, cmdDef := range legacyCommandDefs {
		fmt.Printf("   %-16s %s\n", cmdDef.name, cmdDef.summary)
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
		fieldName := supportedCommands.Type().Field(i).Name
		cmdName := strings.ToLower(fieldName[0:1]) + fieldName[1:]
		cmdArgs, ok := supportedCommands.Field(i).Interface().(command)
		if ok != true {
			panic(fmt.Sprintf("Property %v of type cliCommands is not an instance of command", fieldName))
		}
		if cmdArgs.IsEnabled() { // TODO: can we read Enabled property directly, thanks to reflection ?
			cmdDef, _ := flagSet.FindSubset(cmdName)
			return cmdArgs, cmdDef
		}
	}
	// no command was recognized in args
	return nil, nil
	// flagSet.Help(false) // display usage information, with list of supported commands
	// os.Exit(1)
}
