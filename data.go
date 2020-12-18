package main

import (
	"errors"

	flag "github.com/cosiner/flag"

	"github.com/signaux-faibles/opensignauxfaibles/lib/base"
	"github.com/signaux-faibles/opensignauxfaibles/lib/engine"
)

type reduceHandler struct {
	Enable   bool     // set to true by cosiner/flag if the user is running this command
	BatchKey string   `names:"--until-batch" arglist:"batch_key" desc:"Identifiant du batch jusqu'auquel calculer (ex: 1802, pour Février 2018)"`
	Key      string   `names:"--key" desc:"Numéro SIRET or SIREN d'une entité à calculer exclusivement"`
	From     string   `names:"--from"`                                                                                                // TODO: à définir et tester
	To       string   `names:"--to"`                                                                                                  // TODO: à définir et tester
	Types    []string `names:"--type" arglist:"all|apart" desc:"Sélection des types de données qui vont être calculés ou recalculés"` // Valeurs autorisées pour l'instant: "apart", "all"
}

func (params reduceHandler) Documentation() flag.Flag {
	return flag.Flag{
		Usage: "Calcule les variables destinées à la prédiction",
		Desc: `
		Alimente la collection Features en calculant les variables avec le traitement mapreduce demandé dans la propriété "features".
		Le traitement remplace les objets similaires en sortie du calcul dans la collection Features, les objets non concernés par le traitement ne seront ainsi pas remplacés, de sorte que si un seul siret est demandé le calcul ne remplacera qu'un seul objet.
		Ces traitements ne prennent en compte que les objets déjà compactés.
		Répond "ok" dans la sortie standard, si le traitement s'est bien déroulé.
		`,
	}
}

func (params reduceHandler) IsEnabled() bool {
	return params.Enable
}

func (params reduceHandler) Validate() error {
	if params.BatchKey == "" {
		return errors.New("paramètre `until-batch` obligatoire")
	}
	return nil
}

func (params reduceHandler) Run() error {

	batch, err := engine.GetBatch(params.BatchKey)
	if err != nil {
		return errors.New("Batch inexistant: " + err.Error())
	}

	if params.Key == "" && params.From == "" && params.To == "" {
		err = engine.Reduce(batch, params.Types)
	} else {
		err = engine.ReduceOne(batch, params.Key, params.From, params.To, params.Types)
	}

	if err != nil {
		return err
	}

	printJSON("Traitement effectué")
	return nil
}

type publicHandler struct {
	Enable   bool   // set to true by cosiner/flag if the user is running this command
	BatchKey string `names:"--until-batch" arglist:"batch_key" desc:"Identifiant du batch jusqu'auquel calculer (ex: 1802, pour Février 2018)"`
	Key      string `names:"--key" desc:"Numéro SIRET or SIREN d'une entité à calculer exclusivement"`
}

func (params publicHandler) Documentation() flag.Flag {
	return flag.Flag{
		Usage: "Génère les données destinées au site web",
		Desc: `
		Alimente la collection Public avec les objets calculés pour le batch cité en paramètre, à partir de la collection RawData.
		Le traitement prend en paramètre la clé du batch (obligatoire) et un SIREN (optionnel). Lorsque le SIREN n'est pas précisé, tous les objets lié au batch sont traités, à conditions qu'ils soient dans le périmètre de scoring "algo2".
		Cette collection sera ensuite accédée par les utilisateurs pour consulter les données des entreprises.
		Des niveaux d'accéditation fins (ligne ou colonne) pour la consultation de ces données peuvent être mis en oeuvre.
		Ces filtrages sont effectués grace à la notion de scope. Les objets et les utilisateurs disposent d'un ensemble de tags et les objets partageant au moins un tag avec les utilisateurs peuvent être consultés par ceux-ci.
		Ces tags sont exploités pour traiter la notion de région (ligne) mais aussi les permissions (colonne).
		Répond "ok" dans la sortie standard, si le traitement s'est bien déroulé.
		`,
	}
}

func (params publicHandler) IsEnabled() bool {
	return params.Enable
}

