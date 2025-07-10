package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

var update = flag.Bool("update", false, "Update the expected test values in golden file")

const (
	tmpDir     = "tests/tmp-test-execution-files"
	outputFile = "test-cli.output.txt" // filename within tmpDir
	goldenFile = "tests/output-snapshots/test-cli.golden.txt"
)

func TestCLI(t *testing.T) {

	// Setup temporary directory
	err := os.MkdirAll(tmpDir, 0755)
	require.NoError(t, err)

	t.Cleanup(func() {
	})

	var output bytes.Buffer

	testCases := []struct {
		name string
		args []string
	}{
		{"sfdata", []string{}},
		{"sfdata unknown_command", []string{"unknown_command"}},
		{"sfdata --help", []string{"--help"}},
		{"sfdata check --help", []string{"check", "--help"}},
		{"sfdata import --help", []string{"import", "--help"}},
		{"sfdata parseFile --help", []string{"parseFile", "--help"}},
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
			output.WriteString(fmt.Sprintf("$ %s\n", cmdStr))

			if stdout.Len() > 0 {
				output.WriteString(stdout.String())
			}
			if stderr.Len() > 0 {
				output.WriteString(stderr.String())
			}

			output.WriteString("---\n")

			// Log for debugging
			t.Logf("- %s", cmdStr)
			if err != nil {
				t.Logf("Command failed with: %v", err)
			}
		})
	}

	// Write output to temp file for easy diffing
	outputFilePath := filepath.Join(tmpDir, outputFile)
	err = os.WriteFile(outputFilePath, []byte(output.String()), 0644)
	require.NoError(t, err)
	t.Logf("💾 Output written to: %s", outputFilePath)

	// Handle golden file comparison/update
	if *update {
		err := updateGoldenFile(goldenFile, output.String())
		require.NoError(t, err)
		t.Log("✅ Golden master file updated")
	} else {
		err := compareWithGoldenFile(goldenFile, output.String())
		require.NoError(t, err)
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
