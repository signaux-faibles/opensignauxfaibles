//go:build e2e

package main

import (
	"context"
	"opensignauxfaibles/lib/db"
	"opensignauxfaibles/lib/engine"
	"opensignauxfaibles/lib/filter"
	"opensignauxfaibles/lib/parsing/effectif"
	"opensignauxfaibles/lib/parsing/sirene"
	sireneul "opensignauxfaibles/lib/parsing/sirene_ul"
	"opensignauxfaibles/lib/registry"
	"opensignauxfaibles/lib/sinks"
	"sort"
	"testing"
	"time"

	"slices"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
)

var period, _ = time.Parse("2006-01-02", "2025-01-01")
var effectifContent = effectif.MakeEffectifEntCSV(
	[]time.Time{period},
	map[string][]int{"000000000": {5}, "111111111": {20}},
)

const (
	filterContent = `siren
111111111`
	sirenOut = "000000000"
	sirenIn  = "111111111"
)

var (
	SireneUl = engine.NewBatchFile("lib/parsing/sirene_ul/testData/sireneULTestData.csv")
	Debit    = engine.NewBatchFile("lib/parsing/urssaf/testData/debitTestData.csv")
	Effectif = engine.NewBatchFile("tests/testData/effectifEntTestData.csv")
)

// importWithDiscardData is a test helper that executes an import with default
// test behaviors for a given batch
// All data is discarded.
func importWithDiscardData(t *testing.T, batch engine.AdminBatch) error {
	t.Helper()
	return executeBatchImport(
		batch,
		[]engine.ParserType{}, // empty means all parsers
		registry.DefaultParsers,
		defaultFilterResolver(batch),
		&engine.DiscardSinkFactory{},
		&engine.DiscardReportSink{},
	)
}

func importWithDB(t *testing.T, batch engine.AdminBatch) error {
	t.Helper()
	return executeBatchImport(
		batch,
		[]engine.ParserType{}, // empty means all parsers
		registry.DefaultParsers,
		defaultFilterResolver(batch),
		sinks.NewPostgresSinkFactory(db.DB),
		engine.NewPostgresReportSink(db.DB),
	)
}

