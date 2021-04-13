package main

import (
	"bytes"
	"go/format"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

// Reads all .js files in the current folder
// and encodes them as strings maps in jsFunctions.go
func bundleJsFunctions(jsRootDir string) {
	folders, err := ioutil.ReadDir(jsRootDir)
	if err != nil {
		log.Fatal(err)
	}

	var out bytes.Buffer
	out.Write([]byte("package engine\n"))
	out.Write([]byte("import \"errors\"\n"))
	out.Write([]byte("import \"github.com/globalsign/mgo/bson\"\n"))
	out.Write([]byte("type functions = map[string]string\n"))
	out.Write([]byte("type functionGetter = func (bson.M) (functions, error)\n"))
	out.Write([]byte("var jsFunctions = map[string]functionGetter{\n"))

	// For each folder
	for _, folder := range folders {
		if folder.IsDir() &&
			folder.Name() != "node_modules" && // skip node/npm dependencies cache
			folder.Name() != "typings" && // skip typescript types for javascript dependencies (e.g. concordance)
			!strings.HasPrefix(folder.Name(), ".") && // skip hidden directories, e.g. `.nyc_output`
			!strings.HasPrefix(folder.Name(), "test") {

			out.Write([]byte(`"` + folder.Name() + `"` + ": func (params bson.M) (functions, error) {\n"))

			globals, err := getTypeScriptGlobals(jsRootDir, folder.Name())
			if err != nil {
				log.Fatal(err)
			}
			for _, globalParam := range globals {
				out.Write([]byte("if _, ok := params[\"" + globalParam + "\"]; !ok {\n"))
				out.Write([]byte("return nil, errors.New(\"missing required parameter: " + globalParam + "\")\n"))
				out.Write([]byte("}\n"))
			}
			out.Write([]byte("return functions{\n"))

			files, err := ioutil.ReadDir(filepath.Join(jsRootDir, folder.Name()))
			if err != nil {
				log.Print(err)
			}

			// For each file in folder
			for _, file := range files {
				if shouldInclude(file) {
					out.Write([]byte(
						`"` + strings.TrimSuffix(file.Name(), ".js") + `"` +
							": `"))

					function, err := ioutil.ReadFile(filepath.Join(jsRootDir, folder.Name(), file.Name()))
					if err != nil {
						log.Fatal(err)
					}
					stringFunction := string(function)
					stringFunction = strings.Replace(stringFunction, "`", "` + \"`\" + `", -1) // escape nested "backticks" quotes
					stringFunction = strings.Trim(stringFunction, "\n")

					out.Write([]byte(stringFunction))
					out.Write([]byte("`,\n"))
				}
			}
			out.Write([]byte("}, nil \n },\n"))
		}
	}
	out.Write([]byte("}\n"))

	formatted, err := format.Source(out.Bytes())
	if err != nil {
		log.Fatal(err)
	}
	fileOut, err := os.Create("jsFunctions.go")
	if err != nil {
		log.Fatal(err)
	}
	fileOut.Write(formatted)
}

func main() {
	jsRootDir := filepath.Join("..", "..", "js")
	TranspileTsFunctions(jsRootDir)  // convert *.ts files to .js
	bundleJsFunctions(jsRootDir)     // bundle *.js files to jsFunctions.go
	DeleteTranspiledFiles(jsRootDir) // delete the *.js files
}

func shouldInclude(file os.FileInfo) bool {
	return file.Name() != "functions.js" &&
		(strings.HasSuffix(file.Name(), ".js") ||
			strings.HasSuffix(file.Name(), ".json"))
}

func shouldTranspile(filePath string) bool {
	return !strings.Contains(filePath, "node_modules") &&
		!strings.Contains(filePath, "test") &&
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
func DeleteTranspiledFiles(jsRootDir string) {
	tsFiles := ListTsFiles(jsRootDir)
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
	cmd := exec.Command("bash", "generate-javascript.sh") // output: .js files
	cmd.Dir = jsRootDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

// getTypeScriptGlobals liste les variables globales utilisées par les fichiers TypeScript de subDir.
func getTypeScriptGlobals(jsRootDir string, sudDir string) ([]string, error) {
	cmd := exec.Command("bash", "./get-globals.sh", sudDir+"/*.ts") // output: comma-separated list of globals
	cmd.Dir = jsRootDir
	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return nil, err
	}
	commaSeparatedList := strings.TrimSpace(output.String())
	return strings.Split(commaSeparatedList, ","), nil
}
