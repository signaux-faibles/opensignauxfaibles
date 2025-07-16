//go:build e2e

package main

import (
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseFileEndToEnd(t *testing.T) {
	t.Log("💎 Parsing data...")

	outputFile := "test-parseFile.output.txt"
	goldenFile := "test-parseFile.golden.txt"

	cmd := exec.Command("./sfdata", "parseFile", "--parser", "apconso", "--file", "./lib/apconso/testData/apconsoTestData.csv")
	cmd.Env = append(os.Environ(), "NO_DB=1")

	output, err := cmd.CombinedOutput()

	assert.NoError(t, err, "Command failed: %s", string(output))

	// Compare with or update golden file
	compareWithGoldenFileOrUpdate(t, goldenFile, string(output), outputFile)
}
