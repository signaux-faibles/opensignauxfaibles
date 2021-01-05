package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"testing"

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
	viper.Set("DB_DIAL", fmt.Sprintf("mongodb://localhost:%v", mongoPort))
	viper.Set("DB", mongoDatabase)

	os.Args = []string{"sfdata", "etablissements"}
	mainLogic() // n'appelle pas os.Exit() => le cleanup du test pourra avoir lieu
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
