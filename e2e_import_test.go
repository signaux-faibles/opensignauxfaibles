//go:build e2e

package main

import (
	"context"
	"fmt"
	"opensignauxfaibles/lib/base"
	"opensignauxfaibles/lib/engine"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
)

func TestImportEndToEnd(t *testing.T) {
	cleanDB := setupDBTest(t)
	defer cleanDB()

	createImportTestBatch(t)
	t.Run("Create batch and run import", func(t *testing.T) {

		exitCode := runCLI("sfdata", "import", "--batch", "1910", "--no-filter", "--batch-config", path.Join(tmpDir, "batch.json"))
		assert.Equal(t, 0, exitCode, "sfdata import should succeed")
	})

	t.Run("Verify exported Postgres reports", func(t *testing.T) {
		verifyReports(t)
	})

	t.Run("Verify exported CSV files", func(t *testing.T) {
		verifyExportedCSVFiles(t)
	})

	t.Run("Verify exported Postgres files", func(t *testing.T) {
		verifyPostgresExport(t)
	})

	t.Run("Run with --dry-run", func(t *testing.T) {
		exitCode := runCLI("sfdata", "import", "--dry-run", "--batch", "1910", "--no-filter", "--batch-config", path.Join(tmpDir, "batch.json"))
		assert.Equal(t, 0, exitCode, "sfdata import should succeed with --dry-run")
	})
}

const Dummy base.ParserType = "dummy"

func createImportTestBatch(t *testing.T) {

	batch := base.AdminBatch{
		Key: "1910",
		Files: map[base.ParserType][]base.BatchFile{
			Dummy:            {},
			base.Filter:      {},
			base.Apconso:     {base.NewBatchFile("lib/apconso/testData/apconsoTestData.csv")},
			base.Apdemande:   {base.NewBatchFile("lib/apdemande/testData/apdemandeTestData.csv")},
			base.Sirene:      {base.NewBatchFile("lib/sirene/testData/sireneTestData.csv")},
			base.SireneUl:    {base.NewBatchFile("lib/sirene_ul/testData/sireneULTestData.csv")},
			base.AdminUrssaf: {base.NewBatchFile("lib/urssaf/testData/comptesTestData.csv")},
			base.Debit:       {base.NewBatchFile("lib/urssaf/testData/debitTestData.csv")},
			base.Ccsf:        {base.NewBatchFile("lib/urssaf/testData/ccsfTestData.csv")},
			base.Cotisation:  {base.NewBatchFile("lib/urssaf/testData/cotisationTestData.csv")},
			base.Delai:       {base.NewBatchFile("lib/urssaf/testData/delaiTestData.csv")},
			base.Effectif:    {base.NewBatchFile("lib/urssaf/testData/effectifTestData.csv")},
			base.EffectifEnt: {base.NewBatchFile("lib/urssaf/testData/effectifEntTestData.csv")},
			base.Procol:      {base.NewBatchFile("lib/urssaf/testData/procolTestData.csv")},
		},
		Params: base.AdminBatchParams{
			DateDebut: time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC),
			DateFin:   time.Date(2019, time.February, 1, 0, 0, 0, 0, time.UTC),
		},
	}
	writeBatchConfig(t, batch)
}

func verifyReports(t *testing.T) {
	t.Log("💎 Verifying exported reports...")

	conn, err := pgxpool.New(context.Background(), suite.PostgresURI)
	if err != nil {
		t.Errorf("Unable to connect to test database: %s", err)
	}

	table := engine.ReportTable
	query := fmt.Sprintf("SELECT * FROM %s ORDER BY parser", table)
	output := getTableContents(t, conn, query)
	goldenFile := fmt.Sprintf("test-import.sql.%s.golden.txt", table)
	tmpOutputFile := fmt.Sprintf("test-import.sql.%s.output.txt", table)
	compareWithGoldenFileOrUpdate(t, goldenFile, output, tmpOutputFile)
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
		output := fmt.Sprintf("==== %s ====\n%s", baseName, string(content))

		compareWithGoldenFileOrUpdate(t, goldenFile, output, tmpOutputFile)
	}
}

func verifyPostgresExport(t *testing.T) {
	t.Log("🗃️ Verifying postgresql export...")

	conn, err := pgxpool.New(context.Background(), suite.PostgresURI)
	if err != nil {
		t.Errorf("Unable to connect to test database: %s", err)
	}

	tables := getAllTables(t, conn)

	hasMigrationTable := false
	hasEventTable := false

	for _, table := range tables {
		if table == engine.VersionTable {
			hasMigrationTable = true
			continue
		}

		if table == engine.ReportTable {
			hasEventTable = true
			continue
		}

		query := fmt.Sprintf("SELECT * FROM %s", table)
		output := getTableContents(t, conn, query)
		goldenFile := fmt.Sprintf("test-import.sql.%s.golden.txt", table)
		tmpOutputFile := fmt.Sprintf("test-import.sql.%s.output.txt", table)
		compareWithGoldenFileOrUpdate(t, goldenFile, output, tmpOutputFile)
	}
	assert.True(t, hasMigrationTable, "Expecting the migration table to be present")
	assert.True(t, hasEventTable, "Expecting the event table to be present")
}

