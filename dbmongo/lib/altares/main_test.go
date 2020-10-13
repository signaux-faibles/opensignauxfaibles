package altares

import (
	"flag"
	"path/filepath"
	"testing"

	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/marshal"
)

var update = flag.Bool("update", false, "Update the expected test values in golden file")

var testData = filepath.Join("testData", "altaresTestData.csv") // ce fichier définit 3 entreprises: 000111222, 000111223 et 000111224

func TestAltaresOutput(t *testing.T) {
	var golden = filepath.Join("testData", "expectedAltaresOutput.json")
	marshal.TestParserOutput(t, Parser, marshal.NewCache(), "altares", testData, golden, *update)
}
