package engine

import (
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

func shouldTranspile(filePath string) bool {
	return !strings.Contains(filePath, "node_modules") &&
		!strings.Contains(filePath, "_tests.ts") &&
		path.Ext(filePath) == ".ts"
}

// ListTsFiles retourne la liste des fichiers TypeScript transpilable en JavaScript
// en cherchant récursivement depuis le répertoire jsRootDir.
func ListTsFiles(jsRootDir string) []string {
	var files []string
	err := filepath.Walk(jsRootDir, func(filePath string, info os.FileInfo, err error) error {
		if err == nil && shouldTranspile(filePath) {
			files = append(files, filePath)
		}
		return err
	})
	if err != nil {
		log.Fatal(err)
	}
	return files
}

// DeleteTranspiledFiles supprime les fichiers JavaScript résultant de la
// transpilation des fichiers TypeScript listés dans tsFiles.
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

// TranspileTsFunctions convertit les fichiers TypeScript listés dans tsFiles
// au format JavaScript.
func TranspileTsFunctions(tsFiles []string) {
	cmd := exec.Command("npx", append([]string{"typescript", "--listFiles", "--lib", "es5", "--skipLibCheck", "--noImplicitUseStrict"}, tsFiles...)...) // output: .js files
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}
