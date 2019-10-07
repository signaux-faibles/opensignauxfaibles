package diane

import (
	"flag"
	"opensignauxfaibles/dbmongo/lib/engine"
	"opensignauxfaibles/dbmongo/lib/marshal"
	"path/filepath"
	"testing"
)

var update = flag.Bool("update", false, "Update the expected test values in golden file")

func TestDiane(t *testing.T) {
	var golden = filepath.Join("testData", "expectedDianeMD5.csv")
	var testData = filepath.Join("testData", "dianeTestData.csv")
	marshal.TestParserTupleOutput(t, Parser, engine.NewCache(), "diane", testData, golden, *update)
}
