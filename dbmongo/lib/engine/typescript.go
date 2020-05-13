package engine

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func TranspileTsFunctions(jsRootDir string) {
	// TODO: also transpile any other TS files
	tsFiles := []string{
		filepath.Join(jsRootDir, "common", "raison_sociale.ts"),
		filepath.Join(jsRootDir, "reduce.algo2", "fraisFinancier.ts"),
	}
	cmd := exec.Command("npx", append([]string{"typescript", "--listFiles", "--lib", "es5", "--skipLibCheck", "--noImplicitUseStrict"}, tsFiles...)...) // output: .js files
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}
