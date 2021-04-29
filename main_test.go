package main

import (
	"flag"
	"fmt"
	"log"
	"os/exec"
	"testing"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/signaux-faibles/opensignauxfaibles/lib/engine"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var _ = flag.Bool("update", false, "Update the expected test values in golden file") // please keep this line until https://github.com/kubernetes-sigs/service-catalog/issues/2319#issuecomment-425200065 is fixed

const (
	mongoImage     = "mongo:4.2@sha256:1c2243a5e21884ffa532ca9d20c221b170d7b40774c235619f98e2f6eaec520a"
	mongoContainer = "sf-mongodb"
	mongoPort      = 27016
	mongoDatabase  = "signauxfaibles"
)

func TestMain(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping in short mode")
	}

	startMongoContainer(t) // the test will fail in case of error
	t.Cleanup(stopMongoContainer)

	viper.AddConfigPath(".")
	viper.SetConfigType("toml")
	viper.SetConfigName("config-sample") // => config will be loaded from ./config-sample.toml

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
				"paydex":       []string{"/../lib/paydex/testData/paydexTestData.csv"},
			},
			"param": bson.M{
				"date_debut": time.Date(2019, 0, 1, 0, 0, 0, 0, time.UTC), // ISODate("2019-01-01T00:00:00.000+0000"),
				"date_fin":   time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC), // ISODate("2019-02-01T00:00:00.000+0000"),
			},
		})
		// Exécution des commandes et vérification qu'elles s'achèvent toutes avec un exit code nul    
    assert.Equal(t, 0, runCLI("sfdata", "check", "--batch=1910"))
    assert.Equal(t, 0, runCLI("sfdata", "import", "--batch=1910", "--no-filter"))
    assert.Equal(t, 0, runCLI("sfdata", "validate", "--collection=ImportedData"))
    assert.Equal(t, 0, runCLI("sfdata", "compact", "--since-batch=1910"))
    assert.Equal(t, 0, runCLI("sfdata", "public", "--until-batch=1910"))
    assert.Equal(t, 0, runCLI("sfdata", "public", "--until-batch=1910", "--key=012345678")) // pour couvrir PublicOne()
    assert.Equal(t, 0, runCLI("sfdata", "reduce", "--until-batch=1910"))
    assert.Equal(t, 0, runCLI("sfdata", "reduce", "--until-batch=1910", "--key=012345678")) // pour couvrir ReduceOne()
    assert.Equal(t, 0, runCLI("sfdata", "etablissements"))
    assert.Equal(t, 0, runCLI("sfdata", "entreprises"))
    assert.Equal(t, 0, runCLI("sfdata", "purgeNotCompacted", "--i-understand-what-im-doing")) // pour couvrir PurgeNotCompacted()
    assert.Equal(t, 0, runCLI("sfdata", "purge", "--since-batch=1910", "--i-understand-what-im-doing"))
    assert.Equal(t, 0, runCLI("sfdata", "purge", "--since-batch=1910", "--i-understand-what-im-doing", "--debug-for-key=012345678")) // pour couvrir PurgeBatchOne()
    assert.Equal(t, 0, runCLI("sfdata", "parseFile", "--parser=diane", "--file=lib/diane/testData/dianeTestData.txt"))
    assert.Equal(t, 2, runCLI("sfdata", "check"))                         // => "Erreur: paramètre `batch` obligatoire."
    assert.Equal(t, 3, runCLI("sfdata", "pruneEntities", "--batch=1910")) // => "Erreur: Ce batch ne spécifie pas de filtre"
	})

	t.Run("Les données importées sont validées par MongoDB", func(t *testing.T) {
		// Création d'une collection associée au schéma de validation de données JSON de ImportedData
		colName := "FakeImportedData"
		err := engine.CreateImportedDataCollection(db, colName)
		if assert.NoError(t, err) {
			coll := db.C(colName)
			// L'insertion d'un document valide ne doit pas provoquer d'erreur
			assert.NoError(t, coll.Insert(bson.M{
				"value": bson.M{
					"scope": "entreprise",
					"key":   "000000002",
					"batch": bson.M{
						"2002_2": bson.M{
							"paydex": bson.M{
								"afafafafafafaf": bson.M{"date_valeur": time.Now(), "nb_jours": 4},
							},
						},
					},
				},
			}))
			// L'insertion d'un document invalide doit provoquer une erreur
			assert.EqualError(t, coll.Insert(bson.M{"a": 1}), "Document failed validation")
		}
	})
}

func startMongoContainer(t *testing.T) {
	t.Log("Starting MongoDB in Docker container...")
	exec.Command("docker", "stop", mongoContainer).Run()
	exec.Command("docker", "rm", mongoContainer).Run()
	portMapping := fmt.Sprintf("%v:27017", mongoPort)
	err := exec.Command("docker", "run", "--rm", "-d", "-p", portMapping, "--name", mongoContainer, mongoImage).Run()
	if err != nil {
		t.Fatalf("docker run: %v", err)
	}
}

func stopMongoContainer() {
	if err := exec.Command("docker", "stop", mongoContainer).Run(); err != nil {
		log.Println(err) // affichage à titre informatif
	}
}
