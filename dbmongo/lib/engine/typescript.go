package engine

import (
	"log"
	"os"
	"os/exec"
)

func TranspileTsFunctions(jsRootDir string) {
	// TODO: also transpile any other TS files
	cmd := exec.Command("npx", "typescript", jsRootDir+"/common/raison_sociale.ts") // output: dbmongo/js/common/raison_sociale.js
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}
