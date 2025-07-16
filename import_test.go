package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestImportEndToEnd(t *testing.T) {

	// Setup
	setupImportTest(t)
	t.Cleanup(cleanupImportTest)

	// Configure viper
	viper.AddConfigPath(".")
	viper.SetConfigType("toml")
	viper.SetConfigName("config-sample") // => config will be loaded from ./config-sample.toml
	viper.Set("export.path", "tests/tmp-test-execution-files")

	mongoURI := fmt.Sprintf("mongodb://localhost:%v", mongoPort)
	viper.Set("DB_DIAL", mongoURI)
	viper.Set("DB", mongoDatabase)

	mongodb, err := mgo.Dial(mongoURI)
	assert.NoError(t, err)
	defer mongodb.Close()
	db := mongodb.DB(mongoDatabase)

	t.Run("Insert test data and run import", func(t *testing.T) {
		insertImportTestBatch(t, db)

		exitCode := runCLI("sfdata", "import", "--batch", "1910", "--no-filter")
		assert.Equal(t, 0, exitCode, "sfdata import should succeed")
	})

	t.Run("Verify Journal reports", func(t *testing.T) {
		verifyJournalReports(t, db)
	})

	t.Run("Verify exported CSV files", func(t *testing.T) {
		verifyExportedCSVFiles(t)
	})
}

func setupImportTest(t *testing.T) {
	t.Log("Setting up import test...")

	// Stop any existing container
	exec.Command("docker", "stop", mongoContainer).Run()
	exec.Command("docker", "rm", mongoContainer).Run()

	startMongoContainer(t)

	// Create temp directory
	tmpDir := "tests/tmp-test-execution-files"
	os.RemoveAll(tmpDir)
	err := os.MkdirAll(tmpDir, 0755)
	assert.NoError(t, err)

	// Give MongoDB time to start
	time.Sleep(2 * time.Second)
}

func cleanupImportTest() {
	// Stop MongoDB container
	stopMongoContainer()

	// Clean up temp directory
	os.RemoveAll("tests/tmp-test-execution-files")

	// Clean up environment variables
	os.Unsetenv("TMP_DIR")
	os.Unsetenv("MONGODB_PORT")
}

func insertImportTestBatch(t *testing.T, db *mgo.Database) {
	t.Log("📝 Inserting test data...")

	batch := bson.M{
		"_id": bson.M{
			"key":  "1910",
			"type": "batch",
		},
		"files": bson.M{
			"dummy":        []string{},
			"filter":       []string{},
			"apconso":      []string{"/../lib/apconso/testData/apconsoTestData.csv"},
			"apdemande":    []string{"/../lib/apdemande/testData/apdemandeTestData.csv"},
			"sirene":       []string{"/../lib/sirene/testData/sireneTestData.csv"},
			"sirene_ul":    []string{"/../lib/sirene_ul/testData/sireneULTestData.csv"},
			"admin_urssaf": []string{"/../lib/urssaf/testData/comptesTestData.csv"},
			"debit":        []string{"/../lib/urssaf/testData/debitTestData.csv"},
			"ccsf":         []string{"/../lib/urssaf/testData/ccsfTestData.csv"},
			"cotisation":   []string{"/../lib/urssaf/testData/cotisationTestData.csv"},
			"delai":        []string{"/../lib/urssaf/testData/delaiTestData.csv"},
			"effectif":     []string{"/../lib/urssaf/testData/effectifTestData.csv"},
			"effectif_ent": []string{"/../lib/urssaf/testData/effectifEntTestData.csv"},
			"procol":       []string{"/../lib/urssaf/testData/procolTestData.csv"},
		},
		"param": bson.M{
			"date_debut": time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC),
			"date_fin":   time.Date(2019, time.February, 1, 0, 0, 0, 0, time.UTC),
		},
	}

	err := db.C("Admin").Insert(batch)
	assert.NoError(t, err)
}

