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
		!strings.Contains(filePath, ".d.ts") &&
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
		transpiledFile := strings.TrimSuffix(tsFile, ext) + ".js"
		err := os.Remove(transpiledFile)
		if err != nil {
			panic("failed to delete " + transpiledFile)
		}
	}
}

// TranspileTsFunctions convertit les fichiers TypeScript au format JavaScript.
func TranspileTsFunctions(jsRootDir string) {
	cmd := exec.Command("npx", "typescript", "--listFiles", "--p", filepath.Join(jsRootDir, "tsconfig.json")) // output: .js files
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

// GlobalizeJsFunctions retire le préfixe "export" des fonctions, pour les rendre compatibles avec jsc.
func GlobalizeJsFunctions(jsRootDir string) {
	cmd := exec.Command("bash", "globalize-functions.sh") // output: .js files
	cmd.Dir = jsRootDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}
