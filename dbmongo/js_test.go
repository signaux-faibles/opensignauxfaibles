package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"syscall"
	"testing"
)

func Test_js(t *testing.T) {

	scriptNameRegex, _ := regexp.Compile(".*[.]sh")

	files, err := ioutil.ReadDir("js/test/")
	if err != nil {
		t.Errorf("scripts de test inaccessibles: %v", err.Error())
	}

	for _, f := range files {
		if scriptNameRegex.MatchString(f.Name()) {
			t.Run(f.Name(), func(t *testing.T) {
				if os.Getenv("CI") != "" {
					t.Skip("Skipping testing in CI environment")
				}

				cmd := exec.Command("/bin/bash", f.Name())
				cmd.Dir = "js/test/"

				err := cmdTester(t, cmd)
				if err != nil {
					t.Error(err)
				}
			})
		}
	}

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
