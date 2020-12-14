package apconso

import (
	"flag"
	"path/filepath"
	"testing"

	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/marshal"
)

var update = flag.Bool("update", false, "Update the expected test values in golden file")

func TestApconso(t *testing.T) {
	var golden = filepath.Join("testData", "expectedApconso.json")
	var testData = filepath.Join("testData", "apconsoTestData.csv")
	marshal.TestParserOutput(t, Parser, marshal.NewCache(), testData, golden, *update)
}
