package main

import (
	"encoding/json"
	"errors"
	"fmt"

	flag "github.com/cosiner/flag"

	"github.com/signaux-faibles/opensignauxfaibles/lib/base"
	"github.com/signaux-faibles/opensignauxfaibles/lib/engine"
	"github.com/signaux-faibles/opensignauxfaibles/lib/parsing"

	"github.com/globalsign/mgo/bson"
)

type purgeBatchHandler struct {
	Enable                 bool   // set to true by cosiner/flag if the user is running this command
	FromBatchKey           string `names:"--since-batch" arglist:"batch_key" desc:"Identifiant du batch à partir duquel supprimer les données (ex: 1802, pour Février 2018)"`
	Key                    string `json:"debugForKey"` // TODO: populer "debugForKey" (ex: "012345678901234")
	IUnderstandWhatImDoing bool   `names:"--i-understand-what-im-doing" desc:"Nécessaire pour confirmer la suppression de données"`
}

var purgeBatchMetadata = flag.Flag{
	Usage: "Supprime une partie des données compactées",
	Desc: `
		/!\ ce traitement est destructif et irréversible /!\
		Supprime les données dans les objets de la collection RawData pour les batches suivant le numéro de batch donné.
		La propriété "debugForKey" permet de traiter une entreprise en fournissant son siren, le résultat n'impacte pas la collection RawData mais est déversé dans purgeBatch_debug à des fins de vérifications.
		Lorsque "key" n'est pas fourni, le traitement s'exécute sur l'ensemble de la base, et dans ce cas la propriété IUnderstandWhatImDoing doit être fournie à la valeur "true" sans quoi le traitement refusera de se lancer.
		Répond "ok" dans la sortie standard, si le traitement s'est bien déroulé.
		/!\ ce traitement est destructif et irréversible /!\
		`,
}

func (params purgeBatchHandler) IsEnabled() bool {
	return params.Enable
}

func (params purgeBatchHandler) Validate() error {
	if params.FromBatchKey == "" {
		return errors.New("paramètre `since-batch` obligatoire")
	}
	return nil
}

func (params purgeBatchHandler) Run() error {
	var batch base.AdminBatch
	err := engine.Load(&batch, params.FromBatchKey)
	if err != nil {
		return errors.New("le batch " + params.FromBatchKey + " n'est pas accessible: " + err.Error())
	}

	if params.Key != "" {
		err = engine.PurgeBatchOne(batch, params.Key)
		if err != nil {
			return errors.New("erreur pendant le MapReduce: " + err.Error())
		}
	} else {
		if !params.IUnderstandWhatImDoing {
			return errors.New("pour une purge de la base complète, IUnderstandWhatImDoing doit être `true`")
		}
		err = engine.PurgeBatch(batch)
		if err != nil {
			return errors.New("(✖╭╮✖) le traitement n'a pas abouti: " + err.Error())
		}
	}
	printJSON("ok")
	return nil
}

type importBatchParams struct {
	BatchKey string   `json:"batch"`
	Parsers  []string `json:"parsers"`
	NoFilter bool     `json:"noFilter"`
}

// importBatchHandler traite les demandes d'import par l'API
// on peut demander l'exécution de tous les parsers sans fournir d'option
// ou demander l'exécution de parsers particuliers en fournissant une liste de leurs codes.
func importBatchHandler(params importBatchParams) error {
	batch := base.AdminBatch{}
	err := engine.Load(&batch, params.BatchKey)
	if err != nil {
		return errors.New("Batch inexistant: " + err.Error())
	}

	parsers, err := parsing.ResolveParsers(params.Parsers)
	if err != nil {
		return err
	}

	dataChan := engine.InsertIntoImportedData(engine.Db.DB)
	err = engine.ImportBatch(batch, parsers, params.NoFilter, dataChan)
	if err != nil {
		return err
	}

	printJSON("ok")
	return nil
}

type checkBatchParams struct {
	BatchKey string   `json:"batch"`
	Parsers  []string `json:"parsers"`
}

func checkBatchHandler(params checkBatchParams) error {

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

	printJSON(bson.M{"reports": reports})
	return nil
}

func printJSON(object interface{}) {
	res, _ := json.Marshal(object)
	fmt.Println(string(res))
}

func purgeNotCompactedHandler() error {
	return engine.PurgeNotCompacted()
}

type pruneEntitiesParams struct {
	BatchKey string `json:"batch"`
	Delete   bool   `json:"delete"`
}

// Count – then delete – companies from RawData that should have been
// excluded by the SIREN Filter.
func pruneEntitiesHandler(params pruneEntitiesParams) error {
	count, err := engine.PruneEntities(params.BatchKey, params.Delete)
	if err == nil {
		printJSON(bson.M{"count": count})
	}
	return err
}
