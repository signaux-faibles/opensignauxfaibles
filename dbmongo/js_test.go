package main

import (
	"fmt"
	"io/ioutil"
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
			cmd := exec.Command("/bin/sh", f.Name())
			cmd.Dir = "js/test/"

			err := cmdTester(cmd)
			if err != nil {
				t.Errorf("erreur levée par %v: "+err.Error(), f.Name())
			} else {
				t.Logf("execution de %v ok", f.Name())
			}
		}
	}

}

func cmdTester(cmd *exec.Cmd) error {
	if err := cmd.Run(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				return fmt.Errorf("Error status: %v", status.ExitStatus())
			}
		} else {
			return fmt.Errorf("cmd.Run: %v", err.Error())
		}
	}

	return nil
}
