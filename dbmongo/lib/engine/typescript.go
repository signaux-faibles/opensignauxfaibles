package engine

import (
	"log"
	"os"
	"os/exec"
)

func TranspileTsFunctions(jsRootDir string) {
	// TODO: also transpile any other TS files
	tsFiles := []string{
		jsRootDir + "/common/raison_sociale.ts",
		jsRootDir + "/reduce.algo2/fraisFinancier.ts",
	}
	cmd := exec.Command("npx", append([]string{"typescript", "--listFiles", "--lib", "es5", "--skipLibCheck"}, tsFiles...)...) // output: .js files
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}
