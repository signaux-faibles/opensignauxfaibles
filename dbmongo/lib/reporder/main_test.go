package reporder

import (
	"flag"
	"opensignauxfaibles/dbmongo/lib/engine"
	"opensignauxfaibles/dbmongo/lib/marshal"
	"path/filepath"
	"testing"
)

var update = flag.Bool("update", false, "Update the expected test values in golden file")

func TestReporder(t *testing.T) {
	var golden = filepath.Join("testData", "expectedReporderMD5.csv")
	var testData = filepath.Join("testData", "reporderTestData.csv")
	marshal.TestParserTupleOutput(t, Parser, engine.NewCache(), "reporder", testData, golden, *update)
}
