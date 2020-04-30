package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"
	"syscall"
	"testing"
)

const SKIP_ON_CI = "SKIP_ON_CI"

var update = flag.Bool("update", false, "Update the expected test values in golden file")

func Test_js(t *testing.T) {

	scriptNameRegex, _ := regexp.Compile(".*[.]sh")
	testdir := path.Join("js", "test")
	files, err := ioutil.ReadDir(testdir)
	if err != nil {
		t.Errorf("scripts de test inaccessibles: %v", err.Error())
	}

	for _, f := range files {
		if scriptNameRegex.MatchString(f.Name()) {
			t.Run(f.Name(), func(t *testing.T) {
				filepath := path.Join(testdir, f.Name())
				if os.Getenv("CI") != "" && shouldSkipOnCi(t, filepath) {
					t.Skip("Skipping testing in CI environment")
				}

				cmd := exec.Command("/bin/bash", f.Name())
				cmd.Dir = testdir

				err := cmdTester(t, cmd)
				if err != nil {
					t.Error(err)
				}
			})
		}
	}
}

func shouldSkipOnCi(t *testing.T, filepath string) bool {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		t.Error(err)
	}
	file := string(data)
	return strings.Contains(file, SKIP_ON_CI)
}

func cmdTester(t *testing.T, cmd *exec.Cmd) error {
	t.Helper()

	var cmdOutput bytes.Buffer
	var cmdError bytes.Buffer
	cmd.Stdout = &cmdOutput
	cmd.Stderr = &cmdError

	if err := cmd.Run(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				return fmt.Errorf("Error status: %v\nstderr: %v\nstdout: %v", status.ExitStatus(), cmdError.String(), cmdOutput.String())
			}
		} else {
			return fmt.Errorf("cmd.Run: %v", err.Error())
		}
	}

	return nil
}
