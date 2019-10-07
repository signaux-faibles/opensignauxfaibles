package sirene_ul

import (
	"flag"
	"opensignauxfaibles/dbmongo/lib/engine"
	"opensignauxfaibles/dbmongo/lib/marshal"
	"path/filepath"
	"testing"
)

var update = flag.Bool("update", false, "Update the expected test values in golden file")

var golden = filepath.Join("testData", "expectedMD5.csv")
var testData = filepath.Join("testData", "sireneULTestData.csv")

func TestSirene(t *testing.T) {
	marshal.TestParserTupleOutput(t, Parser, engine.NewCache(), "sirene_ul", testData, golden, *update)
}
