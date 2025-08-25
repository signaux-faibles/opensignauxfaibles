//go:build e2e

package main

import (
	"context"
	"fmt"
	"opensignauxfaibles/lib/base"
	"opensignauxfaibles/lib/engine"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/globalsign/mgo"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
)

func TestImportEndToEnd(t *testing.T) {

	mongodb, err := mgo.Dial(suite.MongoURI)
	assert.NoError(t, err)
	defer mongodb.Close()

	db := mongodb.DB(mongoDatabase)

	t.Run("Insert test data and run import", func(t *testing.T) {
		cleanDatabase(t, db)

		createImportTestBatch(t)

		exitCode := runCLI("sfdata", "import", "--batch", "1910", "--no-filter")
		assert.Equal(t, 0, exitCode, "sfdata import should succeed")
	})

	t.Run("Verify Events reports", func(t *testing.T) {
		verifyEventsReports(t, db)
	})

	t.Run("Verify exported CSV files", func(t *testing.T) {
		verifyExportedCSVFiles(t)
	})

	t.Run("Verify exported Postgres files", func(t *testing.T) {
		verifyPostgresExport(t)

	})
}

func createImportTestBatch(t *testing.T) {

	batch := base.AdminBatch{
		ID: base.AdminID{
			Key:  "1910",
			Type: "batch",
		},
		Files: map[string][]base.BatchFile{
			"dummy":        {},
			"filter":       {},
			"apconso":      {base.NewBatchFile("./lib/apconso/testData/apconsoTestData.csv")},
			"apdemande":    {base.NewBatchFile("./lib/apdemande/testData/apdemandeTestData.csv")},
			"sirene":       {base.NewBatchFile("./lib/sirene/testData/sireneTestData.csv")},
			"sirene_ul":    {base.NewBatchFile("./lib/sirene_ul/testData/sireneULTestData.csv")},
			"admin_urssaf": {base.NewBatchFile("./lib/urssaf/testData/comptesTestData.csv")},
			"debit":        {base.NewBatchFile("./lib/urssaf/testData/debitTestData.csv")},
			"ccsf":         {base.NewBatchFile("./lib/urssaf/testData/ccsfTestData.csv")},
			"cotisation":   {base.NewBatchFile("./lib/urssaf/testData/cotisationTestData.csv")},
			"delai":        {base.NewBatchFile("./lib/urssaf/testData/delaiTestData.csv")},
			"effectif":     {base.NewBatchFile("./lib/urssaf/testData/effectifTestData.csv")},
			"effectif_ent": {base.NewBatchFile("./lib/urssaf/testData/effectifEntTestData.csv")},
			"procol":       {base.NewBatchFile("./lib/urssaf/testData/procolTestData.csv")},
		},
		Params: base.AdminBatchParams{
			DateDebut: time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC),
			DateFin:   time.Date(2019, time.February, 1, 0, 0, 0, 0, time.UTC),
		},
	}
	writeBatchConfig(t, batch)
}

func verifyEventsReports(t *testing.T) {
	t.Log("💎 Verifying Events reports...")

	conn, err := pgxpool.New(context.Background(), suite.PostgresURI)
	if err != nil {
		t.Errorf("Unable to connect to test database: %s", err)
	}

	table := engine.EventTable
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

		if table == engine.EventTable {
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
	for _, fd := range fieldDescriptions {
		headers = append(headers, fmt.Sprintf("%-20s", fd.Name)[:20])
	}
	result.WriteString(strings.Join(headers, "\t") + "\n")

	for rows.Next() {
		values, err := rows.Values()
		if err != nil {
			t.Errorf("failed to get row values: %s", err)
		}

		var strValues []string
		for _, v := range values {
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
