package createfilter

import (
	"bytes"
	"log"
	"os"

	"github.com/pkg/errors"
)

// DiffWithGoldenFile compares the output of the execution of a test with the expected output.
func DiffWithGoldenFile(filename string, updateGoldenFile bool, cmdOutput bytes.Buffer) []byte {

	if updateGoldenFile {
		_ = os.WriteFile(filename, cmdOutput.Bytes(), 0644)
	}
	expected, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	return expected
}

// ReadGoldenFile retourne le fichier passé en paramètre sous forme de chaine de caractère
func ReadGoldenFile(filename string) string {
	expected, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(errors.Wrap(err, "Erreur de lecture du Golden File"))
	}
	return string(expected)
}
