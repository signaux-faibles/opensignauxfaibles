package main

import (
	"errors"
	"strconv"
	"time"

	"github.com/signaux-faibles/opensignauxfaibles/lib/base"
	"github.com/signaux-faibles/opensignauxfaibles/lib/engine"
)

type reduceParams struct {
	BatchKey string   `json:"batch"`
	Key      string   `json:"key"`
	From     string   `json:"from"`
	To       string   `json:"to"`
	Types    []string `json:"types"`
	// Sélection des types de données qui vont être calculés ou recalculés.
	// Valeurs autorisées pour l'instant: "apart", "all"
}

func reduceHandler(params reduceParams) error {

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

type publicParams struct {
	BatchKey string `json:"batch"`
	Key      string `json:"key"`
}

func publicHandler(params publicParams) error {
	if params.BatchKey == "" {
		return errors.New("batch vide")
	}

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

type compactParams struct {
	FromBatchKey string `json:"fromBatchKey"`
}

func compactHandler(params compactParams) error {
	err := engine.Compact(params.FromBatchKey)
	if err == nil {
		printJSON("ok")
	}
	return err
}

func getTimestamp() string {
	return strconv.FormatInt(time.Now().Unix(), 10)
}

type exportParams struct {
	Key string `json:"key"`
}

func getKeyParam(params exportParams) (string, error) {
	if !(len(params.Key) == 9 || len(params.Key) == 0) {
		return "", errors.New("si fourni, key doit être un numéro SIREN (9 chiffres)")
	}
	return params.Key, nil
}

func exportEtablissementsHandler(params exportParams) error {
	key, err := getKeyParam(params)
	if err != nil {
		return err
	}
	return engine.ExportEtablissements(key)
}

func exportEntreprisesHandler(params exportParams) error {
	key, err := getKeyParam(params)
	if err != nil {
		return err
	}
	return engine.ExportEntreprises(key)
}

type validateParams struct {
	Collection string `json:"collection"`
}

func validateHandler(params validateParams) error {

	if params.Collection != "RawData" && params.Collection != "ImportedData" {
		return errors.New("le paramètre collection doit valoir RawData ou ImportedData")
	}

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
