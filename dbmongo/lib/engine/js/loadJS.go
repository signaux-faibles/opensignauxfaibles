package main

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// Reads all .js files in the current folder
// and encodes them as strings maps in jsFunctions.go
func main() {
	// TODO: use filepath
	jsRootDir := "../../js/"
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
		if folder.IsDir() && folder.Name() != "tests" {

			out.Write([]byte(`"` + folder.Name() + `"` + ":{\n"))

			files, err := ioutil.ReadDir(jsRootDir + folder.Name())
			if err != nil {
				log.Print(err)
			}

			// For each file in folder
			for _, file := range files {
				if strings.HasSuffix(file.Name(), ".js") {
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
