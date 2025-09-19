package main

import (
	"flag"
	"log"

	"github.com/pkg/errors"
	"github.com/spf13/viper"

	"opensignauxfaibles/lib/base"
	"opensignauxfaibles/lib/prepare-import/prepareimport"
)

// Implementation of the prepare-import command.
func main() {
	var path = flag.String("path", viper.GetString("APP_DATA"), "Chemin d'accès au répertoire des batches")

	var batchKey = flag.String(
		"batch",
		"",
		"Clé du batch à importer au format AAMM (année + mois + suffixe optionnel)\n"+
			"Exemple: 1802_1",
	)
	var configFile = flag.String("configFile", "./batch.toml", "Chemin du fichier où est écrit la configuration\n"+
		"Exemple: ./batch.toml")

	flag.Parse()
	adminObject, err := prepare(*path, *batchKey)
	if err != nil {
		panic(err)
	}
	saveAdminObject(adminObject, *configFile)
}

func prepare(path, batchKey string) (base.AdminBatch, error) {
	validBatchKey, err := base.NewBatchKey(batchKey)
	if err != nil {
		return base.AdminBatch{}, errors.Wrap(err, "erreur lors de la création de la clé de batch")
	}
	adminObject, err := prepareimport.PrepareImport(path, validBatchKey)
	if _, ok := err.(prepareimport.UnsupportedFilesError); ok {
		return adminObject, err
	} else if err != nil {
		return base.AdminBatch{}, errors.Wrap(err, "erreur inattendue pendant la préparation de l'import : ")
	}
	return adminObject, nil
}

func saveAdminObject(toSave base.AdminBatch, configFile string) {
	err := prepareimport.SaveToFile(toSave, configFile)

	if err != nil {
		log.Fatal("Erreur inattendue pendant la sauvegarde de l'import : ", err)
	}
}
