//go:build e2e

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/globalsign/mgo"
	"github.com/stretchr/testify/assert"
)

type TestSuite struct {
	TmpDir         string
	GoldenFilesDir string
	MongoURI       string
}

var suite *TestSuite

const (
	mongoImage     = "mongo:4.2@sha256:1c2243a5e21884ffa532ca9d20c221b170d7b40774c235619f98e2f6eaec520a"
	mongoContainer = "sf-mongodb"
	mongoPort      = 27016
	mongoDatabase  = "signauxfaibles"
)

var update = flag.Bool("update", false, "Update the expected test values in golden file")

func TestMain(m *testing.M) {
	var err error

	suite, err = setupSuite()
	if err != nil {
		panic(err)
	}

	code := m.Run()

	teardownSuite()
	os.Exit(code)
}

func setupSuite() (*TestSuite, error) {
	log.Println("Setting up e2e tests...")

	log.Println("  Setting up MongoDB")

	exec.Command("docker", "stop", mongoContainer).Run()
	exec.Command("docker", "rm", mongoContainer).Run()

	startMongoContainer()
	mongoURI := fmt.Sprintf("mongodb://localhost:%v", mongoPort)

	log.Println("  Setting up configuration")

	tmpDir := filepath.Join("tests", "tmp-test-execution-files")

	// When running the command with cmd.Exec, the viper config is lost, so we
	// use environment variables instead
	os.Setenv("DB_DIAL", mongoURI)
	os.Setenv("DB", mongoDatabase)
	os.Setenv("APP_DATA", ".")
	os.Setenv("LOG_LEVEL", "error")
	os.Setenv("EXPORT_PATH", tmpDir)

	// When testing "runCli" directly, the config is not initialized, so we do
	// it here
	initConfig()

	log.Println("  Setting up temporary directory")

	os.RemoveAll(tmpDir)
	err := os.MkdirAll(tmpDir, 0755)

	if err != nil {
		panic(err)
	}

	// Give the DB time to start
	time.Sleep(2 * time.Second)

	return &TestSuite{
		TmpDir:         tmpDir,
		MongoURI:       mongoURI,
		GoldenFilesDir: filepath.Join("tests", "output-snapshots"),
	}, nil
}

func teardownSuite() {
	log.Println("Teardown db containers and temporary directory...")
	stopMongoContainer()
	deleteTempFolder()
}

// compareWithGoldenFileOrUpdate compares the "actualOutput" string with the
// contents of a golden file at "goldenPath", or updates the golden file with
// the output if the flag "--update" is provided.
//
// Any difference makes the test fail. A file is then written at "outputPath"
// for inspection.
func compareWithGoldenFileOrUpdate(t *testing.T, goldenFile, actualOutput, outputFile string) {

	goldenPath := filepath.Join(suite.GoldenFilesDir, goldenFile)
	outputPath := filepath.Join(suite.TmpDir, outputFile)

	if *update {
		err := updateGoldenFile(goldenPath, actualOutput)
		assert.NoError(t, err)

		t.Log("✅ Golden master file updated")

	} else {

		err := compareWithGoldenFile(goldenPath, actualOutput)

		if err != nil {
			// Write output to temp file for easy diffing
			t.Logf("❌Output different from golden file, writing output for inspection to %s", outputPath)
			_ = os.WriteFile(outputPath, []byte(actualOutput), 0644)
		} else {
			_ = os.Remove(outputPath)
		}

		assert.NoError(t, err)
	}
}

// compareWithGoldenFile compares the output with the golden file
func compareWithGoldenFile(goldenPath, actualOutput string) error {
	expected, err := os.ReadFile(goldenPath)
	if err != nil {
		return fmt.Errorf("failed to read golden file %s: %w", goldenPath, err)
	}

	if string(expected) != actualOutput {
		return fmt.Errorf("output doesn't match golden file %s.\nExpected:\n%s\nActual:\n%s",
			goldenPath, string(expected), actualOutput)
	}

	return nil
}

// updateGoldenFile writes the output to the golden file
func updateGoldenFile(goldenPath, actualOutput string) error {
	dir := filepath.Dir(goldenPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	return os.WriteFile(goldenPath, []byte(actualOutput), 0644)
}

func startMongoContainer() {
	exec.Command("docker", "stop", mongoContainer).Run()
	exec.Command("docker", "rm", mongoContainer).Run()
	portMapping := fmt.Sprintf("%v:27017", mongoPort)
	startMongoCommand := exec.Command("docker", "run", "--rm", "-d", "-p", portMapping, "--name", mongoContainer, mongoImage)
	err := startMongoCommand.Run()
	if err != nil {
		panic(err)
	}
}

func cleanDatabase(t *testing.T, db *mgo.Database) {
	t.Log("🧹 Cleaning database...")

	// Drop all collections
	collections := []string{"Admin", "Journal"}
	for _, collection := range collections {
		err := db.C(collection).DropCollection()
		if err != nil && err.Error() != "ns not found" {
			t.Logf("Warning: could not drop collection %s: %v", collection, err)
		}
	}
}

func stopMongoContainer() {
	if err := exec.Command("docker", "stop", mongoContainer).Run(); err != nil {
		log.Println(err) // affichage à titre informatif
	}
}

func deleteTempFolder() {
	os.RemoveAll(suite.TmpDir)
}
