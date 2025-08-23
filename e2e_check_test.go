//go:build e2e

package main

import (
	"encoding/json"
	"opensignauxfaibles/lib/base"
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
	createCheckTestBatch(t)

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

			cmd := exec.Command("./sfdata", tc.args...)
			cmd.Env = os.Environ()

			output, err := cmd.CombinedOutput()

			assert.NoError(t, err, "%s failed: %s", tc.name, string(output))
			compareWithGoldenFileOrUpdate(t, tc.goldenFile, string(output), tc.tmpFile)
		})
	}

	verifyCheckJournalReports(t, db)
}

func createCheckTestBatch(t *testing.T) {

	batch := base.AdminBatch{
		ID: base.AdminID{
			Key:  "1910",
			Type: "batch",
		},
		Files: base.BatchFiles{
			"admin_urssaf": {base.NewBatchFile("./lib/urssaf/testData/comptesTestData.csv")},
			"debit":        {base.NewBatchFile("./lib/urssaf/testData/debitCorrompuTestData.csv")},
		},
		Params: base.AdminBatchParams{
			DateDebut: time.Date(2001, time.January, 1, 0, 0, 0, 0, time.UTC),
			DateFin:   time.Date(2019, time.February, 1, 0, 0, 0, 0, time.UTC),
		},
	}
	writeBatchConfig(t, batch)

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

	output, err := json.MarshalIndent(transformedEntries, "", "  ")
	assert.NoError(t, err)

	goldenFilePath := "test-check.journal.golden.txt"
	tmpOutputPath := "test-check.journal.output.txt"

	compareWithGoldenFileOrUpdate(t, goldenFilePath, string(output), tmpOutputPath)
}
