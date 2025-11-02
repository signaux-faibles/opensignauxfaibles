//go:build e2e

package main

import (
	"opensignauxfaibles/lib/engine"
	"path"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	SireneUl = engine.NewBatchFile("lib/parsing/sirene_ul/testData/sireneULTestData.csv")
	Debit    = engine.NewBatchFile("lib/parsing/urssaf/testData/debitTestData.csv")
	Effectif = engine.NewBatchFile("lib/parsing/effectif/testData/effectifTestData.csv")
)

func TestFilter(t *testing.T) {
	cleanDB := setupDBTest(t)
	defer cleanDB()

	t.Run("Import without filter should fail when filter tables are empty", func(t *testing.T) {
		// Create a batch with only Debit file, no filter provided
		batch := engine.AdminBatch{
			Key: "1902",
			Files: map[engine.ParserType][]engine.BatchFile{
				engine.Debit: {Debit},
			},
			Params: engine.AdminBatchParams{
				DateDebut: time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC),
				DateFin:   time.Date(2019, time.February, 1, 0, 0, 0, 0, time.UTC),
			},
		}
		writeBatchConfig(t, batch)

		// Run import without --no-filter flag
		exitCode := runCLI("sfdata", "import", "--batch", "1902", "--batch-config", path.Join(tmpDir, "batch.json"))

		// Should fail because filter tables are empty
		assert.NotEqual(t, 0, exitCode, "sfdata import should fail when no filter is provided and filter tables are empty")
	})
}

// func createImportTestBatch(t *testing.T) {

// 	batch := engine.AdminBatch{
// 		Key: "1910",
// 		Files: map[engine.ParserType][]engine.BatchFile{
// 			Dummy:              {},
// 			engine.Filter:      {},
// 			engine.Apconso:     {engine.NewBatchFile("lib/parsing/apconso/testData/apconsoTestData.csv")},
// 			engine.Apdemande:   {engine.NewBatchFile("lib/parsing/apdemande/testData/apdemandeTestData.csv")},
// 			engine.Sirene:      {engine.NewBatchFile("lib/parsing/sirene/testData/sireneTestData.csv")},
// 			engine.SireneUl:    {engine.NewBatchFile("lib/parsing/sirene_ul/testData/sireneULTestData.csv")},
// 			engine.AdminUrssaf: {engine.NewBatchFile("lib/parsing/urssaf/testData/comptesTestData.csv")},
// 			engine.Debit:       {engine.NewBatchFile("lib/parsing/urssaf/testData/debitTestData.csv")},
// 			engine.Ccsf:        {engine.NewBatchFile("lib/parsing/urssaf/testData/ccsfTestData.csv")},
// 			engine.Cotisation:  {engine.NewBatchFile("lib/parsing/urssaf/testData/cotisationTestData.csv")},
// 			engine.Delai:       {engine.NewBatchFile("lib/parsing/urssaf/testData/delaiTestData.csv")},
// 			engine.Effectif:    {engine.NewBatchFile("lib/parsing/effectif/testData/effectifTestData.csv")},
// 			engine.EffectifEnt: {engine.NewBatchFile("lib/parsing/effectif/testData/effectifEntTestData.csv")},
// 			engine.Procol:      {engine.NewBatchFile("lib/parsing/urssaf/testData/procolTestData.csv")},
// 		},
// 		Params: engine.AdminBatchParams{
// 			DateDebut: time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC),
// 			DateFin:   time.Date(2019, time.February, 1, 0, 0, 0, 0, time.UTC),
// 		},
// 	}
// 	writeBatchConfig(t, batch)
// }

// func verifyReports(t *testing.T) {
// 	t.Log("💎 Verifying exported reports...")

// 	conn, err := pgxpool.New(context.Background(), suite.PostgresURI)
// 	if err != nil {
// 		t.Errorf("Unable to connect to test database: %s", err)
// 	}

// 	table := engine.ReportTable
// 	query := fmt.Sprintf("SELECT * FROM %s ORDER BY parser", table)
// 	output := getTableContents(t, conn, query)
// 	goldenFile := fmt.Sprintf("test-import.sql.%s.golden.txt", table)
// 	tmpOutputFile := fmt.Sprintf("test-import.sql.%s.output.txt", table)
// 	compareWithGoldenFileOrUpdate(t, goldenFile, output, tmpOutputFile)
// }
