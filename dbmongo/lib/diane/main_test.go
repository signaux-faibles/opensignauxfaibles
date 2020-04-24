package diane

import (
	"flag"
	"path/filepath"
	"testing"

	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/engine"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/marshal"
)

var update = flag.Bool("update", false, "Update the expected test values in golden file")

func TestDiane(t *testing.T) {
	var golden = filepath.Join("testData", "expectedDiane.json")
	var testData = filepath.Join("testData", "dianeTestData.csv")
	marshal.TestParserTupleOutput(t, Parser, engine.NewCache(), "diane", testData, golden, *update)
}
