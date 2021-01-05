package main

import (
	"bytes"
	"fmt"
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

// Code from https://github.com/niilo/golib/blob/master/test/dockertest/docker.go

func stopMongoContainer() {
	if err := exec.Command("docker", "stop", mongoContainer).Run(); err != nil {
		log.Println(err)
	}
}

// startMongoContainer sets up a real MongoDB instance for testing purposes,
// using a Docker container. It makes the test fail on error.
func startMongoContainer(t *testing.T) {
	checkDockerImage(t, mongoImage)
	log.Println("Starting mongodb container...")
	if _, err := run("--rm", "-d", "-p", "27017:27017", "--name", mongoContainer, mongoImage); err != nil {
		t.Fatalf("docker run: %v", err)
	}
}

func haveImage(name string) (ok bool, err error) {
	out, err := exec.Command("docker", "images", "--no-trunc").Output()
	if err != nil {
		return
	}
	return bytes.Contains(out, []byte(name)), nil
}

// Pull retrieves the docker image with 'docker pull'.
func Pull(image string) error {
	out, err := exec.Command("docker", "pull", image).CombinedOutput()
	if err != nil {
		err = fmt.Errorf("%v: %s", err, out)
	}
	return err
}

// check all conditions to run a docker container based on image.
func checkDockerImage(t *testing.T, image string) {
	if _, err := exec.LookPath("docker"); err != nil {
		t.Error("'docker' command not found")
	}
	if ok, err := haveImage(image); !ok || err != nil {
		if err != nil {
			t.Errorf("Error running docker to check for %s: %v", image, err)
		}
		log.Printf("Pulling docker image %s ...", image)
		if err := Pull(image); err != nil {
			t.Errorf("Error pulling %s: %v", image, err)
		}
	}
}

func run(args ...string) (containerID string, err error) {
	log.Println(append([]string{"docker", "run"}, args...))
	cmd := exec.Command("docker", append([]string{"run"}, args...)...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	log.Println(stdout.String())
	log.Println(stderr.String())
	if err = cmd.Run(); err != nil {
		return "", err // fmt.Errorf("%v%v", stderr.String(), err)
	}
	containerID = "" // strings.TrimSpace(stdout.String())
	// if containerID == "" {
	// 	return "", errors.New("unexpected empty output from `docker run`")
	// }
	return
}
