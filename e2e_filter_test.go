//go:build e2e

package main

import (
	"opensignauxfaibles/lib/db"
	"opensignauxfaibles/lib/engine"
	"opensignauxfaibles/lib/filter"
	"opensignauxfaibles/lib/registry"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	effectifContent = `compte;siret;rais_soc;ape_ins;dep;eff202501;base;UR_EMET
000000000000000000;00000000000000;ENTREPRISE_A;1234Z;75;5;116;075077
111111111111111111;11111111111111;ENTREPRISE_B;5678Z;92;20;116;075077`
	filterContent = `siren
111111111`
	sirenOut = "000000000"
	sirenIn  = "111111111"
)

var (
	SireneUl = engine.NewBatchFile("lib/parsing/sirene_ul/testData/sireneULTestData.csv")
	Debit    = engine.NewBatchFile("lib/parsing/urssaf/testData/debitTestData.csv")
	Effectif = engine.NewBatchFile("lib/parsing/effectif/testData/effectifTestData.csv")
)

// importWithDefaults is a test helper that executes an import with default
// behaviors for a given batch
func importWithDefaults(t *testing.T, batch engine.AdminBatch) error {
	t.Helper()
	return executeBatchImport(
		batch,
		[]engine.ParserType{}, // empty means all parsers
		registry.DefaultParsers,
		defaultFilterReader(batch),
		defaultFilterWriter(),
		&engine.DiscardSinkFactory{},
		&engine.DiscardReportSink{},
	)
}

func TestFilter(t *testing.T) {
	cleanDB := setupDBTest(t)
	defer cleanDB()

	t.Run("Import without filter should fail when filter tables are empty, and no explicit filter is provided", func(t *testing.T) {
		// Create a batch with only Debit file, no explicitely filter provided
		batch := engine.AdminBatch{
			Key: "1902",
			Files: map[engine.ParserType][]engine.BatchFile{
				engine.Debit: {Debit},
			},
		}

		err := importWithDefaults(t, batch)

		assert.Error(t, err, "should fail to import when filter tables are empty and no explicit filter is provided")
	})

	t.Run("Import with explicit filter file should succeed", func(t *testing.T) {
		// Create a batch with Debit file and an explicit filter file
		batch := engine.AdminBatch{
			Key: "1902",
			Files: map[engine.ParserType][]engine.BatchFile{
				engine.Debit:  {Debit},
				engine.Filter: {engine.NewMockBatchFile(filterContent)},
			},
		}

		// Run import with the filter
		err := importWithDefaults(t, batch)

		assert.NoError(t, err, "should succeed to import when an explicit filter file is provided")
	})

	t.Run("Import with effectif file should succeed", func(t *testing.T) {

		// Create a batch with Debit file and an explicit filter file
		batch := engine.AdminBatch{
			Key: "1902",
			Files: map[engine.ParserType][]engine.BatchFile{
				engine.Effectif: {engine.NewMockBatchFile(effectifContent)},
			},
		}

		err := importWithDefaults(t, batch)

		assert.NoError(t, err, "should succeed to import when an effectif file is provided")

		// Check that the filter has been properly updated
		filter, err := defaultFilterReader(batch).Read()
		assert.NoError(t, err)
		assert.True(t, filter.ShouldSkip(sirenOut))
		assert.False(t, filter.ShouldSkip(sirenIn))
	})
	t.Run("When filter exists, new import with effectif updates the filter", func(t *testing.T) {

		// Create a batch with Debit file and an explicit filter file
		batch1 := engine.AdminBatch{
			Key: "1902",
			Files: map[engine.ParserType][]engine.BatchFile{
				engine.Effectif: {engine.NewMockBatchFile(effectifContent)},
			},
		}
		newEffectifContent := `compte;siret;rais_soc;ape_ins;dep;eff202501;eff202502;base;UR_EMET
000000000000000000;00000000000000;ENTREPRISE_A;1234Z;75;5;20;116;075077
111111111111111111;11111111111111;ENTREPRISE_B;5678Z;92;20;20;116;075077`

		batch2 := engine.AdminBatch{
			Key: "1903",
			Files: map[engine.ParserType][]engine.BatchFile{
				engine.Effectif: {engine.NewMockBatchFile(newEffectifContent)},
			},
		}

		err := importWithDefaults(t, batch1)
		assert.NoError(t, err, "should succeed to import when an effectif file is provided")

		err = importWithDefaults(t, batch2)
		assert.NoError(t, err, "should succeed to import again when filter exists")

		// Check that the filter has been properly updated
		filter, err := defaultFilterReader(batch2).Read()
		assert.NoError(t, err)
		// The new effectif should include former "sirenOut" inside the perimeter.
		assert.False(t, filter.ShouldSkip(sirenOut))
		assert.False(t, filter.ShouldSkip(sirenIn))
	})

	t.Run("Filter created in first import is saved to be reused in subsequent imports", func(t *testing.T) {

		// A first batch creates the filter
		batch1 := engine.AdminBatch{
			Key: "1902",
			Files: map[engine.ParserType][]engine.BatchFile{
				engine.Effectif: {engine.NewMockBatchFile(effectifContent)},
			},
		}

		err := importWithDefaults(t, batch1)
		assert.NoError(t, err) // tested in test above already

		// A second batch has no effectif or filter file, but should reuse
		// existing filter in DB

		batch2 := engine.AdminBatch{
			Key: "1903",
			Files: map[engine.ParserType][]engine.BatchFile{
				engine.Debit: {Debit},
			},
		}
		assert.NoError(t, err, "should succeed to import when a filter has been created in DB")

		// Check that the filter has been left unchanged
		filter, err := defaultFilterReader(batch2).Read()
		assert.NoError(t, err)
		assert.True(t, filter.ShouldSkip("000000000"))
		assert.False(t, filter.ShouldSkip("111111111"))
	})
}

// Default FilterReader (without the --no-filter option)
func defaultFilterReader(batch engine.AdminBatch) engine.FilterReader {
	return &filter.Reader{Batch: &batch, DB: db.DB}
}

// Default FilterWriter (without the --no-filter and --dry-run options)
func defaultFilterWriter() engine.FilterWriter {
	return &filter.DBWriter{DB: db.DB}
}
