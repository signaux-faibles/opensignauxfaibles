//go:build e2e

package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
	"testing"
)

func TestCLI(t *testing.T) {
	t.Log("Checking sfdata command's output")
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

			// Format output
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

				re := regexp.MustCompile(`app\.sha1=[a-f0-9]+ `)
				reproducibleStderr := re.ReplaceAllString(stderr.String(), "")

				output.WriteString(reproducibleStderr)
			}

			output.WriteString("---\n")

			compareWithGoldenFileOrUpdate(t, tc.goldenFile, output.String(), tc.tmpFile)
		})
	}
}
