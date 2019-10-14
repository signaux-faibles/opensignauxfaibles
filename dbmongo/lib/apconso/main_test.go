package apconso

import (
	"flag"
	"path/filepath"
	"testing"

	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/engine"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/marshal"
)

var update = flag.Bool("update", false, "Update the expected test values in golden file")

func TestApconso(t *testing.T) {
	var golden = filepath.Join("testData", "expectedApconsoMD5.csv")
	var testData = filepath.Join("testData", "apconsoTestData.csv")
	marshal.TestParserTupleOutput(t, Parser, engine.NewCache(), "apconso", testData, golden, *update)
}
