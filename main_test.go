//go:build e2e

package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestPrincipal(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping in short mode")
	}

	viper.AddConfigPath(".")
	viper.SetConfigType("toml")
	viper.SetConfigName("config-sample") // => config will be loaded from ./config-sample.toml
	viper.Set("export.path", "tmp")

	mongoURI := fmt.Sprintf("mongodb://localhost:%v", mongoPort)
	viper.Set("DB_DIAL", mongoURI)
	viper.Set("DB", mongoDatabase)

	mongodb, err := mgo.Dial(mongoURI)
	if err != nil {
		t.Fatal(err)
	}
	// mongodb.SetSocketTimeout(72000 * time.Second)
	db := mongodb.DB(mongoDatabase)

	t.Run("Toutes les commandes (CLI) doivent fonctionner à condition qu'un batch valide soit fourni", func(t *testing.T) {
		// Insertion d'un Batch valide
		db.C("Admin").Insert(bson.M{
			"_id": bson.M{
				"key":  "1910",
				"type": "batch",
			},
			"files": bson.M{
				"admin_urssaf": []string{"/../lib/urssaf/testData/comptesTestData.csv"},
				"apconso":      []string{"/../lib/apconso/testData/apconsoTestData.csv"},
			},
			"param": bson.M{
				"date_debut": time.Date(2019, 0, 1, 0, 0, 0, 0, time.UTC), // ISODate("2019-01-01T00:00:00.000+0000"),
				"date_fin":   time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC), // ISODate("2019-02-01T00:00:00.000+0000"),
			},
		})
		// Exécution des commandes et vérification qu'elles s'achèvent toutes avec un exit code nul
		assert.Equal(t, 0, runCLI("sfdata", "check", "--batch=1910"))
		assert.Equal(t, 0, runCLI("sfdata", "import", "--batch=1910", "--no-filter"))
		assert.Equal(t, 0, runCLI("sfdata", "parseFile", "--parser=apconso", "--file=lib/apconso/testData/apconsoTestData.csv"))
		assert.Equal(t, 2, runCLI("sfdata", "check"))                  // => "Erreur: paramètre `batch` obligatoire."
		assert.Equal(t, 3, runCLI("sfdata", "import", "--batch=1910")) // => "Erreur: Ce batch ne spécifie pas de filtre"
	})
}
