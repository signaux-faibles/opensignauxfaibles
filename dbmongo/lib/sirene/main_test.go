package sirene

import (
	"flag"
	"path/filepath"
	"testing"

	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/marshal"
)

var update = flag.Bool("update", false, "Update the expected test values in golden file")

var golden = filepath.Join("testData", "expectedSirene.json")
var testData = filepath.Join("testData", "sireneTestData.csv")

func TestSirene(t *testing.T) {
	marshal.TestParserOutput(t, Parser, marshal.NewCache(), "sirene", testData, golden, *update)
}
