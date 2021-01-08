package main

import (
	"fmt"
	"log"
	"os/exec"
	"testing"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/spf13/viper"
)

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
	db.C("Admin").Insert(bson.M{
		"_id": bson.M{
			"key":  "1910",
			"type": "batch",
		},
		"files": bson.M{
			"paydex": []string{"/../lib/paydex/testData/paydexTestData.csv"},
		},
		"param": bson.M{
			"date_debut": time.Date(2019, 0, 1, 0, 0, 0, 0, time.UTC), // ISODate("2019-01-01T00:00:00.000+0000"),
			"date_fin":   time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC), // ISODate("2019-02-01T00:00:00.000+0000"),
		},
	})

	runCLI("sfdata", "import", "--batch=1910", "--no-filter")
	runCLI("sfdata", "compact", "--since-batch=1910")
	runCLI("sfdata", "public", "--until-batch=1910")
	runCLI("sfdata", "reduce", "--until-batch=1910")
	runCLI("sfdata", "etablissements")
	runCLI("sfdata", "entreprises")

	// var firstBatch base.AdminBatch
	// db.C("Admin").Find(bson.M{}).One(&firstBatch)
	// log.Println(firstBatch)
}

func startMongoContainer(t *testing.T) {
	t.Log("Starting MongoDB in Docker container...")
	portMapping := fmt.Sprintf("%v:27017", mongoPort)
	err := exec.Command("docker", "run", "--rm", "-d", "-p", portMapping, "--name", mongoContainer, mongoImage).Run()
	if err != nil {
		t.Fatalf("docker run: %v", err)
	}
}

func stopMongoContainer() {
	if err := exec.Command("docker", "stop", mongoContainer).Run(); err != nil {
		log.Println(err)
	}
}
