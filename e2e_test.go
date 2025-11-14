//go:build e2e

package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"log/slog"
	"opensignauxfaibles/lib/db"
	"opensignauxfaibles/lib/engine"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

type TestSuite struct {
	TmpDir         string
	GoldenFilesDir string
	PostgresURI    string
}

var suite *TestSuite

const (
	pgImage     = "postgres:17.5@sha256:30fa5c5e240b7b2ff2c31adf5a4c5ccacf22dae1d7760fea39061eb8af475854"
	pgContainer = "test-postgres"
	pgPort      = 5432
	pgDatabase  = "testdb"
	pgUser      = "testuser"
	pgPassword  = "testpass"
)

var tmpDir = filepath.Join("tests", "tmp-test-execution-files")
var update = flag.Bool("update", false, "Update the expected test values in golden file")

func TestMain(m *testing.M) {

	var err error

	suite, err = setupSuite()

	// Do not print logs in tests. Should be called after `setupSuite()` which
	// inits a default logger
	slog.SetDefault(slog.New(slog.NewTextHandler(io.Discard, nil)))

	if err != nil {
		panic(err)
	}

	code := m.Run()

	teardownSuite()
	os.Exit(code)
}

func setupSuite() (*TestSuite, error) {
	log.Println("Setting up e2e tests...")

	log.Println("  Setting up Postgresql")
	startPostgresContainer()

	postgresURI := fmt.Sprintf(
		"postgres://%s:%s@localhost:%v/%s?sslmode=disable",
		pgUser,
		pgPassword,
		pgPort,
		pgDatabase,
	)

	log.Println("  Setting up configuration")

	// When running the command with cmd.Exec, the viper config is lost, so we
	// use environment variables instead
	os.Setenv("APP_DATA", ".")
	os.Setenv("BATCH_CONFIG_FILE", path.Join(tmpDir, "batch.json"))
	os.Setenv("EXPORT_PATH", tmpDir)
	os.Setenv("POSTGRES_DB_URL", postgresURI)

	// Allow to set a different log level with LOG_LEVEL environment variable
	// This may break the tests, which expect an "error" log level,
	// but it makes it easier to debug a test that fails.
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "error"
	}
	os.Setenv("LOG_LEVEL", logLevel)

	// When testing "runCli" directly, the config is not initialized, so we do
	// it here
	initConfig()

	log.Println("  Setting up temporary directory")

	os.RemoveAll(tmpDir)
	err := os.MkdirAll(tmpDir, 0755)

	if err != nil {
		panic(err)
	}

	time.Sleep(3 * time.Second)
	return &TestSuite{
		TmpDir:         tmpDir,
		PostgresURI:    postgresURI,
		GoldenFilesDir: filepath.Join("tests", "output-snapshots"),
	}, nil
}

// setupDBTest initializes the database connection, runs migrations, and returns
// a cleanup function that truncates tables and refreshes materialized views.
func setupDBTest(t *testing.T) func() {
	t.Helper()

	// Initialize db.DB if not already done, with migrations
	if db.DB == nil {
		err := db.Init(true)
		assert.NoError(t, err)
	}

	postgresURI := suite.PostgresURI

	return func() {
		ctx := context.Background()
		conn, err := pgx.Connect(ctx, postgresURI)
		assert.NoError(t, err)
		defer conn.Close(ctx)

		tables := db.AllTables()
		materializedViews := db.AllMaterializedViews()

		// Teardown: truncate tables
		_, err = conn.Exec(
			context.Background(),
			fmt.Sprintf(
				"TRUNCATE %s CASCADE;",
				strings.Join(tables, ", "),
			),
		)
		assert.NoError(t, err)

		// Teardown: refresh views
		for _, view := range materializedViews {
			_, err = conn.Exec(
				context.Background(),
				fmt.Sprintf("REFRESH MATERIALIZED VIEW %s;", view),
			)
			assert.NoError(t, err)
		}
	}
}

func teardownSuite() {
	log.Println("Teardown db containers and temporary directory...")
	rmPostgresContainer()
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

func startPostgresContainer() {
	rmPostgresContainer()

	portMapping := fmt.Sprintf("%v:5432", pgPort)
	startPgCommand := exec.Command(
		"docker",
		"run",
		"--rm",
		"-d",
		"-p",
		portMapping,
		"--name",
		pgContainer,
		"-e",
		fmt.Sprintf("POSTGRES_DB=%s", pgDatabase),
		"-e",
		fmt.Sprintf("POSTGRES_USER=%s", pgUser),
		"-e",
		fmt.Sprintf("POSTGRES_PASSWORD=%s", pgPassword),
		pgImage,
	)

	out, err := startPgCommand.CombinedOutput()
	if err != nil {
		log.Fatalf("Failed to start Postgres container: %v\nOutput: %s", err, string(out))
	}
}

func rmPostgresContainer() {
	exec.Command("docker", "stop", pgContainer).Run()
	exec.Command("docker", "rm", "-f", pgContainer).Run()
}

func deleteTempFolder() {
	os.RemoveAll(suite.TmpDir)
}

func writeBatchConfig(t *testing.T, batch engine.AdminBatch) {
	t.Log("📝 Writing test batch config...")

	bytes, err := json.Marshal(batch)
	assert.NoError(t, err)
	err = os.WriteFile(viper.GetString("BATCH_CONFIG_FILE"), bytes, 0644)
	assert.NoError(t, err)
}
