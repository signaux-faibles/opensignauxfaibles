package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

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

var cmds = map[string]*commandDefinition{
	"purge": {
		"Supprime une partie des données compactées",
		/**
		/!\ ce traitement est destructif et irréversible /!\
		Supprime les données dans les objets de la collection RawData pour les batches suivant le numéro de batch donné.
		La propriété `debugForKey` permet de traiter une entreprise en fournissant son siren, le résultat n'impacte pas la collection RawData mais est déversé dans purgeBatch_debug à des fins de vérifications.
		Lorsque `key` n'est pas fourni, le traitement s'exécute sur l'ensemble de la base, et dans ce cas la propriété IUnderstandWhatImDoing doit être fournie à la valeur `true` sans quoi le traitement refusera de se lancer.
		Répond "ok" dans la sortie standard, si le traitement s'est bien déroulé.
		/!\ ce traitement est destructif et irréversible /!\
		*/
		func(args []string) error {
			var params purgeBatchParams
			// TODO: populer "debugForKey" (ex: "012345678901234")
			flag.StringVar(&params.FromBatchKey, "since-batch", "", "Identifiant du batch à partir duquel supprimer les données (ex: `1802`, pour Février 2018)")
			flag.BoolVar(&params.IUnderstandWhatImDoing, "i-understand-what-im-doing", false, "Nécessaire pour confirmer la suppression de données")
			flag.CommandLine.Parse(args)
			connectDb()
			return purgeBatchHandler(params) // [x] écrit dans Journal
		}},
	"check": {
		"Vérifie la validité d'un batch avant son importation",
		/**
		Vérifie la validité du batch sur le point d'être importé et des fichiers qui le constituent.
		Pour exécuter tous les parsers, il faut ne pas spécifier la propriété parsers ou lui donner la valeur null.
		Répond avec un propriété JSON "reports" qui contient les rapports textuels de parsing de chaque fichier.
		*/
		func(args []string) error {
			var parsers string
			params := checkBatchParams{}
			flag.StringVar(&params.BatchKey, "batch", "", "Identifiant du batch à vérifier (ex: `1802`, pour Février 2018)")
			flag.StringVar(&parsers, "parsers", "", "Parseurs à employer (ex: `altares,cotisation`)")
			flag.CommandLine.Parse(args)
			params.Parsers = strings.Split(parsers, ",")
			connectDb()
			return checkBatchHandler(params) // [x] écrit dans Journal
		}},
	"pruneEntities": {
		"Compte/supprime les entités hors périmètre",
		/**
		Compte puis supprime dans la collection `RawData` les entités (établissements et entreprises)
		non listées dans le filtre de périmètre du batch spécifié.
		Répond avec un propriété JSON "count" qui vaut le nombre d'entités hors périmètre comptées ou supprimées.
		*/
		func(args []string) error {
			params := pruneEntitiesParams{}
			flag.StringVar(&params.BatchKey, "batch", "", "Identifiant du batch à nettoyer (ex: `1802`, pour Février 2018)")
			flag.BoolVar(&params.Delete, "delete", false, "Nécessaire pour confirmer la suppression de données")
			flag.CommandLine.Parse(args)
			connectDb()
			return pruneEntitiesHandler(params) // [x] écrit dans Journal
		}},
	"import": {
		"Importe des fichiers",
		/**
		Effectue l'import de tous les fichiers du batch donné en paramètre.
		Pour exécuter tous les parsers, il faut ne pas spécifier la propriété parsers ou lui donner la valeur null.
		Répond "ok" dans la sortie standard, si le traitement s'est bien déroulé.
		*/
		func(args []string) error {
			params := importBatchParams{}
			// TODO: populer "Parsers", documentation: "Parseurs à employer (ex: `altares,cotisation`)"
			flag.StringVar(&params.BatchKey, "batch", "", "Identifiant du batch à importer (ex: `1802`, pour Février 2018)")
			flag.BoolVar(&params.NoFilter, "no-filter", false, "Pour procéder à l'importation même si aucun filtre n'est fourni")
			flag.CommandLine.Parse(args)
			connectDb()
			return importBatchHandler(params) // [x] écrit dans Journal
		}},
	"validate": {
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
		}},
	"compact": {
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
		}},
	// "purgeNotCompacted": {"TODO - summary", func(args []string) error {
	// 	return purgeNotCompactedHandler() // TODO: écrire rapport dans Journal ?
	// }},
	"reduce": {
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
		}},
	"public": {
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
		}},
	"etablissements": {
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
		}},
	"entreprises": {
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

type commandDefinition struct {
	summary string
	run     func(args []string) error
}

func runCommand(args []string) error {
	if len(args) < 1 {
		printSupportedCommands()
		return errors.New("Error: You must pass a command")
	}

	command := os.Args[1]
	commandDef := cmds[command]
	if commandDef != nil {
		return commandDef.run(os.Args[2:])
	}

	printSupportedCommands()
	return fmt.Errorf("Unknown command: %s", command)
}

func printSupportedCommands() {
	fmt.Println("Supported commands:")
	for cmd, cmdDef := range cmds {
		fmt.Printf(" - %s (%s)\n", cmd, cmdDef.summary)
	}
}
