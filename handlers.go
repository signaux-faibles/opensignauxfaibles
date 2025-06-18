package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"

	"github.com/cosiner/flag"
	"github.com/globalsign/mgo/bson"

	"opensignauxfaibles/lib/base"
	"opensignauxfaibles/lib/engine"
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
		return errors.New("Batch inexistant: " + err.Error())
	}

	parsers, err := parsing.ResolveParsers(params.Parsers)
	if err != nil {
		return err
	}

	dataChan := engine.InsertIntoCSVs()
	err = engine.ImportBatch(batch, parsers, params.NoFilter, dataChan)
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

	reports, err := engine.CheckBatch(batch, parsers)
	if err != nil {
		return errors.New("Erreurs détectées: " + err.Error())
	}

	sort.Strings(reports) // to make sure that parsed files are always listed in the same order
	printJSON(bson.M{"reports": reports})
	return nil
}

func printJSON(object interface{}) {
	res, _ := json.Marshal(object)
	fmt.Println(string(res))
}
