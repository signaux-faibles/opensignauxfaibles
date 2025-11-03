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

var (
	SireneUl = engine.NewBatchFile("lib/parsing/sirene_ul/testData/sireneULTestData.csv")
	Debit    = engine.NewBatchFile("lib/parsing/urssaf/testData/debitTestData.csv")
	Effectif = engine.NewBatchFile("lib/parsing/effectif/testData/effectifTestData.csv")
)

func TestFilter(t *testing.T) {
	cleanDB := setupDBTest(t)
	defer cleanDB()

	t.Run("Import without filter should fail when filter tables are empty, and no explicit filter is provided", func(t *testing.T) {
		// Create a batch with only Debit file, no filter provided
		batch := engine.AdminBatch{
			Key: "1902",
			Files: map[engine.ParserType][]engine.BatchFile{
				engine.Debit: {Debit},
			},
		}

		// Try to read filter - should fail because tables are empty
		filterProvider := &filter.Reader{Batch: &batch, DB: db.DB}
		_, err := filterProvider.Read()

		assert.Error(t, err, "should fail to read filter when filter tables are empty and no explicit filter is provided")
	})

	t.Run("Import with explicit filter file should succeed", func(t *testing.T) {
		// Create a mock filter file with inline data
		filterData := "siren\n111111111\n222222222"
		mockFilter := engine.NewMockBatchFile(filterData)

		// Create a batch with Debit file and an explicit filter file
		batch := engine.AdminBatch{
			Key: "1902",
			Files: map[engine.ParserType][]engine.BatchFile{
				engine.Debit:  {Debit},
				engine.Filter: {mockFilter},
			},
		}

		// Get filter from the explicit filter file
		filterProvider := &filter.Reader{Batch: &batch, DB: db.DB}
		sirenFilter, err := filterProvider.Read()
		assert.NoError(t, err, "should succeed to read filter from explicit filter file")

		// Run import with the filter
		err = engine.ImportBatch(
			batch,
			[]engine.ParserType{}, // empty means all parsers
			registry.DefaultParsers,
			sirenFilter,
			&engine.DiscardSinkFactory{},
			&engine.DiscardReportSink{},
		)

		assert.NoError(t, err, "should succeed to import when an explicit filter file is provided")
	})

	t.Run("Import with explicit filter file should succeed", func(t *testing.T) {
		// Create a batch with Debit file and an explicit filter file
		batch := engine.AdminBatch{
			Key: "1902",
			Files: map[engine.ParserType][]engine.BatchFile{
				engine.Effectif: {Effectif},
			},
		}

	})
}
