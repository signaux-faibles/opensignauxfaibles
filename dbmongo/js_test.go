package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"syscall"
	"testing"

	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/engine"
)

var update = flag.Bool("update", false, "Update the expected test values in golden file")

// TestMain sera exécuté avant les tests
func TestMain(m *testing.M) {
	fmt.Println("Transpilation des fonctions JS depuis TypeScript...")
	jsRootDir := filepath.Join("js") // chemin vers les fichiers TS et JS (sous-répertoire)
	tsFiles := engine.ListTsFiles(jsRootDir)
	engine.TranspileTsFunctions(jsRootDir) // convert *.ts files to .js
	engine.GlobalizeJsFunctions(jsRootDir) // remove "export" prefixes from JS functions, for jsc compatibility
	code := m.Run()
	engine.DeleteTranspiledFiles(tsFiles) // delete the *.js files
	os.Exit(code)
}

func Test_js(t *testing.T) {

	scriptNameRegex, _ := regexp.Compile(".*[.]sh")
	testdir := path.Join("js", "test")
	files, err := ioutil.ReadDir(testdir)
	if err != nil {
		t.Errorf("scripts de test inaccessibles: %v", err.Error())
	}

	if *update {
		fmt.Println("Les golden files vont être mis à jour")
	}

	for _, f := range files {
		if scriptNameRegex.MatchString(f.Name()) {
			t.Run(f.Name(), func(t *testing.T) {
				var cmd *exec.Cmd
				if *update {
					cmd = exec.Command("/bin/bash", f.Name(), "--update")
				} else {
					cmd = exec.Command("/bin/bash", f.Name())
				}
				cmd.Dir = testdir

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