func getTableContents(t *testing.T, conn *pgxpool.Pool, query string) string {

	rows, err := conn.Query(context.Background(), query)
	if err != nil {
		t.Errorf("failed to query table: %s", err)
	}

	defer rows.Close()

	var result strings.Builder

	fieldDescriptions := rows.FieldDescriptions()

	var headers []string

	// timestamp columns should be skipped for reproducible output
	var includeColumn []bool
	for _, fd := range fieldDescriptions {
		isTimestamp := fd.Name == "start_date"
		includeColumn = append(includeColumn, !isTimestamp)
		if !isTimestamp {
			headers = append(headers, fmt.Sprintf("%-20s", fd.Name)[:20])
		}
	}
	result.WriteString(strings.Join(headers, "\t") + "\n")

	for rows.Next() {
		values, err := rows.Values()
		if err != nil {
			t.Errorf("failed to get row values: %s", err)
		}

		var strValues []string
		for i, v := range values {
			if !includeColumn[i] {
				continue
			}

			if v == nil {
				strValues = append(strValues, fmt.Sprintf("%-20s", "NULL"))
			} else {
				formatted := fmt.Sprintf("%-20s", fmt.Sprintf("%v", v))[:20]
				strValues = append(strValues, formatted)
			}
		}

		result.WriteString(strings.Join(strValues, "\t") + "\n")
	}

	if err := rows.Err(); err != nil {
		t.Errorf("error iterating rows: %s", err)
	}

	return result.String()
}

func getAllTables(t *testing.T, conn *pgxpool.Pool) []string {
	query := `SELECT tablename
  FROM pg_catalog.pg_tables
  WHERE schemaname = 'public'
  `

	rows, err := conn.Query(context.Background(), query)
	if err != nil {
		t.Errorf("failed to query tables: %s", err)
	}
	defer rows.Close()

	var tables []string

	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			t.Errorf("failed to scan table name: %s", err)
		}
		tables = append(tables, tableName)
	}

	if err := rows.Err(); err != nil {
		t.Errorf("error iterating rows: %s", err)
	}

	return tables
}

func TestFilter(t *testing.T) {
	t.Run("Empty filter triggers an error if skipFilter=false", func(t *testing.T) {
		batch := base.AdminBatch{
			Key: "1910",
			Files: map[base.ParserType][]base.BatchFile{
				base.Apconso: {base.NewBatchFile("lib/apconso/testData/apconsoTestData.csv")},
			},
		}

		err := engine.ImportBatch(
			base.BasicBatchProvider{Batch: batch},
			[]base.ParserType{},
			false,
			// Should not write anything to DB
			&engine.FailSinkFactory{},
			engine.DiscardReportSink{},
		)

		t.Log(err)
		assert.ErrorContains(t, err, "n'a pas été initialisé")
		assert.ErrorContains(t, err, "import d'un fichier 'effectif'")
	})

	t.Run("Empty filter triggers no error if skipFilter=true", func(t *testing.T) {
		cleanDB := setupDBTest(t)
		defer cleanDB()

		batch := base.AdminBatch{
			Key: "1910",
			Files: map[base.ParserType][]base.BatchFile{
				base.Apconso: {base.NewBatchFile("lib/apconso/testData/apconsoTestData.csv")},
			},
		}

		err := engine.ImportBatch(
			base.BasicBatchProvider{Batch: batch},
			[]base.ParserType{},
			true,
			// Should not write anything to DB
			&engine.DiscardSinkFactory{},
			engine.DiscardReportSink{},
		)

		assert.NoError(t, err)
	})
}

func TestNonEmptyFilter(t *testing.T) {
	cleanDB := setupDBTest(t)
	defer cleanDB()

	db, err := pgx.Connect(context.Background(), suite.PostgresURI)
	assert.NoError(t, err)

	_, err = db.Exec(
		context.Background(),
		`INSERT INTO stg_effectif (siret, periode, effectif) VALUES('43362355000020', '2018-01-01', 759);`,
	)
	assert.NoError(t, err)
	_, err = db.Exec(
		context.Background(),
		`REFRESH MATERIALIZED VIEW filter;`,
	)
	assert.NoError(t, err)
	_, err = db.Exec(
		context.Background(),
		`TRUNCATE stg_effectif;`,
	)
	assert.NoError(t, err)

	var count int
	err = db.QueryRow(
		context.Background(),
		`SELECT count(*) FROM filter;`,
	).Scan(&count)
	assert.NoError(t, err)

	assert.Equal(t, 1, count)

}