func (params publicHandler) Validate() error {
	if params.BatchKey == "" {
		return errors.New("paramètre `until-batch` obligatoire")
	}
	if len(params.Key) < 9 {
		return errors.New("la clé fait moins de 9 caractères (siren)")
	}
	return nil
}

func (params publicHandler) Run() error {

	batch := base.AdminBatch{}
	err := engine.Load(&batch, params.BatchKey)
	if err != nil {
		return errors.New("batch non trouvé")
	}

	if params.Key == "" {
		err = engine.Public(batch)
	} else if len(params.Key) >= 9 {
		err = engine.PublicOne(batch, params.Key[0:9])
	} else {
		return errors.New("la clé fait moins de 9 caractères (siren)")
	}

	if err == nil {
		printJSON("ok")
	}
	return err
}

type compactHandler struct {
	Enable       bool   // set to true by cosiner/flag if the user is running this command
	FromBatchKey string `names:"--since-batch" arglist:"batch_key" desc:"Identifiant du batch à partir duquel compacter (ex: 1802, pour Février 2018)"`
}

func (params compactHandler) Documentation() flag.Flag {
	return flag.Flag{
		Usage: "Compacte la base de données",
		Desc: `
		Ce traitement permet le compactage de la base de données.
		Ce compactage a pour effet de réduire tous les objets en clé uniques comportant dans la même arborescence toutes les données en rapport avec ces clés.
		Ce traitement est nécessaire avant l'usage des commandes "reduce" et "public", après chaque import de données.
		Répond "ok" dans la sortie standard, si le traitement s'est bien déroulé.
		`,
	}
}

func (params compactHandler) IsEnabled() bool {
	return params.Enable
}

func (params compactHandler) Validate() error {
	if params.FromBatchKey == "" {
		return errors.New("paramètre `since-batch` obligatoire")
	}
	return nil
}

func (params compactHandler) Run() error {
	err := engine.Compact(params.FromBatchKey)
	if err == nil {
		printJSON("ok")
	}
	return err
}

type exportEtablissementsHandler struct {
	Enable bool   // set to true by cosiner/flag if the user is running this command
	Key    string `names:"--key" desc:"Numéro SIREN à utiliser pour filtrer les résultats (ex: 012345678)"`
}

func (params exportEtablissementsHandler) Documentation() flag.Flag {
	return flag.Flag{
		Usage: "Exporte la liste des établissements",
		Desc: `
		Exporte la liste des établissements depuis la collection Public.
		Répond dans la sortie standard une ligne JSON par établissement.
		`,
	}
}

func (params exportEtablissementsHandler) IsEnabled() bool {
	return params.Enable
}

func (params exportEtablissementsHandler) Validate() error {
	if !(len(params.Key) == 9 || len(params.Key) == 0) {
		return errors.New("si fourni, paramètre `key` doit être un numéro SIREN (9 chiffres)")
	}
	return nil
}

func (params exportEtablissementsHandler) Run() error {
	return engine.ExportEtablissements(params.Key)
}

type exportParams struct {
	Key string `json:"key"`
}

func exportEntreprisesHandler(params exportParams) error {
	return engine.ExportEntreprises(params.Key)
}

type validateHandler struct {
	Enable     bool   // set to true by cosiner/flag if the user is running this command
	Collection string `names:"--collection" arglist:"RawData|ImportedData" desc:"Nom de la collection à valider"`
}

func (params validateHandler) Documentation() flag.Flag {
	return flag.Flag{
		Usage: "Liste les entrées de données invalides",
		Desc: `
		Vérifie la validité des entrées de données contenues dans les documents de la collection RawData ou ImportedData.
		Répond en listant dans la sortie standard les entrées invalides au format JSON.
		`,
	}
}

func (params validateHandler) IsEnabled() bool {
	return params.Enable
}

func (params validateHandler) Validate() error {
	if params.Collection != "RawData" && params.Collection != "ImportedData" {
		return errors.New("le paramètre collection doit valoir RawData ou ImportedData")
	}
	return nil
}

func (params validateHandler) Run() error {

	jsonSchema, err := engine.LoadJSONSchemaFiles()
	if err != nil {
		return err
	}

	err = engine.ValidateDataEntries(jsonSchema, params.Collection)
	if err != nil {
		return err
	}

	return nil
}
