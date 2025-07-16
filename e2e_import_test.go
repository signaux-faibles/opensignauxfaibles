//go:build e2e

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/stretchr/testify/assert"
)

func TestImportEndToEnd(t *testing.T) {

	mongodb, err := mgo.Dial(suite.MongoURI)
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
	var transformedEntries []map[string]any
	for _, doc := range journalEntries {
		transformed := make(map[string]any)

		if event, hasEvent := doc["event"]; hasEvent {
			eventMap := event.(bson.M)
			transformed["event"] = map[string]any{
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
	goldenFilePath := "test-import.journal.golden.txt"
	tmpOutputPath := "test-import.journal.output.txt"

	compareWithGoldenFileOrUpdate(t, goldenFilePath, output, tmpOutputPath)
}

func verifyExportedCSVFiles(t *testing.T) {
	t.Log("📁 Verifying exported CSV files...")

	exportDir := filepath.Join(suite.TmpDir, "1910")

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

		content, err := os.ReadFile(file)
		assert.NoError(t, err)

		baseName := filepath.Base(file)
		parserType := strings.TrimSuffix(baseName, ".csv")

		goldenFile := fmt.Sprintf("test-import.%s.golden.txt", parserType)
		tmpOutputFile := fmt.Sprintf("test-import.%s.output.txt", parserType)

		// Format output with output csv filename as header
		output := fmt.Sprintf("==== %s ====\n\n%s\n", baseName, string(content))

		compareWithGoldenFileOrUpdate(t, goldenFile, output, tmpOutputFile)
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
