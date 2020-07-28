package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/engine"
)

// Reads all .js files in the current folder
// and encodes them as strings maps in jsFunctions.go
func bundleJsFunctions(jsRootDir string) {
	folders, err := ioutil.ReadDir(jsRootDir)
	if err != nil {
		log.Fatal(err)
	}

	out, err := os.Create("jsFunctions.go")
	if err != nil {
		log.Fatal(err)
	}
	out.Write([]byte("package engine \n\n var jsFunctions = map[string]map[string]string{\n"))

	// For each folder
	for _, folder := range folders {
		if folder.IsDir() &&
			folder.Name() != "node_modules" && // skip node/npm dependencies cache
			!strings.HasPrefix(folder.Name(), ".") && // skip hidden directories, e.g. `.nyc_output`
			!strings.HasPrefix(folder.Name(), "test") {

			out.Write([]byte(`"` + folder.Name() + `"` + ":{\n"))

			files, err := ioutil.ReadDir(filepath.Join(jsRootDir, folder.Name()))
			if err != nil {
				log.Print(err)
			}

			// For each file in folder
			for _, file := range files {
				if strings.HasSuffix(file.Name(), ".js") && !strings.HasSuffix(file.Name(), "_test.js") {
					out.Write([]byte(
						`"` + strings.TrimSuffix(file.Name(), ".js") + `"` +
							": `"))

					function, err := ioutil.ReadFile(filepath.Join(jsRootDir, folder.Name(), file.Name()))
					if err != nil {
						log.Fatal(err)
					}
					stringFunction := string(function)
					exportsDefRegex := regexp.MustCompile(`(?m)^Object.defineProperty\(exports.*$`)
					stringFunction = exportsDefRegex.ReplaceAllLiteralString(stringFunction, "")
					finalExportRegex := regexp.MustCompile(`(?m)^exports..*$`)
					skipLineRegex := regexp.MustCompile(`(?m)^.*DO_NOT_INCLUDE_IN_JSFUNCTIONS_GO.*$`)
					stringFunction = skipLineRegex.ReplaceAllLiteralString(stringFunction, "")
					stringFunction = finalExportRegex.ReplaceAllLiteralString(stringFunction, "")
					stringFunction = strings.Replace(stringFunction, "`", "` + \"`\" + `", -1) // escape nested "backticks" quotes
					stringFunction = strings.Trim(stringFunction, "\n")

					out.Write([]byte(stringFunction))
					out.Write([]byte("`,\n"))
				}
			}
			out.Write([]byte("},\n"))
		}
	}
	out.Write([]byte("}\n"))
}

func main() {
	jsRootDir := filepath.Join("..", "..", "js")
	engine.TranspileTsFunctions(jsRootDir)  // convert *.ts files to .js
	engine.GlobalizeJsFunctions(jsRootDir)  // remove "export" prefixes from JS functions, for jsc compatibility
	bundleJsFunctions(jsRootDir)            // bundle *.js files to jsFunctions.go
	engine.DeleteTranspiledFiles(jsRootDir) // delete the *.js files
}
