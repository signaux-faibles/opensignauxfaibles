package diane

import (
	"bytes"
	"flag"
	"io/ioutil"
	"log"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/base"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/marshal"
	"github.com/stretchr/testify/assert"
)

var update = flag.Bool("update", false, "Update the expected test values in golden file")

func TestDiane(t *testing.T) {

	t.Run("Diane parser (JSON output)", func(t *testing.T) {
		var golden = filepath.Join("testData", "expectedDiane.json")
		var testData = filepath.Join("testData", "dianeTestData.txt")
		marshal.TestParserTupleOutput(t, Parser, base.NewCache(), "diane", testData, golden, *update)
	})

	t.Run("Diane converter (CSV output)", func(t *testing.T) {
		var golden = filepath.Join("testData", "expectedDianeConvert.csv")
		var testData = filepath.Join("testData", "dianeTestData.txt")

		cmd := exec.Command("bash", "./convert_diane.sh", testData)
		var cmdOutput bytes.Buffer
		cmd.Stdout = &cmdOutput
		// var cmdError bytes.Buffer
		// cmd.Stderr = &cmdError
		err := cmd.Run()
		if err != nil {
			log.Fatal(err)
		}
		expectedOutput := diffWithGoldenFile(golden, *update, cmdOutput)
		assert.Equal(t, string(expectedOutput), cmdOutput.String())
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
