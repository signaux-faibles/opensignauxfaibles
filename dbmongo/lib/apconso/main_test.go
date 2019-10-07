package apconso

import (
	"flag"
	"opensignauxfaibles/dbmongo/lib/engine"
	"opensignauxfaibles/dbmongo/lib/marshal"
	"path/filepath"
	"testing"
)

var update = flag.Bool("update", false, "Update the expected test values in golden file")

func TestApconso(t *testing.T) {
	var golden = filepath.Join("testData", "expectedApconsoMD5.csv")
	var testData = filepath.Join("testData", "apconsoTestData.csv")
	marshal.TestParserTupleOutput(t, Parser, engine.NewCache(), "apconso", testData, golden, *update)
}
