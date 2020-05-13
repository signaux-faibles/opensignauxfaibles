package engine

import (
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
)

func ListTsFiles(jsRootDir string) []string {
	// TODO: fetch the list of TS files by walking the file hierarchy
	return []string{
		filepath.Join(jsRootDir, "common", "raison_sociale.ts"),
		filepath.Join(jsRootDir, "reduce.algo2", "fraisFinancier.ts"),
	}
}

func DeleteTranspiledFiles(tsFiles []string) {
	for _, tsFile := range tsFiles {
		ext := path.Ext(tsFile)
		if ext != ".ts" {
			panic("expected a .ts file, found: " + tsFile)
		}
		transpiledFile := tsFile[0:len(tsFile)-len(ext)] + ".js"
		err := os.Remove(transpiledFile)
		if err != nil {
			panic("failed to delete " + transpiledFile)
		}
	}
}

func TranspileTsFunctions(tsFiles []string) {
	cmd := exec.Command("npx", append([]string{"typescript", "--listFiles", "--lib", "es5", "--skipLibCheck", "--noImplicitUseStrict"}, tsFiles...)...) // output: .js files
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}