func TestImportFilter(t *testing.T) {
	cleanDB := setupDBTest(t)

	t.Run(`Import without filter should fail when
     1. filter tables are empty,
     2. no explicit filter is provided and
     3. no "effectif_ent" file is provided`, func(t *testing.T) {

		// Create a batch with only Debit file, no explicitely filter provided
		defer cleanDB()
		batch := engine.AdminBatch{
			Key: "1902",
			Files: map[engine.ParserType][]engine.BatchFile{
				engine.Debit: {Debit},
			},
		}

		err := importWithDiscardData(t, batch)

		assert.Error(t, err)
	})

	t.Run("Import with explicit filter file should succeed", func(t *testing.T) {
		// Create a batch with Debit file and an explicit filter file
		defer cleanDB()
		batch := engine.AdminBatch{
			Key: "1902",
			Files: map[engine.ParserType][]engine.BatchFile{
				engine.Debit:  {Debit},
				engine.Filter: {engine.NewMockBatchFile(filterContent)},
			},
		}

		// Run import with the filter
		err := importWithDiscardData(t, batch)

		assert.NoError(t, err)
	})

	t.Run("Import with \"effectif_ent\" file should succeed", func(t *testing.T) {
		defer cleanDB()

		// Create a batch with Debit file and an explicit filter file
		batch := engine.AdminBatch{
			Key: "1902",
			Files: map[engine.ParserType][]engine.BatchFile{
				engine.EffectifEnt: {engine.NewMockBatchFile(effectifContent)},
			},
		}

		err := importWithDiscardData(t, batch)

		assert.NoError(t, err, "should succeed to import when an \"effectif_ent\" file is provided")

		// Check that the filter has been properly updated
		filter, err := readFilter(batch)
		assert.NoError(t, err)
		assert.True(t, filter.ShouldSkip(sirenOut))
		assert.False(t, filter.ShouldSkip(sirenIn))
	})

	t.Run("When filter exists, new import with effectif updates the filter", func(t *testing.T) {
		defer cleanDB()

		// Create a batch with Debit file and an explicit filter file
		batch1 := engine.AdminBatch{
			Key: "1902",
			Files: map[engine.ParserType][]engine.BatchFile{
				engine.EffectifEnt: {engine.NewMockBatchFile(effectifContent)},
			},
		}

		period2 := period.AddDate(0, 1, 0)
		newEffectifContent := effectif.MakeEffectifEntCSV(
			[]time.Time{period, period2},
			map[string][]int{"000000000": {5, 20}, "111111111": {20, 20}},
		)

		batch2 := engine.AdminBatch{
			Key: "1903",
			Files: map[engine.ParserType][]engine.BatchFile{
				engine.EffectifEnt: {engine.NewMockBatchFile(newEffectifContent)},
			},
		}

		err := importWithDiscardData(t, batch1)
		assert.NoError(t, err, "should succeed to import when an effectif file is provided")

		err = importWithDiscardData(t, batch2)
		assert.NoError(t, err, "should succeed to import again when filter exists")

		// Check that the filter has been properly updated
		filter, err := readFilter(batch2)
		assert.NoError(t, err)
		// The new effectif should include former "sirenOut" inside the perimeter.
		assert.False(t, filter.ShouldSkip(sirenOut))
		assert.False(t, filter.ShouldSkip(sirenIn))
	})

	t.Run("Filter created in first import is saved to be reused in subsequent imports", func(t *testing.T) {
		defer cleanDB()

		// A first batch creates the filter
		batch1 := engine.AdminBatch{
			Key: "1902",
			Files: map[engine.ParserType][]engine.BatchFile{
				engine.EffectifEnt: {engine.NewMockBatchFile(effectifContent)},
			},
		}

		err := importWithDiscardData(t, batch1)
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
		filter, err := readFilter(batch2)
		assert.NoError(t, err)
		assert.True(t, filter.ShouldSkip("000000000"))
		assert.False(t, filter.ShouldSkip("111111111"))
	})
}

// defaultFilterResolver creates a StandardFilterResolver for testing (without --no-filter flag)
func defaultFilterResolver(batch engine.AdminBatch) engine.FilterResolver {
	return &filter.StandardFilterResolver{
		Reader: &filter.StandardReader{Batch: &batch, DB: db.DB},
		Writer: &filter.DBWriter{DB: db.DB},
	}
}

// readFilter reproduces the default filter reading strategy
func readFilter(batch engine.AdminBatch) (engine.SirenFilter, error) {
	reader := &filter.StandardReader{Batch: &batch, DB: db.DB}
	return reader.Read()
}

