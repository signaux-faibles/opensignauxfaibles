//go:build e2e

package main

import (
	"opensignauxfaibles/lib/base"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/globalsign/mgo"
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
		})
	}
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
