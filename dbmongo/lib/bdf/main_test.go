package bdf

import (
	"flag"
	"path/filepath"
	"testing"

	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/marshal"
)

var update = flag.Bool("update", false, "Update the expected test values in golden file")

var testData = filepath.Join("testData", "bdfTestData.csv")

func TestBdf(t *testing.T) {
	var golden = filepath.Join("testData", "expectedBdf.json")
	marshal.TestParserTupleOutput(t, Parser, marshal.NewCache(), "bdf", testData, golden, *update)
}

func TestBdfOutput(t *testing.T) {
	var golden = filepath.Join("testData", "expectedBdfOutput.json")
	marshal.TestParserOutput(t, Parser, marshal.NewCache(), "bdf", testData, golden, *update)
}