func TestCleanFilter(t *testing.T) {
	cleanDB := setupDBTest(t)
	t.Run("Test du périmètre selon les imports effectif_ent, sirene et sireneul", func(t *testing.T) {
		// Les données d'effectif sont filtrées, mais pas les données de Sirene
		// (clean_sirene et clean_sirene_ul), selon les besoins du front-end

		const (
			siren0 = "000000000"
			siren1 = "111111111"
			siren2 = "222222222"
			siren3 = "333333333"
			siren4 = "444444444"
		)
		effectifContent := effectif.MakeEffectifEntCSV(
			[]time.Time{period},
			map[string][]int{siren0: {5}, siren1: {20}, siren2: {20}, siren3: {20}},
		)

		effectifContentNewCompany := effectif.MakeEffectifEntCSV(
			[]time.Time{period},
			map[string][]int{siren0: {5}, siren1: {20}, siren2: {20}, siren3: {20}, siren4: {20}},
		)

		sireneUlContent := sireneul.MakeSireneULCSV(
			[]sireneul.SireneULEntry{
				{Siren: siren0, APE: "62.02A", CategorieJuridique: "5499"}, // private
				{Siren: siren1, APE: "62.01Z", CategorieJuridique: "4110"}, // public entity
				{Siren: siren2, APE: "62.02A", CategorieJuridique: "5499"}, // private
				{Siren: siren3, APE: "62.02A", CategorieJuridique: "5499"}, // private
				{Siren: siren4, APE: "62.02A", CategorieJuridique: "5499"}, // private
			},
		)

		sireneContent := sirene.MakeSireneCSV(
			[]sirene.SireneEntry{
				{Siret: siren0 + "00001", Siege: true, Etranger: false},
				{Siret: siren1 + "00001", Siege: true, Etranger: false},
				{Siret: siren2 + "00001", Siege: false, Etranger: false}, // not headquarters
				{Siret: siren2 + "00002", Siege: true, Etranger: true},   // headquarters abroad
				{Siret: siren3 + "00001", Siege: true, Etranger: false},
				{Siret: siren4 + "00001", Siege: true, Etranger: false},
			},
		)

		testCases := []struct {
			name              string
			batches           []engine.AdminBatch
			expectedPerimeter []string
		}{
			{
				name: "only effectif import, data can be imported but no clean_* layer yet (sirene and sirene_ul are missing)",
				batches: []engine.AdminBatch{
					{
						Key: "1902",
						Files: map[engine.ParserType][]engine.BatchFile{
							engine.EffectifEnt: {engine.NewMockBatchFile(effectifContent)},
						},
					},
				},
				expectedPerimeter: []string{},
			},
			{
				name: "effectif, sirene and sireneul simultanous import",
				batches: []engine.AdminBatch{
					{
						Key: "1902",
						Files: map[engine.ParserType][]engine.BatchFile{
							engine.EffectifEnt: {engine.NewMockBatchFile(effectifContent)},
							engine.SireneUl:    {engine.NewMockBatchFile(sireneUlContent)},
							engine.Sirene:      {engine.NewMockBatchFile(sireneContent)},
						},
					},
				},
				expectedPerimeter: []string{siren3},
			},
			{
				name: "effectif import, then sireneul and sirene import",
				batches: []engine.AdminBatch{
					{
						Key: "1902",
						Files: map[engine.ParserType][]engine.BatchFile{
							engine.EffectifEnt: {engine.NewMockBatchFile(effectifContent)},
						},
					},
					{
						Key: "1903",
						Files: map[engine.ParserType][]engine.BatchFile{
							engine.SireneUl: {engine.NewMockBatchFile(sireneUlContent)},
							engine.Sirene:   {engine.NewMockBatchFile(sireneContent)},
						},
					},
				},
				expectedPerimeter: []string{siren3},
			},
			{
				name: "new company appears in effectif : as sirene is not filtered, company is included right away",
				batches: []engine.AdminBatch{
					{
						Key: "1902",
						Files: map[engine.ParserType][]engine.BatchFile{
							engine.EffectifEnt: {engine.NewMockBatchFile(effectifContent)},
						},
					},
					{
						Key: "1903",
						Files: map[engine.ParserType][]engine.BatchFile{
							engine.SireneUl: {engine.NewMockBatchFile(sireneUlContent)},
							engine.Sirene:   {engine.NewMockBatchFile(sireneContent)},
						},
					},
					{
						Key: "1904",
						Files: map[engine.ParserType][]engine.BatchFile{
							// siren4 appears for the first time in the effectif file
							engine.EffectifEnt: {engine.NewMockBatchFile(effectifContentNewCompany)},
						},
					},
				},
				expectedPerimeter: []string{siren3, siren4},
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				defer cleanDB()

				for _, batch := range tc.batches {
					err := importWithDB(t, batch)
					assert.NoError(t, err)
				}

				rows, err := db.DB.Query(context.Background(), "SELECT siren FROM clean_filter")
				assert.NoError(t, err)
				actualPerimeter, err := pgx.CollectRows(rows, pgx.RowTo[string])
				assert.NoError(t, err)

				sort.Strings(actualPerimeter)
				expected := slices.Clone(tc.expectedPerimeter)
				sort.Strings(expected)
				assert.Equal(t, expected, actualPerimeter)
			})
		}
	})
}
