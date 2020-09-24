package bdf

import (
	"flag"
	"path/filepath"
	"testing"

	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/marshal"
)

var update = flag.Bool("update", false, "Update the expected test values in golden file")

var golden = filepath.Join("testData", "expectedBdf.json")
var testData = filepath.Join("testData", "bdfTestData.csv")

func TestBdf(t *testing.T) {
	marshal.TestParserTupleOutput(t, Parser, marshal.NewCache(), "bdf", testData, golden, *update)
}
