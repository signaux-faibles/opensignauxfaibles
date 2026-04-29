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

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
)

const exportDir = "tests/export"

func TestExportEndToEnd(t *testing.T) {
	cleanDB := setupDBTest(t)
	defer cleanDB()

	// First, import test data so clean views have content
	createImportTestBatch(t)
	exitCode := runCLI("sfdata", "import", "--batch", "1910", "--no-filter")
	assert.Equal(t, 0, exitCode, "sfdata import should succeed")

	// Clear any previous export files
	os.RemoveAll(exportDir)
	os.MkdirAll(exportDir, 0755)

	t.Run("Run export command", func(t *testing.T) {
		exitCode := runCLI("sfdata", "export", "--path", exportDir)
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
		query := fmt.Sprintf("SELECT * FROM %s", view)
		output := getTableContents(t, conn, query)
		goldenFile := fmt.Sprintf("test-export.%s.golden.txt", view)
		tmpOutputFile := fmt.Sprintf("test-export.%s.output.txt", view)
		compareWithGoldenFileOrUpdate(t, goldenFile, output, tmpOutputFile)
	}
}

