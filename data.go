package main

import (
	"errors"

	flag "github.com/cosiner/flag"

	"opensignauxfaibles/lib/engine"
)

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
