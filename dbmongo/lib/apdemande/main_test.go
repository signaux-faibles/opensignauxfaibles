package apdemande

import (
	"flag"
	"opensignauxfaibles/dbmongo/lib/engine"
	"opensignauxfaibles/dbmongo/lib/marshal"
	"path/filepath"
	"testing"
)

var update = flag.Bool("update", false, "Update the expected test values in golden file")

func TestApdemande(t *testing.T) {
	var golden = filepath.Join("testData", "expectedApdemandeMD5.csv")
	var testData = filepath.Join("testData", "apdemandeTestData.csv")
	marshal.TestParserTupleOutput(t, Parser, engine.NewCache(), "apdemande", testData, golden, *update)
}
