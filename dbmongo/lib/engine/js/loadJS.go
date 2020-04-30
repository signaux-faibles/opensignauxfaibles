package main

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Reads all .js files in the current folder
// and encodes them as strings maps in jsFunctions.go
func bundleJsFunctions() {
	jsRootDir := filepath.Join("..", "..", "js")
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
		if folder.IsDir() && !strings.HasPrefix(folder.Name(), "test") { // TODO: remove this

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
					file, err := os.Open(filepath.Join(jsRootDir, folder.Name(), file.Name()))
					if err != nil {
						log.Print(err)
					}
					io.Copy(out, file)
					out.Write([]byte("`,\n"))
				}
			}
			out.Write([]byte("},\n"))
		}
	}
	out.Write([]byte("}\n"))
}

func transpileTsFunctions() {
	// TODO: also transpile any other TS files
	cmd := exec.Command("npx", "typescript", "../../js/common/raison_sociale.ts") // output: dbmongo/js/common/raison_sociale.js
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	transpileTsFunctions() // convert *.ts files to .js
	bundleJsFunctions()    // bundle *.js files to jsFunctions.go
}
