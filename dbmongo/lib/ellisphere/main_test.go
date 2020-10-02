package apdemande

import (
	"flag"
	"path/filepath"
	"testing"

	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/marshal"
)

var update = flag.Bool("update", false, "Update the expected test values in golden file")

func TestApdemande(t *testing.T) {
	var golden = filepath.Join("testData", "expectedEllisphere.json")
	var testData = filepath.Join("testData", "ellisphereTestData.xlsx")
	marshal.TestParserTupleOutput(t, Parser, marshal.NewCache(), "ellisphere", testData, golden, *update)
}
