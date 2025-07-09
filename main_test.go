package main

import (
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
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

const (
	postgresImage     = "postgres:17@sha256:fe3f571d128e8efadcd8b2fde0e2b73ebab6dbec33f6bfe69d98c682c7d8f7bd"
	postgresContainer = "sf-postgres"
	postgresPort      = 5432
	postgresDatabase  = "testdb"
	postgresUser      = "testuser"
	postgresPassword  = "testpass"
)

func TestPrincipal(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping in short mode")
	}

	startMongoContainer(t) // the test will fail in case of error
	startPostgresContainer(t)
	t.Cleanup(stopMongoContainer)
	t.Cleanup(stopPostgresContainer)
	t.Cleanup(deleteTempFolder)

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
		assert.Equal(t, 4, runCLI("sfdata", "import", "--batch=1910")) // => "Erreur: Ce batch ne spécifie pas de filtre"
	})
}

func startMongoContainer(t *testing.T) {
	t.Log("Starting MongoDB in Docker container...")
	exec.Command("docker", "stop", mongoContainer).Run()
	exec.Command("docker", "rm", mongoContainer).Run()
	portMapping := fmt.Sprintf("%v:27017", mongoPort)
	startMongoCommand := exec.Command("docker", "run", "--rm", "-d", "-p", portMapping, "--name", mongoContainer, mongoImage)
	slog.Info("starting mongo", slog.Any("command", startMongoCommand.Args))
	err := startMongoCommand.Run()
	if err != nil {
		t.Fatalf("docker run: %v", err)
	}
}

func stopMongoContainer() {
	if err := exec.Command("docker", "stop", mongoContainer).Run(); err != nil {
		log.Println(err) // affichage à titre informatif
	}
}

func startPostgresContainer(t *testing.T) {
	t.Log("Starting PostgreSQL in Docker container...")
	startPostgresCommand := exec.Command("bash", "./tests/helpers/postgres-container.sh", "start")
	slog.Info("starting mongo", slog.Any("command", startPostgresCommand.Args))
	err := startPostgresCommand.Run()
	if err != nil {
		t.Fatalf("postgresql docker run: %v", err)
	}
}

func stopPostgresContainer() {
	if err := exec.Command("docker", "stop", postgresContainer).Run(); err != nil {
		log.Println(err) // affichage à titre informatif
	}
}

func deleteTempFolder() {
	os.RemoveAll(viper.GetString("export.path"))
}
