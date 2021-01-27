package diane

import (
	"bytes"
	"encoding/csv"
	"flag"
	"io/ioutil"
	"log"
	"path/filepath"
	"testing"

	"github.com/signaux-faibles/opensignauxfaibles/lib/marshal"
	"github.com/stretchr/testify/assert"
)

var update = flag.Bool("update", false, "Update the expected test values in golden file")

func TestDiane(t *testing.T) {

	t.Run("openFile() doit retourner un flux csv avec un en-tête sans duplication de caractères d'espacement", func(t *testing.T) {
		var testData = filepath.Join("testData", "dianeTestData.txt")
		_, reader, err := openFile(testData)
		if assert.NoError(t, err) {
			csvReader := csv.NewReader(*reader)
			csvReader.Comma = ';'
			csvReader.LazyQuotes = true
			header, err := csvReader.Read() // Discard header
			if assert.NoError(t, err) {
				for _, field := range header {
					assert.NotContains(t, field, "  ")
				}
			}
		}
	})

	t.Run("openFile() doit produire un fichier csv intermédiaire conforme", func(t *testing.T) {
		var golden = filepath.Join("testData", "expectedDianeConvert.csv")
		var testData = filepath.Join("testData", "dianeTestData.txt")
		_, reader, err := openFile(testData)
		if assert.NoError(t, err) {
			buf := new(bytes.Buffer)
			buf.ReadFrom(*reader)
			expected := diffWithGoldenFile(golden, *update, *buf)
			assert.Equal(t, expected, buf)
		}
	})

	t.Run("Diane parser (JSON output)", func(t *testing.T) {
		var golden = filepath.Join("testData", "expectedDiane.json")
		var testData = filepath.Join("testData", "dianeTestData.txt")
		marshal.TestParserOutput(t, Parser, marshal.NewCache(), testData, golden, *update)
	})
}

func diffWithGoldenFile(filename string, updateGoldenFile bool, cmdOutput bytes.Buffer) []byte {

	if updateGoldenFile {
		ioutil.WriteFile(filename, cmdOutput.Bytes(), 0644)
	}
	expected, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	return expected
}
