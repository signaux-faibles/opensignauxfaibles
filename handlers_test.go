package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var update = flag.Bool("update", false, "Update the expected test values in golden file")

const (
	tmpDir         = "tests/tmp-test-execution-files"
	goldenFilesDir = "tests/output-snapshots"
)

func TestCLI(t *testing.T) {

	// Setup temporary directory
	err := os.MkdirAll(tmpDir, 0755)
	assert.NoError(t, err)

	t.Cleanup(func() {
	})

	testCases := []struct {
		name       string
		args       []string
		goldenFile string
		tmpFile    string
	}{
		{
			"sfdata",
			[]string{},
			"test-cli.1.golden.txt",
			"test-cli.1.output.txt",
		},
		{
			"sfdata --help",
			[]string{"--help"},
			"test-cli.2.golden.txt",
			"test-cli.2.output.txt",
		},
		{
			"sfdata unknown_command",
			[]string{"unknown_command"},
			"test-cli.unknown.golden.txt",
			"test-cli.unknown.output.txt",
		},
		{
			"sfdata check --help",
			[]string{"check", "--help"},
			"test-cli.check.golden.txt",
			"test-cli.check.output.txt",
		},
		{
			"sfdata import --help",
			[]string{"import", "--help"},
			"test-cli.import.golden.txt",
			"test-cli.import.output.txt",
		},
		{
			"sfdata parseFile --help",
			[]string{"parseFile", "--help"},
			"test-cli.parseFile.golden.txt",
			"test-cli.parseFile.output.txt",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Execute command and capture output
			cmd := exec.Command("./sfdata", tc.args...)
			var stdout, stderr bytes.Buffer
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr

			err := cmd.Run()

			// Format output similar to your bash script
			cmdStr := fmt.Sprintf("./sfdata %s", strings.Join(tc.args, " "))
			t.Logf("- %s", cmdStr)
			if err != nil {
				t.Logf("Command failed with: %v", err)
			}

			var output bytes.Buffer
			output.WriteString(fmt.Sprintf("$ %s\n", cmdStr))

			if stdout.Len() > 0 {
				output.WriteString(stdout.String())
			}
			if stderr.Len() > 0 {
				output.WriteString("--- stderr capture\n")
				output.WriteString(stderr.String())
			}

			output.WriteString("---\n")

			goldenFilePath := path.Join(goldenFilesDir, tc.goldenFile)

			// Handle golden file comparison/update
			if *update {
				err := updateGoldenFile(goldenFilePath, output.String())
				assert.NoError(t, err)

				t.Log("✅ Golden master file updated")

			} else {

				err := compareWithGoldenFile(goldenFilePath, output.String())
				if err != nil {
					// Write output to temp file for easy diffing
					outputFilePath := filepath.Join(tmpDir, tc.tmpFile)
					_ = os.WriteFile(outputFilePath, output.Bytes(), 0644)
					t.Logf("💾 Output written to: %s", outputFilePath)
				}

				assert.NoError(t, err)
			}
		})
	}

	// Only if all tests passes, otherwise we want to keep the tmp files for
	// inspection
	os.RemoveAll(tmpDir)
}

// updateGoldenFile writes the output to the golden file
func updateGoldenFile(goldenPath, content string) error {
	dir := filepath.Dir(goldenPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	return os.WriteFile(goldenPath, []byte(content), 0644)
}

// compareWithGoldenFile compares the output with the golden file
func compareWithGoldenFile(goldenPath, actual string) error {
	expected, err := os.ReadFile(goldenPath)
	if err != nil {
		return fmt.Errorf("failed to read golden file %s: %w", goldenPath, err)
	}

	if string(expected) != actual {
		return fmt.Errorf("output doesn't match golden file %s.\nExpected:\n%s\nActual:\n%s",
			goldenPath, string(expected), actual)
	}

	return nil
}
