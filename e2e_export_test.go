//go:build e2e

package main

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

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

	t.Run("Verify exported parquet files", func(t *testing.T) {
		verifyExportedFiles(t)
	})

	os.RemoveAll(exportDir)
}

func verifyExportedFiles(t *testing.T) {
	t.Log("Verifying exported parquet files from PostgreSQL...")

	files, err := filepath.Glob(filepath.Join(exportDir, "*.parquet"))
	assert.NoError(t, err)
	assert.NotEmpty(t, files, "Expected at least one exported parquet file")

	sort.Strings(files)

	for _, file := range files {
		content, err := os.ReadFile(file)
		assert.NoError(t, err)

		baseName := filepath.Base(file)
		viewName := strings.TrimSuffix(baseName, ".parquet")

		goldenFile := "test-export." + viewName + ".golden.txt"
		tmpOutputFile := "test-export." + viewName + ".output.txt"

		compareWithGoldenFileOrUpdate(t, goldenFile, string(content), tmpOutputFile)
	}
}
