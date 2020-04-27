package sireneul

import (
	"flag"
	"path/filepath"
	"testing"

	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/engine"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/marshal"
)

var update = flag.Bool("update", false, "Update the expected test values in golden file")

var golden = filepath.Join("testData", "expectedSireneUL.json")
var testData = filepath.Join("testData", "sireneULTestData.csv")

func TestSirene(t *testing.T) {
	marshal.TestParserTupleOutput(t, Parser, engine.NewCache(), "sirene_ul", testData, golden, *update)
}
