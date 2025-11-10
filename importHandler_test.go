package main

import (
	"encoding/json"
	"opensignauxfaibles/lib/engine"
	"os"
	"path"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestDryRunWithoutDatabase verifies that an import with --dry-run and --no-filter
// succeeds even when no database is available.
func TestDryRunWithoutDatabase(t *testing.T) {
	// Save original POSTGRES_DB_URL and clear it for this test
	originalDBURL := os.Getenv("POSTGRES_DB_URL")
	os.Setenv("POSTGRES_DB_URL", "")
	defer func() {
		// Restore original value
		if originalDBURL != "" {
			os.Setenv("POSTGRES_DB_URL", originalDBURL)
		} else {
			os.Unsetenv("POSTGRES_DB_URL")
		}
	}()

	// Create a temporary directory for the batch
	tmpDir := t.TempDir()

	// Create a minimal batch configuration
	batch := engine.AdminBatch{
		Key: "1910",
		Files: map[engine.ParserType][]engine.BatchFile{
			engine.Apconso:   {engine.NewBatchFile("lib/parsing/apconso/testData/apconsoTestData.csv")},
			engine.Apdemande: {engine.NewBatchFile("lib/parsing/apdemande/testData/apdemandeTestData.csv")},
			engine.Sirene:    {engine.NewBatchFile("lib/parsing/sirene/testData/sireneTestData.csv")},
		},
		Params: engine.AdminBatchParams{
			DateDebut: time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC),
			DateFin:   time.Date(2019, time.February, 1, 0, 0, 0, 0, time.UTC),
		},
	}

	// Write batch configuration
	batchConfigPath := path.Join(tmpDir, "batch.json")
	batchBytes, err := json.Marshal(batch)
	assert.NoError(t, err, "failed to marshal batch config")
	err = os.WriteFile(batchConfigPath, batchBytes, 0644)
	assert.NoError(t, err, "failed to write batch config")

	t.Run("dry-run with no-filter should succeed without database", func(t *testing.T) {
		exitCode := runCLI(
			"sfdata",
			"import",
			"--dry-run",
			"--no-filter",
			"--batch", "1910",
			"--batch-config", batchConfigPath,
		)
		assert.Equal(t, 0, exitCode, "sfdata import with --dry-run and --no-filter should succeed without database")
	})

	t.Run("dry-run without no-filter should fail gracefully without database", func(t *testing.T) {
		// Without --no-filter, the import should fail because it needs to read the filter
		// from the database, which is not available
		exitCode := runCLI(
			"sfdata",
			"import",
			"--dry-run",
			"--batch", "1910",
			"--batch-config", batchConfigPath,
		)
		assert.NotEqual(t, 0, exitCode, "sfdata import with --dry-run but without --no-filter should fail when database is not available")
	})
}
