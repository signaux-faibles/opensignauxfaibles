//go:build e2e

package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"testing"

	"opensignauxfaibles/lib/export"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
)

const exportDir = "tests/export"

func TestExportEndToEnd(t *testing.T) {
	cleanDB := setupDBTest(t)
	defer cleanDB()

	// First, import test data so clean views have content
	createImportTestBatch(t)
	exitCode := runCLI("sfdata", "import", "--schema", "sfdata", "--batch", "1910", "--no-filter")
	assert.Equal(t, 0, exitCode, "sfdata import should succeed")

	// Clear any previous export files
	os.RemoveAll(exportDir)
	os.MkdirAll(exportDir, 0755)

	t.Run("Run export command", func(t *testing.T) {
		exitCode := runCLI("sfdata", "export", "--schema", "sfdata", "--path", exportDir)
		assert.Equal(t, 0, exitCode, "sfdata export should succeed")
	})

	t.Run("Verify exported parquet files exist", func(t *testing.T) {
		verifyExportedParquetFilesExist(t)
	})

	t.Run("Verify exported view contents", func(t *testing.T) {
		verifyExportedViewContents(t)
	})

	os.RemoveAll(exportDir)
}

// verifyExportedParquetFilesExist checks that parquet files were created for each exported view
func verifyExportedParquetFilesExist(t *testing.T) {
	t.Log("Verifying exported parquet files exist...")

	files, err := filepath.Glob(filepath.Join(exportDir, "*.parquet"))
	assert.NoError(t, err)
	assert.NotEmpty(t, files, "Expected at least one exported parquet file")

	for _, file := range files {
		info, err := os.Stat(file)
		assert.NoError(t, err)
		assert.Greater(t, info.Size(), int64(0), "Parquet file %s should not be empty", file)
	}
}

// TestExportSucceedsOnFreshSchema reproduces the bug where a `WITH NO DATA`
// materialized view (e.g. clean_procol_at_date) makes `COPY ... TO STDOUT`
// fail with SQLSTATE 55000, which used to cascade through errgroup and leave
// most parquet files at 0 bytes.
//
// We run export on a fresh schema (only migrations, no import) and check that
// every parquet file got fully written by pg_parquet (non-zero size — even an
// empty result set produces the parquet header/footer).
func TestExportSucceedsOnFreshSchema(t *testing.T) {
	const schema = "export_fresh_schema_test"

	ctx := context.Background()
	conn, err := pgx.Connect(ctx, suite.PostgresURI)
	if err != nil {
		t.Fatalf("Unable to connect to test database: %s", err)
	}
	_, err = conn.Exec(ctx, fmt.Sprintf("DROP SCHEMA IF EXISTS %s CASCADE", schema))
	assert.NoError(t, err)
	conn.Close(ctx)

	freshExportDir := filepath.Join(suite.TmpDir, "export-fresh")
	os.RemoveAll(freshExportDir)
	os.MkdirAll(freshExportDir, 0755)
	defer os.RemoveAll(freshExportDir)

	exitCode := runCLI("sfdata", "export", "--schema", schema, "--path", freshExportDir)
	assert.Equal(t, 0, exitCode, "sfdata export should succeed on a fresh schema")

	files, err := filepath.Glob(filepath.Join(freshExportDir, "*.parquet"))
	assert.NoError(t, err)
	assert.NotEmpty(t, files, "Expected parquet files on a fresh schema")
	for _, file := range files {
		info, err := os.Stat(file)
		assert.NoError(t, err)
		assert.Greater(t, info.Size(), int64(0),
			"Parquet file %s is empty — unpopulated MV likely cascaded through errgroup", file)
	}
}

// verifyExportedViewContents queries each exported view via SQL and compares with golden files
func verifyExportedViewContents(t *testing.T) {
	t.Log("Verifying exported view contents via SQL...")

	conn, err := pgxpool.New(context.Background(), suite.PostgresURI)
	if err != nil {
		t.Fatalf("Unable to connect to test database: %s", err)
	}
	defer conn.Close()

	views := export.ViewsToExport()

	sort.Strings(views)

	for _, view := range views {
		query := fmt.Sprintf("SELECT * FROM %s ORDER BY 1", view)
		output := getTableContents(t, conn, query)
		goldenFile := fmt.Sprintf("test-export.%s.golden.txt", view)
		tmpOutputFile := fmt.Sprintf("test-export.%s.output.txt", view)
		compareWithGoldenFileOrUpdate(t, goldenFile, output, tmpOutputFile)
	}
}

