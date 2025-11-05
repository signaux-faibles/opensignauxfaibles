//go:build e2e

package main

import (
	"context"
	"opensignauxfaibles/lib/db"
	"opensignauxfaibles/lib/engine"
	"opensignauxfaibles/lib/filter"
	"opensignauxfaibles/lib/registry"
	"opensignauxfaibles/lib/sinks"
	"testing"

	"github.com/jackc/pgx/v5"
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

// importWithDiscardData is a test helper that executes an import with default
// behaviors for a given batch
func importWithDiscardData(t *testing.T, batch engine.AdminBatch) error {
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

func importWithDB(t *testing.T, batch engine.AdminBatch) error {
	t.Helper()
	return executeBatchImport(
		batch,
		[]engine.ParserType{}, // empty means all parsers
		registry.DefaultParsers,
		defaultFilterReader(batch),
		defaultFilterWriter(),
		sinks.NewPostgresSinkFactory(db.DB),
		engine.NewPostgresReportSink(db.DB),
	)
}

func TestImportFilter(t *testing.T) {
	cleanDB := setupDBTest(t)

	t.Run("Import without filter should fail when filter tables are empty, and no explicit filter is provided", func(t *testing.T) {
		// Create a batch with only Debit file, no explicitely filter provided
		defer cleanDB()
		batch := engine.AdminBatch{
			Key: "1902",
			Files: map[engine.ParserType][]engine.BatchFile{
				engine.Debit: {Debit},
			},
		}

		err := importWithDiscardData(t, batch)

		assert.Error(t, err, "should fail to import when filter tables are empty and no explicit filter is provided")
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

		assert.NoError(t, err, "should succeed to import when an explicit filter file is provided")
	})

	t.Run("Import with effectif file should succeed", func(t *testing.T) {
		defer cleanDB()

		// Create a batch with Debit file and an explicit filter file
		batch := engine.AdminBatch{
			Key: "1902",
			Files: map[engine.ParserType][]engine.BatchFile{
				engine.Effectif: {engine.NewMockBatchFile(effectifContent)},
			},
		}

		err := importWithDiscardData(t, batch)

		assert.NoError(t, err, "should succeed to import when an effectif file is provided")

		// Check that the filter has been properly updated
		filter, err := defaultFilterReader(batch).Read()
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

		err := importWithDiscardData(t, batch1)
		assert.NoError(t, err, "should succeed to import when an effectif file is provided")

		err = importWithDiscardData(t, batch2)
		assert.NoError(t, err, "should succeed to import again when filter exists")

		// Check that the filter has been properly updated
		filter, err := defaultFilterReader(batch2).Read()
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
				engine.Effectif: {engine.NewMockBatchFile(effectifContent)},
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

func TestCleanFilter(t *testing.T) {
	cleanDB := setupDBTest(t)
	t.Run("Après l'import d'effectif et de sirene_ul, les vues préfixées par \"clean_\" sont correctement filtrées", func(t *testing.T) {
		defer cleanDB()

		effectifContentTwoIns := effectifContent + "\n" +
			"222222222222222222;22222222222222;ENTREPRISE_C;5678Z;95;20;116;075077`"

		sireneUlContent := `siren,statutDiffusionUniteLegale,unitePurgeeUniteLegale,dateCreationUniteLegale,sigleUniteLegale,sexeUniteLegale,prenom1UniteLegale,prenom2UniteLegale,prenom3UniteLegale,prenom4UniteLegale,prenomUsuelUniteLegale,pseudonymeUniteLegale,identifiantAssociationUniteLegale,trancheEffectifsUniteLegale,anneeEffectifsUniteLegale,dateDernierTraitementUniteLegale,nombrePeriodesUniteLegale,categorieEntreprise,anneeCategorieEntreprise,dateDebut,etatAdministratifUniteLegale,nomUniteLegale,nomUsageUniteLegale,denominationUniteLegale,denominationUsuelle1UniteLegale,denominationUsuelle2UniteLegale,denominationUsuelle3UniteLegale,categorieJuridiqueUniteLegale,activitePrincipaleUniteLegale,nomenclatureActivitePrincipaleUniteLegale,nicSiegeUniteLegale,economieSocialeSolidaireUniteLegale,caractereEmployeurUniteLegale
111111111,O,,2000-01-01,,,,,,,,,,,,2020-01-01T00:00:00,1,PME,2020,2000-01-01,A,,,ENTREPRISE PUBLIQUE,,,,4110,62.01Z,NAFRev2,00001,,O
222222222,O,,2010-01-01,,,,,,,,,,,,2020-01-01T00:00:00,1,PME,2020,2010-01-01,A,,,ENTREPRISE PRIVEE,,,,5499,62.02A,NAFRev2,00001,,O`

		batch := engine.AdminBatch{
			Key: "1903",
			Files: map[engine.ParserType][]engine.BatchFile{
				engine.Effectif: {engine.NewMockBatchFile(effectifContentTwoIns)},
				engine.SireneUl: {engine.NewMockBatchFile(sireneUlContent)},
			},
		}
		err := importWithDB(t, batch)
		assert.NoError(t, err)

		// Vérifier que 222222222 (entreprise privée) est présent dans clean_effectif
		rows, err := db.DB.Query(context.Background(), "SELECT siret FROM clean_effectif WHERE LEFT(siret, 9) = '222222222'")
		assert.NoError(t, err)
		siretsFor222, err := pgx.CollectRows(rows, pgx.RowTo[string])
		assert.NoError(t, err)
		assert.Greater(t, len(siretsFor222), 0, "L'entreprise 222222222 (privée) devrait être présente dans clean_effectif")

		// Vérifier que 111111111 (organisation publique) n'est PAS présent dans clean_effectif
		rows, err = db.DB.Query(context.Background(), "SELECT siret FROM clean_effectif WHERE LEFT(siret, 9) = '111111111'")
		assert.NoError(t, err)
		siretsFor111, err := pgx.CollectRows(rows, pgx.RowTo[string])
		assert.NoError(t, err)
		assert.Equal(t, 0, len(siretsFor111), "L'entreprise 111111111 (publique) ne devrait PAS être présente dans clean_effectif")
	})
}
