package main

import (
	"log"
	"os"
	"os/exec"
	"testing"

	"github.com/spf13/viper"
)

const (
	mongoImage     = "mongo:4.2@sha256:1c2243a5e21884ffa532ca9d20c221b170d7b40774c235619f98e2f6eaec520a"
	mongoContainer = "dockertestmongodb"
	mongoURI       = "mongodb://localhost:27017" // TODO: switch to 27016
	mongoDatabase  = "signauxfaibles"
)

func TestMain(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping in short mode")
	}

	startMongoContainer(t) // may skip or fatal the test
	t.Cleanup(stopMongoContainer)

	viper.AddConfigPath(".")
	viper.SetConfigType("toml")
	viper.SetConfigName("config-sample") // => config will be loaded from ./config-sample.toml
	viper.Set("DB_DIAL", mongoURI)
	viper.Set("DB", mongoDatabase)

	os.Args = []string{"sfdata", "etablissements"}
	mainLogic() // n'appelle pas os.Exit() => le cleanup du test pourra avoir lieu
}

// startMongoContainer sets up a real MongoDB instance for testing purposes,
// using a Docker container. It makes the test fail on error.
func startMongoContainer(t *testing.T) {
	err := exec.Command("docker", "run", "--rm", "-d", "-p", "27017:27017", "--name", mongoContainer, mongoImage).Run()
	if err != nil {
		t.Fatalf("docker run: %v", err)
	}
}

func stopMongoContainer() {
	if err := exec.Command("docker", "stop", mongoContainer).Run(); err != nil {
		log.Println(err)
	}
}
