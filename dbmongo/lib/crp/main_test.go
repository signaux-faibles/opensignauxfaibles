package crp

import (
	"flag"
	"path/filepath"
	"testing"

	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/marshal"
)

var update = flag.Bool("update", false, "Update the expected test values in golden file")

func TestCrpOutput(t *testing.T) {
	var testData = filepath.Join("testData", "crpTestData.csv")
	var golden = filepath.Join("testData", "expectedCrpOutput.json")
	marshal.TestParserOutput(t, Parser, marshal.NewCache(), "crp", testData, golden, *update)
}
