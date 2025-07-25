//go:build e2e

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"opensignauxfaibles/lib/engine"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
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

		insertImportTestBatch(t, db)

		exitCode := runCLI("sfdata", "import", "--batch", "1910", "--no-filter")
		assert.Equal(t, 0, exitCode, "sfdata import should succeed")
	})

	t.Run("Verify Journal reports", func(t *testing.T) {
		verifyJournalReports(t, db)
	})

	t.Run("Verify exported CSV files", func(t *testing.T) {
		verifyExportedCSVFiles(t)
	})

	t.Run("Verify exported Postgres files", func(t *testing.T) {
		verifyPostgresExport(t)

	})
}

func insertImportTestBatch(t *testing.T, db *mgo.Database) {
	t.Log("📝 Inserting test data...")

	batch := bson.M{
		"_id": bson.M{
			"key":  "1910",
			"type": "batch",
		},
		"files": bson.M{
			"dummy":        []string{},
			"filter":       []string{},
			"apconso":      []string{"./lib/apconso/testData/apconsoTestData.csv"},
			"apdemande":    []string{"./lib/apdemande/testData/apdemandeTestData.csv"},
			"sirene":       []string{"./lib/sirene/testData/sireneTestData.csv"},
			"sirene_ul":    []string{"./lib/sirene_ul/testData/sireneULTestData.csv"},
			"admin_urssaf": []string{"./lib/urssaf/testData/comptesTestData.csv"},
			"debit":        []string{"./lib/urssaf/testData/debitTestData.csv"},
			"ccsf":         []string{"./lib/urssaf/testData/ccsfTestData.csv"},
			"cotisation":   []string{"./lib/urssaf/testData/cotisationTestData.csv"},
			"delai":        []string{"./lib/urssaf/testData/delaiTestData.csv"},
			"effectif":     []string{"./lib/urssaf/testData/effectifTestData.csv"},
			"effectif_ent": []string{"./lib/urssaf/testData/effectifEntTestData.csv"},
			"procol":       []string{"./lib/urssaf/testData/procolTestData.csv"},
		},
		"param": bson.M{
			"date_debut": time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC),
			"date_fin":   time.Date(2019, time.February, 1, 0, 0, 0, 0, time.UTC),
		},
	}

	err := db.C("Admin").Insert(batch)
	assert.NoError(t, err)
}

func verifyJournalReports(t *testing.T, db *mgo.Database) {
	t.Log("💎 Verifying Journal reports...")

	var journalEntries []bson.M
	err := db.C("Journal").Find(nil).Sort("reportType", "parserCode").All(&journalEntries)
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
		transformed["hasCommitHash"] = doc["commitHash"] != nil
		transformed["hasDate"] = doc["date"] != nil
		transformed["hasStartDate"] = doc["startDate"] != nil

		transformedEntries = append(transformedEntries, transformed)
	}

	// Convert to string for comparison
	output, err := json.MarshalIndent(transformedEntries, "", "  ")
	assert.NoError(t, err)

	goldenFilePath := "test-import.journal.golden.txt"
	tmpOutputPath := "test-import.journal.output.txt"

	compareWithGoldenFileOrUpdate(t, goldenFilePath, string(output), tmpOutputPath)
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

	for _, table := range tables {
		if table == engine.VersionTable {
			continue
		}
		output := getTableContents(t, conn, table)
		goldenFile := fmt.Sprintf("test-import.sql.%s.golden.txt", table)
		tmpOutputFile := fmt.Sprintf("test-import.sql.%s.output.txt", table)
		compareWithGoldenFileOrUpdate(t, goldenFile, output, tmpOutputFile)
	}
}

func getTableContents(t *testing.T, conn *pgxpool.Pool, tableName string) string {
	query := fmt.Sprintf("SELECT * FROM %s", tableName)

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
				strValues = append(strValues, fmt.Sprintf("%-20v", v)[:20])
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
