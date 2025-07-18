//go:build e2e

package main

import (
	"encoding/json"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/stretchr/testify/assert"
)

func TestCheckEndToEnd(t *testing.T) {
	mongodb, err := mgo.Dial(suite.MongoURI)
	assert.NoError(t, err)
	defer mongodb.Close()

	db := mongodb.DB(mongoDatabase)

	cleanDatabase(t, db)
	insertCheckTestBatch(t, db)

	testCases := []struct {
		name       string
		args       []string
		goldenFile string
		tmpFile    string
	}{
		{
			"sfdata check --batch=1910 --parsers=debit",
			[]string{"check", "--batch=1910", "--parsers=debit"},
			"test-check.1.golden.txt",
			"test-check.1.output.txt",
		},
		{
			"sfdata check --batch=1910",
			[]string{"check", "--batch=1910"},
			"test-check.2.golden.txt",
			"test-check.2.output.txt",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Log("💎 Parsing data...")

			// Run check command with single parser (debit)
			cmd := exec.Command("./sfdata", tc.args...)
			cmd.Env = os.Environ()

			output, err := cmd.CombinedOutput()

			assert.NoError(t, err, "%s failed: %s", tc.name, string(output))
			compareWithGoldenFileOrUpdate(t, tc.goldenFile, string(output), tc.tmpFile)

		})
	}

	verifyCheckJournalReports(t, db)
}

func insertCheckTestBatch(t *testing.T, db *mgo.Database) {
	t.Log("📝 Inserting test data...")

	batch := bson.M{
		"_id": bson.M{
			"key":  "1910",
			"type": "batch",
		},
		"files": bson.M{
			"admin_urssaf": []string{"./lib/urssaf/testData/comptesTestData.csv"},
			"debit":        []string{"./lib/urssaf/testData/debitCorrompuTestData.csv"},
		},
		"param": bson.M{
			"date_debut": time.Date(2001, time.January, 1, 0, 0, 0, 0, time.UTC),
			"date_fin":   time.Date(2019, time.February, 1, 0, 0, 0, 0, time.UTC),
		},
	}

	err := db.C("Admin").Insert(batch)
	assert.NoError(t, err)
}

func verifyCheckJournalReports(t *testing.T, db *mgo.Database) {
	t.Log("💎 Verifying Journal reports...")

	// Query Journal collection with sorting
	var journalEntries []bson.M
	err := db.C("Journal").Find(nil).Sort("-reportType", "parserCode").All(&journalEntries)
	assert.NoError(t, err)

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
		transformed["hasDate"] = doc["date"] != nil
		transformed["hasStartDate"] = doc["startDate"] != nil

		transformedEntries = append(transformedEntries, transformed)
	}

	// Format output to match the original script
	output, err := json.MarshalIndent(transformedEntries, "", "  ")
	assert.NoError(t, err)

	// Handle golden file comparison/update
	goldenFilePath := "test-check.journal.golden.txt"
	tmpOutputPath := "test-check.journal.output.txt"

	compareWithGoldenFileOrUpdate(t, goldenFilePath, string(output), tmpOutputPath)
}