func verifyJournalReports(t *testing.T, db *mgo.Database) {
	t.Log("💎 Verifying Journal reports...")

	// Query Journal collection
	var journalEntries []bson.M
	err := db.C("Journal").Find(nil).Sort("reportType", "parserCode").All(&journalEntries)
	assert.NoError(t, err)

	// Transform the data similar to the MongoDB query in the bash script
	var transformedEntries []map[string]interface{}
	for _, doc := range journalEntries {
		transformed := make(map[string]interface{})

		if event, hasEvent := doc["event"]; hasEvent {
			eventMap := event.(bson.M)
			transformed["event"] = map[string]interface{}{
				"headRejected": eventMap["headRejected"],
				"headFatal":    eventMap["headFatal"],
				"linesSkipped": eventMap["linesSkipped"],
				"summary":      eventMap["summary"],
				"batchKey":     eventMap["batchKey"],
			}
		}

		transformed["reportType"] = doc["reportType"]
		transformed["parserCode"] = doc["parserCode"]
		transformed["hasCommitHash"] = doc["commitHash"] != nil
		transformed["hasDate"] = doc["date"] != nil
		transformed["hasStartDate"] = doc["startDate"] != nil

		transformedEntries = append(transformedEntries, transformed)
	}

	// Convert to string for comparison
	output := formatJournalOutput(transformedEntries)

	// Handle golden file comparison/update
	goldenFilePath := "tests/output-snapshots/test-import.journal.golden.txt"
	if *update {
		err := updateGoldenFile(goldenFilePath, output)
		assert.NoError(t, err)

		t.Log("✅ Golden master file updated")
	} else {

		err := compareWithGoldenFile(t, goldenFilePath, output)
		outputFilePath := filepath.Join(tmpDir, "test-import.journal.output.txt")

		if err != nil {
			// Write output to temp file for easy diffing
			_ = os.WriteFile(outputFilePath, []byte(output), 0644)
		} else {
			_ = os.Remove(outputFilePath)
		}
		assert.NoError(t, err)
	}
}

func verifyExportedCSVFiles(t *testing.T) {
	t.Log("📁 Verifying exported CSV files...")

	exportDir := "tests/tmp-test-execution-files/1910"

	// Find all CSV files in the export directory
	files, err := filepath.Glob(filepath.Join(exportDir, "*"))
	assert.NoError(t, err)

	// Sort files for consistent ordering
	sort.Strings(files)

	for _, file := range files {
		if !strings.HasSuffix(file, ".csv") {
			continue
		}

		fileInfo, err := os.Stat(file)
		assert.NoError(t, err)

		if fileInfo.IsDir() {
			continue
		}

		// Read file content
		content, err := os.ReadFile(file)
		assert.NoError(t, err)

		// Create golden file name based on CSV file name
		baseName := filepath.Base(file)
		goldenFileName := fmt.Sprintf("test-import.%s.golden.txt",
			strings.TrimSuffix(baseName, ".csv"))
		goldenFilePath := filepath.Join("tests/output-snapshots", goldenFileName)

		// Format output with header
		output := fmt.Sprintf("==== %s ====\n\n%s\n", baseName, string(content))

		// Compare with golden file
		compareWithGoldenFile(t, goldenFilePath, output)
		if *update {
			err := updateGoldenFile(goldenFilePath, output)
			assert.NoError(t, err)

			t.Log("✅ Golden master file updated")
		} else {

			err := compareWithGoldenFile(t, goldenFilePath, output)
			outputFilePath := filepath.Join(tmpDir, "test-import.journal.output.txt")

			if err != nil {
				// Write output to temp file for easy diffing
				_ = os.WriteFile(outputFilePath, []byte(output), 0644)
			} else {
				_ = os.Remove(outputFilePath)
			}
			assert.NoError(t, err)
		}
	}
}

func formatJournalOutput(entries []map[string]any) string {
	var output strings.Builder
	output.WriteString("// Reports from db.Journal:\n")

	for _, entry := range entries {
		// Convert to JSON-like format (simplified)
		output.WriteString(fmt.Sprintf("%+v\n", entry))
	}

	return output.String()
}
