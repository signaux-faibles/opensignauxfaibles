package altares

import (
	"flag"
	"path/filepath"
	"testing"

	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/marshal"
)

var update = flag.Bool("update", false, "Update the expected test values in golden file")

func TestAltaresOutput(t *testing.T) {
	var testData = filepath.Join("testData", "altaresTestData.csv")
	var golden = filepath.Join("testData", "expectedAltaresOutput.json")
	marshal.TestParserOutput(t, Parser, marshal.NewCache(), "altares", testData, golden, *update)
}
